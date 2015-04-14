package main

import (
	"testing"
)

const sleepyTime = 500

func BenchmarkCompatabilityChecking(b *testing.B) {
	body := "hello, world"
	realServer := makeFakeDownstreamServer(body, sleepyTime)
	checker, err := makeChecker(multipleEndpointYAML)

	if err != nil {
		b.Fatalf("Unable to create checker from YAML %v", err)
	}

	for i := 0; i < b.N; i++ {
		checker.CheckCompatability(realServer.URL)
	}
}

const multipleEndpointYAML = `
 - name: Test endpoint 1
   request:
     uri: /hello1
     method: GET
   response:
     code: 200
     body: 'hello, world'

 - name: Test endpoint 2
   request:
     uri: /hello2
     method: GET
   response:
     code: 200
     body: 'hello, world'

 - name: Test endpoint 1
   request:
     uri: /hello3
     method: GET
   response:
     code: 200
     body: 'hello, world'
`
