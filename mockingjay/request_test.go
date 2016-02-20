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
