package main

import (
	"github.com/quii/mockingjay-server/mockingjay"
	"log"
	"net/http"
	"testing"
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

	checker := NewCompatabilityChecker()

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
