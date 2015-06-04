package main

import (
	"github.com/quii/mockingjay-server/mockingjay"
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

func IgnoreTestItLaunchesServersAndIsCompatibleWithItsOwnConfig(t *testing.T) {
	endpoints, err := mockingjay.NewFakeEndpoints([]byte(yaml))

	if err != nil {
		t.Log(err)
		t.Fatal("Couldnt make endpoints from YAML")
	}

	server := mockingjay.NewServer(endpoints)
	http.Handle("/", server)
	go http.ListenAndServe(":9092", nil)

	checker := NewCompatabilityChecker()

	if !checker.CheckCompatability(endpoints, "http://localhost:9092") {
		t.Log("Endpoints were not seen as compatible and they should've been.")
		t.Fail()
	}
}
