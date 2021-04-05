#!/bin/bash

set -e

cov(){
	go test -v -race -coverprofile=coverage.txt ./...
	go tool cover -html=coverage.txt -o cover.html
	open cover.html
}

fmt(){
	echo fmt...
	gofmt -l -w -s .
}

lint(){
	echo lint...
	golint ./...
}

cyclo(){
	echo cyclo...
	gocyclo -top 10 .
}

all(){
	fmt && lint && cyclo
}


if [[ $# -eq 0 ]]; then
	cmd=cov
else
	cmd=${1:-cov} && shift
fi
$cmd "$@"