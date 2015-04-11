# mockingjay server

Mockingjay server creates HTTP web servers for integration tests from YAML configuration and then compares them against real servers to ensure they are valid to test against.

Mockingjay is fast, requires no coding and is better than other solutions because it will ensure your mock servers and real integration points are consistent.

## Installation

     go get github.com/quii/mockingjay-server

### Non go-centric install

    TODO - Need to make cross platform binaries and stick them somewhere.

## Running a fake server

    mockingjay-server -config=example.yaml -port=1234

## Check configuration is compatible with a real server

    mockingjay-server -config=example.yaml -realURL=http://some-real-api.com

## Building

### Requirements

- Go 1.3+ installed ($GOPATH set, et al)
- godep https://github.com/tools/godep
- golint https://github.com/golang/lint


    ./build.sh

### TODO

- Pretty diagrams explaining it all (interations between fakes, CDCs et al.)
- Currently the CDC part only checks status codes. Now need to check structure of content.
	- JSON
	- XML
- Get the binaries built somewhere
- Code is a mess and hacky right now, needs tests.
- Performance tests needed too. Pretty sure it can be sped up by checking the endpoints concurrenty.

### Things to figure out

- Still somewhat reliaint on golden data. i.e fake /user/2 - does user 2 exist on the real server? What can be done?
- Is it possible to do chains of requests for more complicated tests but still keep it nice and simple. Should we?
