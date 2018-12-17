#!/usr/bin/env bash

pushd $GOPATH/src/github.com/kubevault/cli/hack/gendocs
go run main.go
popd
