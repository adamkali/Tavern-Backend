#!/bin/bash

# Get the commit message
COMMIT_MESSAGE=$1

# Setup a progress bar
PUR='\033[0;35m'
BLU='\033[0;34m'
GREEN='\033[0;32m'
NC='\033[0m' # No Color

STAGE0 = "Pulling main brach from git"
STAGE1 = "Updating VERSION"
STAGE2 = "Logging into AWS ECR"
STAGE3 = "Building docker image"
STAGE4 = "Tagging docker image"
STAGE5 = "Pushing docker image to ECR"
STAGE6 = "COMPLETE!"

PROG0 = "[$BLU=>$PUR--------------------------------------$NCR] 0%"
PROG1 = "[$BLU######$PUR=>--------------------------------$NCR] 20%"
PROG2 = "[$BLU##############$PUR=>------------------------$NCR] 40%"
PROG3 = "[$BLU######################$PUR=>----------------$NCR] 60%"
PROG4 = "[$BLU##############################$PUR=>--------$NCR] 80%"
PROG5 = "[$BLU######################################$PUR=>$NCR] 100%"

echo -e "$PUR$STAGE0$NC$PROG0"
# git add -A
# git commit -m "$COMMIT_MESSAGE"
# then checkout the Beor and throw away output to avoid printing it
git checkout Beor 2>&1 > /dev/null
git merge main 2>&1 > /dev/null

git add -A 2>&1 > /dev/null
git commit -m "$COMMIT_MESSAGE" 2>&1 > /dev/null

# git push origin beor
git push origin Beor 2>&1 > /dev/null

echo -e "$PUR$STAGE1$NC$PROG1"

# check if there is a VERSION.yaml file
# if not make one throw away any output to the terminal
if [ ! -f ./VERSION.yaml ]; then
    echo "major: 0" >  ./VERSION.yaml 2>&1 > /dev/null
    echo "minor: 1" >> ./VERSION.yaml 2>&1 > /dev/null
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

# update the VERSION.yaml file and throw away any output to the terminal
echo "major: $MAJOR" >  ./VERSION.yaml 2>&1 > /dev/null
echo "minor: $MINOR" >> ./VERSION.yaml 2>&1 > /dev/null

# build the docker image and throw away any output to the terminal
echo -e "$PUR$STAGE2$NC$PROG2"
aws ecr get-login-password --region us-east-1 | docker login --username AWS --password-stdin 739810740537.dkr.ecr.us-east-1.amazonaws.com 2>&1 > /dev/null

echo -e "$PUR$STAGE3$NC$PROG3"
docker build -t tavern-profile-beor . 2>&1 > /dev/null

echo -e "$PUR$STAGE4$NC$PROG4"
docker tag tavern-profile-beor:latest 739810740537.dkr.ecr.us-east-1.amazonaws.com/tavern-profile-beor:$MAJOR.$MINOR 2>&1 > /dev/null

echo -e "$PUR$STAGE5$NC$PROG5"
docker push 739810740537.dkr.ecr.us-east-1.amazonaws.com/tavern-profile-beor:$MAJOR.$MINOR 2>&1 > /dev/null


# git checkout main
git checkout main 2>&1 > /dev/null

echo -e "$PUR$STAGE6$NC"