// Package monkey lets you wrap bad behaviour around your http.Handlers, such as request delays, incorrect responses, bodies and returning garbage from configuration. By using monkey you can simulate the unpredictable nature of calling services over HTTP.
package monkey

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

// server wraps around a http.Handler and adds destructive behaviour (monkey business) based on the behaviours passed in
type server struct {
	delegate   http.Handler
	behaviours []behaviour
	randomiser randomiser
}

// NewServerFromYAML creates a http.Handler which wraps monkey business defined from YAML around it, to return a new http.Handler. If the YAML is invalid, it will return an error.
func NewServerFromYAML(server http.Handler, YAML []byte) (http.Handler, error) {
	behaviours, err := monkeyConfigFromYAML(YAML)

	if err != nil {
		return nil, errors.New(fmt.Sprintf("Problem occured when trying to parse the config file: %v", err))
	}

	log.Println("Monkey config loaded")
	for _, b := range behaviours {
		log.Println(b)
	}

	return newServerFromBehaviour(server, behaviours), nil
}

// NewServer creates a http.Handler which wraps it's monkey business around it, to return a new http.Handler. If no behaviours are defined in the config it will return the original handler, otherwise an error
func NewServer(server http.Handler, configPath string) (http.Handler, error) {
	if configPath == "" {
		return server, nil
	}

	data, err := ioutil.ReadFile(configPath)

	if err != nil {
		return nil, errors.New(fmt.Sprintf("Problem occured when trying to read the config file: %v", err))
	}

	return NewServerFromYAML(server, data)

}

func newServerFromBehaviour(degegate http.Handler, behaviours []behaviour) http.Handler {
	return &server{degegate, behaviours, new(defaultRandomiser)}
}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	var responseWriter http.ResponseWriter
	if chosenBehaviour := getBehaviour(s.behaviours, s.randomiser); chosenBehaviour != nil {
		responseWriter = monkeyWriter{w, chosenBehaviour}
	} else {
		responseWriter = w
	}

	s.delegate.ServeHTTP(responseWriter, r)
}

type monkeyWriter struct {
	http.ResponseWriter
	behaviour *behaviour
}

func (w monkeyWriter) Write(data []byte) (int, error) {

	time.Sleep(w.behaviour.Delay * time.Millisecond)

	if w.behaviour.Garbage > 0 {
		content := []byte{}
		for i := 0; i < w.behaviour.Garbage; i++ {
			content = append(content, byte('a'))
		}
		return w.ResponseWriter.Write(content)
	}

	if len(w.behaviour.Body) > 0 {
		return w.ResponseWriter.Write([]byte(w.behaviour.Body))
	}
	return w.ResponseWriter.Write(data)
}

func (w monkeyWriter) WriteHeader(code int) {
	if w.behaviour.Status != 0 {
		w.ResponseWriter.WriteHeader(w.behaviour.Status)
	} else {
		w.ResponseWriter.WriteHeader(code)
	}
}
