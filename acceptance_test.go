package main

import (
	"bytes"
	"github.com/quii/mockingjay-server/mockingjay"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

const testYAMLPath = "examples/example.yaml"

func TestHeadersArentTooStrict(t *testing.T) {
	app, _ := buildMJ(t, "examples/issue42.yaml")
	failSvr, _ := app.CreateServer("", false)
	svr := httptest.NewServer(failSvr)
	defer svr.Close()

	reqBody := []byte(`{"query":{"match_all":{}}}`)
	req, _ := http.NewRequest("POST", svr.URL+"/profile/validate-query", bytes.NewBuffer(reqBody))
	req.Header["Content-Type"] = []string{"application/json"}
	req.Header["Content-Butt"] = []string{"application/json"}

	client := http.Client{}
	res, err := client.Do(req)

	if err != nil {
		t.Fatal(err)
	}

	defer res.Body.Close()

	body, _ := ioutil.ReadAll(res.Body)

	if res.StatusCode != http.StatusOK {
		t.Error("WTF", res, string(body))
	}
}

func TestItLaunchesServersAndIsCompatibleWithItsOwnConfig(t *testing.T) {
	app, mjServer := buildMJ(t, testYAMLPath)
	svr := httptest.NewServer(mjServer)
	defer svr.Close()
	cdcErrors := app.CheckCompatibility(svr.URL)
	assert.Empty(t, cdcErrors, "There should be no CDC errors with itself")
}

func TestItCanReadConfigFromAURL(t *testing.T) {
	config, _ := ioutil.ReadFile(testYAMLPath)
	consumerCDCSvr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write(config)
	}))

	app, mjServer := buildMJ(t, consumerCDCSvr.URL)
	producerSvr := httptest.NewServer(mjServer)

	defer producerSvr.Close()
	defer consumerCDCSvr.Close()

	cdcErrors := app.CheckCompatibility(producerSvr.URL)
	assert.Nil(t, cdcErrors, "There should be no CDC errors with itself")
}

func TestItListsRequestsItHasReceived(t *testing.T) {
	_, mjServer := buildMJ(t, testYAMLPath)

	svr := httptest.NewServer(mjServer)
	defer svr.Close()

	http.Get(svr.URL + "/hello")

	res, err := http.Get(svr.URL + "/requests")

	assert.Nil(t, err)
	assert.Equal(t, res.StatusCode, http.StatusOK)
}

func TestANewEndpointCanBeAdded(t *testing.T) {
	_, mjServer := buildMJ(t, testYAMLPath)

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

func buildMJ(t *testing.T, config string) (app *application, mjServer http.Handler) {
	t.Helper()
	app = defaultApplication(log.New(ioutil.Discard, "", 0), mockingjay.DefaultHTTPTimeoutSeconds, config)
	svr, err := app.CreateServer("", false)

	if err != nil {
		t.Fatal("Couldn't load MJ config from ", config)
	}

	mjServer = svr
	return
}
