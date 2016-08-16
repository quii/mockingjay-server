package mockingjay

import (
	"io/ioutil"
	"net/http"
)

// ExampleNewServer is an example as to how to make a fake server. The mockingjay server implements what is needed to mount it as a standard web server.
func ExampleNewServer() {

	// Create a fake server from YAML
	testYAML := `
---
 - name: Test endpoint
   request:
     uri: /hello
     method: GET
     headers:
       content-type: application/json
     body: foobar
   response:
     code: 200
     body: hello, world
     headers:
       content-type: text/plain

 - name: Test endpoint 2
   request:
     uri: /world
     method: DELETE
   response:
     code: 200
     body: hello, world

 - name: Failing endpoint
   request:
     uri: /card
     method: POST
     body: Greetings
   response:
     code: 500
     body: Oh bugger

 - name: Endpoint not used for CDC
   cdcdisabled: true
   request:
     uri: /burp
     method: POST
     body: Belch
   response:
     code: 500
     body: Oh no
 `

	endpoints, _ := NewFakeEndpoints([]byte(testYAML))
	server := NewServer(endpoints, false, ioutil.Discard)

	// Mount it just like any other server
	http.Handle("/", server)
	http.ListenAndServe(":9090", nil)

}
