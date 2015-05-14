package monkey

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"math/rand"
	"time"
)

type behaviour struct {
	Delay     time.Duration
	Frequency float64
	Status    int
	Body      string
	Garbage   int
}

type randomiser interface {
	getFloat() float64
}

type defaultRandomiser struct{}

func (d *defaultRandomiser) getFloat() float64 {
	return rand.Float64()
}

func getBehaviour(behaviours []behaviour, randomiser randomiser) *behaviour {
	randnum := randomiser.getFloat()
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

func monkeyConfigFromYAML(data []byte) (result []behaviour, err error) {
	err = yaml.Unmarshal(data, &result)
	return
}
