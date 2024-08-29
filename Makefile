# Makefile for ecombase project

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod
GOLINT=golangci-lint

# Main package path
MAIN_PATH=./cmd/server

# Binary name
BINARY_NAME=goecom

# Build the project
build:
	$(GOBUILD) -o $(BINARY_NAME) -v $(MAIN_PATH)

# Run the project
run:
	$(GOBUILD) -o $(BINARY_NAME) -v $(MAIN_PATH)
	./$(BINARY_NAME)

# Clean build files
clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)

# Run tests
test:
	$(GOTEST) -v ./...

# Run tests with coverage
test-coverage:
	$(GOTEST) -v -coverprofile=coverage.out ./...
	$(GOCMD) tool cover -html=coverage.out -o coverage.html

# Run linter
lint:
	$(GOLINT) run

# Download dependencies
deps:
	$(GOGET) -v -t -d ./...
	$(GOMOD) tidy

# Update dependencies
update-deps:
	$(GOGET) -u -v -t -d ./...
	$(GOMOD) tidy

# Build for multiple platforms
build-all:
	GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(BINARY_NAME)-linux-amd64 $(MAIN_PATH)
	GOOS=windows GOARCH=amd64 $(GOBUILD) -o $(BINARY_NAME)-windows-amd64.exe $(MAIN_PATH)
	GOOS=darwin GOARCH=amd64 $(GOBUILD) -o $(BINARY_NAME)-darwin-amd64 $(MAIN_PATH)

setup-dev:
	cp .env.example .env

# Help command
help:
	@echo "Available commands:"
	@echo "  make setup-dev     - Setup development environment"
	@echo "  make build         - Build the project"
	@echo "  make run           - Run the project"
	@echo "  make clean         - Clean build files"
	@echo "  make test          - Run tests"
	@echo "  make test-coverage - Run tests with coverage"
	@echo "  make lint          - Run linter"
	@echo "  make deps          - Download dependencies"
	@echo "  make update-deps   - Update dependencies"
	@echo "  make build-all     - Build for multiple platforms"

.PHONY: build run clean test test-coverage lint deps update-deps build-all help