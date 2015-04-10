package mockingjay

import (
	"net/http"
)

// ExampleNewServer is an example as to how to make a fake server. The mockingjay server implements what is needed to mount it as a standard web server.
func ExampleNewServer() {

	// Create a fake server from json
	testJSON := `
[
	{
		"Name": "A descriptive name useful for when stuff is logged",
		"Request":{
	    	"URI" : "/hello",
	    	"Method": "GET"
		},
		"Response":{
			"Code": 200,
			"Body": "hello, world"
		}
	},
	{
		"Name": "Blah blah",
		"Request":{
	    	"URI" : "/world",
	    	"Method": "DELETE"
		},
		"Response":{
			"Code": 200,
			"Body": "hello, world"
		}
	}
]`

	endpoints, _ := NewFakeEndpoints(testJSON)
	server := NewServer(endpoints)

	// Mount it just like any other server
	http.Handle("/", server)
	http.ListenAndServe(":9090", nil)

}
