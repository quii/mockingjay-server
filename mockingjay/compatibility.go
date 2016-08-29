package mockingjay

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/johnmuth/xmlcompare"
	"github.com/quii/jsonequaliser"
	"net"
)

// CompatibilityChecker is responsible for checking endpoints are compatible
type CompatibilityChecker struct {
	client http.RoundTripper
	logger *log.Logger
}

// DefaultHTTPTimeoutSeconds is the default http timeout for compatability checks
const DefaultHTTPTimeoutSeconds = 5

// NewCompatabilityChecker creates a new CompatabilityChecker. The httpTimeout refers to the http timeout when making requests to the real server
func NewCompatabilityChecker(logger *log.Logger, httpTimeout time.Duration) *CompatibilityChecker {
	c := new(CompatibilityChecker)

	c.client = &http.Transport{
		Dial: (&net.Dialer{
			Timeout: httpTimeout * time.Second,
		}).Dial,
	}

	c.logger = logger
	return c
}

// CheckCompatibility checks the endpoints against a "real" URL
func (c *CompatibilityChecker) CheckCompatibility(endpoints []FakeEndpoint, realURL string) bool {

	numberOfEndpoints := len(endpoints)

	results := make(chan bool, numberOfEndpoints)

	for _, endpoint := range endpoints {

		if endpoint.CDCDisabled {
			c.logger.Println("! IGNORED", endpoint.Name)
			results <- true
			continue
		}

		go c.compatabilityWorker(endpoint, realURL, results)
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

func (c *CompatibilityChecker) compatabilityWorker(endpoint FakeEndpoint, realURL string, results chan<- bool) {
	errorMessages := c.check(&endpoint, realURL)

	if len(errorMessages) > 0 {
		c.logger.Println(fmt.Sprintf("✗ %s is incompatible with %s", endpoint.String(), realURL))
		for _, msg := range errorMessages {
			c.logger.Println(msg)
		}
		results <- false
	} else {
		c.logger.Println(fmt.Sprintf("✓ %s is compatible with %s", endpoint.String(), realURL))
		results <- true
	}
}

func (c *CompatibilityChecker) check(endpoint *FakeEndpoint, realURL string) (errors []string) {

	request, err := endpoint.Request.AsHTTPRequest(realURL)

	if err != nil {
		errors = append(errors, "Unable to create request from config, maybe try again?")
		return
	}

	response, err := c.client.RoundTrip(request)

	if err != nil {
		errMsg := fmt.Sprintf("Couldn't reach real server: %v", err)
		errors = append(errors, errMsg)
		return
	}

	defer response.Body.Close()

	if response.StatusCode != endpoint.Response.Code {
		errors = append(errors, fmt.Sprintf("Got %d expected %d", response.StatusCode, endpoint.Response.Code))
	}

	body, err := ioutil.ReadAll(response.Body)

	if err != nil {
		errors = append(errors, fmt.Sprintf("Couldn't read response body [%s]", err))
		return
	}

	errors = append(errors, checkBody(string(body), endpoint.Response.Body)...)
	errors = append(errors, findMissingHeaders(endpoint.Response.Headers, response)...)

	return
}

func findMissingHeaders(expectedHeaders map[string]string, response *http.Response) (missing []string) {
	for name, value := range expectedHeaders {
		actualResponseHeader := response.Header.Get(name)
		if actualResponseHeader == "" || actualResponseHeader != value {
			missing = append(missing, fmt.Sprintf("Missing or incorrect header value for %s", name))
		}
	}
	return
}

var (
	msgNotJSON          = "Expected JSON to be returned"
	msgXMLNotCompatible = "XML is not compatible"
	msgExactMatchFailed = "Exact body match did not pass (i did this because the body doesn't look like JSON or XML to me)"
)

func checkBody(downstreamBody string, expectedBody string) (errors []string) {

	if expectedBody == "*" {
		return
	}

	if isJSON(expectedBody) {

		if !isJSON(downstreamBody) {
			return []string{msgNotJSON}
		}
		errMessages, err := jsonequaliser.IsCompatible(expectedBody, downstreamBody)

		if err != nil {
			errors = append(errors, err.Error())
		}

		for k, v := range errMessages {
			errors = append(errors, fmt.Sprintf("JSON err on field '%s' : %s", k, v))
		}

		return
	}

	if isXML(expectedBody) {
		compatible, err := xmlcompare.IsCompatible(expectedBody, downstreamBody)
		if err != nil {
			return []string{err.Error()}
		}
		if !compatible {
			return []string{msgXMLNotCompatible}
		}
		return
	}

	if !strings.Contains(downstreamBody, expectedBody) {
		return []string{msgExactMatchFailed}
	}

	return
}

func isJSON(s string) bool {
	var js interface{}
	return json.Unmarshal([]byte(s), &js) == nil
}

func isXML(s string) bool {
	var x interface{}
	return xml.Unmarshal([]byte(s), &x) == nil
}
