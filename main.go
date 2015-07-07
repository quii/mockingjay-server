package main

import (
	"flag"
	"log"
	"os"
	"strconv"
)

var logger *log.Logger

func init() {
	logger = log.New(os.Stdout, "mocking-jay: ", log.Ldate|log.Ltime)
}

func main() {

	envPort := 9090

	if i, err := strconv.Atoi(os.Getenv("PORT")); err == nil {
		envPort = i
	}

	var port = flag.Int("port", envPort, "Port to listen on")
	var configPath = flag.String("config", "", "Path to config YAML")
	var realURL = flag.String("realURL", "", "Optional: Set this to a URL to check your config against a real server for compatibility")
	var monkeyConfigPath = flag.String("monkeyConfig", "", "Optional: Set this to add some monkey business")

	flag.Parse()

	if *configPath == "" {
		flag.Usage()
		os.Exit(1)
	}

	app := defaultApplication(logger)

	err := app.Run(*configPath, *port, *realURL, *monkeyConfigPath)

	if err != nil {
		log.Fatal(err)
	}
}
