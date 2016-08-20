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
	if reqError := f.Request.errors(); reqError != nil {
		return reqError
	}
	if !f.Response.isValid() {
		return errResponseInvalid
	}
	return nil
}

var (
	errResponseInvalid = errors.New("Response is not configured correctly")
)

func errDuplicateRequestsError(duplicates []string) error {
	return fmt.Errorf("There were duplicated requests in YAML %v", duplicates)
}

// NewFakeEndpoints returns an array of Endpoints from a YAML byte array. Returns an error if YAML cannot be parsed or there are validation concerns
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

	if duplicates := findDuplicates(endpoints); len(duplicates) > 0 {
		return nil, errDuplicateRequestsError(duplicates)
	}

	return
}

func findDuplicates(endpoints []FakeEndpoint) []string {
	requests := make(map[string]int)

	for _, e := range endpoints {
		requests[e.Request.hash()] = requests[e.Request.hash()] + 1
	}

	var duplicates []string

	for k, v := range requests {
		if v > 1 {
			duplicates = append(duplicates, k)
		}
	}

	return duplicates
}
