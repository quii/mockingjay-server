#!/bin/bash

docker build -t mjtest .

docker run --rm=true -p 8080:8080 -v ${PWD}/examples/example.yaml:/issue42.yaml mjtest -config=/issue42.yaml -realURL=https://facebook.com -port=8080
