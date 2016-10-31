#!/bin/sh

set -e

export GOPATH="$(pwd)/gopath"

echo '***'
id
pwd
which apt-get
echo '***'

cd "$(dirname "$0")"/..

sudo apt-get install libsystemd-dev

go get .
go test ./...
