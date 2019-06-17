package mockingjay

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"reflect"
	"regexp"
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
	return fmt.Sprintf(fakeEndpointStringerFormat, f.Name, f.Request.String())
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
	return fmt.Errorf("There were duplicated requests found %v", duplicates)
}

// NewFakeEndpoints returns an array of Endpoints from a YAML byte array. Returns an error if YAML cannot be parsed or there are validation concerns
func NewFakeEndpoints(data io.ReadCloser) (endpoints []FakeEndpoint, err error) {
	defer data.Close()
	err = mjYAMLDecoder(data).Decode(&endpoints)

	if err != nil {
		return nil, fmt.Errorf(
			"The structure of the supplied config is wrong, please refer to https://github.com/quii/mockingjay-server for an example [%v]",
			err)
	}

	err = validateEndpoints(endpoints)

	return
}

// NewFakeEndpointsFromJSON returns an array of Endpoints from a JSON byte array. Returns an error if JSON cannot be parsed or there are validation concerns
func NewFakeEndpointsFromJSON(data io.ReadCloser) (endpoints []FakeEndpoint, err error) {
	defer data.Close()
	err = json.NewDecoder(data).Decode(&endpoints)

	if jsonErr, isJsonErr := err.(*json.InvalidUnmarshalError); isJsonErr {
		return nil, fmt.Errorf("problem unmarshalling JSON %v %v", jsonErr.Type, jsonErr.Error())
	}

	if err != nil {
		return
	}

	err = validateEndpoints(endpoints)

	return
}

func validateEndpoints(endpoints []FakeEndpoint) error {
	for _, endPoint := range endpoints {
		if endpointErr := endPoint.isValid(); endpointErr != nil {
			return endpointErr
		}
	}

	if duplicates := findDuplicates(endpoints); len(duplicates) > 0 {
		return errDuplicateRequestsError(duplicates)
	}

	return nil
}

func findDuplicates(endpoints []FakeEndpoint) []string {
	request := make(map[string]string)
	requestHashes := make(map[string]int)

	for _, e := range endpoints {
		request[e.Request.Hash()] = e.Request.String()
		requestHashes[e.Request.Hash()] = requestHashes[e.Request.Hash()] + 1
	}

	var duplicates []string

	for k, v := range requestHashes {
		if v > 1 {
			duplicates = append(duplicates, request[k])
		}
	}

	return duplicates
}

// Generate creates a random endpoint typically used for testing
func (f FakeEndpoint) Generate(rand *rand.Rand, size int) reflect.Value {
	randomMethod := httpMethods[rand.Intn(len(httpMethods))]

	//todo: Creation of random URL and corresponding regex is a bit naff, needs improvement
	uri, _ := randomURL(uint16(rand.Intn(100) + 10))
	quotedURI := regexp.QuoteMeta(uri)
	regexURI, err := regexp.Compile(quotedURI)
	if err != nil {
		return reflect.ValueOf(err)
	}

	req := Request{
		Method:   randomMethod,
		URI:      uri,
		RegexURI: &RegexField{regexURI},
	}

	res := response{
		Code: rand.Intn(599-300) + 300,
	}

	return reflect.ValueOf(FakeEndpoint{
		Name:     "Generated",
		Request:  req,
		Response: res,
	})
}
