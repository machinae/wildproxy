# Git tag, if any and commit hash
VERSION := $(shell git describe --tags --long 2>/dev/null || git rev-parse --short HEAD)

# Directory for output, compiled files will go in $DIR/bin
BIN_DIR := "$(shell pwd)/bin"
LDFLAGS=-ldflags "-X main.Version=$(VERSION)"

# Go commands
GOINSTALL=go install -v $(LDFLAGS) ./...

.PHONY: all get build clean install

all: build

# Get dependencies
get:
	go get -t ./...

build: get
	GOBIN=$(BIN_DIR) $(GOINSTALL)

clean:
	rm -f $(BIN_DIR)/*
	rmdir $(BIN_DIR)

test: get
	go test -v ./...

install: get
	$(GOINSTALL)
