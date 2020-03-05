package main

import (
	"crypto/md5"
	_ "expvar"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/quii/mockingjay-server/mockingjay"
	"github.com/quii/mockingjay-server/monkey"
)

var (
	// ErrCDCFail describes when a fake server is not compatible with a given URL
	ErrCDCFail = fmt.Errorf("at least one endpoint was incompatible with the real URL supplied")
)

type mockingjayLoader func(closer io.ReadCloser) ([]mockingjay.FakeEndpoint, error)

type compatabilityChecker interface {
	CheckCompatibility(endpoints []mockingjay.FakeEndpoint, realURL string) bool
}

type serverMaker func([]mockingjay.FakeEndpoint, bool, io.Writer) *mockingjay.Server
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

func defaultApplication(logger *log.Logger, httpTimeout time.Duration, configPath string) (app *application) {
	app = new(application)
	app.configPath = configPath

	if isURL(configPath) {
		app.configLoader = urlLoader{}
	} else {
		app.configLoader = globFileLoader{}
	}

	app.mockingjayLoader = mockingjay.NewFakeEndpoints
	app.compatabilityChecker = mockingjay.NewCompatabilityChecker(logger, httpTimeout)
	app.mockingjayServerMaker = mockingjay.NewServer
	app.monkeyServerMaker = monkey.NewServer
	app.logger = logger

	return
}

func (a *application) PollConfig() {
	if _, pollable := a.configLoader.(pollable); pollable {
		for range time.Tick(time.Millisecond * 500) {
			endpoints, err := a.loadConfig()
			if err != nil {
				log.Println(err)
			} else {
				a.mjServer.Endpoints = endpoints
			}
		}
	} else {
		a.logger.Println("config loader is not pollable, restart app to update config")
	}
}

// CreateServer will create a fake server from the configuration found in configPath with optional performance constraints from configutation found in monkeyConfigPath
func (a *application) CreateServer(monkeyConfigPath string, debugMode bool, disablePolling bool) (server http.Handler, err error) {
	a.monkeyConfigPath = monkeyConfigPath
	endpoints, err := a.loadConfig()

	if err != nil || len(endpoints) == 0 {
		return
	}

	return a.createFakeServer(endpoints, debugMode, disablePolling)
}

// CheckCompatibility will run a MJ config against a realURL to see if it's compatible
func (a *application) CheckCompatibility(realURL string) error {
	endpoints, err := a.loadConfig()

	if err != nil {
		return err
	}

	if a.compatabilityChecker.CheckCompatibility(endpoints, realURL) {
		a.logger.Println("All endpoints are compatible")
		return nil
	}

	return ErrCDCFail
}

func (a *application) loadConfig() (endpoints []mockingjay.FakeEndpoint, err error) {

	configs, _, err := a.configLoader.Load(a.configPath)

	if err != nil {
		return
	}

	for _, conf := range configs {
		mjEndpoint, err := a.mockingjayLoader(conf)

		if err != nil {
			return nil, err
		}

		conf.Close()

		endpoints = append(endpoints, mjEndpoint...)
	}

	return
}

/*
Giovanni Bajo [5:09 PM]
it doesn't respect the semantics of io.Writer

[5:09]
like you wouldn’t be able to compose it with a gzip.Writer

[5:10]
so I’m not sure you’re doing yourself a favor in implementing the io.Writer interface, it’s prone to mistakes
*/
type fileUpdater struct {
	path string
}

func (fu *fileUpdater) Write(p []byte) (n int, err error) {
	f, err := os.Create(fu.path)

	if err != nil {
		return 0, err
	}

	n, err = f.Write(p)
	_ = f.Sync()
	return
}

func (a *application) createFakeServer(endpoints []mockingjay.FakeEndpoint, debugMode bool, disablePolling bool) (server http.Handler, err error) {
	if !disablePolling {
		go a.PollConfig()
	}

	configFile := fileUpdater{a.configPath}

	a.mjServer = a.mockingjayServerMaker(endpoints, debugMode, &configFile)
	monkeyServer, err := a.monkeyServerMaker(a.mjServer, a.monkeyConfigPath)

	if err != nil {
		return nil, err
	}

	router := http.NewServeMux()

	router.Handle("/", monkeyServer)

	return router, nil
}
