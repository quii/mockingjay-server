# mockingjay server

[![Build Status](https://travis-ci.org/quii/mockingjay-server.svg?branch=master)](https://travis-ci.org/quii/mockingjay-server)[![Coverage Status](https://coveralls.io/repos/quii/mockingjay-server/badge.svg?branch=master)](https://coveralls.io/r/quii/mockingjay-server?branch=master)[![GoDoc](https://godoc.org/github.com/quii/mockingjay-server?status.svg)](https://godoc.org/github.com/quii/mockingjay-server)

![mj example](http://i.imgur.com/ZtI1Q39.gif)

Mockingjay lets you define the contract between a consumer and producer and with just a configuration file you get:

- A fast to launch fake server for your integration tests
 - Configurable to simulate the eratic nature of calling other services
- [Consumer driven contracts (CDCs)](http://martinfowler.com/articles/consumerDrivenContracts.html) to run against your real downstream services.

**Mockingjay makes it really easy to check your HTTP integration points**. It's fast, requires no coding and is better than other solutions because it will ensure your mock servers and real integration points are consistent so that you never have a green build but failing software.

- [Installation](https://github.com/quii/mockingjay-server/wiki/Installing) - [Download a binary](https://github.com/quii/mockingjay-server/releases/latest), [use a Docker image](https://hub.docker.com/r/quii/mockingjay-server/) or `go get`
- [Rationale](https://github.com/quii/mockingjay-server/wiki/Rationale)
- [See how mockingjay can easily fit into your workflow to make integration testing really easy and robust](https://github.com/quii/mockingjay-server/wiki/Suggested-workflow)


## Running a fake server

````yaml
---
 - name: My very important integration point
   request:
     uri: /hello
     method: POST
     body: "Chris" # * matches any body
   response:
     code: 200
     body: '{"message": "hello, Chris"}'   # * matches any body
     headers:
       content-type: application/json

# define as many as you need...
````

````bash
$ mockingjay-server -config=example.yaml -port=1234 &
2015/04/13 14:27:54 Serving 3 endpoints defined from example.yaml on port 1234
$ curl http://localhost:1234/hello
{"message": "hello, world"}
````

## Check configuration is compatible with a real server

````bash
$ mockingjay-server -config=example.yaml -realURL=http://some-real-api.com
2015/04/13 21:06:06 Test endpoint (GET /hello) is incompatible with http://some-real-api - Couldn't reach real server
2015/04/13 21:06:06 Test endpoint 2 (DELETE /world) is incompatible with http://some-real-api - Couldn't reach real server
2015/04/13 21:06:06 Failing endpoint (POST /card) is incompatible with http://some-real-api - Couldn't reach real server
2015/04/13 21:06:06 At least one endpoint was incompatible with the real URL supplied
````
This ensures your integration test is working against a *reliable* fake.

### Inspect what requests mockingjay has received

     http://{mockingjayhost}:{port}/requests

Calling this will return you a JSON list of requests

## Make your fake server flaky

Mockingjay has an annoying friend, a monkey. Given a monkey configuration you can make your fake service misbehave. This can be useful for performance tests where you want to simulate a more realistic scenario (i.e all integration points are painful).

````yaml
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

# Write 10,000,000 garbage bytes 9% of the time
- garbage: 10000000
  frequency: 0.09
````

````bash
$ mockingjay-server -config=examples/example.yaml -monkeyConfig=examples/monkey-business.yaml
2015/04/17 14:19:53 Serving 3 endpoints defined from examples/example.yaml on port 9090
2015/04/17 14:19:53 Monkey config loaded
2015/04/17 14:19:53 50% of the time | Body: This is wrong :(
2015/04/17 14:19:53 20% of the time | Delay: 1s
2015/04/17 14:19:53 30% of the time | Status: 404
2015/04/17 14:19:53  9% of the time | Garbage bytes: 10000000
````

## Building

### Requirements

- Go 1.3+ installed ($GOPATH set, et al)
- Node 4.2.x (for the frontend)
- golint https://github.com/golang/lint
- `$ go get github.com/jteeuwen/go-bindata/...`
- `$ go get github.com/elazarl/go-bindata-assetfs/...`

### Build application

````bash
$ go get github.com/quii/mockingjay-server
$ cd $GOPATH/src/github.com/quii/mockingjay-server
$ ./build.sh
````

### Working with the frontend
````bash
$ cd ui/
$ npm install
$ npm run dev
````
This will rebuild the frontend assets on file changes. Then in another terminal:
````bash
$ ENV=LOCAL mockingjay-server -config=$PATH-TO-CONFIG
````

The env local ensures the web server serves your constantly building assets, rather than the bundled ones made on `build.sh`.

MIT license
