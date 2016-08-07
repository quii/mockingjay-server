package mockingjay

import (
	"encoding/json"
	"log"
)

type response struct {
	Code    int
	Body    string
	Headers map[string]string
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
	return &response{404, string(j), nil}
}
