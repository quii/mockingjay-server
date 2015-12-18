package mockingjay

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"reflect"
)

// Server allows you to configure a request and a corresponding response. It implements http.ServeHTTP so you can bind it to a port.
type Server struct {
	endpoints []FakeEndpoint
	requests  []Request
}

// NewServer creates a new Server instance
func NewServer(endpoints []FakeEndpoint) *Server {
	s := new(Server)
	s.endpoints = endpoints
	s.requests = make([]Request, 0)
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

	for _, endpoint := range s.endpoints {
		if requestMatches(endpoint.Request, r) {
			return &endpoint.Response
		}
	}

	return newNotFound(r.Method, r.URI, r.Body, r.Headers, s.endpoints)
}

func requestMatches(a Request, b Request) bool {

	headersOk := !(a.Headers != nil && !reflect.DeepEqual(a.Headers, b.Headers))
	bodyOk := a.Body == "*" || a.Body == b.Body
	urlOk := a.URI == b.URI
	methodOk := a.Method == b.Method

	return bodyOk && urlOk && methodOk && headersOk
}

func (s *Server) serveEndpoints(w http.ResponseWriter) {
	endpointsBody, err := json.Marshal(s.endpoints)

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

	s.endpoints = append(s.endpoints, newEndpoint)

	w.WriteHeader(http.StatusCreated)
}
