# mockingjay server

Mockingjay server creates HTTP web servers for integration tests from JSON configuration and then compares them against real servers to ensure they are valid to test against.

Mockingjay is fast, requires no coding and is better than other solutions because it will ensure your mock servers and real integration points are consistent.

## Installation

     go get github.com/quii/mockingjay-server

### Non go-centric install

    TODO - Need to make cross platform binaries and stick them somewhere.

## Running a fake server

    mockingjay-server -config=example.json -port=1234

## Comparing it against a real server

    mockingjay-server -config=example.json -realURL=http://localhost:1234

## Building

### Requirements

- Go 1.3+ installed ($GOPATH set, et al)
- godep https://github.com/tools/godep
- golint https://github.com/golang/lint

    ./build.sh	

### TODO

- Pretty diagrams explaining it all (interations between fakes, CDCs et al.)
- Take a base url and compare them with the fake
- Get the binaries built somewhere
- Code is a mess and hacky right now, needs tests.
- Performance tests needed too. Pretty sure it can be sped up by checking the endpoints concurrenty.