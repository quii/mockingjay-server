#!/bin/bash

go list ./... | grep -v /vendor/ | xargs -L1 golint
go list ./... | grep -v /vendor/ | xargs -L1 go vet
go list ./... | grep -v /vendor/ | xargs -L1 go test -cover
go install
