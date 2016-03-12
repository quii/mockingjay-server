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

var (
	// ErrCDCFail describes when a fake server is not compatible with a given URL
	ErrCDCFail = fmt.Errorf("At least one endpoint was incompatible with the real URL supplied")
)

type configLoader func(string) ([]byte, error)
type mockingjayLoader func([]byte) ([]mockingjay.FakeEndpoint, error)

type compatabilityChecker interface {
	CheckCompatability(endpoints []mockingjay.FakeEndpoint, realURL string) bool
}

type serverMaker func([]mockingjay.FakeEndpoint) *mockingjay.Server
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

func defaultApplication(logger *log.Logger) (app *application) {
	app = new(application)
	app.configLoader = ioutil.ReadFile
	app.mockingjayLoader = mockingjay.NewFakeEndpoints
	app.compatabilityChecker = NewCompatabilityChecker(logger)
	app.mockingjayServerMaker = mockingjay.NewServer
	app.monkeyServerMaker = monkey.NewServer
	app.logger = logger

	return
}

func (a *application) PollConfig() {
	ticker := time.NewTicker(time.Millisecond * 500)

	for range ticker.C {
		endpoints, err := a.loadConfig()
		if err != nil {
			log.Println(err)
		} else if len(endpoints) > 0 {
			a.logger.Println("Reloaded config")
			a.mjServer.Endpoints = endpoints
		}
	}
}

// Run will create a fake server from the configuration found in configPath with optional performance constraints from configutation found in monkeyConfigPath. If the realURL is supplied then it will not launch as a server and instead will check the config against the URL to see if it is structurally compatible (CDC mode)
func (a *application) Run(configPath string, realURL string, monkeyConfigPath string) (server http.Handler, err error) {
	a.configPath = configPath
	a.monkeyConfigPath = monkeyConfigPath
	endpoints, err := a.loadConfig()

	if err != nil || len(endpoints) == 0 {
		return
	}

	inCheckCompatabilityMode := realURL != ""

	if inCheckCompatabilityMode {
		err = a.checkCompatability(endpoints, realURL)
		return nil, err
	}

	return a.createFakeServer(endpoints)
}

func (a *application) checkCompatability(endpoints []mockingjay.FakeEndpoint, realURL string) error {
	if a.compatabilityChecker.CheckCompatability(endpoints, realURL) {
		a.logger.Println("All endpoints are compatible")
		return nil
	}
	return ErrCDCFail
}

func (a *application) loadConfig() (endpoints []mockingjay.FakeEndpoint, err error) {

	configData, err := a.configLoader(a.configPath)

	if err != nil {
		return
	}

	if newMD5 := md5.Sum(configData); newMD5 != a.yamlMD5 {
		a.yamlMD5 = newMD5
		endpoints, err = a.mockingjayLoader(configData)
	}
	return
}

func (a *application) createFakeServer(endpoints []mockingjay.FakeEndpoint) (server http.Handler, err error) {
	go a.PollConfig()
	a.mjServer = a.mockingjayServerMaker(endpoints)
	monkeyServer, err := a.monkeyServerMaker(a.mjServer, a.monkeyConfigPath)

	if err != nil {
		return nil, err
	}

	router := http.NewServeMux()
	router.Handle("/", monkeyServer)

	return router, nil
}
