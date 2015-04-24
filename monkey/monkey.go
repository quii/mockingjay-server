package monkey

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

// Server wraps around a http.Handler and adds destructive behaviour (monkey business) based on the behaviours passed in
type Server struct {
	delegate   http.Handler
	behaviours []behaviour
	randomiser randomiser
}

// NewServer creates http.Handler which wraps it's monkey business around it, to return a new http.Handler. If no behaviours are defined in the config it will return the original handler, otherwise an error
func NewServer(server http.Handler, configPath string) (http.Handler, error) {
	if configPath == "" {
		return server, nil
	}

	config, err := ioutil.ReadFile(configPath)

	if err != nil {
		return nil, errors.New(fmt.Sprintf("Problem occured when trying to read the config file: %v", err))
	}

	behaviours, err := monkeyConfigFromYAML(config)

	if err != nil {
		return nil, errors.New(fmt.Sprintf("Problem occured when trying to parse the config file: %v", err))
	}

	log.Println("Monkey config loaded")
	for _, b := range behaviours {
		log.Println(b)
	}

	return newServerFromBehaviour(server, behaviours), nil
}

func newServerFromBehaviour(server http.Handler, behaviours []behaviour) http.Handler {
	return &Server{server, behaviours, new(defaultRandomiser)}
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	var responseWriter http.ResponseWriter
	if chosenBehaviour := getBehaviour(s.behaviours, s.randomiser); chosenBehaviour != nil {
		s.misbehave(*chosenBehaviour, w)
		responseWriter = monkeyWriter{w, []byte(chosenBehaviour.Body), chosenBehaviour.Garbage}
	} else {
		responseWriter = w
	}

	s.delegate.ServeHTTP(responseWriter, r)
}

func (s *Server) misbehave(behaviour behaviour, w http.ResponseWriter) {
	time.Sleep(behaviour.Delay * time.Millisecond)
	if behaviour.Status != 0 {
		w.WriteHeader(behaviour.Status)
	}
}

type monkeyWriter struct {
	http.ResponseWriter
	newBody      []byte
	garbageCount int
}

func (w monkeyWriter) Write(data []byte) (int, error) {

	if w.garbageCount > 0 {
		content := []byte{}
		for i := 0; i < w.garbageCount; i++ {
			content = append(content, byte('a'))
		}
		return w.ResponseWriter.Write(content)
	}

	if len(w.newBody) > 0 {
		return w.ResponseWriter.Write(w.newBody)
	}
	return w.ResponseWriter.Write(data)
}
