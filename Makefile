BINARY_NAME=terraform-provider-mailu
VERSION=0.1.0
OS_ARCH=darwin_arm64
GO ?= go
GOCACHE ?= /tmp/go-build
GOPATH ?= /tmp/go
GORELEASER ?= goreleaser
GITLEAKS_VERSION ?= v8.28.0

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
	terraform fmt -recursive examples

.PHONY: fmt-check
fmt-check:
	test -z "$$(gofmt -l .)"
	terraform fmt -check -recursive examples

.PHONY: docs
docs:
	GOCACHE=$(GOCACHE) GOPATH=$(GOPATH) go run github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs generate --provider-name mailu

.PHONY: secret-scan
secret-scan:
	GOCACHE=$(GOCACHE) GOPATH=$(GOPATH) $(GO) run github.com/zricethezav/gitleaks/v8@$(GITLEAKS_VERSION) git --redact --verbose

.PHONY: release-snapshot
release-snapshot:
	GOCACHE=$(GOCACHE) GOPATH=$(GOPATH) $(GORELEASER) release --snapshot --clean

.PHONY: testacc
testacc:
	TF_ACC=1 GOCACHE=$(GOCACHE) GOPATH=$(GOPATH) $(GO) test ./...

.PHONY: install-local
install-local: build
	mkdir -p ~/.terraform.d/plugins/registry.terraform.io/cetic-group/mailu/$(VERSION)/$(OS_ARCH)
	cp $(BINARY_NAME) ~/.terraform.d/plugins/registry.terraform.io/cetic-group/mailu/$(VERSION)/$(OS_ARCH)/
