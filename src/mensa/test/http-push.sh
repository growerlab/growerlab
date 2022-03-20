#!/bin/bash

#
# for test push
#

ROOT=$GOPATH/src/github.com/growerlab/mensa/test
TEMP_PATH=$ROOT/temp

git clone http://localhost:8080/moli/test.git $TEMP_PATH

cd $TEMP_PATH || exit

touch push.txt
echo "for push" >>push.txt

touch push2.txt
echo "for push" >>push2.txt

git add push.txt push2.txt

git commit -m 'for push commit'

git push origin master

rm -rf $TEMP_PATH
echo "done"
