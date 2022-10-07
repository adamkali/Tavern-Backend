#!/bin/bash

# Get the commit message
COMMIT_MESSAGE=$1

# if there is no comit message, throw an error and quit out
if [ -z "$COMMIT_MESSAGE" ]; then
    echo "No commit message provided"
    exit 1
fi

# create a quit function
quit() {
    echo "Build failed"
    exit 1
}

# create a git function and make sure there 
# are no errors
# $1 is the commit message
# throw away any output
gitstep() {
    # git add -A
    # git commit -m "$COMMIT_MESSAGE"
    # then checkout the Beor and throw away output to avoid printing it
    git checkout Beor                   &> /dev/null
    git merge main                      &> /dev/null

    git add -A                          &> /dev/null
    git commit -m "$COMMIT_MESSAGE"     &> /dev/null 

    # git push origin beor
    git push origin Beor                &> /dev/null 
}

versionstep() {
# check if there is a VERSION.yaml file
# if not make one throw away any output to the terminal
    if [ ! -f ./cmd/VERSION.yaml ]; then
        echo "major: 0" >  ./cmd/VERSION.yaml       &> /dev/null
        echo "minor: 1" >> ./cmd/VERSION.yaml       &> /dev/null
    fi
    
    # check if -M is set
    if [ "$2" = "-M" ]; then
        # increment major version
        echo "Incrementing major version"
        # Get the current major version and increment it
        MAJOR=$(cat ./cmd/VERSION.yaml | grep major | cut -d " " -f 2)
        # Increment the major version
        MAJOR=$((MAJOR+1))
        # Set the minor version to 0
        MINOR=0
    else
        # increment minor version
        echo "Incrementing minor version"
        #save the major version
        MAJOR=$(cat ./cmd/VERSION.yaml | grep major | cut -d " " -f 2)
        # Get the current minor version and increment it
        MINOR=$(cat ./cmd/VERSION.yaml | grep minor | cut -d " " -f 2)
        # Increment the minor version
        MINOR=$((MINOR+1))
    fi
    
    # update the VERSION.yaml file and throw away any output to the terminal
    echo "major: $MAJOR" >  ./cmd/VERSION.yaml      &> /dev/null
    echo "minor: $MINOR" >> ./cmd/VERSION.yaml      &> /dev/null
}

# Setup a progress bar
PUR='\033[0;35m'
BLU='\033[0;34m'
GRE='\033[0;32m'
NCR='\033[0m' # No Color

STAGE0="Pulling main brach from git"
STAGE1="Updating VERSION"
STAGE2="Logging into AWS ECR"
STAGE3="Building docker image"
STAGE4="Tagging docker image"
STAGE5="Pushing docker image to ECR"
STAGE6="COMPLETE!"

PROG0="[${BLU}=>${PUR}--------------------------------------${NCR}] 0%"
PROG1="[${BLU}######${PUR}=>--------------------------------${NCR}] 20%"
PROG2="[${BLU}##############${PUR}=>------------------------${NCR}] 40%"
PROG3="[${BLU}######################${PUR}=>----------------${NCR}] 60%"
PROG4="[${BLU}##############################${PUR}=>--------${NCR}] 80%"
PROG5="[${BLU}######################################${PUR}=>${NCR}] 100%"

echo -e "${PUR}${STAGE0}${NC}${PROG0}\r"
# do the git stuff and if there is an error, quit
gitstep || quit
echo -e "${PUR}${STAGE1}${NCR}${PROG1}\r"
versionstep || quit
# build the docker image and throw away any output to the terminal
echo -e "${PUR}${STAGE2}${NCR}${PROG2}\r"
aws ecr get-login-password --region us-east-1 | docker login --username AWS --password-stdin 739810740537.dkr.ecr.us-east-1.amazonaws.com &> /dev/null || quit

echo -e "${PUR}${STAGE3}${NCR}${PROG3}\r"
# docker build -t tavern-profile-beor . &> /dev/null || quit
docker build -t tavern-profile-beor .|| quit

echo -e "${PUR}${STAGE4}${NCR}${PROG4}\r"
docker tag tavern-profile-beor:latest 739810740537.dkr.ecr.us-east-1.amazonaws.com/tavern-profile-beor:$MAJOR.$MINOR &> /dev/null || quit

echo -e "${PUR}${STAGE5}${NCR}${PROG5}\r"
docker push 739810740537.dkr.ecr.us-east-1.amazonaws.com/tavern-profile-beor:$MAJOR.$MINOR &> /dev/null || quit

# git checkout main
git checkout main &> /dev/null

echo -e "${GRE}$STAGE6${NCR}"