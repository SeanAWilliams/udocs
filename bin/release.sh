#!/bin/bash

set -e -x

git checkout master
version=$(head ./version)
git add version
git commit -m "Bumped version tag to $version"
git push 
git tag $version
git push origin $version