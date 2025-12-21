.PHONY: all build run test clean fmt lint docker-build docker-run help

# Variables
BINARY_NAME=transactions-api
DOCKER_IMAGE=accountstack/api-transactions
DOCKER_TAG=latest
GO=go
GOFLAGS=-v

all: clean fmt build

## build: Build the application binary
build:
	@echo "Building..."
	@$(GO) build $(GOFLAGS) -o bin/$(BINARY_NAME) cmd/server/main.go

## run: Run the application
run: build
	@echo "Running..."
	@./bin/$(BINARY_NAME)

## run-dev: Run the application with development settings
run-dev:
	@echo "Running in development mode..."
	@export LOG_LEVEL=debug && \
	export DATA_PATH=../../data/seed && \
	$(GO) run cmd/server/main.go

## test: Run tests
test:
	@echo "Running tests..."
	@$(GO) test -v -race -coverprofile=coverage.out ./...

## test-coverage: Run tests with coverage report
test-coverage: test
	@$(GO) tool cover -html=coverage.out

## fmt: Format code
fmt:
	@echo "Formatting code..."
	@$(GO) fmt ./...

## lint: Run linter
lint:
	@echo "Running linter..."
	@golangci-lint run

## clean: Clean build artifacts
clean:
	@echo "Cleaning..."
	@rm -rf bin/
	@rm -f coverage.out
	@$(GO) clean

## deps: Download dependencies
deps:
	@echo "Downloading dependencies..."
	@$(GO) mod download
	@$(GO) mod tidy

## docker-build: Build Docker image
docker-build:
	@echo "Building Docker image..."
	@docker build -t $(DOCKER_IMAGE):$(DOCKER_TAG) .

## docker-run: Run Docker container
docker-run:
	@echo "Running Docker container..."
	@docker run -p 8002:8002 \
		-e PORT=8002 \
		-e LOG_LEVEL=info \
		-e CLOUDBEES_FM_API_KEY=$(CLOUDBEES_FM_API_KEY) \
		-v $(PWD)/../../data/seed:/data/seed \
		-e DATA_PATH=/data/seed \
		$(DOCKER_IMAGE):$(DOCKER_TAG)

## help: Show this help message
help:
	@echo "Usage: make [target]"
	@echo ""
	@echo "Targets:"
	@sed -n 's/^##//p' $(MAKEFILE_LIST) | column -t -s ':' | sed -e 's/^/ /'
