package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/quii/mockingjay-server/mockingjay"
	"github.com/stretchr/testify/assert"
)

var (
	checker = NewCompatabilityChecker(log.New(os.Stdout, "mocking-jay: ", log.Ldate|log.Ltime), defaultHTTPTimeoutSeconds)
)

func TestItIgnoresEndpointsNotSetToCDC(t *testing.T) {
	yaml := `
---
 - name: Endpoint doesnt matter so much as its ignored
   cdcdisabled: true
   request:
     uri: /hello
     method: GET
   response:
     code: 200
     body: 'ok'
`
	endpoints, _ := mockingjay.NewFakeEndpoints([]byte(yaml))

	realServerThatsNotCompatible := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, ":(")
	}))

	compatible := checker.CheckCompatibility(endpoints, realServerThatsNotCompatible.URL)
	assert.True(t, compatible, `Checker shouldve found downstream server to be "compatible" as the endpoint is ignored`)
}

func TestItMatchesHeaders(t *testing.T) {
	yaml := `
---
 - name: Endpoint with response headers
   request:
     uri: /hello
     method: GET
   response:
     code: 200
     body: 'ok'
     headers:
       content-type: text/json
`
	endpoints, err := mockingjay.NewFakeEndpoints([]byte(yaml))

	assert.Nil(t, err)

	realServerWithoutHeaders := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "ok")
	}))

	realServerWithHeaders := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-TYPE", "text/json")
		fmt.Fprint(w, "ok")
	}))

	assert.False(t,
		checker.CheckCompatibility(endpoints, realServerWithoutHeaders.URL),
		"Checker should've found downstream server to be incompatible as it did not include the response headers we expected")

	assert.True(t, checker.CheckCompatibility(endpoints, realServerWithHeaders.URL))
}

func TestItMatchesBodiesAsStrings(t *testing.T) {
	body := "Chris"
	downstreamBody := "Christopher"
	realServer := makeFakeDownstreamServer(downstreamBody, noSleep)
	endpoints := makeEndpoints(body)

	assert.True(t, checker.CheckCompatibility(endpoints, realServer.URL))
}

func TestItFailsWhenStringsDontMatch(t *testing.T) {
	body := "Chris"
	downstreamBody := "Pixies"
	realServer := makeFakeDownstreamServer(downstreamBody, noSleep)
	endpoints := makeEndpoints(body)

	assert.False(t, checker.CheckCompatibility(endpoints, realServer.URL))
}

func TestItChecksAValidEndpointsJSON(t *testing.T) {
	body := `{"foo":"bar"}`
	realServer := makeFakeDownstreamServer(body, noSleep)
	endpoints := makeEndpoints(body)

	assert.True(t, checker.CheckCompatibility(endpoints, realServer.URL))
}

func TestItChecksAValidEndpointsXML(t *testing.T) {
	body := `<foo><bar>x</bar></foo>`
	realServerBody := `<foo><bar>y</bar></foo>`
	realServer := makeFakeDownstreamServer(realServerBody, noSleep)
	endpoints := makeEndpoints(body)

	assert.True(t, checker.CheckCompatibility(endpoints, realServer.URL))
}

func TestItChecksInvaidXML(t *testing.T) {
	body := `<foo><bar>x</bar></foo>`
	realServerBody := `not xml`
	realServer := makeFakeDownstreamServer(realServerBody, noSleep)
	endpoints := makeEndpoints(body)

	assert.False(t, checker.CheckCompatibility(endpoints, realServer.URL))
}

func TestItFlagsDifferentJSONToBeIncompatible(t *testing.T) {
	serverResponseBody := `{"foo": "bar"}`
	fakeResponseBody := `{"baz": "boo"}`

	realServer := makeFakeDownstreamServer(serverResponseBody, noSleep)
	endpoints := makeEndpoints(fakeResponseBody)

	assert.False(t, checker.CheckCompatibility(endpoints, realServer.URL))
}

func TestItChecksStatusCodes(t *testing.T) {
	body := `{"foo": "bar"}`

	realServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(body))
	}))

	endpoints := makeEndpoints(body)

	assert.False(t, checker.CheckCompatibility(endpoints, realServer.URL))
}

func TestIfItsNotJSONItDoesAStringMatch(t *testing.T) {
	serverResponseBody := "hello world"
	fakeResponseBody := "hello bob"

	realServer := makeFakeDownstreamServer(serverResponseBody, noSleep)
	endpoints := makeEndpoints(fakeResponseBody)

	assert.False(t, checker.CheckCompatibility(endpoints, realServer.URL))
}

func TestItAllowsAnyBodyOnWildcard(t *testing.T) {
	body := "hello world"

	realServer := makeFakeDownstreamServer(body, noSleep)
	endpoints := makeEndpoints("*")

	assert.True(t, checker.CheckCompatibility(endpoints, realServer.URL))
}

func TestItIsIncompatibleWhenRealServerIsntReachable(t *testing.T) {
	body := "doesnt matter"
	endpoints := makeEndpoints(body)

	assert.False(t, checker.CheckCompatibility(endpoints, "http://localhost:12344"))
}

func TestItHandlesBadURLsInConfig(t *testing.T) {
	yaml := fmt.Sprintf(yamlFormat, "not a real url", "foobar")
	fakeEndPoints, _ := mockingjay.NewFakeEndpoints([]byte(yaml))

	assert.False(t, checker.CheckCompatibility(fakeEndPoints, "also not a real url"))
}

func TestItFailsWhenExpectedJSONButGotSomethingElse(t *testing.T) {
	serverResponseBody := `not json`
	fakeResponseBody := `{"isJSON": true}`

	realServer := makeFakeDownstreamServer(serverResponseBody, noSleep)
	endpoints := makeEndpoints(fakeResponseBody)

	assert.False(t, checker.CheckCompatibility(endpoints, realServer.URL))
}

// https://github.com/quii/mockingjay-server/issues/3
func TestWhitespaceSensitivity(t *testing.T) {
	y := `
---
- name: Testing whitespace sensitivity
  request:
    uri: /hello
    method: POST
    body: '{"email":"foo@bar.com","password":"xxx"}'
  response:
    code: 200
    body: '{"token": "1234abc"}'
        `
	fakeEndPoints, _ := mockingjay.NewFakeEndpoints([]byte(y))
	realServer := makeFakeDownstreamServer(`{"token":    "1234abc"}`, noSleep)

	assert.True(t, checker.CheckCompatibility(fakeEndPoints, realServer.URL))
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
