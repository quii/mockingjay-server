package mockingjay

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"testing/quick"

	"gopkg.in/yaml.v2"
)

func TestItMarshalsAndUnmarshalsJSONCorrectly(t *testing.T) {
	assertion := func(endpoint FakeEndpoint) bool {
		data, err := json.Marshal(endpoint)

		if err != nil {
			t.Log("Had a problem marshalling config into JSON", err)
			return false
		}

		var parsed FakeEndpoint
		err = json.Unmarshal(data, &parsed)

		if err != nil {
			t.Log(string(data))
			t.Log(parsed)
			t.Log("Couldn't re-parse JSON into the config", err)
		}

		return reflect.DeepEqual(endpoint, parsed)
	}

	config := quick.Config{
		MaxCount: 1000,
	}

	if err := quick.Check(assertion, &config); err != nil {
		t.Error(err)
	}
}

func TestItMarshalsAndUnmarshalsYAMLCorrectly(t *testing.T) {
	assertion := func(endpoint FakeEndpoint) bool {
		data, err := yaml.Marshal(endpoint)

		if err != nil {
			t.Log("Had a problem marshalling config into JSON", err)
			return false
		}

		var parsed FakeEndpoint
		err = yaml.Unmarshal(data, &parsed)

		if err != nil {
			t.Log("Couldn't re-parse JSON into the config", err)
		}

		return reflect.DeepEqual(endpoint, parsed)
	}

	config := quick.Config{
		MaxCount: 1000,
	}

	if err := quick.Check(assertion, &config); err != nil {
		t.Error(err)
	}
}

func TestItIsAlwaysCompatibleWithItself(t *testing.T) {

	compatabilityChecker := NewCompatabilityChecker(noopLogger, 1)

	assertion := func(endpoint FakeEndpoint) bool {

		// Start an MJ server with the random configuration
		mjSvr := NewServer([]FakeEndpoint{endpoint}, false, ioutil.Discard)
		svr := httptest.NewServer(http.HandlerFunc(mjSvr.ServeHTTP))
		defer svr.Close()

		// Run CDC against "itself". An MJ server should always be compatible with itself.
		errors := compatabilityChecker.check(&endpoint, svr.URL)

		if len(errors) > 0 {
			t.Logf("Not compatible with itself %+v", endpoint)
			for _, err := range errors {
				t.Log(err)
			}
		}

		return len(errors) == 0
	}

	config := quick.Config{
		MaxCount: 1000,
	}

	if err := quick.Check(assertion, &config); err != nil {
		t.Error(err)
	}
}
