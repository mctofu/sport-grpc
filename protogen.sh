#!/bin/bash -e

docker-compose run --rm protoc /protobuf/bin/protoc -I/protobuf --go_out=plugins=grpc:/sportgrpc --proto_path=/sportgrpc /sportgrpc/sport.proto
