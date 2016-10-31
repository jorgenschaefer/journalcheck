#!/bin/sh

cd "$(dirname "$0")"/..
export GOPATH="$(pwd)/../../../../gopath"

go get .
go test ./...
