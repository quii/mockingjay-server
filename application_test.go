package main

import (
	"errors"
	"log"
	"net/http"
	"testing"

	"github.com/quii/mockingjay-server/mockingjay"
	"github.com/stretchr/testify/assert"
)

const someMonkeyConfigString = "Hello, world"

func TestItFailsWhenTheConfigFileCantBeLoaded(t *testing.T) {
	app := testApplication()
	app.configLoader = failingIOUtil

	_, err := app.CreateServer("mockingjay config path", "")

	assert.NotNil(t, err)
	assert.Equal(t, err, errIOError)
}

func TestItFailsWhenTheConfigIsInvalid(t *testing.T) {
	app := testApplication()
	app.mockingjayLoader = failingMockingjayLoader

	_, err := app.CreateServer("mockingjay config path", "")

	assert.NotNil(t, err, "Didnt get an error when the mockingjay config failed to load")
	assert.Equal(t, err, errMJLoaderError)
}

func TestItFailsWhenTheMonkeyConfigIsInvalid(t *testing.T) {
	app := testApplication()

	_, err := app.CreateServer("mockingjay config path", "monkey config path")

	assert.NotNil(t, err, "Didnt get an error when the monkey config failed to load")
	assert.Equal(t, err, errMonkeyLoadError)
}

func TestItReturnsCDCErrIfCompatabilityFails(t *testing.T) {
	app := new(application)
	app.configLoader = passingIOUtil
	app.mockingjayLoader = passingMockingjayLoader

	app.compatabilityChecker = fakeCompatabilityChecker{passes: false}

	err := app.CheckCompatibility("mj config path", "http://someurl")

	assert.NotNil(t, err, "Didn't get an error when compatability fails")
	assert.Equal(t, err, ErrCDCFail)
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
