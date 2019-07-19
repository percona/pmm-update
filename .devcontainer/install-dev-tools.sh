#!/bin/bash

set -o errexit
set -o xtrace

yum install -y golang mc tmux
go env

go get -u -v golang.org/x/tools/cmd/gopls

#curl https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $GOPATH/bin
