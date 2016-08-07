package mockingjay

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const testYAML = `
---
 - name: Test endpoint
   request:
     uri: /hello/chris
     method: GET
     regexuri: "\\/hello\\/[a-z]+"
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
     form:
       age: 10
       name: Hudson
     body: Greetings
   response:
     code: 500
     body: Oh bugger
 `

func TestItCreatesAServerConfigFromYAML(t *testing.T) {
	endpoints, err := NewFakeEndpoints([]byte(testYAML))

	assert.Nil(t, err, "Shouldn't have got an error for valid YAML")
	assert.Len(t, endpoints, 3, "There should be 3 endpoints found in YAML")

	firstEndpoint := endpoints[0]

	assert.Equal(t, firstEndpoint.Name, "Test endpoint")
	assert.Equal(t, firstEndpoint.Request.URI, "/hello/chris")
	assert.Equal(t, firstEndpoint.Request.Headers["content-type"], "application/json")
	assert.Equal(t, firstEndpoint.Request.Method, "GET")
	assert.Equal(t, firstEndpoint.Request.Body, "foobar")
	assert.Equal(t, firstEndpoint.Response.Headers["content-type"], "text/json")
	assert.Equal(t, firstEndpoint.Response.Code, 200)
	assert.Equal(t, firstEndpoint.Response.Body, `{"message": "hello, world"}`)
	assert.False(t, firstEndpoint.CDCDisabled)
	assert.NotNil(t, firstEndpoint.Request.RegexURI)

	endpoint2 := endpoints[1]

	assert.True(t, endpoint2.CDCDisabled)

	failingEndpoint := endpoints[2]

	assert.Equal(t, failingEndpoint.Request.Form["age"], "10")
	assert.Equal(t, failingEndpoint.Request.Form["name"], "Hudson")
}

func TestItReturnsAnErrorWhenNotValidYAML(t *testing.T) {
	_, err := NewFakeEndpoints([]byte("not real YAML"))
	assert.NotNil(t, err, "Expected an error to be returned because the YAML is bad")
	assert.Contains(t, err.Error(), "The structure of the supplied YAML is wrong")
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
	assert.NotNil(t, err, "Expected an error to be returned because the YAML is bad")
	assert.Equal(t, err, errEmptyURI)
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
	assert.NotNil(t, err, "Expected an error to be returned because the YAML has missing fields")
	assert.Equal(t, err, errResponseInvalid)
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

func TestItReturnsErrorWhenRequestsAreDuplicated(t *testing.T) {
	_, err := NewFakeEndpoints([]byte(duplicatedRequest))
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "duplicated")
}

const badRegex = `
---
-
  name: "This doesn't make sense"
  request:
    method: GET
    regexuri: "\\/hello\\/[a-z]+"
    uri: /goodbye/chris
  response:
    body: WOOT
    code: 200
`

func TestItReturnsErrorWhenRegexDoesntMatchURI(t *testing.T) {
	_, err := NewFakeEndpoints([]byte(badRegex))
	assert.NotNil(t, err)
	assert.Equal(t, err, errBadRegex)
}
