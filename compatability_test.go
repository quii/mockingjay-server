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
	realServer := makeRealServer(body, noSleep)

	fakeEndPoints, err := mockingjay.NewFakeEndpoints([]byte(testYAML(body)))

	if err != nil {
		t.Fatalf("Couldn't make mockingjay endpoints, is your data correct? [%v]", err)
	}

	checker := NewCompatabilityChecker(fakeEndPoints)

	if !checker.CheckCompatability(realServer.URL) {
		t.Error("Checker should've found this endpoint to be correct")
	}
}

func TestItFlagsDifferentJSONToBeIncompatible(t *testing.T) {
	serverResponseBody := `{"foo": "bar"}`
	fakeResponseBody := `{"baz": "boo"}`

	realServer := makeRealServer(serverResponseBody, noSleep)

	fakeEndPoints, err := mockingjay.NewFakeEndpoints([]byte(testYAML(fakeResponseBody)))

	if err != nil {
		t.Fatalf("Couldn't make mockingjay endpoints, is your data correct? [%v]", err)
	}

	checker := NewCompatabilityChecker(fakeEndPoints)

	if checker.CheckCompatability(realServer.URL) {
		t.Error("Checker should've found this endpoint to be incorrect")
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

func makeRealServer(responseBody string, sleepTime time.Duration) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(sleepTime * time.Millisecond)
		if r.URL.RequestURI() == defaultRequestURI {
			fmt.Fprint(w, responseBody)
		} else {
			http.Error(w, "Not found", http.StatusNotFound)
		}
	}))
}

func testYAML(responseBody string) string {
	return fmt.Sprintf(yamlFormat, defaultRequestURI, responseBody)
}
