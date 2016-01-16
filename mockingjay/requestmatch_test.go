package mockingjay

import (
	"regexp"
	"testing"
)

var (
	incomingRequest = Request{
		URI:    "/hello/chris",
		Method: "GET",
	}
)

func TestMatchingWithRegex(t *testing.T) {

	uriPathRegex, err := regexp.Compile(`\/hello\/[a-z]+`)

	if err != nil {
		t.Fatal(err)
	}

	serverConfig := Request{
		URI:      "/hello/world",
		RegexURI: uriPathRegex,
		Method:   "GET",
	}

	if !requestMatches(serverConfig, incomingRequest) {
		t.Error("Requests didnt match when we expected them to", incomingRequest, serverConfig)
	}
}

func TestItMatchesOnURL(t *testing.T) {
	serverConfig := Request{
		URI:    "/hello/bob",
		Method: "GET",
	}

	if requestMatches(serverConfig, incomingRequest) {
		t.Error("Should not match", serverConfig, incomingRequest)
	}
}

func TestItMatchesWildcardBodies(t *testing.T) {
	incomingRequest := Request{
		URI:    "/x",
		Method: "POST",
		Body:   "Doesn't matter",
	}

	serverConfig := Request{
		URI:    "/x",
		Method: "POST",
		Body:   "*",
	}

	if !requestMatches(serverConfig, incomingRequest) {
		t.Error("Expected wildcards to match", incomingRequest, serverConfig)
	}
}
