#!/bin/bash

set -e -x

if [ ! -d "${GOPATH}/bin/go-bindata" ]; then
  go get -v github.com/jteeuwen/go-bindata/...
fi

cd static
go-bindata -pkg static ./*/* 
cd -

go fmt ./cli/...
go vet -v ./cli/...
go test -v ./cli/...
go install -v 

# build Linux image for Dockerfile
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -v -a -installsuffix cgo -o bin/udocs -ldflags "-X main.buildNumber=${1}"
