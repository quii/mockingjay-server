//go:generate go-bindata-assetfs -ignore=node_modules -prefix "ui/src/client/" -pkg $GOPACKAGE ui/src/client/public/...
package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
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
		svr, err := app.CreateServer(config.configPath, config.monkeyConfigPath, config.debugMode, getUIServer())

		if err != nil {
			log.Fatal(err)
		} else {
			config.logger.Printf("Listening on port %d", config.port)
			config.logger.Printf("Admin on http://localhost:%d/mj-admin", config.port)

			err = http.ListenAndServe(fmt.Sprintf(":%d", config.port), svr)
			if err != nil {
				msg := fmt.Sprintf("There was a problem starting the mockingjay server on port %d: %s", config.port, err.Error())
				config.logger.Fatal(msg)
			}
		}
	}
}

func getUIServer() http.Handler {
	if os.Getenv("ENV") == "LOCAL" {
		log.Println("Detected local dev mode, serving files from /ui")
		return http.FileServer(http.Dir("./ui/src/client/public"))
	}
	return http.FileServer(assetFS())
}
