package main

import (
	"errors"
	"log"
	"net/http"
	"testing"

	"fmt"
	"github.com/quii/mockingjay-server/mockingjay"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http/httptest"
)

const someMonkeyConfigString = "Hello, world"

func TestCompatabilityWithWildcards(t *testing.T) {

	wildcardPath := "examples/issue40/1.yaml"

	app := defaultApplication(log.New(ioutil.Discard, "", log.Ldate|log.Ltime), 1)
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Log("Got request", r.URL)
		if r.URL.String() == "/hello" {
			fmt.Fprint(w, "world")
		}
	}))

	defer svr.Close()
	cdcErrors, err := app.CheckCompatibility(wildcardPath, svr.URL)
	assert.NoError(t, err)
	assert.Empty(t, cdcErrors)
}

func TestItFailsWhenTheConfigFileCantBeLoaded(t *testing.T) {
	app := testApplication()
	app.configLoader = failingIOUtil

	configPath := "mockingjay config path"
	_, err := app.CreateServer(configPath, "", false)

	assert.NotNil(t, err)
	assert.Equal(t, err, errIOError)

	_, err = app.CheckCompatibility(configPath, "some url")
	assert.NotNil(t, err)
	assert.Equal(t, err, errIOError)
}

func TestItFailsWhenTheConfigIsInvalid(t *testing.T) {
	app := testApplication()
	app.mockingjayLoader = failingMockingjayLoader

	_, err := app.CreateServer("mockingjay config path", "", false)

	assert.NotNil(t, err, "Didnt get an error when the mockingjay config failed to load")
	assert.Equal(t, err, errMJLoaderError)
}

func TestCompatFailsWhenConfigIsInvalid(t *testing.T) {
	app := testApplication()
	app.mockingjayLoader = failingMockingjayLoader

	_, err := app.CheckCompatibility("mockingjay config path", "some url")

	assert.NotNil(t, err, "Didnt get an error when the mockingjay config failed to load")
	assert.Equal(t, err, errMJLoaderError)
}

func TestItFailsWhenTheMonkeyConfigIsInvalid(t *testing.T) {
	app := testApplication()

	_, err := app.CreateServer("mockingjay config path", "monkey config path", false)

	assert.NotNil(t, err, "Didnt get an error when the monkey config failed to load")
	assert.Equal(t, err, errMonkeyLoadError)
}

func TestItReturnsCDCErrorIfCompatabilityFails(t *testing.T) {
	app := new(application)
	app.configLoader = passingIOUtil
	app.mockingjayLoader = passingMockingjayLoader

	app.compatabilityChecker = fakeCompatabilityChecker{passes: false}

	cdcError, _ := app.CheckCompatibility("mj config path", "http://someurl")

	assert.NotNil(t, cdcError, "Didn't get an error when compatability fails")
	assert.IsType(t, []CDCFailError{}, cdcError)
}

func testApplication() *application {
	app := new(application)
	app.configLoader = passingIOUtil
	app.mockingjayLoader = passingMockingjayLoader
	app.mockingjayServerMaker = mockingjay.NewServer
	app.monkeyServerMaker = failingMonkeyServerMaker
	return app
}

func testMockingJayConfig() []mockingjay.FakeEndpoint {

	m, err := mockingjay.NewFakeEndpoints([]byte(testYAML("hello, world")))

	if err != nil {
		log.Fatal(err)
	}

	return m
}

func passingIOUtil(path string) ([]byte, error) {
	return []byte(someMonkeyConfigString), nil
}

var errIOError = errors.New("Couldn't load err from FS")

func failingIOUtil(path string) ([]byte, error) {
	return nil, errIOError
}

var errMJLoaderError = errors.New("Couldnt load mj file")

func failingMockingjayLoader([]byte) ([]mockingjay.FakeEndpoint, error) {
	return nil, errMJLoaderError
}

func passingMockingjayLoader([]byte) ([]mockingjay.FakeEndpoint, error) {
	return testMockingJayConfig(), nil
}

var errMonkeyLoadError = errors.New("Couldn't load monkey file")

func failingMonkeyServerMaker(http.Handler, string) (http.Handler, error) {
	return nil, errMonkeyLoadError
}

type fakeCompatabilityChecker struct {
	passes bool
}

func (f fakeCompatabilityChecker) CheckCompatibility(endpoints []mockingjay.FakeEndpoint, realURL string) bool {
	return f.passes
}
