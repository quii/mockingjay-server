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
