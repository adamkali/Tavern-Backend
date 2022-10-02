# run go clean
go clean
# run go build
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "-s -w" -o ./TavernProfile ./main.go

