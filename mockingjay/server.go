package mockingjay

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

// Server allows you to configure a HTTP server for a slice of fake endpoints
type Server struct {
	Endpoints      []FakeEndpoint
	requests       []Request
	requestMatcher func(a, b Request) bool
}

// NewServer creates a new Server instance
func NewServer(endpoints []FakeEndpoint) *Server {
	s := new(Server)
	s.Endpoints = endpoints
	s.requests = make([]Request, 0)
	s.requestMatcher = requestMatches
	return s
}

const requestsURL = "/requests"
const endpointsURL = "/mj-endpoints"
const newEndpointURL = "/mj-new-endpoint"

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.URL.String() {
	case endpointsURL:
		s.serveEndpoints(w)
	case newEndpointURL:
		s.createEndpoint(w, r)
	case requestsURL:
		s.listAvailableRequests(w)
	default:
		mjRequest := NewRequest(r)
		log.Println("Got request", mjRequest)
		s.requests = append(s.requests, mjRequest)

		cannedResponse := s.getResponse(mjRequest)

		for name, value := range cannedResponse.Headers {
			w.Header().Set(name, value)
		}

		w.WriteHeader(cannedResponse.Code)
		w.Write([]byte(cannedResponse.Body))
	}
}

func (s *Server) listAvailableRequests(w http.ResponseWriter) {
	payload, err := json.Marshal(s.requests)

	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	} else {
		w.WriteHeader(http.StatusOK)
		w.Write(payload)
	}
}

func (s *Server) getResponse(r Request) *response {

	for _, endpoint := range s.Endpoints {
		if s.requestMatcher(endpoint.Request, r) {
			return &endpoint.Response
		}
	}

	return newNotFound(r, s.Endpoints)
}

func (s *Server) serveEndpoints(w http.ResponseWriter) {
	endpointsBody, err := json.Marshal(s.Endpoints)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Write(endpointsBody)
	w.Header().Set("Content-Type", "application/json")
}

func (s *Server) createEndpoint(w http.ResponseWriter, r *http.Request) {
	var newEndpoint FakeEndpoint

	defer r.Body.Close()

	endpointBody, err := ioutil.ReadAll(r.Body)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	err = json.Unmarshal(endpointBody, &newEndpoint)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	s.Endpoints = append(s.Endpoints, newEndpoint)

	w.WriteHeader(http.StatusCreated)
}
