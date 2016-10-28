#!/bin/bash

set -e -x

# clear out older binaries
rm -rf ./bin/udocs*

if [[ $1 =~ ^[0-9]+\.[0-9]+\.[0-9]+$ ]]; then
  echo "v${1}" > ./version
fi
version=$(head ./version)

go fmt ./cli/...
go vet -v ./cli/...
go test -v ./cli/...
go install # install locally to keep up-to-date

# build for Linux Docker image
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -v -a -installsuffix cgo -o bin/udocs-docker-$version -ldflags "-X main.version=$version"

# build for Mac OS X 
GOOS=darwin GOARCH=amd64 go build -o bin/udocs-osx-$version -ldflags "-X main.version=$version"

# build for Linux
GOOS=linux GOARCH=amd64 go build -o bin/udocs-linux-$version -ldflags "-X main.version=$version"
