package mockingjay

import (
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestItCreatesHTTPRequests(t *testing.T) {
	headers := make(map[string]string)
	headers["foo"] = "bar"
	uri := "/hello"
	method := "PUT"
	body := "Body body body"
	baseURL := "http://localhost:1234"

	mockingJayRequest := Request{
		URI:     uri,
		Method:  method,
		Headers: headers,
		Body:    body}

	httpRequest, _ := mockingJayRequest.AsHTTPRequest(baseURL)

	assert.Equal(t, httpRequest.URL.String(), httpRequest.URL.String())
	assert.Equal(t, httpRequest.Method, method)
	assert.Equal(t, httpRequest.Header.Get("foo"), "bar")

	requestBody, _ := ioutil.ReadAll(httpRequest.Body)

	assert.Equal(t, string(requestBody), body)
}

func TestItValidatesRequests(t *testing.T) {
	noURIRequest := Request{
		URI:    "",
		Method: "POST"}

	assert.Equal(t, noURIRequest.errors(), errEmptyURI)

	noMethodRequest := Request{
		URI:    "/",
		Method: ""}

	assert.Equal(t, noMethodRequest.errors(), errEmptyMethod)

	validRequest := Request{
		URI:    "/",
		Method: "POST",
	}

	assert.Nil(t, validRequest.errors())

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

	assert.True(t, requestMatches(config, incomingRequest))
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

	assert.True(t, requestMatches(expectedRequest, incomingRequest))
}
