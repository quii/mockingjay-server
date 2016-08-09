package mockingjay

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

const debugModeOff = false

type mjLogger interface {
	Println(...interface{})
}

// Server allows you to configure a HTTP server for a slice of fake endpoints
type Server struct {
	Endpoints      []FakeEndpoint
	requests       []Request
	requestMatcher func(a, b Request, endpointName string) bool
	logger         mjLogger
}

// NewServer creates a new Server instance
func NewServer(endpoints []FakeEndpoint, debugMode bool) *Server {
	s := new(Server)
	s.Endpoints = endpoints
	s.requests = make([]Request, 0)

	if debugMode {
		s.logger = log.New(os.Stdout, "mocking-jay", log.Ldate|log.Ltime)
	} else {
		s.logger = log.New(ioutil.Discard, "mocking-jay", log.Ldate|log.Ltime)
	}

	s.requestMatcher = func(a, b Request, endpointName string) bool {
		return requestMatches(a, b, endpointName, s.logger)
	}

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

		s.logger.Println("Trying to match a request")
		s.logger.Println(mjRequest.String())

		s.requests = append(s.requests, mjRequest)

		writeToHTTP(s.getResponse(mjRequest), w)
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
		if s.requestMatcher(endpoint.Request, r, endpoint.Name) {
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
