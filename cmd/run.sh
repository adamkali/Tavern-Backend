#!/bin/zsh
# run go clean
go clean
# run go build
go build -ldflags "-s -w" -o ./TavernProfile
# make sure the binary is executable
chmod +x ./TavernProfile
# run the binary with prod argument
./TavernProfile prod    


