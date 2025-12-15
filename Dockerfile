FROM golang:1.25

RUN apt-get update && apt-get install -y unzip wget

ARG PROTOBUF_VERSION=3.11.4

RUN mkdir -p /protobuf && \
  mkdir -p /tools && \
  wget "https://github.com/protocolbuffers/protobuf/releases/download/v$PROTOBUF_VERSION/protoc-$PROTOBUF_VERSION-linux-x86_64.zip" && \
  unzip "protoc-$PROTOBUF_VERSION-linux-x86_64.zip" -d "/protobuf" && \
  rm "protoc-$PROTOBUF_VERSION-linux-x86_64.zip"

WORKDIR /tools

RUN go mod init tools && \
  go get github.com/golang/protobuf/protoc-gen-go@v1.3.5 && \
  go get google.golang.org/grpc@v1.28.1

RUN mkdir -p /sportgrpc

WORKDIR /sportgrpc
