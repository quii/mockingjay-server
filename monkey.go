package main

import (
	"net/http"
	"time"
)

// MonkeyServer wraps around a http.Handler and adds destructive behaviour (monkey business) based on the behaviours passed in
type MonkeyServer struct {
	delegate   http.Handler
	behaviours []behaviour
	randomiser randomiser
}

// NewMonkeyServer creates http.Handler which wraps it's monkey business around it, to return a new http.Handler
func NewMonkeyServer(server http.Handler, behaviours []behaviour) http.Handler {
	return &MonkeyServer{server, behaviours, new(defaultRandomiser)}
}

func (s *MonkeyServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	var responseWriter http.ResponseWriter
	if chosenBehaviour := getBehaviour(s.behaviours, s.randomiser); chosenBehaviour != nil {
		s.misbehave(*chosenBehaviour, w)
		responseWriter = monkeyWriter{w, []byte(chosenBehaviour.Body), chosenBehaviour.Garbage}
	} else {
		responseWriter = w
	}

	s.delegate.ServeHTTP(responseWriter, r)
}

func (s *MonkeyServer) misbehave(behaviour behaviour, w http.ResponseWriter) {
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
