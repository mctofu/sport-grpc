#!/bin/bash -e

docker compose run --rm protoc /protobuf/bin/protoc -I/protobuf --go_out=. --go-grpc_out=. --go_opt=paths=source_relative --go-grpc_opt=paths=source_relative --proto_path=/sportgrpc /sportgrpc/sport.proto
