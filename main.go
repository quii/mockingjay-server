package main

import (
	"flag"
	"fmt"
	"github.com/quii/mockingjay"
	"io/ioutil"
	"log"
	"net/http"
)

func main() {
	var port = flag.Int("port", 9090, "Port to listen on")
	var configPath = flag.String("config", "", "Path to config json")
	var realURL = flag.String("realURL", "", "Optional: Set this to a URL to check your config against a real server for compatibility")

	flag.Parse()

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
		makeFakeServer(endpoints, *port)
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

func makeFakeServer(endpoints []mockingjay.FakeEndpoint, port int) {
	server := mockingjay.NewServer(endpoints)
	http.Handle("/", server)
	http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
}
