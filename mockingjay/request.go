package mockingjay

import (
	"crypto/sha1"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/moul/http2curl"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"reflect"
	"strconv"
	"strings"
)

// Request is a simplified version of a http.Request
type Request struct {
	URI      string            `yaml:"uri"`
	RegexURI *RegexField       `yaml:"regexuri,omitempty" json:"RegexURI,omitempty"`
	Method   string            `yaml:"method"`
	Headers  map[string]string `yaml:"headers,omitempty"`
	Body     string            `yaml:"body,omitempty"`
	Form     map[string]string `yaml:"form,omitempty"`
}

var (
	errBadRegex    = errors.New("A regex defined in the request does not pass against it's defined URI")
	errEmptyURI    = errors.New("Cannot have an empty URI")
	errBadURI      = errors.New("URI should start with a slash")
	errEmptyMethod = errors.New("Cannot have an empty HTTP method")
	errBadHeaders  = errors.New("Headers are bad (no spaces please, each key must have a value")
	errBodyAndForm = errors.New("Cannot have a request with both a body and a form defined")
)

func (r Request) errors() error {
	regexPassed := r.RegexURI == nil || r.RegexURI.MatchString(r.URI)

	if !regexPassed {
		return errBadRegex
	}

	if r.URI == "" {
		return errEmptyURI
	}

	if r.URI[0] != '/' {
		return errBadURI
	}

	if r.Method == "" {
		return errEmptyMethod
	}

	if !httpHeadersValid(r.Headers) {
		return errBadHeaders
	}

	if r.Form != nil && len(r.Body) > 0 {
		return errBodyAndForm
	}

	return nil
}

// AsCURL returns a string which is the curl command to match the reqquest
func (r Request) AsCURL(baseURL string) (string, error) {
	asHTTPReq, err := r.AsHTTPRequest(baseURL)

	if err != nil {
		return "", err
	}

	curl, err := http2curl.GetCurlCommand(asHTTPReq)
	if err != nil {
		return "", err
	}

	return curl.String(), err
}

// AsHTTPRequest tries to create a http.Request from a given baseURL
func (r Request) AsHTTPRequest(baseURL string) (req *http.Request, err error) {

	body := r.Body
	if r.Form != nil {
		form := url.Values{}
		for formKey, formValue := range r.Form {
			form.Add(formKey, formValue)
		}
		body = form.Encode()
	}

	req, err = http.NewRequest(r.Method, baseURL+r.URI, strings.NewReader(body))

	if err != nil {
		return
	}

	for headerName, headerValue := range r.Headers {
		req.Header.Add(headerName, headerValue)
	}

	if r.Form != nil {
		req.Header.Add("content-type", "application/x-www-form-urlencoded")
	}

	return
}

// NewRequest creates a mockingjay request from a http request
func NewRequest(httpRequest *http.Request) (req Request) {

	req.URI = httpRequest.URL.String()
	req.Method = httpRequest.Method

	req.Headers = make(map[string]string)
	for header, values := range httpRequest.Header {
		req.Headers[strings.ToLower(header)] = strings.Join(values, ",")
	}

	if httpRequest.Header.Get("content-type") == "application/x-www-form-urlencoded" {
		err := httpRequest.ParseForm()
		if err == nil {
			req.Form = make(map[string]string)
			for key, val := range httpRequest.PostForm {
				req.Form[key] = val[0] // bit naughty
			}
		} else {
			log.Println("Problem parsing http form", err)
		}
	} else if httpRequest.Body != nil {
		reqBody, err := ioutil.ReadAll(httpRequest.Body)
		if err != nil {
			log.Println(err)
		} else {
			req.Body = string(reqBody)
		}
	}

	return
}

func (r Request) String() string {
	base := fmt.Sprintf("%s %s", r.Method, r.URI)

	if r.Headers != nil {
		headersAsString := stringifyMap(r.Headers)
		base = fmt.Sprintf("%s Headers: [%s]", base, headersAsString)
	}

	if r.Form != nil {
		headersAsString := stringifyMap(r.Form)
		base = fmt.Sprintf("%s Form: [%s]", base, headersAsString)
	}

	if len(r.Body) > 0 && len(r.Body) > 50 {
		base = fmt.Sprintf("%s Body: [%s...]", base, r.Body[:50])
	} else if len(r.Body) > 0 {
		base = fmt.Sprintf("%s Body: [%s]", base, r.Body)
	}

	return base
}

