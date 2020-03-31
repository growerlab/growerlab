#!/bin/bash

if [ ! -f ./router/router ]
then
  cd ./router || exit 1
  GOOS=linux go build -o ./router ./router.go
  cd - || exit 1
fi

