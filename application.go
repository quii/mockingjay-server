package main

import (
	"crypto/md5"
	_ "expvar"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/quii/mockingjay-server/mockingjay"
	"github.com/quii/mockingjay-server/monkey"
)

// CDCFailError provides a summary of errors
type CDCFailError struct {
	ConfigPath string
	URL        string
}

func (c CDCFailError) Error() string {
	return fmt.Sprintf("Config defined at %s has at least one endpoint that was incompatible with %s", c.ConfigPath, c.URL)
}

type configLoader func(string) ([]byte, error)
type mockingjayLoader func([]byte) ([]mockingjay.FakeEndpoint, error)

type compatabilityChecker interface {
	CheckCompatibility(endpoints []mockingjay.FakeEndpoint, realURL string) bool
}

type serverMaker func([]mockingjay.FakeEndpoint, bool) *mockingjay.Server
type monkeyServerMaker func(http.Handler, string) (http.Handler, error)

type application struct {
	configLoader          configLoader
	mockingjayLoader      mockingjayLoader
	compatabilityChecker  compatabilityChecker
	mockingjayServerMaker serverMaker
	monkeyServerMaker     monkeyServerMaker
	logger                *log.Logger
	mjServer              *mockingjay.Server
	configPath            string
	monkeyConfigPath      string
	yamlMD5               [md5.Size]byte
}

func defaultApplication(logger *log.Logger, httpTimeout time.Duration) (app *application) {
	app = new(application)
	app.configLoader = ioutil.ReadFile
	app.mockingjayLoader = mockingjay.NewFakeEndpoints
	app.compatabilityChecker = NewCompatabilityChecker(logger, httpTimeout)
	app.mockingjayServerMaker = mockingjay.NewServer
	app.monkeyServerMaker = monkey.NewServer
	app.logger = logger

	return
}

func (a *application) PollConfig() {
	for range time.Tick(time.Millisecond * 500) {
		endpoints, err := a.loadConfig()
		if err != nil {
			log.Println(err)
		} else if len(endpoints) > 0 {
			a.logger.Println("Reloaded config")
			a.mjServer.Endpoints = endpoints
		}
	}
}

// CreateServer will create a fake server from the configuration found in configPath with optional performance constraints from configutation found in monkeyConfigPath
func (a *application) CreateServer(configPath string, monkeyConfigPath string, debugMode bool) (server http.Handler, err error) {
	a.configPath = configPath
	a.monkeyConfigPath = monkeyConfigPath
	endpoints, err := a.loadConfig()

	if err != nil || len(endpoints) == 0 {
		return
	}

	return a.createFakeServer(endpoints, debugMode)
}

// CheckCompatibility will run a MJ config against a realURL to see if it's compatible
func (a *application) CheckCompatibility(configPath string, realURL string) ([]CDCFailError, error) {
	a.configPath = configPath
	endpoints, err := a.loadConfig()

	if err != nil {
		return nil, err
	}

	if a.compatabilityChecker.CheckCompatibility(endpoints, realURL) {
		a.logger.Println("All endpoints are compatible")
		return nil, nil
	}

	error := CDCFailError{
		ConfigPath: configPath,
		URL:        realURL,
	}
	return []CDCFailError{error}, nil
}

func (a *application) loadConfig() (endpoints []mockingjay.FakeEndpoint, err error) {
	configData, err := a.configLoader(a.configPath)

	if err != nil {
		return nil, err
	}

	if newMD5 := md5.Sum(configData); newMD5 != a.yamlMD5 {
		a.yamlMD5 = newMD5
		endpoints, err = a.mockingjayLoader(configData)
	}
	return
}

func (a *application) createFakeServer(endpoints []mockingjay.FakeEndpoint, debugMode bool) (server http.Handler, err error) {
	go a.PollConfig()
	a.mjServer = a.mockingjayServerMaker(endpoints, debugMode)
	monkeyServer, err := a.monkeyServerMaker(a.mjServer, a.monkeyConfigPath)

	if err != nil {
		return nil, err
	}

	router := http.NewServeMux()
	router.Handle("/", monkeyServer)

	return router, nil
}
