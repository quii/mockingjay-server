package main

import (
	"fmt"
	"github.com/quii/jsonequaliser"
	"github.com/quii/mockingjay"
	"io/ioutil"
	"log"
	"net/http"
)

// CompatabilityChecker is responsible for checking endpoints are compatible
type CompatabilityChecker struct {
	client    *http.Client
	endpoints []mockingjay.FakeEndpoint
}

// NewCompatabilityChecker creates a new CompatabilityChecker
func NewCompatabilityChecker(endpoints []mockingjay.FakeEndpoint) *CompatabilityChecker {
	c := new(CompatabilityChecker)
	c.endpoints = endpoints
	c.client = &http.Client{}
	return c
}

// CheckCompatability checks the endpoints against a "real" URL
func (c *CompatabilityChecker) CheckCompatability(realURL string) bool {

	results := make(chan bool, len(c.endpoints))

	for _, endpoint := range c.endpoints {
		go func() {
			msg, compatible := c.check(&endpoint, realURL)
			log.Println(msg)
			results <- compatible
		}()
	}

	allCompatible := true
	for i := 0; i < len(c.endpoints); i++ {
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
		return fmt.Sprintf("%s - Got %d expected %d", errorMsg, response.StatusCode, endpoint.Response.Code), false
	}

	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)

	if err != nil {
		return fmt.Sprintln("Couldn't read response body [%s]", err), false
	}

	bodyCompatible, err := checkBody(string(body), endpoint.Response.Body)

	if err != nil {
		return fmt.Sprintf("There was a problem checking the compatibility of the body", err), false
	}

	if !bodyCompatible {
		return fmt.Sprintf("Body [%s] was not compatible with config body [%s]", string(body), endpoint.Response.Body), false
	}

	return fmt.Sprintf("%s Tentatively compatible", endpoint), true
}

func checkBody(downstreamBody string, expectedBody string) (bool, error) {
	if isJSON(downstreamBody) {
		return jsonequaliser.IsCompatible(expectedBody, downstreamBody)
	}
	return true, nil
}

//todo: This is clearly flaky and stupid
func isJSON(x string) bool {
	return x[0] == '{'
}
