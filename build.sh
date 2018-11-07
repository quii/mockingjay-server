#!/bin/bash

set -e

echo "Building application"
go list ./... | grep -v /vendor/ | grep -v bindata_assetfs.go | xargs -L1 go vet
go list ./... | grep -v /vendor/ | grep -v bindata_assetfs.go | xargs -L1 go vet
go fmt ./...
go test ./... --cover
go install

