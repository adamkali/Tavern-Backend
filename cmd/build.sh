#!/bin/bash

# Get the commit message
COMMIT_MESSAGE=$1

# check if there is a VERSION.yaml file
# if not make one
if [ ! -f ./VERSION.yaml ]; then
    echo "major: 0" > ./VERSION.yaml
    echo "minor: 1" >> ./VERSION.yaml
fi

# check if -M is set
if [ "$2" = "-M" ]; then
    # increment major version
    echo "Incrementing major version"
    # Get the current major version
    MAJOR=$(cat ./VERSION.yaml | grep -oP '(?<=major: ).*')
    # Increment the major version
    MAJOR=$((MAJOR+1))
    # Set the minor version to 0
    MINOR=0
else
    # increment minor version
    echo "Incrementing minor version"
    # Get the current major version
    MAJOR=$(cat ./VERSION.yaml | grep -oP '(?<=major: ).*')
    # Get the current minor version
    MINOR=$(cat ./VERSION.yaml | grep -oP '(?<=minor: ).*')
    # Increment the minor version
    MINOR=$((MINOR+1))
fi

# git add -A
# git commit -m "$COMMIT_MESSAGE"
# then checkout the Beor
git checkout Beor
git merge main

git add -A
git commit -m "$COMMIT_MESSAGE"

# git push origin beor
git push origin Beor

# git checkout main
git checkout main