package mockingjay

import "encoding/json"

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

func newNotFound(method string, url string, body string, headers map[string]string, endpoints []FakeEndpoint) *response {
	notFound := notFoundResponse{"Endpoint not found", Request{url, method, headers, body}, endpoints}
	j, _ := json.Marshal(notFound)
	return &response{404, string(j), nil}
}
