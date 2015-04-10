# mockingjay

[![GoDoc](https://godoc.org/github.com/quii/mockingjay?status.svg)](https://godoc.org/github.com/quii/mockingjay)

Create a server based on json configuration. Useful for integration tests, hopefully.

## Example

```go
package main

import (
	"github.com/quii/mockingjay"
	"log"
	"net/http"
)

func main() {
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
		"Name": "Amazing endpoint",
		"Request":{
	    	"URI" : "/world",
	    	"Method": "POST",
	    	"Headers":
	    		{
	    			"Content-Type": "application/json"
	    		}
		},
		"Response":{
			"Code": 201,
			"Body": "hello, world"
		}
	}
]`
	endpoints, err := mockingjay.NewFakeEndpoints(testJSON)

	if err != nil {
		log.Fatal(err)
	}

	server := mockingjay.NewServer(endpoints)

	// Mount it just like any other server
	http.Handle("/", server)
	http.ListenAndServe(":9090", nil)
}
```


## Todo

- Although it supports request/response headers, it only supports one value per header (http allows you to set multiple values)
- Tests for stuff inside request.go
