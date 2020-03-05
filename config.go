package main

import (
	"flag"
	"log"
	"os"
	"strconv"
	"time"
)

type appConfig struct {
	logger           *log.Logger
	port             int
	configPath       string
	monkeyConfigPath string
	realURL          string
	httpTimeout      time.Duration
	debugMode        bool
	disablePolling   bool
}

func loadConfig() *appConfig {
	port := flag.Int("port", 9090, "Port to listen on")
	debug := flag.Bool("debug", false, "Print debug statements")
	disablePolling := flag.Bool("disable-polling", false, "Disable file change polling")
	configPath := flag.String("config", "", "Path to config YAML")
	httpTimeout := flag.Int("timeout", 5, "Optional: HTTP timeout when performing CDC")
	monkeyConfigPath := flag.String("monkeyConfig", "", "Optional: Set this to add some monkey business")
	realURL := flag.String("realURL", "", "Optional: Set this to a URL to check your config against a real server for compatibility")

	flag.Parse()

	config := &appConfig{
		logger:           log.New(os.Stdout, "mocking-jay: ", log.Ldate|log.Ltime),
		port:             *port,
		configPath:       *configPath,
		monkeyConfigPath: *monkeyConfigPath,
		realURL:          *realURL,
		httpTimeout:      time.Duration(*httpTimeout),
		debugMode:        *debug,
		disablePolling:   *disablePolling,
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