func (r Request) Hash() string {
	base := fmt.Sprintf("%s %s", r.Method, r.URI)

	if r.Headers != nil {
		headersAsString := stringifyMap(r.Headers)
		base = fmt.Sprintf("%s Headers: [%s]", base, headersAsString)
	}

	if r.Form != nil {
		headersAsString := stringifyMap(r.Form)
		base = fmt.Sprintf("%s Form: [%s]", base, headersAsString)
	}

	if len(r.Body) > 0 {
		base = fmt.Sprintf("%s Body: [%s]", base, r.Body)
	}

	return fmt.Sprintf("%x", sha1.Sum([]byte(base)))
}

func stringifyMap(m map[string]string) string {
	var headerStrings []string
	for k, v := range m {
		headerStrings = append(headerStrings, fmt.Sprintf("%s->%s", k, v))
	}
	return strings.Join(headerStrings, ", ")
}

func requestMatches(expected, incoming Request, endpointName string, logger mjLogger) bool {

	headersOk := matchHeaders(expected.Headers, incoming.Headers)
	bodyOk := expected.Body == "*" || expected.Body == incoming.Body || matchJSON(expected.Body, incoming.Body)
	urlOk := matchURI(expected.URI, expected.RegexURI, incoming.URI)
	methodOk := expected.Method == incoming.Method
	formOK := matchForm(expected.Form, incoming.Form)
	matches := bodyOk && urlOk && methodOk && headersOk && formOK

	if matches {
		logger.Println("Matched on", endpointName)
	} else {
		debugMsg := fmt.Sprintf(
			"%s | URL OK? %s | Body OK? %s | Method OK? %s | Form OK? %s | Headers OK? %s",
			fixedLengthString(20, endpointName),
			fixedLengthString(5, strconv.FormatBool(urlOk)),
			fixedLengthString(5, strconv.FormatBool(bodyOk)),
			fixedLengthString(5, strconv.FormatBool(methodOk)),
			fixedLengthString(5, strconv.FormatBool(formOK)),
			fixedLengthString(5, strconv.FormatBool(headersOk)),
		)

		logger.Println(debugMsg)
	}

	return matches
}

func matchForm(expected map[string]string, incoming map[string]string) bool {
	if expected == nil {
		return true
	}

	for expectedKey, expectedValue := range expected {
		if incoming[expectedKey] != expectedValue {
			return false
		}
	}

	return true
}

func matchJSON(a string, b string) bool {

	var aJSON map[string]interface{}
	var bJSON map[string]interface{}

	err := json.Unmarshal([]byte(a), &aJSON)
	err = json.Unmarshal([]byte(b), &bJSON)

	if err != nil {
		return false
	}

	return reflect.DeepEqual(aJSON, bJSON)
}

func matchHeaders(expected, incoming map[string]string) bool {
	incominglowercased := lowercaseMapKeys(incoming)
	expectedLowercased := lowercaseMapKeys(expected)

	for key, expectedValue := range expectedLowercased {
		if value, exists := incominglowercased[key]; exists {
			if value != expectedValue {
				return false
			}
		} else {
			return false
		}
	}
	return true
}

func matchURI(serverURI string, serverRegex *RegexField, incomingURI string) bool {
	if serverURI == incomingURI {
		return true
	} else if serverRegex != nil {
		return serverRegex.MatchString(incomingURI)
	} else {
		inURL, err := url.Parse(incomingURI)
		expectedURL, err := url.Parse(serverURI)

		if err != nil {
			return false
		}

		if reflect.DeepEqual(inURL.Query(), expectedURL.Query()) {
			return inURL.Path == expectedURL.Path
		}
	}
	return false
}

func lowercaseMapKeys(upperCasedMap map[string]string) map[string]string {
	lowerCasedMap := make(map[string]string)

	for key, value := range upperCasedMap {
		lowerCasedMap[strings.ToLower(key)] = value
	}

	return lowerCasedMap
}

func fixedLengthString(length int, str string) string {
	verb := fmt.Sprintf("%%%d.%ds", length, length)
	return fmt.Sprintf(verb, str)
}
