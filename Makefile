
.PHONY: all
all:
	@$(MAKE) format
	@$(MAKE) tidy
	@$(MAKE) build
	@$(MAKE) lint
	@$(MAKE) test

.PHONY: format
format:
	go fmt ./...

.PHONY: tidy
tidy:
	go mod tidy

.PHONY: build
build:
	go generate ./... # Needed?
	CGO_ENABLED=0 GOOS=linux go build -o build/cosi-driver -buildvcs=false ./cmd

.PHONY: image
image: build
	docker build -t cloudian-cosi-driver:v0.0.0 .

.PHONY: lint
lint:
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run --timeout 10m0s ./... --max-same-issues 0; \
	else \
		docker run --rm -v $(shell pwd):/app -w /app golangci/golangci-lint:v1.62.0 golangci-lint run --timeout 10m0s ./... --max-same-issues 0; \
	fi

.PHONY: test
test:
	go test ./pkg/...
