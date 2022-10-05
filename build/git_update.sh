#!/bin/bash

VERSION=""

MAJOR=""
MINOR=""

# Get version from ../cmd/VERSION.yaml
# if not found just use "latest" as version
if [ -f ../cmd/VERSION.yaml ]; then
    VERSION=$(cat ../cmd/VERSION.yaml)

    # Major is a string like "1"
    MAJOR=$(echo $VERSION | grep -oP '(?<=major: ).*')
    # Minor is a string like "0"
    MINOR=$(echo $VERSION | grep -oP '(?<=minor: ).*')

    VERSION="$MAJOR.$MINOR"
else
    VERSION="latest"
fi  

echo ::set-output name=git-tag::$VERSION    
