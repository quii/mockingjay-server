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
)

var (
	checker = NewCompatabilityChecker(log.New(os.Stdout, "mocking-jay: ", log.Ldate|log.Ltime))
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

	if !checker.CheckCompatability(endpoints, realServerThatsNotCompatible.URL) {
		t.Error(`Checker shouldve found downstream server to be "compatible" as the endpoint is ignored`)
	}
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

	if err != nil {
		t.Fatal(err)
	}

	realServerWithoutHeaders := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "ok")
	}))

	realServerWithHeaders := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-TYPE", "text/json")
		fmt.Fprint(w, "ok")
	}))

	if checker.CheckCompatability(endpoints, realServerWithoutHeaders.URL) {
		t.Error("Checker should've found downstream server to be incompatible as it did not include the response headers we expected")
	}

	if !checker.CheckCompatability(endpoints, realServerWithHeaders.URL) {
		t.Error("Checker shouldve found downstream server to be compatible")
	}
}

func TestItMatchesBodiesAsStrings(t *testing.T) {
	body := "Chris"
	downstreamBody := "Christopher"
	realServer := makeFakeDownstreamServer(downstreamBody, noSleep)
	endpoints := makeEndpoints(body)

	if !checker.CheckCompatability(endpoints, realServer.URL) {
		t.Error("Checker should've found this endpoint to be correct")
	}
}

func TestItFailsWhenStringsDontMatch(t *testing.T) {
	body := "Chris"
	downstreamBody := "Pixies"
	realServer := makeFakeDownstreamServer(downstreamBody, noSleep)
	endpoints := makeEndpoints(body)

	if checker.CheckCompatability(endpoints, realServer.URL) {
		t.Error("Checker should've found this endpoint to be incompatible")
	}
}

func TestItChecksAValidEndpointsJSON(t *testing.T) {
	body := `{"foo":"bar"}`
	realServer := makeFakeDownstreamServer(body, noSleep)
	endpoints := makeEndpoints(body)

	if !checker.CheckCompatability(endpoints, realServer.URL) {
		t.Error("Checker should've found this endpoint to be correct")
	}
}

func TestItChecksAValidEndpointsXML(t *testing.T) {
	body := `<foo><bar>x</bar></foo>`
	realServerBody := `<foo><bar>y</bar></foo>`
	realServer := makeFakeDownstreamServer(realServerBody, noSleep)
	endpoints := makeEndpoints(body)

	if !checker.CheckCompatability(endpoints, realServer.URL) {
		t.Error("Checker should've found this endpoint to be correct")
	}
}

func TestItFlagsDifferentJSONToBeIncompatible(t *testing.T) {
	serverResponseBody := `{"foo": "bar"}`
	fakeResponseBody := `{"baz": "boo"}`

	realServer := makeFakeDownstreamServer(serverResponseBody, noSleep)
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

	endpoints := makeEndpoints(body)

	if checker.CheckCompatability(endpoints, realServer.URL) {
		t.Error("Checker should've found this endpoint to be incorrect")
	}
}

func TestIfItsNotJSONItDoesAStringMatch(t *testing.T) {
	serverResponseBody := "hello world"
	fakeResponseBody := "hello bob"

	realServer := makeFakeDownstreamServer(serverResponseBody, noSleep)
	endpoints := makeEndpoints(fakeResponseBody)

	if checker.CheckCompatability(endpoints, realServer.URL) {
		t.Error("Checker should've found this endpoint to be incorrect")
	}
}

func TestIfItsNotJSONItKnowsItsCompatable(t *testing.T) {
	body := "hello world"

	realServer := makeFakeDownstreamServer(body, noSleep)
	endpoints := makeEndpoints(body)

	if !checker.CheckCompatability(endpoints, realServer.URL) {
		t.Error("Checker should've found this endpoint to be correct")
	}
}

func TestIfItsNotJSONItAllowsAnyBody(t *testing.T) {
	body := "hello world"

	realServer := makeFakeDownstreamServer(body, noSleep)
	endpoints := makeEndpoints("*")

	if !checker.CheckCompatability(endpoints, realServer.URL) {
		t.Error("Checker should've found this endpoint to be correct")
	}
}

func TestItIsIncompatibleWhenRealServerIsntReachable(t *testing.T) {
	body := "doesnt matter"
	endpoints := makeEndpoints(body)
	if checker.CheckCompatability(endpoints, "http://localhost:12344") {
		t.Error("Checker shouldve found this to be an error as the real server isnt reachable")
	}
}

func TestItHandlesBadURLsInConfig(t *testing.T) {
	yaml := fmt.Sprintf(yamlFormat, "not a real url", "foobar")
	fakeEndPoints, _ := mockingjay.NewFakeEndpoints([]byte(yaml))

	if checker.CheckCompatability(fakeEndPoints, "also not a real url") {
		t.Error("Checker should've found that the URL in the YAML cannot be made into a request")
	}
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

	if !checker.CheckCompatability(fakeEndPoints, realServer.URL) {
		t.Error("Checker should've found that the two JSONs are compatible despite having different whitespace")
	}
}

// this is panicing in goconvey?
func TestErrorReportingOfEmptyJSONArrays(t *testing.T) {
	y := `
---
- name: Testing whitespace sensitivity
  request:
    uri: /hello
    method: POST
    body: '{"email":"foo@bar.com","password":"xxx"}'
  response:
    code: 200
    body: '{"stuff": [{"foo":"bar"}]}'
`

	fakeEndPoints, err := mockingjay.NewFakeEndpoints([]byte(y))

	if err != nil {
		t.Fatal(err)
	}

	realServer := makeFakeDownstreamServer(`{"stuff":[]}`, noSleep)

	if checker.CheckCompatability(fakeEndPoints, realServer.URL) {
		t.Error("Checker shouldn't have found it compatible because its got an empty array downstream")
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
