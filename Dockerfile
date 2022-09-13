FROM golang:1.19-bullseye AS builder

# make the working directory as /app
WORKDIR /app

# copy the source code to the working directory
COPY go.* ./

# RUN the command to download the dependencies
RUN go mod download

# copy the source code to the working directory
COPY . ./

# build the application
RUN go build -o main /TavernProfile

# expose the port 8000
EXPOSE 8000

# command to run on container start
ENTRYPOINT [ "/TavernProfile", "prod" ]
