package mockingjay

import (
	"encoding/json"
	"errors"
	"fmt"

	"gopkg.in/yaml.v2"
)

type response struct {
	Code    int
	Body    string
	Headers map[string]string
}

func (r response) isValid() bool {
	return r.Code != 0
}

// FakeEndpoint represents the information required to listen to a particular request and respond to it
type FakeEndpoint struct {
	Name        string //A description of what this endpoint is.
	CDCDisabled bool   // When set to true it will not be included in the consumer driven contract tests against real server
	Request     Request
	Response    response
}

func (f *FakeEndpoint) String() string {
	return fmt.Sprintf(fakeEndpointStringerFormat, f.Name, f.Request)
}

const fakeEndpointStringerFormat = "%s (%s)"

var (
	invalidConfigError = errors.New("Config YAML structure is invalid")
)

// NewFakeEndpoints returns an array of Endpoints from a YAML byte array. Returns an error if YAML cannot be parsed
func NewFakeEndpoints(data []byte) (endpoints []FakeEndpoint, err error) {
	err = yaml.Unmarshal(data, &endpoints)

	if err != nil {
		return
	}

	for _, endPoint := range endpoints {
		if !endPoint.isValid() {
			return nil, invalidConfigError
		}
	}

	return
}

func (f FakeEndpoint) isValid() bool {
	return f.Request.isValid() && f.Response.isValid()
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
