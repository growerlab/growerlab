#!/bin/bash

#
# for init bare repo
#

ROOT=$GOPATH/src/github.com/growerlab/mensa/test

REPO_PATH=$ROOT/repos/mo/te/moli/test
TEMP_PATH=$ROOT/temp

mkdir -p $REPO_PATH

cd $REPO_PATH

git init --bare

git clone $REPO_PATH $TEMP_PATH

cd $TEMP_PATH

touch a.txt
echo "a" >> a.txt

touch b.txt
echo "b" >> b.txt

git add a.txt b.txt
git commit -m 'init'
git push origin master


rm -rf $TEMP_PATH
echo "done"