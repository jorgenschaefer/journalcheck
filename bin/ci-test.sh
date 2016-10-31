#!/bin/sh

set -e

export GOPATH="$(pwd)/gopath"

cd "$(dirname "$0")"/..

apt-get install libsystemd-dev

go get .
go test ./...
