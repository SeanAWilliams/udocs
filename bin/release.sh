#!/bin/bash

set -e -x

git checkout master
version=$(head ./version)
git add version
git commit -m "Bumped version tag to $version"
git push 
git tag origin $version
git push $version