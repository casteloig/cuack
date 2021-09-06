# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
BINARY_NAME=cuack


GO_OPT= -mod vendor

.PHONY: build
build:
	go build $(GO_OPT) -o ./bin/$(BINARY_NAME) ./cmd/cuack-ctl/

.PHONY: install
install:
	go install $(GO_OPT) ./cmd/cuack-ctl/