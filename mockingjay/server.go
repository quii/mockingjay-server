package mockingjay

import (
	"io/ioutil"
	"net/http"
)

// Server allows you to configure a request and a corresponding response. It implements http.ServeHTTP so you can bind it to a port.
type Server struct {
	endpoints []FakeEndpoint
}

// NewServer creates a new Server instance
func NewServer(endpoints []FakeEndpoint) *Server {
	s := new(Server)
	s.endpoints = endpoints
	return s
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	cannedResponse := s.getResponse(r)
	for name, value := range cannedResponse.Headers {
		w.Header().Set(name, value)
	}
	w.WriteHeader(cannedResponse.Code)
	w.Write([]byte(cannedResponse.Body))
}

func (s *Server) getResponse(r *http.Request) *response {

	requestBody := ""

	if r.Body != nil {
		body, err := ioutil.ReadAll(r.Body)

		if err != nil {
			return newNotFound(r.Method, r.URL.String(), "", s.endpoints)
		}

		requestBody = string(body)
	}

	for _, endpoint := range s.endpoints {
		if requestMatches(endpoint.Request, r, requestBody) {
			return &endpoint.Response
		}
	}

	return newNotFound(r.Method, r.URL.String(), requestBody, s.endpoints)
}

func requestMatches(a request, b *http.Request, receivedBody string) bool {

	for key, value := range a.Headers {
		if b.Header[key] == nil || b.Header.Get(key) != value {
			return false
		}
	}

	bodyOk := a.Body == "*" || receivedBody == a.Body
	urlOk := a.URI == b.URL.String()
	methodOk := a.Method == b.Method

	return bodyOk && urlOk && methodOk
}
