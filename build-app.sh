#!/bin/bash

go generate
go list ./... | grep -v /vendor/ | grep -v binddata.go | xargs -L1 golint
go list ./... | grep -v /vendor/ | grep -v binddata.go | xargs -L1 go vet
go list ./... | grep -v /vendor/ | grep -v binddata.go | xargs -L1 go vet
go fmt ./...
go test -v ./... --cover
go install
