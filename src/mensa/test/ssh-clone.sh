#!/bin/bash

#
# for test clone
#

ROOT=$GOPATH/src/github.com/growerlab/mensa/test
TEST_PATH=$ROOT/test
USER=moli

cd $ROOT

rm -rf $TEST_PATH

git clone ssh://$USER@localhost:8022/moli/test.git