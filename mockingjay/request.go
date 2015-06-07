package mockingjay

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
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
func (r request) AsHTTPRequest(baseURL string) (req *http.Request, err error) {

	req, err = http.NewRequest(r.Method, baseURL, ioutil.NopCloser(bytes.NewBufferString(r.Body)))

	if err != nil {
		return
	}

	req.URL = &url.URL{
		Scheme: req.URL.Scheme,
		Host:   req.URL.Host,
		Opaque: fmt.Sprintf("//%s%s", req.URL.Host, r.URI),
	}

	for headerName, headerValue := range r.Headers {
		req.Header.Add(headerName, headerValue)
	}

	return
}

const stringerFormat = "%s %s"

func (r request) String() string {
	return fmt.Sprintf(stringerFormat, r.Method, r.URI)
}
