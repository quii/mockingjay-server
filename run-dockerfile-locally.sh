#!/bin/bash

docker build -t mjtest .

docker run --rm=true -p 8080:8080 -v ${PWD}/examples/example.yaml:/example.yaml mjtest -config=/example.yaml -realURL=https://facebook.com -port=8080
