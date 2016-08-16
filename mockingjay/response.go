package mockingjay

import (
	"encoding/json"
	"log"
	"net/http"
)

type response struct {
	Code    int
	Body    string
	Headers map[string]string `yaml:"headers,omitempty"`
}

func (r response) isValid() bool {
	return r.Code != 0
}

type notFoundResponse struct {
	Message            string
	Request            Request        `json:"RequestReceived"`
	EndpointsAvailable []FakeEndpoint `json:"EndpointsAvailable"`
}

func newNotFound(req Request, endpoints []FakeEndpoint) *response {
	notFound := notFoundResponse{
		"Endpoint not found",
		req,
		endpoints}
	j, err := json.Marshal(notFound)
	if err != nil {
		log.Println(err)
	}
	return &response{http.StatusNotFound, string(j), nil}
}

func writeToHTTP(res *response, w http.ResponseWriter) {
	for name, value := range res.Headers {
		w.Header().Set(name, value)
	}

	w.WriteHeader(res.Code)
	w.Write([]byte(res.Body))
}
