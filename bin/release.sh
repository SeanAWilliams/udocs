#!/bin/bash

set -e -x

git checkout master
local_version=$(head ./version)
git pull
remote_version=$(head ./version)

if [ "$local_version" != "$remote_version" ]; then
    echo "Local version $local_version does not match remote version $remote_version"
    exit 1
fi

git tag -a -m "Bumped version tag to $remote_version" $remote_version
git push origin $remote_version