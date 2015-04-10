package mockingjay

import (
	"testing"
)

const testJSON = `
[
	{
		"Name": "Test endpoint",
		"Request":{
	    	"URI" : "/hello",
	    	"Method": "GET",
	    	"Headers":
	    		{
	    			"Content-Type": "application/json"
	    		}

		},
		"Response":{
			"Code": 200,
			"Body": "hello, world",
	    	"Headers":
	    		{
	    			"Content-Type": "text/plain"
	    		}
		}
	},
	{
		"Name": "Test endpoint 2",
		"Request":{
	    	"URI" : "/world",
	    	"Method": "DELETE"
		},
		"Response":{
			"Code": 200,
			"Body": "hello, world"
		}
	},
	{
		"Request":{
	    	"URI" : "/card",
	    	"Method": "POST",
	    	"Body": "Greetings"
		},
		"Response":{
			"Code": 500,
			"Body": "Oh bugger"
		}
	}

]
`

func TestItCreatesAServerConfigFromJSON(t *testing.T) {
	endpoints, err := NewFakeEndpoints(testJSON)

	if err != nil {
		t.Fatalf("Shouldn't have got an error for valid JSON %v", err)
	}
	firstEndpoint := endpoints[0]

	if len(endpoints) != 3 {
		t.Error("There should be 3 endpoints found in json")
	}

	if firstEndpoint.Name != "Test endpoint" {
		t.Error("There should be a name set for the endpoint")
	}

	if firstEndpoint.Request.URI != "/hello" {
		t.Error("Request URI was not properly set")
	}

	if firstEndpoint.Request.Headers["Content-Type"] != "application/json" {
		t.Error("Request headers were not parsed")
	}

	if firstEndpoint.Response.Headers["Content-Type"] != "text/plain" {
		t.Errorf("Response headers were not parsed, got %s", firstEndpoint.Response.Headers["Content-Type"])
	}

	if firstEndpoint.Request.Method != "GET" {
		t.Error("Request method was not properly set")
	}

	if firstEndpoint.Response.Code != 200 {
		t.Error("Response code was not properly set")
	}

	if firstEndpoint.Response.Body != "hello, world" {
		t.Error("Response body was not properly set")
	}

	if endpoints[1].Request.Method != "DELETE" {
		t.Error("Request method for second fake was not properly set")
	}

	if endpoints[1].String() != "Test endpoint 2 (DELETE /world)" {
		t.Errorf("Fake didnt have correct Stringer, got %s", endpoints[1].String())
	}

	if endpoints[2].Request.Body != "Greetings" {
		t.Errorf("Request body for third fake was not properly set, got [%s]", endpoints[2].Request.Body)
	}
}

func TestItReturnsAnErrorWhenNotValidJSON(t *testing.T) {
	_, err := NewFakeEndpoints("not real json")
	if err == nil {
		t.Error("Expected an error to be returned because the JSON is bad")
	}

}

const badJSON = `
[
	{
		"Request":{
	    	"URI" : "/hello",
	    	"Method": "GET",
	    	"Headers":
	    		{
	    			"Content-Type": "application/json"
	    		}

		},
		"Response":{
			"Code": 200,
			"Body": "hello, world",
	    	"Headers":
	    		{
	    			"Content-Type": "text/plain"
	    		}
		}
	},
	{
		"foo":{
	    	"URIBABARBABR" : "/world",
	    	"Method": "DELETE"
		},
		"Bummer":{
			"Code": 200,
			"Body": "hello, world"
		}
	}
]
`

func TestItReturnsAnErrorWhenStructureOfJSONIsWrong(t *testing.T) {
	_, err := NewFakeEndpoints(badJSON)
	if err == nil {
		t.Error("Expected an error to be returned because the JSON is bad")
	}
}
