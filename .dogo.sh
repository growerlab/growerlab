#!/bin/bash
set +x

cd src

go run backend/main.go &

go run go-git-grpc/cli/main.go &

go run mensa/main.go &

wait
echo "done"
