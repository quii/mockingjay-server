package mockingjay

import (
	"fmt"
	"log"
	"os"
	"testing"
)

const sleepyTime = 500
const numberOfEndpoints = 3

// 503439778 ns with 3 endpoints
// 522222809 ns with 100 endpoints (surely we will never see a config that big)
func BenchmarkCompatabilityChecking(b *testing.B) {
	body := "hello, world"
	realServer := makeFakeDownstreamServer(body, sleepyTime)
	checker := NewCompatabilityChecker(log.New(os.Stdout, "mocking-jay: ", log.Ldate|log.Ltime), DefaultHTTPTimeoutSeconds)
	endpoints, err := NewFakeEndpoints([]byte(multipleEndpointYAML(numberOfEndpoints)))

	if err != nil {
		b.Fatalf("Unable to create checker from YAML %v", err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		checker.CheckCompatibility(endpoints, realServer.URL)
	}
}

func multipleEndpointYAML(count int) string {

	benchmarkFormat := `
 - name: Test endpoint %d
   request:
     uri: /hello%d
     method: GET
   response:
     code: 200
     body: 'hello, world'

`
	body := `---
  `
	for i := 0; i < count; i++ {
		body = body + fmt.Sprintf(benchmarkFormat, i, i)
	}

	return body
}
