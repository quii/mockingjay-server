#!/bin/bash

go list ./... | grep -v /vendor/ | xargs -L1 golint
go list ./... | grep -v /vendor/ | xargs -L1 go vet
go test ./...
go install
