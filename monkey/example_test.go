package monkey

import (
	"log"
	"net/http"
)

// ExampleWrapMonkeyBusinessAroundAServer shows how use monkey to wrap around a http.Handler and change it's behaviour from configuration
func ExampleWrapMonkeyBusinessAroundAServer() {

	// Create a fake server from YAML
	testYAML := `
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

	// your server you want to monkey with
	var server http.Handler

	// monkey will create you a new server from your YAML
	server, err := NewServerFromYAML(server, []byte(testYAML))

	// err will be returned if the config is bad
	if err != nil {
		log.Fatal(err)
	}

	// Mount it just like any other server
	http.Handle("/", server)
	http.ListenAndServe(":9090", nil)

}
