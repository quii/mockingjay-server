package main

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"time"
)

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

func NewMonkeyServer(server http.Handler, behaviours []behaviour) http.Handler {
	return &MonkeyServer{server, behaviours}
}

type MonkeyServer struct {
	delegate   http.Handler
	behaviours []behaviour
}

func (s *MonkeyServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	chosenBehaviour := getBehaviour(s.behaviours)
	if chosenBehaviour != nil {
		time.Sleep(chosenBehaviour.Delay * time.Millisecond)
		if chosenBehaviour.Status != 0 {
			w.WriteHeader(chosenBehaviour.Status)
		}
	}

	s.delegate.ServeHTTP(w, r)
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

type behaviour struct {
	Delay     time.Duration
	Frequency float64
	Status    int
}
