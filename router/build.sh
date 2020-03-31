#!/bin/bash

GOOS=linux go build -o ./router ./router.go

docker build "$GOPATH"/src/github.com/growerlab/growerlab/router/.
