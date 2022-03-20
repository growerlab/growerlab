#!/bin/bash

#
# for test clone
#

ROOT=$GOPATH/src/github.com/growerlab/mensa/test
TEST_PATH=$ROOT/test

cd $ROOT || exit

rm -rf $TEST_PATH

git clone http://localhost:8080/moli/test.git
