package main

import (
	"github.com/quii/mockingjay"
	"testing"
)

const sleepyTime = 500

func BenchmarkCompatabilityChecking(b *testing.B) {
	body := "hello, world"
	realServer := makeRealServer(body, sleepyTime)

	fakeEndPoints, err := mockingjay.NewFakeEndpoints([]byte(multipleEndpointYAML))

	if err != nil {
		b.Fatalf("Couldn't make mockingjay endpoints, is your data correct? [%v]", err)
	}

	checker := NewCompatabilityChecker(fakeEndPoints)

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
