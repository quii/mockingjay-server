package monkey

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

const alwaysMonkeyingAround = 1.0
const neverMonkeyAround = 0.0
const cannedResponse = "hello, world"

func TestItLoadsFromYAML(t *testing.T) {

	yaml := `
---
# Writes a different body 50% of the time
- body: "This is wrong :( "
  frequency: 0.5

# Delays initial writing of response by a second 20% of the time
- delay: 1000
  frequency: 0.2

# Returns a 404 30% of the time
- status: 404
  frequency: 0.3

# Write 10,000,000 garbage bytes 10% of the time
- garbage: 10000000
  frequency: 0.09
`
	delegate := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	monkeyServer, err := NewServerFromYAML(delegate.Config.Handler, []byte(yaml))

	assert.Nil(t, err, "It didnt return a server from the YAML")
	assert.Len(t, monkeyServer.(*server).behaviours, 4)
}

func TestItMonkeysWithStatusCodesAndBodies(t *testing.T) {
	monkeyBehaviour := new(behaviour)
	monkeyBehaviour.Frequency = alwaysMonkeyingAround
	monkeyBehaviour.Status = http.StatusNotFound
	monkeyBehaviour.Body = "hello, monkey"

	testServer, request := makeTestServerAndRequest()

	monkeyServer := newServerFromBehaviour(testServer.Config.Handler, []behaviour{*monkeyBehaviour})

	w := httptest.NewRecorder()

	monkeyServer.ServeHTTP(w, request)

	assert.Equal(t, w.Code, monkeyBehaviour.Status, "Server shouldve returned a 404 because of monkey")
	assert.Equal(t, w.Body.String(), monkeyBehaviour.Body, "Server should've returned a different body because of monkey override")
}

func TestItReturnsGarbage(t *testing.T) {
	monkeyBehaviour := new(behaviour)
	monkeyBehaviour.Frequency = alwaysMonkeyingAround
	monkeyBehaviour.Garbage = 1984

	testServer, request := makeTestServerAndRequest()

	monkeyServer := newServerFromBehaviour(testServer.Config.Handler, []behaviour{*monkeyBehaviour})

	w := httptest.NewRecorder()

	monkeyServer.ServeHTTP(w, request)

	assert.Len(t, w.Body.Bytes(), monkeyBehaviour.Garbage, "Server shouldve returned garbage")
}

func TestItDoesntMonkeyAroundWhenFrequencyIsNothing(t *testing.T) {
	monkeyBehaviour := new(behaviour)
	monkeyBehaviour.Frequency = neverMonkeyAround
	monkeyBehaviour.Body = "blah blah"

	testServer, request := makeTestServerAndRequest()

	monkeyServer := newServerFromBehaviour(testServer.Config.Handler, []behaviour{*monkeyBehaviour})

	w := httptest.NewRecorder()

	monkeyServer.ServeHTTP(w, request)

	assert.Equal(t, w.Body.String(), cannedResponse, "Response body shouldve been tampered")
}

func makeTestServerAndRequest() (*httptest.Server, *http.Request) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(cannedResponse))
	}))
	request, _ := http.NewRequest("GET", server.URL, nil)

	return server, request
}
