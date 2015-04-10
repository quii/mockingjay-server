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

	mockingJayRequest := request{uri, method, headers, body}

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
	noURIRequest := request{"", "POST", nil, ""}

	if noURIRequest.isValid() {
		t.Error("A request without a URI is seen as valid")
	}

	noMethodRequest := request{"/", "", nil, ""}

	if noMethodRequest.isValid() {
		t.Error("A request without a method is seen as valid")
	}

	validRequest := request{"/", "POST", nil, ""}

	if !validRequest.isValid() {
		t.Error("A valid request is seen as not valid")
	}
}
