#!/bin/bash

MENSA="$GOPATH"/src/github.com/growerlab/mensa
BACKEND="$GOPATH"/src/github.com/growerlab/backend
SVC="$GOPATH"/src/github.com/growerlab/svc


BASE_DIR="$GOPATH"/src/github.com/growerlab/growerlab

cd "$BACKEND"
mkdir -p "$BASE_DIR"/data/services/backend
go build -o "$BASE_DIR"/data/services/backend/backend "$BACKEND"/main.go
cp -R "$BACKEND"/conf "$BASE_DIR"/data/services/backend
cd -

cd "$MENSA"
mkdir -p "$BASE_DIR"/data/services/mensa
go build -o "$BASE_DIR"/data/services/mensa/mensa "$MENSA"/main.go
cp -R "$MENSA"/conf "$BASE_DIR"/data/services/mensa
cd -

cd "$SVC"
mkdir -p "$BASE_DIR"/data/services/svc
go build -o "$BASE_DIR"/data/services/svc/svc "$SVC"/main.go
cp -R "$SVC/.env.example" "$BASE_DIR"/data/services/svc/.env
cd -
