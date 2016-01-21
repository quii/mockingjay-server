package mockingjay

import (
	"io/ioutil"
	"testing"
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

	if httpRequest.URL.String() != baseURL+uri {
		t.Errorf("Request URL is wrong, got %s", httpRequest.URL)
	}

	if httpRequest.Method != method {
		t.Errorf("Request method is wrong, got %s", httpRequest.Method)
	}

	if httpRequest.Header.Get("foo") != "bar" {
		t.Error("Request didnt have a header set")
	}

	requestBody, _ := ioutil.ReadAll(httpRequest.Body)

	if string(requestBody) != body {
		t.Errorf("Request body was not set properly, got %s", requestBody)
	}
}

func TestItValidatesRequests(t *testing.T) {
	noURIRequest := Request{
		URI:    "",
		Method: "POST"}

	if noURIRequest.isValid() != errEmptyURI {
		t.Error("A request without a URI is seen as valid")
	}

	noMethodRequest := Request{
		URI:    "/",
		Method: ""}

	if noMethodRequest.isValid() != errEmptyMethod {
		t.Error("A request without a method is seen as valid")
	}

	validRequest := Request{
		URI:    "/",
		Method: "POST",
	}

	if validRequest.isValid() != nil {
		t.Error("A valid request is seen as not valid")
	}
}
