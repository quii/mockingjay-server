#!/bin/bash

golint .
godep go test  -coverprofile=coverage.out -covermode=count
godep go install
