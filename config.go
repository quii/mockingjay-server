package main

import (
	"flag"
	"log"
	"os"
	"strconv"
)

type appConfig struct {
	logger           *log.Logger
	port             int
	configPath       string
	monkeyConfigPath string
	realURL          string
}

func loadConfig() *appConfig {
	port := flag.Int("port", 9090, "Port to listen on")
	configPath := flag.String("config", "", "Path to config YAML")
	monkeyConfigPath := flag.String("monkeyConfig", "", "Optional: Set this to add some monkey business")
	realURL := flag.String("realURL", "", "Optional: Set this to a URL to check your config against a real server for compatibility")

	flag.Parse()

	config := &appConfig{
		logger:           log.New(os.Stdout, "mocking-jay: ", log.Ldate|log.Ltime),
		port:             *port,
		configPath:       *configPath,
		monkeyConfigPath: *monkeyConfigPath,
		realURL:          *realURL,
	}

	if i, err := strconv.Atoi(os.Getenv("PORT")); err == nil {
		config.port = i
	}

	if config.configPath == "" {
		flag.Usage()
		os.Exit(1)
	}

	return config
}
