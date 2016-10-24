#!/bin/bash

set -e -x

git checkout master
git tag v0.1.0 -a -m "[ci skip] Bumped version tag"
git push --tags