# +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
#                          Stage 1: Build the application
# +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++

FROM golang:1.18-alpine AS builder

# make the directory for the source code
RUN mkdir -p ./go/src/Tavern-Backend/env/ && mkdir -p ./Files/env && mkdir -p ./Files/awslib

# copy go.mod and go.sum to the working directory
COPY ./go.* /go/src/Tavern-Backend/

# set the working directory to be in the GOROOT directory
WORKDIR /go/src/Tavern-Backend

# RUN the command to download the dependencies
RUN go mod download

WORKDIR /

# copy the source code to the working directory
COPY . /go/src/Tavern-Backend/
COPY ./env /Files/env/
COPY ./aws /Files/awslib/

WORKDIR /go/src/Tavern-Backend

# Now it should build
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "-s -w" -o /Files/TavernProfile

# +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
#                          Stage 2: Create the final image
# +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++

# create alpine image to run the application 
FROM alpine:latest

WORKDIR /root/
RUN mkdir -p /env
RUN mkdir -p ./lib/log && mkdir -p ./awslib
# RUN from the build stage list the files in the directory
RUN ls -la
COPY --from=builder /Files/ .
# Make a file to hold logs 
RUN touch ./lib/log/tavern.log
# expose the port 8000
EXPOSE 8000

# command to run on container start
ENTRYPOINT [ "./TavernProfile", "prod" ]
