package main

import (
	"flag"
	"log"
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

	app := defaultApplication()

	err := app.Run(*configPath, *port, *realURL, *monkeyConfigPath)

	if err != nil {
		log.Fatal(err)
	}
}
