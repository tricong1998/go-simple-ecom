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
PATH_PAYMENT= ./cmd/payment
PATH_PRODUCT= ./cmd/product
PATH_ORDER= ./cmd/order
# Main package path
MAIN_PATH_USER=$(PATH_USER)/cmd
MAIN_PATH_ORDER=$(PATH_ORDER)/cmd
MAIN_PATH_PRODUCT=$(PATH_PRODUCT)/cmd
MAIN_PATH_PAYMENT=$(PATH_PAYMENT)/cmd

# Binary name
BINARY_NAME_ORDER=order
BINARY_NAME_PRODUCT=product
BINARY_NAME_USER=user
BINARY_NAME_PAYMENT=payment

# Build the project
build_order:
	$(GOBUILD) -o $(BINARY_NAME_ORDER) -v $(MAIN_PATH_ORDER)

build_product:
	$(GOBUILD) -o $(BINARY_NAME_PRODUCT) -v $(MAIN_PATH_PRODUCT)

build_user:
	$(GOBUILD) -o $(BINARY_NAME_USER) -v $(MAIN_PATH_USER)

build_payment:
	$(GOBUILD) -o $(BINARY_NAME_PAYMENT) -v $(MAIN_PATH_PAYMENT)

build:
	make build_order
	make build_product
	make build_user
	make build_payment

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

run_payment:
	$(GOBUILD) -o $(BINARY_NAME_PAYMENT) -v $(MAIN_PATH_PAYMENT)
	./$(BINARY_NAME_PAYMENT)

run:
	make run_order
	make run_product
	make run_user
	make run_payment

# Clean build files
clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME_ORDER)
	rm -f $(BINARY_NAME_PRODUCT)
	rm -f $(BINARY_NAME_USER)
	rm -f $(BINARY_NAME_PAYMENT)

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
	GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(BINARY_NAME_ORDER)-linux-amd64 $(MAIN_PATH_ORDER)
	GOOS=windows GOARCH=amd64 $(GOBUILD) -o $(BINARY_NAME_ORDER)-windows-amd64.exe $(MAIN_PATH_ORDER)
	GOOS=darwin GOARCH=amd64 $(GOBUILD) -o $(BINARY_NAME_ORDER)-darwin-amd64 $(MAIN_PATH_ORDER)
	GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(BINARY_NAME_PRODUCT)-linux-amd64 $(MAIN_PATH_PRODUCT)
	GOOS=windows GOARCH=amd64 $(GOBUILD) -o $(BINARY_NAME_PRODUCT)-windows-amd64.exe $(MAIN_PATH_PRODUCT)
	GOOS=darwin GOARCH=amd64 $(GOBUILD) -o $(BINARY_NAME_PRODUCT)-darwin-amd64 $(MAIN_PATH_PRODUCT)
	GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(BINARY_NAME_USER)-linux-amd64 $(MAIN_PATH_USER)
	GOOS=windows GOARCH=amd64 $(GOBUILD) -o $(BINARY_NAME_USER)-windows-amd64.exe $(MAIN_PATH_USER)
	GOOS=darwin GOARCH=amd64 $(GOBUILD) -o $(BINARY_NAME_USER)-darwin-amd64 $(MAIN_PATH_USER)
	GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(BINARY_NAME_PAYMENT)-linux-amd64 $(MAIN_PATH_PAYMENT)
	GOOS=windows GOARCH=amd64 $(GOBUILD) -o $(BINARY_NAME_PAYMENT)-windows-amd64.exe $(MAIN_PATH_PAYMENT)
	GOOS=darwin GOARCH=amd64 $(GOBUILD) -o $(BINARY_NAME_PAYMENT)-darwin-amd64 $(MAIN_PATH_PAYMENT)

setup-dev:
	cp .env.example .env



proto-user: 
	rm -f $(PATH_USER)/pkg/pb/*.go
	protoc --proto_path=$(PATH_USER)/proto --go_out=$(PATH_USER)/pkg/pb --go_opt=paths=source_relative \
    --go-grpc_out=$(PATH_USER)/pkg/pb --go-grpc_opt=paths=source_relative \
		--grpc-gateway_out=$(PATH_USER)/pkg/pb --grpc-gateway_opt paths=source_relative \
    $(PATH_USER)/proto/*.proto

proto-payment:
	rm -f $(PATH_PAYMENT)/pkg/pb/*.go
	protoc --proto_path=$(PATH_PAYMENT)/proto --go_out=$(PATH_PAYMENT)/pkg/pb --go_opt=paths=source_relative \
    --go-grpc_out=$(PATH_PAYMENT)/pkg/pb --go-grpc_opt=paths=source_relative \
		--grpc-gateway_out=$(PATH_PAYMENT)/pkg/pb --grpc-gateway_opt paths=source_relative \
    $(PATH_PAYMENT)/proto/*.proto

proto-product:
	rm -f $(PATH_PRODUCT)/pkg/pb/*.go
	protoc --proto_path=$(PATH_PRODUCT)/proto --go_out=$(PATH_PRODUCT)/pkg/pb --go_opt=paths=source_relative \
    --go-grpc_out=$(PATH_PRODUCT)/pkg/pb --go-grpc_opt=paths=source_relative \
		--grpc-gateway_out=$(PATH_PRODUCT)/pkg/pb --grpc-gateway_opt paths=source_relative \
    $(PATH_PRODUCT)/proto/*.proto

generate-admin-account:
	@if [ -z "$(password)" ]; then \
		echo "Error: Password not provided. Usage: make generate-admin-account password=your_password"; \
		exit 1; \
	fi
	@if [ -z "$(admin)" ]; then \
		echo "Error: Admin not provided. Usage: make generate-admin-account admin=your_admin"; \
		exit 1; \
	fi
	$(GOBUILD) -o $(BINARY_NAME_USER) -v $(PATH_USER)/cmd/scripts/generate_admin_account.go
	./$(BINARY_NAME_USER) "$(admin)" "$(password)"

# Help command
help:
	@echo "Available commands:"
	@echo "  make setup-dev     - Setup development environment"
	@echo "  make build         - Build the project"
	@echo "  make build_order   - Build the order project"
	@echo "  make build_product - Build the product project"
	@echo "  make build_user    - Build the user project"
	@echo "  make build_payment - Build the payment project"
	@echo "  make run           - Run the project"
	@echo "  make run_order     - Run the order project"
	@echo "  make run_product   - Run the product project"
	@echo "  make run_user      - Run the user project"
	@echo "  make run_payment   - Run the payment project"
	@echo "  make clean         - Clean build files"
	@echo "  make test          - Run tests"
	@echo "  make test-coverage - Run tests with coverage"
	@echo "  make lint          - Run linter"
	@echo "  make deps          - Download dependencies"
	@echo "  make update-deps   - Update dependencies"
	@echo "  make build-all     - Build for multiple platforms"
	@echo "  make generate-admin-account admin=your_admin password=your_password - Generate admin account with admin and password"
	@echo "  make proto-user    - Generate proto for user"
	@echo "  make proto-payment    - Generate proto for payment"
	@echo "  make proto-product    - Generate proto for product"

.PHONY: build run clean test test-coverage lint deps update-deps build-all help proto-user proto-payment proto-product generate-admin-account
