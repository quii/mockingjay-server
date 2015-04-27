package main

import (
	"fmt"
	"github.com/quii/mockingjay-server/mockingjay"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestItChecksAValidEndpointsJSON(t *testing.T) {
	body := `{"foo":"bar"}`
	realServer := makeFakeDownstreamServer(body, noSleep)
	checker := NewCompatabilityChecker()
	endpoints := makeEndpoints(body)

	if !checker.CheckCompatability(endpoints, realServer.URL) {
		t.Error("Checker should've found this endpoint to be correct")
	}
}

func TestItFlagsDifferentJSONToBeIncompatible(t *testing.T) {
	serverResponseBody := `{"foo": "bar"}`
	fakeResponseBody := `{"baz": "boo"}`

	realServer := makeFakeDownstreamServer(serverResponseBody, noSleep)
	checker := NewCompatabilityChecker()
	endpoints := makeEndpoints(fakeResponseBody)

	if checker.CheckCompatability(endpoints, realServer.URL) {
		t.Error("Checker should've found this endpoint to be incorrect")
	}
}

func TestItChecksStatusCodes(t *testing.T) {
	body := `{"foo": "bar"}`

	realServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(body))
	}))

	checker := NewCompatabilityChecker()
	endpoints := makeEndpoints(body)

	if checker.CheckCompatability(endpoints, realServer.URL) {
		t.Error("Checker should've found this endpoint to be incorrect")
	}
}

func TestIfItsNotJSONItDoesAStringMatch(t *testing.T) {
	serverResponseBody := "hello world"
	fakeResponseBody := "hello bob"

	realServer := makeFakeDownstreamServer(serverResponseBody, noSleep)
	checker := NewCompatabilityChecker()
	endpoints := makeEndpoints(fakeResponseBody)

	if checker.CheckCompatability(endpoints, realServer.URL) {
		t.Error("Checker should've found this endpoint to be incorrect")
	}
}

func TestIfItsNotJSONItKnowsItsCompatable(t *testing.T) {
	body := "hello world"

	realServer := makeFakeDownstreamServer(body, noSleep)
	checker := NewCompatabilityChecker()
	endpoints := makeEndpoints(body)

	if !checker.CheckCompatability(endpoints, realServer.URL) {
		t.Error("Checker should've found this endpoint to be correct")
	}
}


func TestIfItsNotJSONItAllowsAnyBody(t *testing.T) {
    body := "hello world"

    realServer := makeFakeDownstreamServer(body, noSleep)
    checker := NewCompatabilityChecker()
    endpoints := makeEndpoints("*")

    if !checker.CheckCompatability(endpoints, realServer.URL) {
        t.Error("Checker should've found this endpoint to be correct")
    }
}



func TestItIsIncompatibleWhenRealServerIsntReachable(t *testing.T) {
	body := "doesnt matter"

	checker := NewCompatabilityChecker()
	endpoints := makeEndpoints(body)

	if checker.CheckCompatability(endpoints, "http://localhost:12344") {
		t.Error("Checker shouldve found this to be an error as the real server isnt reachable")
	}
}

func TestItHandlesBadURLsInConfig(t *testing.T) {
	yaml := fmt.Sprintf(yamlFormat, "not a real url", "foobar")
	fakeEndPoints, _ := mockingjay.NewFakeEndpoints([]byte(yaml))
	checker := NewCompatabilityChecker()

	if checker.CheckCompatability(fakeEndPoints, "also not a real url") {
		t.Error("Checker should've found that the URL in the YAML cannot be made into a request")
	}
}

const noSleep = 1

const defaultRequestURI = "/hello"

const yamlFormat = `
---
 - name: Test endpoint
   request:
     uri: %s
     method: GET
   response:
     code: 200
     body: '%s'
`

func makeFakeDownstreamServer(responseBody string, sleepTime time.Duration) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(sleepTime * time.Millisecond)
		if r.URL.RequestURI() == defaultRequestURI {
			fmt.Fprint(w, responseBody)
		} else {
			http.Error(w, "Not found", http.StatusNotFound)
		}
	}))
}

func makeEndpoints(body string) []mockingjay.FakeEndpoint {
	e, _ := mockingjay.NewFakeEndpoints([]byte(testYAML(body)))
	return e
}

func testYAML(responseBody string) string {
	return fmt.Sprintf(yamlFormat, defaultRequestURI, responseBody)
}

