#!/bin/bash
set -e 

DIR=$GOPATH/src/github.com/ultimatesoftware
mkdir -p $DIR

cd $DIR

go get github.com/ultimatesoftware/udocs 

cd udocs

# git pull origin master

./bin/install.sh
