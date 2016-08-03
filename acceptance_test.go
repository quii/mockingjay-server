package main

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

var (
	app      *application
	mjServer http.Handler
)

const testYAMLPath = "examples/example.yaml"

func init() {
	app = defaultApplication(log.New(ioutil.Discard, "", 0))
	svr, err := app.CreateServer(testYAMLPath, "")

	if err != nil {
		log.Fatal("Couldn't load MJ config from", testYAMLPath)
	}

	mjServer = svr
}

func TestIssue42(t *testing.T) {
	failApp := defaultApplication(log.New(ioutil.Discard, "", 0))
	failSvr, _ := failApp.CreateServer("examples/issue42.yaml", "")
	svr := httptest.NewServer(failSvr)
	defer svr.Close()

	reqBody := []byte(`{"query":{"match_all":{}}}`)
	req, _ := http.NewRequest("POST", svr.URL+"/profile/validate-query", bytes.NewBuffer(reqBody))
	req.Header["Content-Type"] = []string{"application/json"}
	req.Header["Content-Butt"] = []string{"application/json"}

	client := http.Client{}
	res, _ := client.Do(req)

	defer res.Body.Close()

	body, _ := ioutil.ReadAll(res.Body)

	if res.StatusCode != http.StatusOK {
		t.Error("WTF", res, string(body))
	}
}

func TestItLaunchesServersAndIsCompatibleWithItsOwnConfig(t *testing.T) {
	svr := httptest.NewServer(mjServer)
	defer svr.Close()
	assert.NoError(t, app.CheckCompatibility(testYAMLPath, svr.URL))
}

func TestItListsRequestsItHasReceived(t *testing.T) {
	svr := httptest.NewServer(mjServer)
	defer svr.Close()

	http.Get(svr.URL + "/hello")

	res, err := http.Get(svr.URL + "/requests")

	assert.Nil(t, err)
	assert.Equal(t, res.StatusCode, http.StatusOK)
}

func TestANewEndpointCanBeAdded(t *testing.T) {
	svr := httptest.NewServer(mjServer)
	defer svr.Close()

	newEndpointJSON := `
	{
	  "Name": "Test endpoint",
	  "CDCDisabled": false,
	  "Request": {
	    "URI": "/new-endpoint",
	    "Method": "GET",
	    "Headers": null,
	    "Body": ""
	  },
	  "Response": {
	    "Code": 200,
	    "Body": "{\"message\": \"hello, world\"}",
	    "Headers": {
	      "content-type": "text\/json"
	    }
	  }
	}
	`
	newEndpointURL := svr.URL + "/mj-new-endpoint"
	res, err := http.Post(newEndpointURL, "application/json", strings.NewReader(newEndpointJSON))

	assert.Nil(t, err)
	assert.Equal(t, res.StatusCode, http.StatusCreated)

	newEndPointResponse, err := http.Get(svr.URL + `/new-endpoint`)

	assert.Nil(t, err)
	assert.Equal(t, newEndPointResponse.StatusCode, http.StatusOK)

	newEndpointBody, _ := ioutil.ReadAll(newEndPointResponse.Body)

	assert.Equal(t, string(newEndpointBody), `{"message": "hello, world"}`)
}
