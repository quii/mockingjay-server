# mockingjay server

Mockingjay creates two artifacts from a given configuration file.

- A program to check that a service responds with a compatible response given each request
- A HTTP server which serves the configured request for each response 

````yaml
---
 - name: Test endpoint
   request:
     uri: /hello
     method: GET
   response:
     code: 200
     body: hello, world
     headers:
       content-type: text/plain

# define as many as you need...
````

Mockingjay is fast, requires no coding and is better than other solutions because it will ensure your mock servers and real integration points are consistent.

## Installation

     go get github.com/quii/mockingjay-server

## Running a fake server

    mockingjay-server -config=example.yaml -port=1234

## Check configuration is compatible with a real server

    mockingjay-server -config=example.yaml -realURL=http://some-real-api.

There are example files inside the examples directory.

## Building

### Requirements

- Go 1.3+ installed ($GOPATH set, et al)
- godep https://github.com/tools/godep
- golint https://github.com/golang/lint


    ./build.sh

### TODO

- Pretty diagrams explaining it all (interations between fakes, CDCs et al.)
- Currently the CDC part only checks status codes. 
- Check XML structure when applicable
- Get the binaries built somewhere for non gophers
- Investigate a more standard test output

### Things to figure out

- Still somewhat reliaint on golden data. i.e fake /user/2 - does user 2 exist on the real server? What can be done?
- Is it possible to do chains of requests for more complicated tests but still keep it nice and simple. Should we?
