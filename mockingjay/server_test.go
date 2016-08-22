package mockingjay

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v2"
	"io/ioutil"
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

	server := NewServer([]FakeEndpoint{endpoint, secondEndpoint}, debugModeOff, ioutil.Discard)

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
	server := NewServer([]FakeEndpoint{}, debugModeOff, ioutil.Discard)

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
	server := NewServer([]FakeEndpoint{}, debugModeOff, ioutil.Discard)

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

	server := NewServer([]FakeEndpoint{}, debugModeOff, ioutil.Discard)
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
	server := NewServer([]FakeEndpoint{config}, debugModeOff, ioutil.Discard)

	requestWithDifferentBody, _ := http.NewRequest("POST", testURL, strings.NewReader("This body isnt what we said but it should match"))
	responseReader := httptest.NewRecorder()

	server.ServeHTTP(responseReader, requestWithDifferentBody)

	assert.Len(t, server.requests, 1)
	assert.Equal(t, server.requests[0].Method, "POST")
}

func TestMJEndpoints(t *testing.T) {

	var newConfigBuffer bytes.Buffer

	mjReq := Request{URI: testURL, Method: "GET", Form: nil}
	endpoint := FakeEndpoint{testEndpointName, cdcDisabled, mjReq, response{http.StatusCreated, cannedResponse, nil}}
	server := NewServer([]FakeEndpoint{endpoint}, debugModeOff, &newConfigBuffer)

	t.Run("get endpoints", func(t *testing.T) {
		request, _ := http.NewRequest("GET", endpointsURL, nil)
		responseReader := httptest.NewRecorder()

		server.ServeHTTP(responseReader, request)

		assert.Equal(t, responseReader.Code, http.StatusOK)
		assert.Equal(t, responseReader.HeaderMap["Content-Type"][0], "application/json")

		var endpointResponse []FakeEndpoint
		err := json.Unmarshal(responseReader.Body.Bytes(), &endpointResponse)

		assert.Nil(t, err)
		assert.Equal(t, endpointResponse[0], endpoint, "The endpoint returned doesnt match what the server was set up with")

	})

	t.Run("update endpoints", func(t *testing.T) {
		newEndpointName := testEndpointName + "changed"
		updatedEndpoint := FakeEndpoint{newEndpointName, cdcDisabled, mjReq, response{http.StatusCreated, cannedResponse, nil}}

		editedBody, _ := json.Marshal([]FakeEndpoint{updatedEndpoint})
		updateRequest, _ := http.NewRequest(http.MethodPut, endpointsURL, bytes.NewReader(editedBody))
		updateResponseReader := httptest.NewRecorder()

		server.ServeHTTP(updateResponseReader, updateRequest)

		assert.Equal(t, http.StatusOK, updateResponseReader.Code, "Didnt work!", updateResponseReader.Body.String())
		assert.Equal(t, updateResponseReader.HeaderMap["Content-Type"][0], "application/json")

		var updatedEndpoints []FakeEndpoint

		err := yaml.Unmarshal(newConfigBuffer.Bytes(), &updatedEndpoints)

		assert.NoError(t, err)
		assert.Equal(t, updatedEndpoints, server.Endpoints)
		assert.Equal(t, newEndpointName, server.Endpoints[0].Name)
	})

	t.Run("rejects bad updates", func(t *testing.T) {
		badUpdate := `{"bad": "config"}`
		updateRequest, _ := http.NewRequest(http.MethodPut, endpointsURL, strings.NewReader(badUpdate))
		updateResponseReader := httptest.NewRecorder()

		server.ServeHTTP(updateResponseReader, updateRequest)

		assert.Equal(t, http.StatusBadRequest, updateResponseReader.Code, updateResponseReader.Body.String())
	})

	t.Run("rejects bad JSON", func(t *testing.T) {
		notJSON := `not json`
		updateRequest, _ := http.NewRequest(http.MethodPut, endpointsURL, strings.NewReader(notJSON))
		updateResponseReader := httptest.NewRecorder()

		server.ServeHTTP(updateResponseReader, updateRequest)

		assert.Equal(t, http.StatusBadRequest, updateResponseReader.Code, updateResponseReader.Body.String())
	})
}

func TestItCanCheckCompatability(t *testing.T) {

	mjReq := Request{URI: testURL, Method: "GET", Form: nil}
	endpoint := FakeEndpoint{testEndpointName, cdcDisabled, mjReq, response{http.StatusCreated, cannedResponse, nil}}
	server := NewServer([]FakeEndpoint{endpoint}, debugModeOff, ioutil.Discard)

	t.Run("cdc failing", func(t *testing.T){
		failingServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			http.Error(w, "Not found", http.StatusNotFound)
		}))
		defer failingServer.Close()
		request, _ := http.NewRequest("GET", checkcompatabilityURL+"?url="+failingServer.URL, nil)
		failingResponseReader := httptest.NewRecorder()

		server.ServeHTTP(failingResponseReader, request)

		assert.Equal(t, failingResponseReader.Code, http.StatusOK)

		var result compatCheckResult

		err := json.Unmarshal(failingResponseReader.Body.Bytes(), &result)

		assert.NoError(t, err, "Shouldn't get an error parsing compatability result")
		assert.False(t, result.Passed, "Compatability check should be fail on failing server")
		assert.NotEmpty(t, result.Messages, "Should be some messages about failure")
	})

	t.Run("cdc passing", func(t *testing.T){
		passingServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			t.Log("Got request", r.URL.RequestURI())
			if r.URL.Path == endpoint.Request.URI {
				w.WriteHeader(endpoint.Response.Code)
				w.Write([]byte(endpoint.Response.Body))
			} else {
				http.Error(w, "Not found", http.StatusNotFound)
			}
		}))
		defer passingServer.Close()

		passingRequest, _ := http.NewRequest("GET", checkcompatabilityURL+"?url="+passingServer.URL, nil)
		passingResponseReader := httptest.NewRecorder()
		server.ServeHTTP(passingResponseReader, passingRequest)

		assert.Equal(t, passingResponseReader.Code, http.StatusOK)

		var result compatCheckResult

		err := json.Unmarshal(passingResponseReader.Body.Bytes(), &result)

		assert.NoError(t, err, "Shouldn't get an error parsing compatability result")
		assert.True(t, result.Passed, "Compatability check should be passing on compat server")

	})

}
