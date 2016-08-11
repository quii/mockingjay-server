//go:generate go-bindata-assetfs -ignore=node_modules -pkg $GOPACKAGE ui/...
package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {

	config := loadConfig()
	app := defaultApplication(config.logger, config.httpTimeout)

	if config.realURL != "" {
		err := app.CheckCompatibility(config.configPath, config.realURL)

		if err != nil {
			log.Fatal(err)
		}
	} else {
		svr, err := app.CreateServer(config.configPath, config.monkeyConfigPath, config.debugMode)

		if err != nil {
			log.Fatal(err)
		} else {
			go ui(config)

			config.logger.Printf("Listening on port %d", config.port)
			err = http.ListenAndServe(fmt.Sprintf(":%d", config.port), svr)
			if err != nil {
				msg := fmt.Sprintf("There was a problem starting the mockingjay server on port %d: %s", config.port, err.Error())
				config.logger.Fatal(msg)
			}
		}
	}
}

func ui(config *appConfig) {
	config.logger.Printf("UI served on port %d", config.uiPort)
	svr := http.NewServeMux()
	svr.Handle("/", http.FileServer(assetFS()))
	err := http.ListenAndServe(fmt.Sprintf(":%d", config.uiPort), svr)

	if err != nil {
		log.Fatal(err)
	}
}
