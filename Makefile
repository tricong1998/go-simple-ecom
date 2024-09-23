# Makefile for ecombase project

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod
GOLINT=golangci-lint


# SERVICE PATH
PATH_USER= ./cmd/user

# Main package path
MAIN_PATH_USER=$(PATH_USER)/cmd
MAIN_PATH_ORDER=./cmd/order/cmd
MAIN_PATH_PRODUCT=./cmd/product/cmd

# Binary name
BINARY_NAME_ORDER=order
BINARY_NAME_PRODUCT=product
BINARY_NAME_USER=user

# Build the project
build_order:
	$(GOBUILD) -o $(BINARY_NAME_ORDER) -v $(MAIN_PATH_ORDER)

build_product:
	$(GOBUILD) -o $(BINARY_NAME_PRODUCT) -v $(MAIN_PATH_PRODUCT)

build_user:
	$(GOBUILD) -o $(BINARY_NAME_USER) -v $(MAIN_PATH_USER)

# Run the project
run_order:
	$(GOBUILD) -o $(BINARY_NAME_ORDER) -v $(MAIN_PATH_ORDER)
	./$(BINARY_NAME_ORDER)

run_user:
	$(GOBUILD) -o $(BINARY_NAME_USER) -v $(MAIN_PATH_USER)
	./$(BINARY_NAME_USER)

run_product:
	$(GOBUILD) -o $(BINARY_NAME_PRODUCT) -v $(MAIN_PATH_PRODUCT)
	./$(BINARY_NAME_PRODUCT)	


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



proto-user: 
	rm -f $(PATH_USER)/pkg/pb/*.go
	protoc --proto_path=$(PATH_USER)/proto --go_out=$(PATH_USER)/pkg/pb --go_opt=paths=source_relative \
    --go-grpc_out=$(PATH_USER)/pkg/pb --go-grpc_opt=paths=source_relative \
		--grpc-gateway_out=$(PATH_USER)/pkg/pb --grpc-gateway_opt paths=source_relative \
    $(PATH_USER)/proto/*.proto

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

.PHONY: build run clean test test-coverage lint deps update-deps build-all help proto