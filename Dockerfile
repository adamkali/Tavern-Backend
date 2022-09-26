FROM golang:1.18-bullseye AS builder

# make the working directory as /app
WORKDIR /app

# get the source code from git and copy it to the working directory
RUN git clone https://github.com/adamkali/Tavern-Backend.git && git checkout Bëor

# copy the source code to the working directory
COPY go.* ./

# RUN the command to download the dependencies
RUN go mod download

# copy the source code to the working directory with recursion 
COPY Tavern-Backend .

# build the application
RUN  go build -ldflags "-s -w" -o TavernProfile

# expose the port 8000
EXPOSE 8000

# command to run on container start
ENTRYPOINT [ "/TavernProfile", "prod" ]
