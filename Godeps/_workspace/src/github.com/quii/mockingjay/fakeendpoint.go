package mockingjay

import (
	"encoding/json"
	"errors"
	"fmt"
)

type response struct {
	Code    int
	Body    string
	Headers map[string]string
}

func (r response) isValid() bool {
	return r.Code != 0 && r.Body != ""
}

// FakeEndpoint represents the information required to listen to a particular request and respond to it
type FakeEndpoint struct {
	Name     string
	Request  request
	Response response
}

func (f *FakeEndpoint) String() string {
	return fmt.Sprintf(fakeEndpointStringerFormat, f.Name, f.Request)
}

const fakeEndpointStringerFormat = "%s (%s)"

// NewFakeEndpoints returns an array of Endpoints from a JSON string. Returns an error if JSON cannot be parsed
func NewFakeEndpoints(data string) ([]FakeEndpoint, error) {
	var config []FakeEndpoint
	err := json.Unmarshal([]byte(data), &config)

	if err != nil {
		return nil, err
	}

	for _, endPoint := range config {
		if !endPoint.isValid() {
			return nil, errors.New("JSON was invalid")
		}
	}

	return config, nil
}

func (f FakeEndpoint) isValid() bool {
	return f.Request.isValid() && f.Response.isValid()
}

type notFoundResponse struct {
	Message            string
	Request            request        `json:"Request received"`
	EndpointsAvailable []FakeEndpoint `json:"Endpoints available"`
}

func newNotFound(method string, url string, endpoints []FakeEndpoint) *response {
	notFound := notFoundResponse{"Endpoint not found", request{url, method, nil, ""}, endpoints}
	j, _ := json.Marshal(notFound)
	return &response{404, string(j), nil}
}
