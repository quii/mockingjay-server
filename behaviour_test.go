package main

import (
	"testing"
)

func TestItGetsRandomBehaviours(t *testing.T) {
	behaviour1 := new(behaviour)
	behaviour1.Frequency = 0.15
	behaviour1.Status = 123

	behaviour2 := new(behaviour)
	behaviour2.Frequency = 0.2
	behaviour2.Status = 456

	allBehaviours := []behaviour{*behaviour1, *behaviour2}

	result := getBehaviour(allBehaviours, fakeRandomiser{0.9})

	if result != nil {
		t.Error("There shouldnt have been a behaviour returned, but there was", result)
	}

	result = getBehaviour(allBehaviours, fakeRandomiser{0.13})

	if result.Status != behaviour1.Status {
		t.Error("It shouldve found behaviour 1", result, behaviour1)
	}

	result = getBehaviour(allBehaviours, fakeRandomiser{0.19})

	if result.Status != behaviour2.Status {
		t.Error("It should've found behaviour 2", result)
	}
}

func TestItLoadsFromFile(t *testing.T) {
	config := loadMonkeyConfig("examples/monkey-business.yaml")
	if len(config) != 4 {
		t.Error("It didnt load all the behaviours from YAML")
	}
}

func TestItReturnsAnErrorForBadYAML(t *testing.T) {
	yaml := "lol not yaml"

	_, err := monkeyConfigFromYAML([]byte(yaml))

	if err == nil {
		t.Error("Error was not returned for bad YAML")
	}
}

func TestItParsesYAMLIntoBehaviour(t *testing.T) {
	yaml := `
---
# Writes a different body 50% of the time
- body: "This is wrong :( "
  frequency: 0.5

# Delays initial writing of response by a second 20% of the time
- delay: 1000
  frequency: 0.2

# Returns a 404 30% of the time
- status: 404
  frequency: 0.3

# Write 10,000,000 garbage bytes 10% of the time
- garbage: 10000000
  frequency: 0.09
  `
	behaviours, _ := monkeyConfigFromYAML([]byte(yaml))

	if len(behaviours) != 4 {
		t.Error("It didnt load all the behaviours from YAML")
	}
}

type fakeRandomiser struct {
	value float64
}

func (f fakeRandomiser) getFloat() float64 {
	return f.value
}
