#!/bin/bash

go generate
golint ./... | grep -v bindata_assetfs | grep -v vendor
go list ./... | grep -v /vendor/ | grep -v binddata.go | xargs -L1 go vet
go list ./... | grep -v /vendor/ | grep -v binddata.go | xargs -L1 go vet
go fmt ./...
go test ./... --cover
go install
