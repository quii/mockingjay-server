package mockingjay

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"reflect"
	"strings"
)

// Request is a simplified version of a http.Request
type Request struct {
	URI      string
	RegexURI *RegexYAML
	Method   string
	Headers  map[string]string
	Body     string
	Form     map[string]string
}

var (
	errBadRegex    = errors.New("A regex defined in the request does not pass against it's defined URI")
	errEmptyURI    = errors.New("Cannot have an empty URI")
	errEmptyMethod = errors.New("Cannot have an empty HTTP method")
)

func (r Request) errors() error {
	regexPassed := r.RegexURI == nil || r.RegexURI.MatchString(r.URI)
	if !regexPassed {
		return errBadRegex
	}

	if r.URI == "" {
		return errEmptyURI
	}

	if r.Method == "" {
		return errEmptyMethod
	}
	return nil
}

// AsHTTPRequest tries to create a http.Request from a given baseURL
func (r Request) AsHTTPRequest(baseURL string) (req *http.Request, err error) {

	//todo: Test me with the form stuff

	body := r.Body
	if r.Form != nil {
		form := url.Values{}
		for formKey, formValue := range r.Form {
			form.Add(formKey, formValue)
		}
		body = form.Encode()
	}

	req, err = http.NewRequest(r.Method, baseURL+r.URI, ioutil.NopCloser(bytes.NewBufferString(body)))

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

	//todo: Test me with the form stuff

	req.URI = httpRequest.URL.String()
	req.Method = httpRequest.Method

	req.Headers = make(map[string]string)
	for header, values := range httpRequest.Header {
		req.Headers[header] = strings.Join(values, ",")
	}

	err := httpRequest.ParseForm()

	if httpRequest.Body != nil {
		reqBody, err := ioutil.ReadAll(httpRequest.Body)
		if err != nil {
			log.Println(err)
		} else {
			req.Body = string(reqBody)
		}
	}

	if err == nil {
		req.Form = make(map[string]string)
		for key, val := range httpRequest.PostForm {
			log.Println(val)
			req.Form[key] = val[0] // bit naughty
		}
	} else {
		log.Println("Problem parsing http form", err)
	}

	return
}

const stringerFormat = "%s %s"

func (r Request) String() string {
	return fmt.Sprintf(stringerFormat, r.Method, r.URI)
}

func (r Request) hash() string {
	return fmt.Sprintf("URI: %v | METHOD: %v | HEADERS: %v | BODY: %v", r.URI, r.Method, r.Headers, r.Body)
}

func requestMatches(expected, incoming Request) bool {

	headersOk := matchHeaders(expected.Headers, incoming.Headers)
	bodyOk := expected.Body == "*" || expected.Body == incoming.Body || matchJSON(expected.Body, incoming.Body)
	urlOk := matchURI(expected.URI, expected.RegexURI, incoming.URI)
	methodOk := expected.Method == incoming.Method
	formOK := matchForm(expected.Form, incoming.Form)

	return bodyOk && urlOk && methodOk && headersOk && formOK
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

func matchURI(serverURI string, serverRegex *RegexYAML, incomingURI string) bool {
	if serverURI == incomingURI {
		return true
	} else if serverRegex != nil {
		return serverRegex.MatchString(incomingURI)
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
