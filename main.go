package main

import (
	"flag"
	"fmt"
	"github.com/quii/mockingjay-server/mockingjay"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

func main() {
	var port = flag.Int("port", 9090, "Port to listen on")
	var configPath = flag.String("config", "", "Path to config YAML")
	var realURL = flag.String("realURL", "", "Optional: Set this to a URL to check your config against a real server for compatibility")
	var monkeyConfigPath = flag.String("monkeyConfig", "", "Optional: Set this to add some monkey business")

	flag.Parse()

	if *configPath == "" {
		flag.Usage()
		os.Exit(1)
	}

	config, err := ioutil.ReadFile(*configPath)

	if err != nil {
		log.Fatalf("Problem occured when trying to open the config file: %v", err)
	}

	endpoints, err := mockingjay.NewFakeEndpoints(config)

	if err != nil {
		log.Fatalf("Problem occured when trying to create a server from the config: %v ", err)
	}

	if *realURL != "" {
		checkEndpoints(endpoints, *realURL)
	} else {
		log.Printf("Serving %d endpoints defined from %s on port %d", len(endpoints), *configPath, *port)
		makeFakeServer(endpoints, *port, loadMonkeyConfig(*monkeyConfigPath))
	}

}

func checkEndpoints(endpoints []mockingjay.FakeEndpoint, realURL string) {
	checker := NewCompatabilityChecker(endpoints)

	if checker.CheckCompatability(realURL) {
		log.Println("All endpoints are compatible")
	} else {
		log.Fatal("At least one endpoint was incompatible with the real URL supplied")
	}
}

func makeFakeServer(endpoints []mockingjay.FakeEndpoint, port int, monkeyConfig []behaviour) {

	var server http.Handler
	if len(monkeyConfig) > 0 {
		server = NewMonkeyServer(mockingjay.NewServer(endpoints), monkeyConfig)
	} else {
		server = mockingjay.NewServer(endpoints)
	}

	http.Handle("/", server)
	err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
	if err != nil {
		log.Fatalf("There was a problem starting the mockingjay server on port %d: %v", port, err)
	}
}
