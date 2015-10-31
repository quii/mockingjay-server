package mockingjay

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
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
		Request: request{
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
		Request: request{
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

	server := NewServer([]FakeEndpoint{endpoint, secondEndpoint})

	request, _ := http.NewRequest("GET", testURL, nil)
	responseReader := httptest.NewRecorder()

	server.ServeHTTP(responseReader, request)

	if responseReader.Code != http.StatusCreated {
		t.Error("Expected a 201 (status created)")
	}

	if responseReader.Body.String() != cannedResponse {
		t.Errorf("Expected canned response to be returned, got [%s] ", responseReader.Body.String())
	}

	if responseReader.Header().Get("Content-Type") != "application/json" {
		t.Errorf("Expected response header to be set")
	}

	responseReader = httptest.NewRecorder()
	requestTwo, _ := http.NewRequest("GET", testURL2, nil)

	server.ServeHTTP(responseReader, requestTwo)

	if responseReader.Body.String() != cannedResponse2 {
		t.Errorf("Expected second endpoint's canned response, got %s", responseReader.Body.String())
	}
}

func TestItReturns404WhenUriIsWrong(t *testing.T) {
	endpoint := FakeEndpoint{testEndpointName, cdcDisabled, request{testURL, "GET", nil, ""}, response{http.StatusCreated, cannedResponse, nil}}
	server := NewServer([]FakeEndpoint{endpoint})
	requestBody := "some body"

	request, _ := http.NewRequest("GET", "/bums", strings.NewReader(requestBody))
	responseReader := httptest.NewRecorder()

	server.ServeHTTP(responseReader, request)

	if responseReader.Code != http.StatusNotFound {
		t.Error("Expected to get a 404")
	}

	if !strings.Contains(responseReader.Body.String(), requestBody) {
		t.Errorf("Expected request body to be returned %v", responseReader.Body.String())
	}
}

func TestItReturns404WhenMethodIsWrong(t *testing.T) {
	endpoint := FakeEndpoint{testEndpointName, cdcDisabled, request{testURL, "GET", nil, ""}, response{http.StatusCreated, cannedResponse, nil}}
	server := NewServer([]FakeEndpoint{endpoint})

	request, _ := http.NewRequest("POST", "/hello", nil)
	request.Header.Set("content-type", "application/bob")

	responseReader := httptest.NewRecorder()

	server.ServeHTTP(responseReader, request)

	if responseReader.Code != http.StatusNotFound {
		t.Error("Expected to get a 404")
	}

	if !strings.Contains(responseReader.Body.String(), "application/bob") {
		t.Log(responseReader.Body.String())
		t.Error("Request headers were not added to the response info")
	}
}

func TestItDoesContentNegotiation(t *testing.T) {
	contentTypes := make(map[string]string)
	contentTypes["Content-Type"] = "application/json"

	endpoint := FakeEndpoint{testEndpointName, cdcDisabled, request{testURL, "GET", contentTypes, ""}, response{http.StatusCreated, cannedResponse, nil}}
	server := NewServer([]FakeEndpoint{endpoint})

	requestWithIncorrectHeaderValue, _ := http.NewRequest("GET", testURL, nil)
	requestWithIncorrectHeaderValue.Header.Add("Content-Type", "application/xml")
	responseReader := httptest.NewRecorder()

	server.ServeHTTP(responseReader, requestWithIncorrectHeaderValue)

	if responseReader.Code != http.StatusNotFound {
		t.Error("Expected to get a 404 because we didnt set a content type header when it was expected")
	}

	requestWithNoHeaderAtAll, _ := http.NewRequest("GET", testURL, nil)
	responseReader = httptest.NewRecorder()

	server.ServeHTTP(responseReader, requestWithNoHeaderAtAll)

	if responseReader.Code != http.StatusNotFound {
		t.Error("Expected to get a 404 because we didnt set a content type header when it was expected")
	}

	requestWithDifferentCasedHeader, _ := http.NewRequest("GET", testURL, nil)
	requestWithDifferentCasedHeader.Header.Add("Content-TYPE", "application/json")
	responseReader = httptest.NewRecorder()

	server.ServeHTTP(responseReader, requestWithDifferentCasedHeader)

	if responseReader.Code != http.StatusCreated {
		t.Errorf("Expected request to match even though the header name was differently cased : \n%d", responseReader.Code)
	}
}

func TestItSendsRequestBodies(t *testing.T) {
	body := "some body"
	expectedStatus := http.StatusInternalServerError

	endpoint := FakeEndpoint{testEndpointName, cdcDisabled, request{testURL, "POST", nil, body}, response{expectedStatus, "", nil}}
	server := NewServer([]FakeEndpoint{endpoint})

	requestWithoutBody, _ := http.NewRequest("POST", testURL, nil)
	responseReader := httptest.NewRecorder()

	server.ServeHTTP(responseReader, requestWithoutBody)

	if responseReader.Code != http.StatusNotFound {
		t.Error("Expected to get a 404 because we didnt set a request body when it was expected")
	}

	requestWithBody, _ := http.NewRequest("POST", testURL, strings.NewReader(body))
	responseReader = httptest.NewRecorder()

	server.ServeHTTP(responseReader, requestWithBody)

	if responseReader.Code != expectedStatus {
		t.Error("Expected request to succeed but it didnt")
	}
}

func TestItMatchesWildcardBodies(t *testing.T) {
	wildcardBody := "*"
	expectedStatus := http.StatusOK

	config := FakeEndpoint{testEndpointName, cdcDisabled, request{testURL, "POST", nil, wildcardBody}, response{expectedStatus, "", nil}}
	server := NewServer([]FakeEndpoint{config})

	requestWithDifferentBody, _ := http.NewRequest("POST", testURL, strings.NewReader("This body isnt what we said but it should match"))
	responseReader := httptest.NewRecorder()

	server.ServeHTTP(responseReader, requestWithDifferentBody)

	if responseReader.Code != expectedStatus {
		t.Errorf("Expected code %v but got %v", expectedStatus, responseReader.Code)
	}
}

func TestItRecordsIncomingRequests(t *testing.T) {
	wildcardBody := "*"
	expectedStatus := http.StatusOK

	config := FakeEndpoint{testEndpointName, cdcDisabled, request{testURL, "POST", nil, wildcardBody}, response{expectedStatus, "", nil}}
	server := NewServer([]FakeEndpoint{config})

	requestWithDifferentBody, _ := http.NewRequest("POST", testURL, strings.NewReader("This body isnt what we said but it should match"))
	responseReader := httptest.NewRecorder()

	server.ServeHTTP(responseReader, requestWithDifferentBody)

	if len(server.requests) != 1 {
		t.Fatalf("Expected one request to be recorded but got %d", len(server.requests))
	}

	if server.requests[0].Method != "POST" {
		t.Error("It doesnt look like it recorded the request properly")
	}
}

func TestItRespectsURLencoding(t *testing.T) {
	escapedURL := "/document/10.1007%2Fs00414-006-0114-x"
	body := "some body"
	expectedStatus := http.StatusOK
	config := FakeEndpoint{testEndpointName, cdcDisabled, request{escapedURL, "POST", nil, body}, response{expectedStatus, "", nil}}
	server := NewServer([]FakeEndpoint{config})

	request, _ := http.NewRequest("POST", escapedURL, strings.NewReader(body))
	responseReader := httptest.NewRecorder()

	server.ServeHTTP(responseReader, request)

	if responseReader.Code != expectedStatus {
		t.Errorf("It looks like it didnt respect the escaped url, got a %d", responseReader.Code)
	}
}
