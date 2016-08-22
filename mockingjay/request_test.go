package mockingjay

import (
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
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
	req, _ := http.NewRequest(http.MethodPost, "/foo", strings.NewReader("body"))
	req.Header.Add("foo", "bar")

	mjRequest := NewRequest(req)

	assert.Equal(t, mjRequest.Method, http.MethodPost)
	assert.Equal(t, mjRequest.Headers["foo"], "bar")
	assert.Equal(t, mjRequest.Body, "body")
}

func TestItMapsHTTPRequestsWithFormsToMJRequests(t *testing.T) {
	form := url.Values{}
	form.Add("age", "10")
	form.Add("name", "Hudson")
	body := form.Encode()

	req, _ := http.NewRequest(http.MethodPost, "/foo", strings.NewReader(body))
	req.Header.Add("content-type", "application/x-www-form-urlencoded")

	mjRequest := NewRequest(req)

	assert.Equal(t, mjRequest.Form["age"], "10")
	assert.Equal(t, mjRequest.Form["name"], "Hudson")
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

	//todo: ?? did i just stop bothering with the test here?
}

func TestItHasPrettyString(t *testing.T) {
	mapOfThings := make(map[string]string)
	mapOfThings["A"] = "B"

	longBody := "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Morbi interdum consectetur diam, sed rhoncus tortor dapibus eget. Mauris lacus metus, laoreet in nunc at, ullamcorper tincidunt turpis. Duis maximus cursus mi, a luctus eros posuere a. In laoreet neque sit amet metus vestibulum porta. Nulla quam eros, pretium at scelerisque et, mattis euismod est. Cum sociis natoque penatibus et magnis dis parturient montes, nascetur ridiculus mus. Integer id odio lorem."

	tests := []struct {
		Request        Request
		ExpectedString string
	}{
		{
			Request: Request{
				URI:    "/hello-world",
				Method: http.MethodGet,
			},
			ExpectedString: "GET /hello-world",
		},
		{
			Request: Request{
				URI:     "/hello-world",
				Method:  http.MethodGet,
				Headers: mapOfThings,
			},
			ExpectedString: "GET /hello-world Headers: [A->B]",
		},
		{
			Request: Request{
				URI:    "/hello-world",
				Method: http.MethodGet,
				Form:   mapOfThings,
			},
			ExpectedString: "GET /hello-world Form: [A->B]",
		},
		{
			Request: Request{
				URI:    "/hello-world",
				Method: http.MethodGet,
				Body:   longBody,
			},
			ExpectedString: "GET /hello-world Body: [Lorem ipsum dolor sit amet, consectetur adipiscing...]",
		},
		{
			Request: Request{
				URI:    "/hello-world",
				Method: http.MethodGet,
				Body:   "short stuff",
			},
			ExpectedString: "GET /hello-world Body: [short stuff]",
		},
	}
	for _, test := range tests {
		assert.Equal(t, test.ExpectedString, test.Request.String())
	}
}
