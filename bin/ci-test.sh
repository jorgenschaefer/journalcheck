#!/bin/sh

set -e

export GOPATH="$(pwd)/gopath"

cd "$(dirname "$0")"/..

apt-get update -qq
apt-get install -qq libsystemd-dev
go get .

go test ./...
