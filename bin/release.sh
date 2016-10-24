#!/bin/bash

set -e -x

VERSION_TAG=v1.0.1

git checkout master
git tag ${VERSION_TAG} -a -m "Bumped version tag to ${VERSION_TAG}"
git push --tags