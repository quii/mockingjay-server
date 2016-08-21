package mockingjay

import (
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
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

	assert.Equal(t, httpRequest.URL.String(), httpRequest.URL.String())
	assert.Equal(t, httpRequest.Method, method)
	assert.Equal(t, httpRequest.Header.Get("foo"), "bar")

	requestBody, _ := ioutil.ReadAll(httpRequest.Body)

	assert.Equal(t, string(requestBody), body)
}

func TestItMapsHTTPRequestsToMJRequests(t *testing.T) {
	req, _ := http.NewRequest(http.MethodPost, "/foo", nil)
	mjRequest := NewRequest(req)
	assert.Equal(t, mjRequest.Method, http.MethodPost)
}

func TestItSendsForms(t *testing.T) {
	mjReq := Request{
		URI:    "/cat",
		Form:   make(map[string]string),
		Method: http.MethodPost,
	}

	mjReq.Form["name"] = "Hudson"

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		if r.PostForm.Get("name") != "Hudson" {
			t.Error("Did not get expected form value from request", r.PostForm)
		}
	})

	req, err := mjReq.AsHTTPRequest("/")

	if err != nil {
		t.Fatal("Couldnt create http request from mj request", err)
	}

	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)
}
