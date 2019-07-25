package mockingjay

import (
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
)

type noOpLogger struct{}

func (noOpLogger) Println(...interface{}) {}

var (
	incomingRequest = Request{
		URI:    "/hello/chris",
		Method: "GET",
	}
	testLogger   = noOpLogger{}
	endpointName = "test endpoint"
)

func TestItMatchesRequests(t *testing.T) {
	requiredHeaders := make(map[string]string)
	requiredHeaders["Content-Type"] = "application/json"

	wrongHeaderValues := make(map[string]string)
	wrongHeaderValues["Content-Type"] = "text/html"

	wrongHeaders := make(map[string]string)
	wrongHeaders["Accept"] = "text/html"

	form := make(map[string]string)
	form["name"] = "Hudson"

	config := Request{
		URI:     "/cats",
		Method:  "POST",
		Headers: requiredHeaders,
		Body:    `123`,
		Form:    form,
	}

	failingCases := []struct {
		name            string
		incomingRequest Request
	}{
		{
			"Incorrect URI",
			Request{
				URI:     "/wrong-uri",
				Method:  "POST",
				Headers: requiredHeaders,
				Body:    `123`,
				Form:    form,
			},
		},
		{
			"Incorrect Method",
			Request{
				URI:     "/cats",
				Method:  "GET",
				Headers: requiredHeaders,
				Body:    `123`,
				Form:    form,
			},
		},
		{
			"Incorrect header values",
			Request{
				URI:     "/cats",
				Method:  "POST",
				Headers: wrongHeaderValues,
				Body:    `123`,
				Form:    form,
			},
		},
		{
			"Incorrect headers",
			Request{
				URI:     "/cats",
				Method:  "POST",
				Headers: wrongHeaders,
				Body:    `123`,
				Form:    form,
			},
		},
		{
			"Incorrect body",
			Request{
				URI:     "/cats",
				Method:  "POST",
				Headers: requiredHeaders,
				Body:    `456`,
				Form:    form,
			},
		},
		{
			"No form",
			Request{
				URI:     "/cats",
				Method:  "POST",
				Headers: requiredHeaders,
				Body:    `123`,
				Form:    nil,
			},
		},
	}
	for _, c := range failingCases {
		assert.False(t, requestMatches(config, c.incomingRequest, endpointName, testLogger), c.name)
	}

	assert.True(t, requestMatches(config, config, endpointName, testLogger), "The exact same request should match")
}

func TestItMatchesJSONWithSpaces(t *testing.T) {
	serverConfig := Request{
		URI:    "/hello/world",
		Method: "POST",
		Body:   `{"foo": 2}`,
	}

	incoming := Request{
		URI:    "/hello/world",
		Method: "POST",
		Body:   `{"foo": 2    }`,
	}

	assert.True(t, requestMatches(serverConfig, incoming, endpointName, testLogger))
}

func TestItDoesntMatchWhenJSONValuesAreDifferent(t *testing.T) {
	serverConfig := Request{
		URI:    "/hello/world",
		Method: "POST",
		Body:   `{"foo": 2}`,
	}

	incoming := Request{
		URI:    "/hello/world",
		Method: "POST",
		Body:   `{"foo": 3    }`,
	}

	assert.False(t, requestMatches(serverConfig, incoming, endpointName, testLogger))
}

func TestItDoesntCrashOnNonJSONAndAssumesNotMatch(t *testing.T) {
	serverConfig := Request{
		URI:    "/hello/world",
		Method: "POST",
		Body:   `{"foo": 2}`,
	}

	incoming := Request{
		URI:    "/hello/world",
		Method: "POST",
		Body:   `not json`,
	}

	assert.False(t, requestMatches(serverConfig, incoming, endpointName, testLogger))
}

func TestMatchingWithRegex(t *testing.T) {

	uriPathRegex, err := regexp.Compile(`\/hello\/[a-z]+`)
	regexURI := &RegexField{Regexp: uriPathRegex}

	assert.Nil(t, err)

	serverConfig := Request{
		URI:      "/hello/world",
		RegexURI: regexURI,
		Method:   "GET",
	}

	assert.True(t, requestMatches(serverConfig, incomingRequest, endpointName, testLogger))
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

	assert.True(t, requestMatches(serverConfig, incomingRequest, endpointName, testLogger))
}

func TestItIgnoresExtraHeadersInEqualityCheck(t *testing.T) {
	requiredHeaders := make(map[string]string)
	requiredHeaders["Content-Type"] = "application/json"

	config := Request{
		URI:     "",
		Method:  "POST",
		Headers: requiredHeaders,
	}

	extraHeaders := make(map[string]string)
	extraHeaders["Content-Type"] = "application/json"
	extraHeaders["Content-Size"] = "ten"
	incomingRequest := config
	incomingRequest.Headers = extraHeaders

	assert.True(t, requestMatches(config, incomingRequest, endpointName, testLogger))
}

func TestItIgnoresHeadersKeyCasing(t *testing.T) {
	requiredHeaders := make(map[string]string)
	requiredHeaders["content-type"] = "application/json"

	expectedRequest := Request{
		URI:     "",
		Method:  "POST",
		Headers: requiredHeaders,
	}

	differentlyCasingHeaders := make(map[string]string)
	differentlyCasingHeaders["Content-Type"] = "application/json"

	incomingRequest := expectedRequest
	incomingRequest.Headers = differentlyCasingHeaders

	assert.True(t, requestMatches(expectedRequest, incomingRequest, endpointName, testLogger))
}

func TestItIgnoresOrderOfQueryString(t *testing.T) {
	expectedRequest := Request{
		URI:    "?a=1&b=2",
		Method: "GET",
	}

	incomingRequest := Request{
		URI:    "?b=2&a=1",
		Method: "GET",
	}

	assert.True(t, requestMatches(expectedRequest, incomingRequest, endpointName, testLogger))
}
