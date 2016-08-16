package main

import (
	"bytes"
	"crypto/md5"
	_ "expvar"
	"fmt"
	"github.com/quii/mockingjay-server/mockingjay"
	"github.com/quii/mockingjay-server/monkey"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

var (
	// ErrCDCFail describes when a fake server is not compatible with a given URL
	ErrCDCFail = fmt.Errorf("At least one endpoint was incompatible with the real URL supplied")
)

type configLoader func(string) ([][]byte, []string, error)
type mockingjayLoader func([]byte) ([]mockingjay.FakeEndpoint, error)

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

func defaultApplication(logger *log.Logger, httpTimeout time.Duration) (app *application) {
	app = new(application)
	app.configLoader = globFileLoader
	app.mockingjayLoader = mockingjay.NewFakeEndpoints
	app.compatabilityChecker = mockingjay.NewCompatabilityChecker(logger, httpTimeout)
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
func (a *application) CreateServer(configPath string, monkeyConfigPath string, debugMode bool, ui http.Handler) (server http.Handler, err error) {
	a.configPath = configPath
	a.monkeyConfigPath = monkeyConfigPath
	endpoints, err := a.loadConfig()

	if err != nil || len(endpoints) == 0 {
		return
	}

	return a.createFakeServer(endpoints, debugMode, ui)
}

// CheckCompatibility will run a MJ config against a realURL to see if it's compatible
func (a *application) CheckCompatibility(configPath string, realURL string) error {
	a.configPath = configPath
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

	configs, _, err := a.configLoader(a.configPath)

	if err != nil {
		return
	}

	if newMD5 := md5.Sum(bytes.Join(configs, []byte{})); newMD5 != a.yamlMD5 {
		a.yamlMD5 = newMD5

		for _, configData := range configs {
			mjEndpoint, err := a.mockingjayLoader(configData)

			if err != nil {
				return nil, err
			}

			endpoints = append(endpoints, mjEndpoint...)
		}
	}

	return
}

type fileUpdater struct {
	path string
}

func (fu *fileUpdater) Write(p []byte) (n int, err error) {
	f, err := os.Create(fu.path)

	if err != nil {
		return 0, err
	}

	n, err = f.Write(p)
	f.Sync()
	return
}

func (a *application) createFakeServer(endpoints []mockingjay.FakeEndpoint, debugMode bool, ui http.Handler) (server http.Handler, err error) {
	go a.PollConfig()

	configFile := fileUpdater{a.configPath}

	a.mjServer = a.mockingjayServerMaker(endpoints, debugMode, &configFile)
	monkeyServer, err := a.monkeyServerMaker(a.mjServer, a.monkeyConfigPath)

	if err != nil {
		return nil, err
	}

	router := http.NewServeMux()

	if ui != nil {
		router.Handle("/mj-admin/", http.StripPrefix("/mj-admin", ui))
	}

	router.Handle("/", monkeyServer)

	return router, nil
}

func globFileLoader(path string) (data [][]byte, paths []string, err error) {
	files, err := filepath.Glob(path)

	if err != nil {
		return nil, nil, fmt.Errorf("Failed to get files from file path (glob) %s, %v", path, err)
	}

	if len(files) == 0 {
		return nil, nil, fmt.Errorf("No files found in path %s", path)
	}

	var configs [][]byte
	for _, file := range files {

		configData, err := ioutil.ReadFile(file)

		if err != nil {
			return nil, nil, err
		}

		configs = append(configs, configData)
	}

	return configs, files, nil
}
