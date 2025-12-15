#/bin/sh
GOOS=linux GOARCH=arm64 go build -o build/sportgrpc -v github.com/mctofu/sport-grpc/cmd/server
