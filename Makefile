# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
BINARY_NAME=cuack-ctl
GORELEASER := $(GOPATH)/bin/goreleaser


GO_OPT= -mod vendor

.PHONY: build
build:
	go build $(GO_OPT) -o ./bin/$(BINARY_NAME) ./cmd/cuack-ctl/

.PHONY: install
install:
	go install $(GO_OPT) ./cmd/cuack-ctl/



$(GORELEASER):
	curl -sfL https://install.goreleaser.com/github.com/goreleaser/goreleaser.sh | BINDIR=$(GOPATH)/bin sh

release: $(GORELEASER)
	$(GORELEASER) build --skip-validate --rm-dist
	$(GORELEASER) release --rm-dist

