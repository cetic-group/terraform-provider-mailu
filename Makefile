BINARY_NAME=terraform-provider-mailu
VERSION=0.1.0
OS_ARCH=darwin_arm64

.PHONY: build
build:
	go build -o $(BINARY_NAME) .

.PHONY: test
test:
	go test ./...

.PHONY: tidy
tidy:
	go mod tidy

.PHONY: install-local
install-local: build
	mkdir -p ~/.terraform.d/plugins/registry.terraform.io/cetic-group/mailu/$(VERSION)/$(OS_ARCH)
	cp $(BINARY_NAME) ~/.terraform.d/plugins/registry.terraform.io/cetic-group/mailu/$(VERSION)/$(OS_ARCH)/
