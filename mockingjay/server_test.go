package mockingjay

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	cdcDisabled = false
)

const cannedResponse = "world"
const cannedResponse2 = "worldy world world"
const testURL = "/hello"
const testURL2 = "/hello2"
const testEndpointName = "Test 123"

func TestItReturnsCannedResponses(t *testing.T) {
	responseHeaders := make(map[string]string)
	responseHeaders["Content-Type"] = "application/json"

	endpoint := FakeEndpoint{
		Name:        "Fake 1",
		CDCDisabled: cdcDisabled,
		Request: Request{
			URI:     testURL,
			Method:  "GET",
			Headers: nil,
			Body:    "",
		},
		Response: response{
			Code:    http.StatusCreated,
			Body:    cannedResponse,
			Headers: responseHeaders,
		},
	}
	secondEndpoint := FakeEndpoint{
		Name:        "Fake 2",
		CDCDisabled: cdcDisabled,
		Request: Request{
			URI:     testURL2,
			Method:  "GET",
			Headers: nil,
			Body:    "",
		},
		Response: response{
			Code:    http.StatusCreated,
			Body:    cannedResponse2,
			Headers: nil,
		},
	}

	server := NewServer([]FakeEndpoint{endpoint, secondEndpoint}, debugModeOff)

	request, _ := http.NewRequest("GET", testURL, nil)
	responseReader := httptest.NewRecorder()

	server.ServeHTTP(responseReader, request)

	assert.Equal(t, responseReader.Code, http.StatusCreated)
	assert.Equal(t, responseReader.Body.String(), cannedResponse)
	assert.Equal(t, responseReader.Header().Get("Content-Type"), "application/json")

	responseReader = httptest.NewRecorder()
	requestTwo, _ := http.NewRequest("GET", testURL2, nil)

	server.ServeHTTP(responseReader, requestTwo)

	assert.Equal(t, responseReader.Body.String(), cannedResponse2)
}

func TestItCanCreateNewEndpointsOverHTTP(t *testing.T) {
	server := NewServer([]FakeEndpoint{}, debugModeOff)

	newEndpoint := FakeEndpoint{
		Name: "New endpoint",
		Request: Request{
			URI:    "/foo",
			Method: "GET",
			Body:   "Blah blah",
		},
		Response: response{
			Code: 200,
			Body: "SUPER",
		},
	}

	newEndpointJSON, _ := json.Marshal(newEndpoint)
	req, _ := http.NewRequest("POST", newEndpointURL, bytes.NewReader(newEndpointJSON))
	responseReader := httptest.NewRecorder()

	server.ServeHTTP(responseReader, req)

	assert.Equal(t, responseReader.Code, http.StatusCreated)
	assert.Len(t, server.Endpoints, 1)
	assert.Equal(t, server.Endpoints[0], newEndpoint)
}

func TestItReturnsBadRequestWhenMakingABadNewEndpoint(t *testing.T) {
	server := NewServer([]FakeEndpoint{}, debugModeOff)

	badBody := []byte("blah")
	req, _ := http.NewRequest("POST", newEndpointURL, bytes.NewReader(badBody))
	responseReader := httptest.NewRecorder()

	server.ServeHTTP(responseReader, req)

	assert.Equal(t, responseReader.Code, http.StatusBadRequest)
}

func TestItReturns404WhenRequestCannotBeMatched(t *testing.T) {
	alwaysNotMatching := func(a, b Request, endpointName string) bool {
		return false
	}

	server := NewServer([]FakeEndpoint{}, debugModeOff)
	server.requestMatcher = alwaysNotMatching

	req, _ := http.NewRequest("POST", "doesnt-matter", nil)

	responseReader := httptest.NewRecorder()

	server.ServeHTTP(responseReader, req)

	assert.Equal(t, responseReader.Code, http.StatusNotFound)
}

func TestItRecordsIncomingRequests(t *testing.T) {
	wildcardBody := "*"
	expectedStatus := http.StatusOK

	mjReq := Request{URI: testURL, Method: "POST", Body: wildcardBody, Form: nil}
	config := FakeEndpoint{testEndpointName, cdcDisabled, mjReq, response{expectedStatus, "", nil}}
	server := NewServer([]FakeEndpoint{config}, debugModeOff)

	requestWithDifferentBody, _ := http.NewRequest("POST", testURL, strings.NewReader("This body isnt what we said but it should match"))
	responseReader := httptest.NewRecorder()

	server.ServeHTTP(responseReader, requestWithDifferentBody)

	assert.Len(t, server.requests, 1)
	assert.Equal(t, server.requests[0].Method, "POST")
}

func TestItReturnsListOfEndpointsAndUpdates(t *testing.T) {
	mjReq := Request{URI: testURL, Method: "GET", Form: nil}
	endpoint := FakeEndpoint{testEndpointName, cdcDisabled, mjReq, response{http.StatusCreated, cannedResponse, nil}}
	server := NewServer([]FakeEndpoint{endpoint}, debugModeOff)

	request, _ := http.NewRequest("GET", endpointsURL, nil)
	responseReader := httptest.NewRecorder()

	server.ServeHTTP(responseReader, request)

	assert.Equal(t, responseReader.Code, http.StatusOK)
	assert.Equal(t, responseReader.HeaderMap["Content-Type"][0], "application/json")

	var endpointResponse []FakeEndpoint
	err := json.Unmarshal(responseReader.Body.Bytes(), &endpointResponse)

	assert.Nil(t, err)
	assert.Equal(t, endpointResponse[0], endpoint, "The endpoint returned doesnt match what the server was set up with")

	updateBody := `[{
	"Name": "Test endpoint updated",
	"CDCDisabled": false,
	"Request": {
		"URI": "/hello",
		"RegexURI": null,
		"Method": null,
		"Body": ""
	},
	"Response": {
		"Code": 200,
		"Body": "{\"message\": \"hello, world\"}"
	}
}, {
	"Name": "New endpoint",
	"CDCDisabled": false,
	"Request": {
		"URI": "/world",
		"RegexURI": null,
		"Method": null,
		"Body": ""
	},
	"Response": {
		"Code": 200,
		"Body": "hello, world"
	}

}]`
	updateRequest, _ := http.NewRequest(http.MethodPut, endpointsURL, strings.NewReader(updateBody))
	updateResponseReader := httptest.NewRecorder()

	server.ServeHTTP(updateResponseReader, updateRequest)

	assert.Equal(t, updateResponseReader.Code, http.StatusOK)
	assert.Equal(t, updateResponseReader.HeaderMap["Content-Type"][0], "application/json")

	assert.Len(t, server.Endpoints, 2)
}
