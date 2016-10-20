#!/bin/bash

set -e -x 

UDOCS_SRC="$GOPATH/src/github.com/ultimatesoftware/udocs"

BUILD_NUMBER=0.1.0
PLATFORM="$(uname -s | tr '[A-Z]' '[a-z]')"

rm -rf /tmp/udocs/bin
mkdir -p /tmp/udocs/bin

GODEP='github.com/tools/godep'
if [ ! -d "${GOPATH}/pkg/${PLATFORM}_amd64/${GODEP}" ]; then
  go get -v ${GODEP}
fi

BIN="/tmp/udocs/bin/udocs-${BUILD_NUMBER}_${PLATFORM}_amd64"
GOOS=${PLATFORM} GOARCH=amd64 ${GOPATH}/bin/godep go build -o "${BIN}" -ldflags "-X main.buildNumber=${BUILD_NUMBER}" "${UDOCS_SRC}/main.go"

# the local user directory where udocs content is installed
UDOCS_PATH="$HOME/.udocs"

if [ -d $UDOCS_PATH ]; then
  rm -rf $UDOCS_PATH
fi

mkdir -p $UDOCS_PATH/bin
mkdir -p $UDOCS_PATH/lib
mkdir -p $UDOCS_PATH/docs
mkdir -p $UDOCS_PATH/var

# copy the default config file and envvars
cp ${UDOCS_SRC}/cli/config/udocs.conf $UDOCS_PATH
cp ${UDOCS_SRC}/cli/config/udocs.env $UDOCS_PATH

# copy static files (css, templates, scripts, images, fonts, etc.)
cp -r ${UDOCS_SRC}/cli/udocs/static $UDOCS_PATH/lib

# copy the build and install scripts
cp  -r ${UDOCS_SRC}/bin/ $UDOCS_PATH/bin

# copy docs
cp -r ${UDOCS_SRC}/docs/ $UDOCS_PATH/docs

# copy Dockerfile
# cp -r ${UDOCS_SRC}/ci/Dockerfile $UDOCS_PATH/Dockerfile

# copy the udocs CLI binary, and symlink it into the user's PATH
cp /tmp/udocs/bin/*_${PLATFORM}_amd64 $UDOCS_PATH/bin/udocs
ln -sf $UDOCS_PATH/bin/udocs /usr/local/bin/udocs
chmod +x /usr/local/bin/udocs
