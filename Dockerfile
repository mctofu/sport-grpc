FROM golang:1.25

RUN apt-get update && apt-get install -y unzip wget

ARG PROTOBUF_VERSION=33.2

RUN mkdir -p /protobuf && \
  mkdir -p /tools && \
  wget "https://github.com/protocolbuffers/protobuf/releases/download/v$PROTOBUF_VERSION/protoc-$PROTOBUF_VERSION-linux-x86_64.zip" && \
  unzip "protoc-$PROTOBUF_VERSION-linux-x86_64.zip" -d "/protobuf" && \
  rm "protoc-$PROTOBUF_VERSION-linux-x86_64.zip"

WORKDIR /tools

RUN go mod init tools && \
  go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.36.11 && \
  go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.6.0

RUN mkdir -p /sportgrpc

WORKDIR /sportgrpc
