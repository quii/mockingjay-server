#!/bin/bash

golint .
godep go vet ./...
godep go test ./...  -cover
# go tool cover -html=coverage.out
godep go install
