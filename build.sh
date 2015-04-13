#!/bin/bash

golint .
godep go test  -coverprofile=coverage.out -covermode=count
go tool cover -html=coverage.out
godep go install
