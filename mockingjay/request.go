package mockingjay

import (
	"bytes"
	"fmt"
	"io"
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

	urlParsed, err := url.Parse(baseURL)

	if err != nil {
		return
	}

	req, err = http.NewRequest(r.Method, baseURL, nil)

	if err != nil {
		return
	}

	req.URL = &url.URL{
		Scheme: urlParsed.Scheme,
		Host:   urlParsed.Host,
		Opaque: fmt.Sprintf("//%s%s", urlParsed.Host, r.URI),
	}

	req.Body = nopCloser{bytes.NewBufferString(r.Body)}

	for headerName, headerValue := range r.Headers {
		req.Header.Add(headerName, headerValue)
	}

	return
}

const stringerFormat = "%s %s"

func (r request) String() string {
	return fmt.Sprintf(stringerFormat, r.Method, r.URI)
}

type nopCloser struct {
	io.Reader
}

func (nopCloser) Close() error { return nil }
