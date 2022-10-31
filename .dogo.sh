#!/bin/bash
set +x

if [ ! -d ".repositories" ]; then
    mkdir .repositories
fi

cd src || exit 1

function clean() {
    if [ -f /tmp/growerlab.pid ]; then
        kill "$(cat /tmp/growerlab.pid)"
    fi
}

function run() {
    go run backend/main.go &

    if [ ! -f ./conf/hooks/update ]; then
        echo "building go-git-grpc 'update' hook..."
        cd go-git-grpc && go build -o ../conf/hooks/update ./script/hulk/main.go
        cd - || exit 1
    fi
    go run go-git-grpc/cli/main.go &

    go run mensa/main.go &

    echo $$ >/tmp/growerlab.pid
}

clean
run
wait
echo "done"
