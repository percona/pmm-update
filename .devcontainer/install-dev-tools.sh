#!/bin/bash

set -o errexit
set -o xtrace

yum install -y golang mc tmux which
go env

go get -u -v golang.org/x/tools/cmd/gopls \
    github.com/acroca/go-symbols

#curl https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $GOPATH/bin
