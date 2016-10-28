#!/bin/bash

set -e -x

git checkout master
version=$(head ./version)
git add version
git tag $version -a -m "Bumped version tag to $version"
git push origin $version