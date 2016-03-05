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
	"net/http/httptest"
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
	server *mockingjay.Server
)

func init() {
	endpoints, err := mockingjay.NewFakeEndpoints([]byte(yaml))

	if err != nil {
		log.Fatal("Couldnt set up mockingjay from config", err)
	}

	server = mockingjay.NewServer(endpoints)
}

func TestItLaunchesServersAndIsCompatibleWithItsOwnConfig(t *testing.T) {
	svr := httptest.NewServer(server)
	defer svr.Close()

	checker := NewCompatabilityChecker(log.New(os.Stdout, "mocking-jay: ", log.Ldate|log.Ltime))
	assert.True(t, checker.CheckCompatability(endpoints, svr.URL))
}

func TestItListsRequestsItHasReceived(t *testing.T) {
	svr := httptest.NewServer(server)
	defer svr.Close()

	http.Get(svr.URL + "/hello")

	res, err := http.Get(svr.URL + "/requests")

	assert.Nil(t, err)
	assert.Equal(t, res.StatusCode, http.StatusOK)
}

func TestANewEndpointCanBeAdded(t *testing.T) {
	svr := httptest.NewServer(server)
	defer svr.Close()

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
	newEndpointURL := svr.URL + "/mj-new-endpoint"
	res, err := http.Post(newEndpointURL, "application/json", strings.NewReader(newEndpointJSON))

	assert.Nil(t, err)
	assert.Equal(t, res.StatusCode, http.StatusCreated)

	newEndPointResponse, err := http.Get(svr.URL + `/hello`)

	assert.Nil(t, err)
	assert.Equal(t, newEndPointResponse.StatusCode, http.StatusOK)

	newEndpointBody, _ := ioutil.ReadAll(newEndPointResponse.Body)

	assert.Equal(t, string(newEndpointBody), `{"message": "hello, world"}`)
}
