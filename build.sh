#!/bin/bash

export PATH=/mnt/d/download/go/bin/:$PATH
export GOPATH=/home/fish/dipamkara/
#export GOPATH=/mnt/d/workspace/dipamkara/
export GO111MODULE=on
export CGO_ENABLED=1
export GLIDE_HOME=$GOPATH
export GOPROXY=https://goproxy.io
export PATH=`pwd`/auto/:$PATH
echo "start build...$#"

protoc --go_out=plugins=grpc:. hello/hello.proto
go build .
mv ./hello
