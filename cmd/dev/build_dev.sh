#!/bin/bash

getCloneURL() {
  echo "https://github.com/growerlab/$1.git"
}

FRONTEND=$(getCloneURL "frontend")
SVC=$(getCloneURL "svc")
BACKEND=$(getCloneURL "backend")
MENSA=$(getCloneURL "mensa")

echo $FRONTEND