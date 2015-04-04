package main

import (
	"github.com/quii/mockingjay"
	"log"
	"net/http"
	"os"
)

func main() {
	testJSON := `
[
    {
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

	port := os.Getenv("PORT")

	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Serving %d endpoints on port %s", len(endpoints), port)

	server := mockingjay.NewServer(endpoints)

	// Mount it just like any other server
	http.Handle("/", server)
	http.ListenAndServe(":"+port, nil)
}
