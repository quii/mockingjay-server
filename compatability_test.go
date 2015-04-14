package main

import (
	"fmt"
	"github.com/quii/mockingjay"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestItChecksAValidEndpointsJSON(t *testing.T) {
	body := `{"foo":"bar"}`
	realServer := makeFakeDownstreamServer(body, noSleep)
	checker := makeChecker(body)

	if !checker.CheckCompatability(realServer.URL) {
		t.Error("Checker should've found this endpoint to be correct")
	}
}

func TestItFlagsDifferentJSONToBeIncompatible(t *testing.T) {
	serverResponseBody := `{"foo": "bar"}`
	fakeResponseBody := `{"baz": "boo"}`

	realServer := makeFakeDownstreamServer(serverResponseBody, noSleep)
	checker := makeChecker(fakeResponseBody)

	if checker.CheckCompatability(realServer.URL) {
		t.Error("Checker should've found this endpoint to be incorrect")
	}
}

func TestItIsIncompatibleWhenRealServerIsntReachable(t *testing.T) {
	if !makeChecker("doesn't matter").CheckCompatability("http://localhost:12344") {
		t.Error("Checker shouldve found this to be an error as the real server isnt reachable")
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

func makeChecker(responseBody string) *CompatabilityChecker {
	fakeEndPoints, _ := mockingjay.NewFakeEndpoints([]byte(testYAML(responseBody)))
	return NewCompatabilityChecker(fakeEndPoints)
}

func testYAML(responseBody string) string {
	return fmt.Sprintf(yamlFormat, defaultRequestURI, responseBody)
}
