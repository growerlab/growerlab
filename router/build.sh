#!/bin/bash

if [ ! -f ./router/main ]; then
    cd ./router || exit 1
    go build -o main ./main.go
    cd - || exit 1
fi
