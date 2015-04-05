package main

import (
	"fmt"
	"github.com/quii/mockingjay"
	"net/http"
	"net/http/httptest"
	"testing"
)

const json = `
[
    {
        "Name": "Example hello",
        "Request":{
            "URI" : "/hello",
            "Method": "GET"
        },
        "Response":{
            "Code": 200,
            "Body": "hello, world"
        }
    }
]
`

func TestItChecksAValidEndpoint(t *testing.T) {
	realServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.RequestURI() == "/hello" {
			fmt.Fprint(w, "hello, world")
		} else {
			http.Error(w, "Not found", http.StatusNotFound)
		}
	}))

	fakeEndPoints, _ := mockingjay.NewFakeEndpoints(json)

	checker := NewCompatabilityChecker(fakeEndPoints)

	if !checker.CheckCompatability(realServer.URL) {
		t.Error("Checker should've found this endpoint to be correct")
	}
}
