package main

import (
	"encoding/json"
	"fmt"
	"github.com/quii/jsonequaliser"
	"github.com/quii/mockingjay-server/mockingjay"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

// CompatabilityChecker is responsible for checking endpoints are compatible
type CompatabilityChecker struct {
	client *http.Client
}

// NewCompatabilityChecker creates a new CompatabilityChecker
func NewCompatabilityChecker() *CompatabilityChecker {
	c := new(CompatabilityChecker)
	c.client = &http.Client{}
	c.client.Timeout = 5 * time.Second
	return c
}

// CheckCompatability checks the endpoints against a "real" URL
func (c *CompatabilityChecker) CheckCompatability(endpoints []mockingjay.FakeEndpoint, realURL string) bool {

	numberOfEndpoints := len(endpoints)

	results := make(chan bool, numberOfEndpoints)

	for _, endpoint := range endpoints {
		go func(ep mockingjay.FakeEndpoint) {
			msg, compatible := c.check(&ep, realURL)
			log.Println(msg)
			results <- compatible
		}(endpoint)
	}

	allCompatible := true
	for i := 0; i < numberOfEndpoints; i++ {
		compatible := <-results
		if !compatible {
			allCompatible = false
		}
	}
	return allCompatible
}

func (c *CompatabilityChecker) check(endpoint *mockingjay.FakeEndpoint, realURL string) (string, bool) {

	request, err := endpoint.Request.AsHTTPRequest(realURL)

	errorMsg := fmt.Sprintf("%s is incompatible with %s", endpoint, realURL)

	if err != nil {
		return "Unable to create request from config, maybe try again?", false
	}

	response, err := c.client.Do(request)

	if err != nil {
		return fmt.Sprintf("%s - Couldn't reach real server", errorMsg), false
	}

	if response.StatusCode != endpoint.Response.Code {
		return fmt.Sprintf("%s - Got %d expected %d (%s)", errorMsg, response.StatusCode, endpoint.Response.Code, request.URL), false
	}

	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)

	if err != nil {
		return fmt.Sprintln("%s - Couldn't read response body [%s]", errorMsg, err), false
	}

	bodyCompatible, err := checkBody(string(body), endpoint.Response.Body)

	if err != nil {
		return fmt.Sprintf("%s - There was a problem checking the compatibility of the body", errorMsg, err), false
	}

	if !bodyCompatible {
		return fmt.Sprintf("%s - Body [%s] was not compatible with config body [%s]", errorMsg, string(body), endpoint.Response.Body), false
	}

	return fmt.Sprintf("%s  ✔", endpoint), true
}

func checkBody(downstreamBody string, expectedBody string) (bool, error) {
	if isJSON(expectedBody) {
		return jsonequaliser.IsCompatible(expectedBody, downstreamBody)
	}

	if expectedBody == "*" {
		return true, nil
	}

	if downstreamBody != expectedBody {
		return false, nil
	}
	return true, nil
}

func isJSON(s string) bool {
	var js interface{}
	return json.Unmarshal([]byte(s), &js) == nil
}
