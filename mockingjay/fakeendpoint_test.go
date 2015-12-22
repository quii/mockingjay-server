package mockingjay

import (
	"strings"
	"testing"
)

const testYAML = `
---
 - name: Test endpoint
   request:
     uri: /hello
     method: GET
     headers:
       content-type: application/json
     body: foobar
   response:
     code: 200
     body: '{"message": "hello, world"}'
     headers:
       content-type: text/json

 - name: Test endpoint 2
   cdcdisabled: true
   request:
     uri: /world
     method: DELETE
   response:
     code: 200
     body: ''

 - name: Failing endpoint
   request:
     uri: /card
     method: POST
     body: Greetings
   response:
     code: 500
     body: Oh bugger
 `

func TestItCreatesAServerConfigFromYAML(t *testing.T) {
	endpoints, err := NewFakeEndpoints([]byte(testYAML))

	if err != nil {
		t.Fatalf("Shouldn't have got an error for valid YAML [%v]", err)
	}
	firstEndpoint := endpoints[0]

	if len(endpoints) != 3 {
		t.Fatalf("There should be 3 endpoints found in YAML")
	}

	if firstEndpoint.Name != "Test endpoint" {
		t.Error("There should be a name set for the endpoint")
	}

	if firstEndpoint.Request.URI != "/hello" {
		t.Error("Request URI was not properly set")
	}

	if firstEndpoint.Request.Headers["content-type"] != "application/json" {
		t.Error("Request headers were not parsed")
	}

	if firstEndpoint.Response.Headers["content-type"] != "text/json" {
		t.Errorf("Response headers were not parsed, got %s", firstEndpoint.Response.Headers)
	}

	if firstEndpoint.Request.Method != "GET" {
		t.Error("Request method was not properly set")
	}

	if firstEndpoint.Response.Code != 200 {
		t.Error("Response code was not properly set")
	}

	if firstEndpoint.Response.Body != `{"message": "hello, world"}` {
		t.Errorf("Response body was not properly set got [%s]", firstEndpoint.Response.Body)
	}

	if firstEndpoint.CDCDisabled {
		t.Error("First endpoint doesnt define cdc preference so it should be enabled by default")
	}

	endpoint2 := endpoints[1]

	if endpoint2.Request.Method != "DELETE" {
		t.Error("Request method for second fake was not properly set")
	}

	if endpoint2.Response.Body != "" {
		t.Error("Response body for second fake was not properly set")
	}

	if endpoint2.String() != "Test endpoint 2 (DELETE /world)" {
		t.Errorf("Fake didnt have correct Stringer, got %s", endpoints[1].String())
	}

	if endpoints[2].Request.Body != "Greetings" {
		t.Errorf("Request body for third fake was not properly set, got [%s]", endpoints[2].Request.Body)
	}

	if !endpoint2.CDCDisabled {
		t.Error("Second endpoint should have CDC disabled")
	}
}

func TestItReturnsAnErrorWhenNotValidYAML(t *testing.T) {
	_, err := NewFakeEndpoints([]byte("not real YAML"))

	if err == nil {
		t.Error("Expected an error to be returned because the YAML is bad")
	}

	if !strings.HasPrefix(err.Error(), "yaml: unmarshal errors:") {
		t.Errorf("Expected unmarshal error actual: %v", err.Error())
	}

}

const badYAML = `
---
 - name: Test endpoint
   roquest:
     internet: /hello
     cats: GET
     headers:
       content-type: application/json
     body: foobar
   response:
     code: 200
     body: hello, world
     headers:
       content-type: text/plain

 - name: Test endpoint 2
 `

func TestItReturnsAnErrorWhenStructureOfYAMLIsWrong(t *testing.T) {
	_, err := NewFakeEndpoints([]byte(badYAML))
	if err == nil {
		t.Error("Expected an error to be returned because the YAML is bad")
	}

	if err != invalidConfigError {
		t.Errorf("Expected YAML was invalid error actual: %v", err.Error())
	}
}

const incompleteYAML = `
---
 - name: Test endpoint
   request:
     uri: /world
     method: GET
   response:
     body: 'A body'
 `

func TestItReturnsAnErrorWhenYAMLIsIncomplete(t *testing.T) {
	_, err := NewFakeEndpoints([]byte(incompleteYAML))
	if err == nil {
		t.Error("Expected an error to be returned because the YAML has missing fields")
	}

	if err != invalidConfigError {
		t.Errorf("Expected YAML was incomplete error actual: %v", err.Error())
	}

}

const duplicatedRequest = `
---
 - name: Test endpoint
   request:
     uri: /hello
     method: GET
   response:
     code: 200
     body: '{"message": "hello, world"}'

 - name: Duplicated test endpoint
   request:
     uri: /hello
     method: GET
   response:
     code: 404
 `

func ignoreTestItReturnsErrorWhenRequestsAreDuplicated(t *testing.T) {
	_, err := NewFakeEndpoints([]byte(duplicatedRequest))

	if err == nil {
		t.Error("Expected an error to be returned for duplicated requests")
	}
}
