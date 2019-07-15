#!/usr/bin/env bash

pushd $GOPATH/src/kubevault.dev/cli/hack/gendocs
go run main.go
popd
