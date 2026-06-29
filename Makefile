BINARY_NAME=terraform-provider-mailu
VERSION=0.1.0
OS_ARCH=darwin_arm64
GO ?= go
GOCACHE ?= /tmp/go-build
GOPATH ?= /tmp/go

.PHONY: build
build:
	GOCACHE=$(GOCACHE) GOPATH=$(GOPATH) $(GO) build -o $(BINARY_NAME) .

.PHONY: test
test:
	GOCACHE=$(GOCACHE) GOPATH=$(GOPATH) $(GO) test ./...

.PHONY: tidy
tidy:
	GOCACHE=$(GOCACHE) GOPATH=$(GOPATH) $(GO) mod tidy

.PHONY: fmt
fmt:
	gofmt -w .
	terraform fmt -recursive examples templates

.PHONY: testacc
testacc:
	TF_ACC=1 GOCACHE=$(GOCACHE) GOPATH=$(GOPATH) $(GO) test ./...

.PHONY: install-local
install-local: build
	mkdir -p ~/.terraform.d/plugins/registry.terraform.io/cetic-group/mailu/$(VERSION)/$(OS_ARCH)
	cp $(BINARY_NAME) ~/.terraform.d/plugins/registry.terraform.io/cetic-group/mailu/$(VERSION)/$(OS_ARCH)/
