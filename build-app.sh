#!/bin/bash

set -e

go generate
#golint ./... | grep -v bindata_assetfs | grep -v vendor <- exits with 1 because of vendor, how do i fix?
go list ./... | grep -v /vendor/ | grep -v bindata_assetfs.go | xargs -L1 go vet
go list ./... | grep -v /vendor/ | grep -v bindata_assetfs.go | xargs -L1 go vet
go fmt ./...
go test ./... --cover
go install
