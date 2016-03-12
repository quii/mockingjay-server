package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
)

func main() {

	logger := log.New(os.Stdout, "mocking-jay: ", log.Ldate|log.Ltime)
	envPort := 9090
	port := flag.Int("port", envPort, "Port to listen on")
	configPath := flag.String("config", "", "Path to config YAML")
	realURL := flag.String("realURL", "", "Optional: Set this to a URL to check your config against a real server for compatibility")
	monkeyConfigPath := flag.String("monkeyConfig", "", "Optional: Set this to add some monkey business")

	flag.Parse()

	if i, err := strconv.Atoi(os.Getenv("PORT")); err == nil {
		envPort = i
	}

	if *configPath == "" {
		flag.Usage()
		os.Exit(1)
	}

	app := defaultApplication(logger)
	svr, err := app.Run(*configPath, *realURL, *monkeyConfigPath)

	if err != nil {
		log.Fatal(err)
	} else if svr != nil {
		log.Printf("Listening on port %d", *port)
		err = http.ListenAndServe(fmt.Sprintf(":%d", *port), svr)
		if err != nil {
			log.Fatal("There was a problem starting the mockingjay server on port %d: %v", *port, err.Error())
		}
	}
}
