#!/bin/bash

golint .
go test  -coverprofile=coverage.out -covermode=count
go tool cover -html=coverage.out