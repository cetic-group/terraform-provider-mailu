BINARY_NAME=terraform-provider-mailu

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
	mkdir -p ~/.terraform.d/plugins/registry.terraform.io/cetic-group/mailu/0.0.0/darwin_arm64
	cp $(BINARY_NAME) ~/.terraform.d/plugins/registry.terraform.io/cetic-group/mailu/0.0.0/darwin_arm64/
