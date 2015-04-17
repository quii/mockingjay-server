package main

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"time"
)

// MonkeyServer wraps around a http.Handler and adds destructive behaviour (monkey business) based on the behaviours passed in
type MonkeyServer struct {
	delegate   http.Handler
	behaviours []behaviour
}

// NewMonkeyServer creates http.Handler which wraps it's monkey business around it, to return a new http.Handler
func NewMonkeyServer(server http.Handler, behaviours []behaviour) http.Handler {
	return &MonkeyServer{server, behaviours}
}

func (s *MonkeyServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	chosenBehaviour := getBehaviour(s.behaviours)

	var resp http.ResponseWriter
	if chosenBehaviour != nil {
		s.misbehave(*chosenBehaviour, w)
		resp = monkeyWriter{w, []byte(chosenBehaviour.Body)}
	} else {
		resp = w
	}

	s.delegate.ServeHTTP(resp, r)
}

func (s *MonkeyServer) misbehave(behaviour behaviour, w http.ResponseWriter) {
	time.Sleep(behaviour.Delay * time.Millisecond)
	if behaviour.Status != 0 {
		w.WriteHeader(behaviour.Status)
	}
}

func getBehaviour(behaviours []behaviour) *behaviour {
	randnum := rand.Float64()
	lower := 0.0
	var upper float64
	for _, behaviour := range behaviours {
		upper = lower + behaviour.Frequency
		if randnum > lower && randnum <= upper {
			return &behaviour
		}

		lower = upper
	}

	return nil

}

func loadMonkeyConfig(path string) []behaviour {

	if path == "" {
		return []behaviour{}
	}

	config, err := ioutil.ReadFile(path)

	if err != nil {
		log.Fatalf("Problem occured when trying to read the config file: %v", err)
	}

	var result []behaviour
	err = yaml.Unmarshal([]byte(config), &result)

	if err != nil {
		log.Fatalf("Problem occured when trying to parse the config file: %v", err)
	}

	log.Println(result)

	return result
}

type behaviour struct {
	Delay      time.Duration
	Frequency  float64
	Status     int
	WriteDelay time.Duration
	Body       string
}

type monkeyWriter struct {
	http.ResponseWriter
	newBody []byte
}

func (w monkeyWriter) Write(data []byte) (int, error) {
	if len(w.newBody) > 0 {
		return w.ResponseWriter.Write(w.newBody)
	} else {
		return w.ResponseWriter.Write(data)
	}
}
