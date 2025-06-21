.PHONY: build test clean install release

# Build the application
build:
	go build -o tf-arm ./cmd

# Run tests
test:
	go test -v -race ./...

# Run tests with coverage
test-coverage:
	go test -v -race -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

# Clean build artifacts
clean:
	rm -f tf-arm
	rm -f coverage.out coverage.html
	rm -f tf-arm-*

# Install dependencies
deps:
	go mod download
	go mod tidy

# Format code
fmt:
	go fmt ./...

# Run linter
lint:
	golangci-lint run

# Build for all platforms (local testing)
build-all:
	GOOS=linux GOARCH=amd64 go build -o tf-arm-linux-amd64 ./cmd
	GOOS=linux GOARCH=arm64 go build -o tf-arm-linux-arm64 ./cmd
	GOOS=darwin GOARCH=amd64 go build -o tf-arm-darwin-amd64 ./cmd
	GOOS=darwin GOARCH=arm64 go build -o tf-arm-darwin-arm64 ./cmd
	GOOS=windows GOARCH=amd64 go build -o tf-arm-windows-amd64.exe ./cmd

# Install the binary to GOPATH/bin
install: build
	cp tf-arm $(GOPATH)/bin/

# Run example tests
test-examples: build
	@echo "Testing all example files..."
	@for file in examples/*.json *.json; do \
		if [ -f "$$file" ]; then \
			echo "Testing $$file..."; \
			./tf-arm "$$file" || exit 1; \
			echo ""; \
		fi; \
	done
	@echo "All example tests passed!"

# Development setup
dev-setup:
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install github.com/goreleaser/goreleaser@latest

# Run GoReleaser in snapshot mode (local testing)
snapshot:
	goreleaser release --snapshot --clean

# Help
help:
	@echo "Available targets:"
	@echo "  build         - Build the application"
	@echo "  test          - Run tests"
	@echo "  test-coverage - Run tests with coverage report"
	@echo "  clean         - Clean build artifacts"
	@echo "  deps          - Install dependencies"
	@echo "  fmt           - Format code"
	@echo "  lint          - Run linter"
	@echo "  build-all     - Build for all platforms"
	@echo "  install       - Install binary to GOPATH/bin"
	@echo "  test-examples - Test with example files"
	@echo "  dev-setup     - Setup development tools"
	@echo "  snapshot      - Run GoReleaser in snapshot mode"
	@echo "  help          - Show this help message"