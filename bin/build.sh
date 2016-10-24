#!/bin/bash

set -e -x

if [ ! -d "${GOPATH}/bin/go-bindata" ]; then
  go get -u -v github.com/jteeuwen/go-bindata/...
fi

cd static
go-bindata -pkg static ./*/* 
cd -

go fmt ./cli/...
go vet -v ./cli/...
go test -v ./cli/...
go install -v 

# build for Linux Docker image
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -v -a -installsuffix cgo -o bin/udocs-docker -ldflags "-X main.buildNumber=${1}"

# build for Mac OS X 
GOOS=darwin GOARCH=amd64 go build -o bin/udocs-osx -ldflags "-X main.buildNumber=${1}"

# build for Linux
GOOS=linux GOARCH=amd64 go build -o bin/udocs-linux -ldflags "-X main.buildNumber=${1}"
