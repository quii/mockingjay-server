package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

func main() {

	config := loadConfig()
	app := defaultApplication(config.logger, config.httpTimeout, config.configPath)

	if config.realURL != "" {
		err := app.CheckCompatibility(config.realURL)

		if err != nil {
			log.Fatal(err)
		}
	} else {
		svr, err := app.CreateServer(config.monkeyConfigPath, config.debugMode, config.disablePolling)

		if err != nil {
			log.Fatal(err)
		} else {
			config.logger.Printf("Listening on port %d", config.port)
			webServer := newServer(svr, config.port)
			err = webServer.ListenAndServe()
			if err != nil {
				msg := fmt.Sprintf("There was a problem starting the mockingjay server on port %d: %s", config.port, err.Error())
				config.logger.Fatal(msg)
			}
		}
	}
}

func newServer(router http.Handler, port int) *http.Server {
	return &http.Server{
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       30 * time.Second,
		ReadHeaderTimeout: 20 * time.Second,
		Handler:           http.TimeoutHandler(Recovery(router), 5*time.Second, "Timed out!"),
		Addr:              fmt.Sprintf(":%d", port),
	}
}

func Recovery(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			err := recover()
			if err != nil {
				msg := fmt.Sprintf("Oh dear, mockingjay has had a panic %#v", err)
				log.Println(msg)
				w.WriteHeader(http.StatusInternalServerError)
				fmt.Fprint(w, msg)
			}
		}()
		next.ServeHTTP(w, r)
	})
}
