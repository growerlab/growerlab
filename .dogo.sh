#!/bin/bash
set +x

cd src || exit 1

go run backend/main.go &

if [ ! -f ./go-git-grpc/hooks/update ]; then
    echo "building go-git-grpc 'update' hook..."
    cd go-git-grpc && go build -o ./hooks/update ./script/hulk/main.go
    cd - || exit 1
fi
go run go-git-grpc/cli/main.go &

go run mensa/main.go &

wait
echo "done"
