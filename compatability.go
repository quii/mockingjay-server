package main

import (
	"fmt"
	"github.com/quii/mockingjay"
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
	failure := false
	for _, endpoint := range c.endpoints {
		msg, compatible := c.check(&endpoint, realURL)
		log.Println(msg)
		if !compatible {
			failure = true
		}
	}
	return failure
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

	return fmt.Sprintf("%s Tentatively compatible", endpoint), true
}
