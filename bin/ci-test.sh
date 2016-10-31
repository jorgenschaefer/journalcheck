#!/bin/sh

export GOPATH="$(pwd/gopath)"

cd "$(dirname "$0")"/..

go get .
go test ./...
