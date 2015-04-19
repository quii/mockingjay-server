package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

const alwaysMonkeyingAround = 1.0

func TestItMonkeysWithStatusCodes(t *testing.T) {
	notFoundBehaviour := new(behaviour)
	notFoundBehaviour.Frequency = alwaysMonkeyingAround
	notFoundBehaviour.Status = http.StatusNotFound

	testServer, request := makeTestServerAndRequest()

	monkeyServer := NewMonkeyServer(testServer.Config.Handler, []behaviour{*notFoundBehaviour})

	w := httptest.NewRecorder()

	monkeyServer.ServeHTTP(w, request)

	if w.Code != notFoundBehaviour.Status {
		t.Error("Server shouldve returned a 404 because of monkey override")
	}
}

func makeTestServerAndRequest() (*httptest.Server, *http.Request) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	request, _ := http.NewRequest("GET", server.URL, nil)

	return server, request
}
