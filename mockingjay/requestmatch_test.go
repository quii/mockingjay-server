package mockingjay

import (
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	incomingRequest = Request{
		URI:    "/hello/chris",
		Method: "GET",
	}
)

func TestItMatchesJSONWithSpaces(t *testing.T){
	serverConfig := Request{
		URI:      "/hello/world",
		Method:   "POST",
		Body : `{"foo": 2}`,
	}

	incoming := Request{
		URI:      "/hello/world",
		Method:   "POST",
		Body : `{"foo": 2    }`,
	}

	assert.True(t, requestMatches(serverConfig, incoming))
}

func TestItDoesntMatchWhenJSONValuesAreDifferent(t *testing.T){
	serverConfig := Request{
		URI:      "/hello/world",
		Method:   "POST",
		Body : `{"foo": 2}`,
	}

	incoming := Request{
		URI:      "/hello/world",
		Method:   "POST",
		Body : `{"foo": 3    }`,
	}

	assert.False(t, requestMatches(serverConfig, incoming))
}

func TestItDoesntCrashOnNonJSONAndAssumesNotMatch(t *testing.T){
	serverConfig := Request{
		URI:      "/hello/world",
		Method:   "POST",
		Body : `{"foo": 2}`,
	}

	incoming := Request{
		URI:      "/hello/world",
		Method:   "POST",
		Body : `not json`,
	}

	assert.False(t, requestMatches(serverConfig, incoming))
}

func TestItDoesntMatchUnequalBodies(t *testing.T){
	serverConfig := Request{
		URI:      "/hello/world",
		Method:   "POST",
		Body : `123`,
	}

	incoming := Request{
		URI:      "/hello/world",
		Method:   "POST",
		Body : `456`,
	}

	assert.False(t, requestMatches(serverConfig, incoming))
}

func TestMatchingWithRegex(t *testing.T) {

	uriPathRegex, err := regexp.Compile(`\/hello\/[a-z]+`)
	regexURI := &RegexYAML{Regexp: uriPathRegex}

	assert.Nil(t, err)

	serverConfig := Request{
		URI:      "/hello/world",
		RegexURI: regexURI,
		Method:   "GET",
	}

	assert.True(t, requestMatches(serverConfig, incomingRequest))
}

func TestItMatchesOnURL(t *testing.T) {
	notMatchingURL := Request{
		URI:    "/hello/bob",
		Method: "GET",
	}

	assert.False(t, requestMatches(notMatchingURL, incomingRequest))
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

	assert.True(t, requestMatches(serverConfig, incomingRequest))
}
