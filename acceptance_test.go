package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"testing"

	"github.com/quii/mockingjay-server/mockingjay"
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
		log.Fatal("Couldn't make endpoints for test", endpoints)
	}

	server := mockingjay.NewServer(endpoints)
	http.Handle("/", server)
	go http.ListenAndServe(":9094", nil)
}

func TestItLaunchesServersAndIsCompatibleWithItsOwnConfig(t *testing.T) {

	checker := NewCompatabilityChecker(log.New(os.Stdout, "mocking-jay: ", log.Ldate|log.Ltime))

	if !checker.CheckCompatability(endpoints, "http://localhost:9094") {
		t.Log("Endpoints were not seen as compatible and they should've been.")
		t.Fail()
	}
}

func TestItListsRequestsItHasReceived(t *testing.T) {
	http.Get("http://localhost:9094/hello")

	res, err := http.Get("http://localhost:9094/requests")

	if err != nil {
		t.Fatal(err)
	}

	if res.StatusCode != http.StatusOK {
		t.Error("Expected a 200 but got", res.StatusCode)
	}
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

	if err != nil {
		t.Fatal("Problem calling new endpoint URL", newEndpointURL, err)
	}

	if res.StatusCode != http.StatusCreated {
		t.Error("Create endpoint didnt return created - got", res.Status)
	}

	newEndPointResponse, err := http.Get("http://localhost:9094/hello")

	if err != nil {
		t.Fatal("Problem requesting newly created endpoint", err)
	}

	if newEndPointResponse.StatusCode != http.StatusOK {
		t.Error("Didnt get a 200 from newly created endpoint, got", newEndPointResponse.StatusCode)
	}

	newEndpointBody, _ := ioutil.ReadAll(newEndPointResponse.Body)

	if string(newEndpointBody) != `{"message": "hello, world"}` {
		t.Error("New endpoint didnt return the correct body, got", newEndpointBody)
	}
}
