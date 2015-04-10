package main

import (
	"fmt"
	"github.com/quii/mockingjay"
	"net/http"
	"net/http/httptest"
	"testing"
)

const defaultRequestURI = "/hello"

const jsonFormat = `
[
    {
        "Name": "Example hello",
        "Request":{
            "URI" : %s,
            "Method": "GET"
        },
        "Response":{
            "Code": 200,
            "Body": "%s"
        }
    }
]
`

func testJSON(responseBody string) string {
	return fmt.Sprintf(jsonFormat, defaultRequestURI, responseBody)
}

func TestItChecksAValidEndpointsJSON(t *testing.T) {
	body := `{"foo":"bar"}`
	realServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.RequestURI() == defaultRequestURI {
			fmt.Fprint(w, body)
		} else {
			http.Error(w, "Not found", http.StatusNotFound)
		}
	}))

	fakeEndPoints, _ := mockingjay.NewFakeEndpoints(testJSON(body))

	checker := NewCompatabilityChecker(fakeEndPoints)

	if !checker.CheckCompatability(realServer.URL) {
		t.Error("Checker should've found this endpoint to be correct")
	}
}

func TestItFlagsDifferentJSONToBeIncompatible(t *testing.T) {
	serverResponseBody := `{"foo": "bar"}`
	fakeResponseBody := `{"baz": "boo"}`

	realServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.RequestURI() == defaultRequestURI {
			fmt.Fprint(w, serverResponseBody)
		} else {
			http.Error(w, "Not found", http.StatusNotFound)
		}
	}))

	fakeEndPoints, _ := mockingjay.NewFakeEndpoints(testJSON(fakeResponseBody))

	checker := NewCompatabilityChecker(fakeEndPoints)

	if checker.CheckCompatability(realServer.URL) {
		t.Error("Checker should've found this endpoint to be incorrect")
	}

}
