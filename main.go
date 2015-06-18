package main

import (
	"flag"
	"log"
	"os"
	"strconv"
)

func main() {

	envPort := 9090
	envPortRaw := os.Getenv("PORT")
	if envPortRaw != "" {
		i, err := strconv.Atoi(envPortRaw)

		if err != nil {
			log.Println("Your PORT environment variable isn't an int, defaulting to 9090")
		} else {
			envPort = i
		}
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

	app := defaultApplication()

	err := app.Run(*configPath, *port, *realURL, *monkeyConfigPath)

	if err != nil {
		log.Fatal(err)
	}
}
