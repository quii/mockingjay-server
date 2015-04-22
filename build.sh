#!/bin/bash

golint .
godep go test ./...
godep go install
