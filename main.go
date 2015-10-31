package main

import (
	"flag"
	"log"
	"os"
	"strconv"
)

var (
	logger           *log.Logger
	envPort          = 9090
	port             = flag.Int("port", envPort, "Port to listen on")
	configPath       = flag.String("config", "", "Path to config YAML")
	realURL          = flag.String("realURL", "", "Optional: Set this to a URL to check your config against a real server for compatibility")
	monkeyConfigPath = flag.String("monkeyConfig", "", "Optional: Set this to add some monkey business")
)

func init() {
	logger = log.New(os.Stdout, "mocking-jay: ", log.Ldate|log.Ltime)
	if i, err := strconv.Atoi(os.Getenv("PORT")); err == nil {
		envPort = i
	}
}

func main() {

	flag.Parse()

	if *configPath == "" {
		flag.Usage()
		os.Exit(1)
	}

	app := defaultApplication(logger)

	if err := app.Run(*configPath, *port, *realURL, *monkeyConfigPath); err != nil {
		log.Fatal(err)
	}
}
