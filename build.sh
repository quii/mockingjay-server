#!/bin/bash

golint .
godep go test ./...  -cover
# go tool cover -html=coverage.out
godep go install
