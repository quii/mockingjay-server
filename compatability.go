package main

import (
	"fmt"
	"github.com/quii/mockingjay"
	"net/http"
	"strings"
)

// CheckCompatability determines whether an endpoint is compatible with an equivalent server
func CheckCompatability(endpoint *mockingjay.FakeEndpoint, realURL string) (string, bool) {

	request, err := http.NewRequest(endpoint.Request.Method, realURL+endpoint.Request.URI, makeRequestBody(endpoint))

	for headerName, headerValue := range endpoint.Request.Headers {
		request.Header.Add(headerName, headerValue)
	}

	errorMsg := fmt.Sprintf("Endpoint %s is incompatible with %s", endpoint, realURL)

	if err != nil {
		return "Unable to create request from config, maybe try again?", false
	}

	client := http.Client{}

	response, err := client.Do(request)

	if err != nil {
		return fmt.Sprintf("%s - Couldn't reach real server", errorMsg), false
	}

	if response.StatusCode != endpoint.Response.Code {
		return fmt.Sprintf("%s - Got %d expected %d", errorMsg, response.StatusCode, endpoint.Response.Code), false
	}

	return fmt.Sprintf("%s Tentatively compatible", endpoint), true
}

func makeRequestBody(endpoint *mockingjay.FakeEndpoint) *strings.Reader {
	var body *strings.Reader
	if endpoint.Request.Body != "" {
		body = strings.NewReader(endpoint.Request.Body)
	} else {
		body = strings.NewReader("")
	}
	return body
}
