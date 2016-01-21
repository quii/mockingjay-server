package mockingjay

import (
	"errors"
	"fmt"

	"gopkg.in/yaml.v2"
)

// FakeEndpoint represents the information required to listen to a particular request and respond to it
type FakeEndpoint struct {
	Name        string //A description of what this endpoint is.
	CDCDisabled bool   // When set to true it will not be included in the consumer driven contract tests against real server
	Request     Request
	Response    response
}

const fakeEndpointStringerFormat = "%s (%s)"

func (f *FakeEndpoint) String() string {
	return fmt.Sprintf(fakeEndpointStringerFormat, f.Name, f.Request)
}

func (f FakeEndpoint) isValid() error {
	if reqError := f.Request.isValid(); reqError != nil {
		return reqError
	}
	if !f.Response.isValid() {
		return errResponseInvalid
	}
	return nil
}

var (
	errDuplicateRequestsError = errors.New("There were duplicated requests in YAML")
	errResponseInvalid        = errors.New("Response is not configured correctly")
)

// NewFakeEndpoints returns an array of Endpoints from a YAML byte array. Returns an error if YAML cannot be parsed
func NewFakeEndpoints(data []byte) (endpoints []FakeEndpoint, err error) {
	err = yaml.Unmarshal(data, &endpoints)

	if err != nil {
		return nil, fmt.Errorf(
			"The structure of the supplied YAML is wrong, please refer to https://github.com/quii/mockingjay-server for an example [%v]",
			err)
	}

	for _, endPoint := range endpoints {
		if endpointErr := endPoint.isValid(); endpointErr != nil {
			return nil, endpointErr
		}
	}

	if isDuplicates(endpoints) {
		return nil, errDuplicateRequestsError
	}

	return
}

func isDuplicates(endpoints []FakeEndpoint) bool {
	requests := make(map[string]bool)

	for _, e := range endpoints {
		requests[e.Request.hash()] = true
	}

	return len(requests) != len(endpoints)
}
