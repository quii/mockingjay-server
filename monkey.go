package main

import (
	"fmt"
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

	var responseWriter http.ResponseWriter
	if chosenBehaviour := getBehaviour(s.behaviours); chosenBehaviour != nil {
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

	behaviours := monkeyConfigFromYAML(config)

	log.Println("Monkey config loaded")
	for _, b := range behaviours {
		log.Println(b)
	}

	return behaviours
}

func monkeyConfigFromYAML(data []byte) []behaviour {
	var result []behaviour
	err := yaml.Unmarshal([]byte(data), &result)

	if err != nil {
		log.Fatalf("Problem occured when trying to parse the config file: %v", err)
	}

	return result
}

type behaviour struct {
	Delay     time.Duration
	Frequency float64
	Status    int
	Body      string
	Garbage   int
}

func (b behaviour) String() string {

	frequency := fmt.Sprintf("%2.0f%% of the time |", b.Frequency*100)

	delay := ""
	if b.Delay != 0 {
		delay = fmt.Sprintf("Delay: %v ", b.Delay*time.Millisecond)
	}

	status := ""
	if b.Status != 0 {
		status = fmt.Sprintf("Status: %v ", b.Status)
	}

	body := ""
	if b.Body != "" {
		body = fmt.Sprintf("Body: %v ", b.Body)
	}

	garbage := ""
	if b.Garbage != 0 {
		garbage = fmt.Sprintf("Garbage bytes: %d ", b.Garbage)
	}

	return fmt.Sprintf("%v %v%v%v%v", frequency, delay, status, body, garbage)
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
