#/bin/sh
GOOS=linux GOARCH=arm GOARM=7 go build -o build/sportgrpc -v github.com/mctofu/sport-grpc/cmd/server
