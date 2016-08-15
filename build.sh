#!/bin/bash

echo "Building frontend assets"
cd ui
npm install
npm run build
cd ..

echo "Building application"
go generate
go list ./... | grep -v /vendor/ | grep -v binddata.go | xargs -L1 golint
go list ./... | grep -v /vendor/ | grep -v binddata.go | xargs -L1 go vet
go list ./... | grep -v /vendor/ | grep -v binddata.go | xargs -L1 go vet
go fmt ./...
go test ./...
go install
