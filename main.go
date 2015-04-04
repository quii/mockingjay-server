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
	flag.Parse()

	config, err := ioutil.ReadFile(*configPath)

	if err != nil {
		log.Fatal(err)
	}

	endpoints, err := mockingjay.NewFakeEndpoints(string(config))

	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Serving %d endpoints defined from %s on port %d", len(endpoints), *configPath, *port)

	server := mockingjay.NewServer(endpoints)

	// Mount it just like any other server
	http.Handle("/", server)
	http.ListenAndServe(fmt.Sprintf(":%d", *port), nil)
}
