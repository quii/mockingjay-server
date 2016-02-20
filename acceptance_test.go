package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"testing"

	"github.com/quii/mockingjay-server/mockingjay"
	"github.com/stretchr/testify/assert"
)

const yaml = `
- name: Valid login details
  request:
    uri: /token
    method: POST
    body: '{"email":"foo@domain.com","password":"secret"}'
    headers:
     X-Api-Key: supersecret
  response:
    code: 200
    body: '{"token": "1234abc"}'
`

var (
	endpoints []mockingjay.FakeEndpoint
)

func init() {
	endpoints, err := mockingjay.NewFakeEndpoints([]byte(yaml))

	if err != nil {
		log.Fatal("Couldnt set up mockingjay from config", err)
	}

	server := mockingjay.NewServer(endpoints)
	http.Handle("/", server)
	go http.ListenAndServe(":9094", nil)
}

func TestItLaunchesServersAndIsCompatibleWithItsOwnConfig(t *testing.T) {
	checker := NewCompatabilityChecker(log.New(os.Stdout, "mocking-jay: ", log.Ldate|log.Ltime))
	assert.True(t, checker.CheckCompatability(endpoints, "http://localhost:9094"))
}

func TestItListsRequestsItHasReceived(t *testing.T) {
	http.Get("http://localhost:9094/hello")

	res, err := http.Get("http://localhost:9094/requests")

	assert.Nil(t, err)
	assert.Equal(t, res.StatusCode, http.StatusOK)
}

func TestANewEndpointCanBeAdded(t *testing.T) {
	newEndpointJSON := `
	{
	  "Name": "Test endpoint",
	  "CDCDisabled": false,
	  "Request": {
	    "URI": "/hello",
	    "Method": "GET",
	    "Headers": null,
	    "Body": ""
	  },
	  "Response": {
	    "Code": 200,
	    "Body": "{\"message\": \"hello, world\"}",
	    "Headers": {
	      "content-type": "text\/json"
	    }
	  }
	}
	`
	newEndpointURL := "http://localhost:9094/mj-new-endpoint"
	res, err := http.Post(newEndpointURL, "application/json", strings.NewReader(newEndpointJSON))

	assert.Nil(t, err)
	assert.Equal(t, res.StatusCode, http.StatusCreated)

	newEndPointResponse, err := http.Get("http://localhost:9094/hello")

	assert.Nil(t, err)
	assert.Equal(t, newEndPointResponse.StatusCode, http.StatusOK)

	newEndpointBody, _ := ioutil.ReadAll(newEndPointResponse.Body)

	assert.Equal(t, string(newEndpointBody), `{"message": "hello, world"}`)
}
