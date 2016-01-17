package mockingjay

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
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
}

func (r Request) isValid() bool {
	return r.URI != "" && r.Method != ""
}

// AsHTTPRequest tries to create a http.Request from a given baseURL
func (r Request) AsHTTPRequest(baseURL string) (req *http.Request, err error) {

	req, err = http.NewRequest(r.Method, baseURL+r.URI, ioutil.NopCloser(bytes.NewBufferString(r.Body)))

	if err != nil {
		return
	}

	for headerName, headerValue := range r.Headers {
		req.Header.Add(headerName, headerValue)
	}

	return
}

// NewRequest creates a mockingjay request from a http request
func NewRequest(httpRequest *http.Request) (req Request) {
	req.URI = httpRequest.URL.String()
	req.Method = httpRequest.Method

	req.Headers = make(map[string]string)
	for header, values := range httpRequest.Header {
		req.Headers[header] = strings.Join(values, ",")
	}

	if httpRequest.Body != nil {
		reqBody, err := ioutil.ReadAll(httpRequest.Body)
		if err != nil {
			log.Println(err)
		} else {
			req.Body = string(reqBody)
		}
	}

	return
}

const stringerFormat = "%s %s"

func (r Request) String() string {
	return fmt.Sprintf(stringerFormat, r.Method, r.URI)
}

func (r Request) hash() string {
	return fmt.Sprintf("%v%v%v%v", r.URI, r.Method, r.Headers, r.Body)
}

func requestMatches(a, b Request) bool {

	headersOk := !(a.Headers != nil && !reflect.DeepEqual(a.Headers, b.Headers))
	bodyOk := a.Body == "*" || a.Body == b.Body
	urlOk := matchURI(a.URI, a.RegexURI, b.URI)
	methodOk := a.Method == b.Method

	return bodyOk && urlOk && methodOk && headersOk
}

func matchURI(serverURI string, serverRegex *RegexYAML, incomingURI string) bool {
	if serverURI == incomingURI {
		return true
	} else if serverRegex != nil {
		return serverRegex.MatchString(incomingURI)
	}
	return false
}
