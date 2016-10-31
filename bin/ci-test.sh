#!/bin/sh

cd "$(dirname "$0")"/../..

mkdir gopath
export GOPATH="$(pwd)/gopath"
mkdir -p "$GOPATH/src/github.com/jorgenschaefer"
ln -s "$(pwd)/journalcheck-source" "$GOPATH/src/github.com/jorgenschaefer/"

go get 
