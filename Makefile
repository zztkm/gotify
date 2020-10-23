.DEFAULT_GOAL := help
BIN := gotify
VERSION := $$(make -s show-version)
CURRENT_REVISION := $(shell git rev-parse --short HEAD)
BUILD_LDFLAGS := "-s -w -X main.revision=$(CURRENT_REVISION)"
export GO111MODULE=on

## Setup tools
setup:
	GO111MODULE=off	go get \
	github.com/Songmu/make2help/cmd/make2help

## Run clean & build
all: clean build

## Build binary
build:
	go build -ldflags=$(BUILD_LDFLAGS) -o $(BIN) .

## Clean repository
clean:
	rm -rf $(BIN)
	go clean

## Run test
test: build
	go test -v ./...

## Show help
help:
	@make2help ${MAKEFILE_LIST}

.PHONY: all build clean test help