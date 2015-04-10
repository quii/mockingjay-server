package mockingjay

import (
	"fmt"
	"net/http"
	"strings"
)

type request struct {
	URI     string
	Method  string
	Headers map[string]string
	Body    string
}

func (r request) isValid() bool {
	return r.URI != "" && r.Method != ""
}

// AsHTTPRequest tries to create a http.Request from a given baseURL
func (r request) AsHTTPRequest(baseURL string) (*http.Request, error) {

	var body *strings.Reader
	if r.Body != "" {
		body = strings.NewReader(r.Body)
	} else {
		body = strings.NewReader("")
	}

	request, err := http.NewRequest(r.Method, baseURL+r.URI, body)

	for headerName, headerValue := range r.Headers {
		request.Header.Add(headerName, headerValue)
	}

	return request, err
}

const stringerFormat = "%s %s"

func (r request) String() string {
	return fmt.Sprintf(stringerFormat, r.Method, r.URI)
}
