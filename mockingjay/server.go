package mockingjay

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
)

// Server allows you to configure a request and a corresponding response. It implements http.ServeHTTP so you can bind it to a port.
type Server struct {
	endpoints []FakeEndpoint
	requests  []http.Request
}

// NewServer creates a new Server instance
func NewServer(endpoints []FakeEndpoint) *Server {
	s := new(Server)
	s.endpoints = endpoints
	s.requests = make([]http.Request, 0)
	return s
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	if r.RequestURI == "/requests" {
		s.returnRequests(w)
	} else {
		s.requests = append(s.requests, *r)
		cannedResponse := s.getResponse(r)
		for name, value := range cannedResponse.Headers {
			w.Header().Set(name, value)
		}
		w.WriteHeader(cannedResponse.Code)
		w.Write([]byte(cannedResponse.Body))
	}
}

func (s *Server) returnRequests(w http.ResponseWriter) {
	payload, err := json.Marshal(s.requests)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	} else {
		w.WriteHeader(http.StatusOK)
		w.Write(payload)
	}
}

func (s *Server) getResponse(r *http.Request) *response {

	requestHeaders := make(map[string]string)
	for key, val := range r.Header {
		requestHeaders[key] = strings.Join(val, ",")
	}

	requestBody := ""

	if r.Body != nil {
		body, err := ioutil.ReadAll(r.Body)

		if err != nil {
			return newNotFound(r.Method, r.URL.String(), "", requestHeaders, s.endpoints)
		}

		requestBody = string(body)
	}

	for _, endpoint := range s.endpoints {
		if requestMatches(endpoint.Request, r, requestBody) {
			return &endpoint.Response
		}
	}

	return newNotFound(r.Method, r.URL.String(), requestBody, requestHeaders, s.endpoints)
}

func requestMatches(a request, b *http.Request, receivedBody string) bool {

	for key, value := range a.Headers {
		if b.Header[key] == nil || strings.ToLower(b.Header.Get(key)) != strings.ToLower(value) {
			return false
		}
	}

	aURL, err := url.QueryUnescape(a.URI)
	bURL, err := url.QueryUnescape(b.URL.String())

	if err != nil {
		log.Fatalf("Unescaping the query string failed horribly, crashing and burning %v", err)
	}

	bodyOk := a.Body == "*" || receivedBody == a.Body
	urlOk := aURL == bURL
	methodOk := a.Method == b.Method

	return bodyOk && urlOk && methodOk
}
