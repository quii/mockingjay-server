#!/bin/bash

go test -v -run=Example -coverprofile=c.out && go tool cover -html=c.out