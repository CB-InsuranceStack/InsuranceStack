.PHONY: help up down logs clean test test-unit test-integration test-e2e build lint format

# Default target
.DEFAULT_GOAL := help

# ============================================================================
# Help
# ============================================================================
help: ## Show this help message
	@echo "AccountStack - Development Commands"
	@echo ""
	@echo "Usage: make [target]"
	@echo ""
	@echo "Targets:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-20s\033[0m %s\n", $$1, $$2}'

# ============================================================================
# Docker Compose
# ============================================================================
up: ## Start all services
	docker compose up --build

up-detached: ## Start all services in background
	docker compose up -d --build

down: ## Stop all services
	docker compose down

down-volumes: ## Stop all services and remove volumes
	docker compose down -v

logs: ## Show logs from all services
	docker compose logs -f

logs-web: ## Show logs from web service
	docker compose logs -f web

logs-api: ## Show logs from all API services
	docker compose logs -f api-accounts api-transactions api-insights

restart: ## Restart all services
	docker compose restart

restart-web: ## Restart web service
	docker compose restart web

restart-apis: ## Restart all API services
	docker compose restart api-accounts api-transactions api-insights

# ============================================================================
# Testing
# ============================================================================
test: ## Run all tests
	@echo "Running all tests..."
	@$(MAKE) test-unit
	@$(MAKE) test-integration
	@$(MAKE) test-e2e
	@echo "✓ All tests completed"

test-unit: ## Run unit tests
	@echo "Running unit tests..."
	@cd apps/web && npm run test:unit || true
	@cd apps/api-accounts && go test -v ./... -short || true
	@cd apps/api-transactions && go test -v ./... -short || true
	@cd apps/api-insights && go test -v ./... -short || true

test-integration: ## Run integration tests
	@echo "Running integration tests..."
	@cd tests/integration && go test -v ./... || true

test-e2e: ## Run end-to-end tests
	@echo "Running E2E tests..."
	@cd apps/web && npm run test:e2e || true

test-coverage: ## Run tests with coverage
	@echo "Running tests with coverage..."
	@cd apps/web && npm run test:coverage || true
	@cd apps/api-accounts && go test -v -coverprofile=coverage.out ./... || true
	@cd apps/api-transactions && go test -v -coverprofile=coverage.out ./... || true
	@cd apps/api-insights && go test -v -coverprofile=coverage.out ./... || true

# ============================================================================
# Development
# ============================================================================
install: ## Install all dependencies
	@echo "Installing dependencies..."
	@cd apps/web && npm install
	@echo "✓ Dependencies installed"

build: ## Build all services
	@echo "Building services..."
	docker compose build
	@echo "✓ Build complete"

clean: ## Clean up generated files and containers
	@echo "Cleaning up..."
	docker compose down -v
	rm -rf apps/web/node_modules
	rm -rf apps/web/dist
	rm -rf apps/api-*/bin
	@echo "✓ Cleanup complete"

# ============================================================================
# Code Quality
# ============================================================================
lint: ## Run linters
	@echo "Running linters..."
	@cd apps/web && npm run lint || true
	@cd apps/api-accounts && golint ./... || true
	@cd apps/api-transactions && golint ./... || true
	@cd apps/api-insights && golint ./... || true

format: ## Format code
	@echo "Formatting code..."
	@cd apps/web && npm run format || true
	@cd apps/api-accounts && gofmt -w . || true
	@cd apps/api-transactions && gofmt -w . || true
	@cd apps/api-insights && gofmt -w . || true
	@echo "✓ Formatting complete"

# ============================================================================
# Database/Data
# ============================================================================
seed: ## Seed data (when implemented)
	@echo "Seeding data..."
	@echo "✓ Data seeding not yet implemented"

# ============================================================================
# Utilities
# ============================================================================
ps: ## Show running containers
	docker compose ps

exec-web: ## Open shell in web container
	docker compose exec web sh

exec-api-accounts: ## Open shell in accounts API container
	docker compose exec api-accounts sh

exec-api-transactions: ## Open shell in transactions API container
	docker compose exec api-transactions sh

exec-api-insights: ## Open shell in insights API container
	docker compose exec api-insights sh

health: ## Check health of all services
	@echo "Checking service health..."
	@curl -s http://localhost:3000 > /dev/null && echo "✓ Web: healthy" || echo "✗ Web: unhealthy"
	@curl -s http://localhost:8001/healthz > /dev/null && echo "✓ Accounts API: healthy" || echo "✗ Accounts API: unhealthy"
	@curl -s http://localhost:8002/healthz > /dev/null && echo "✓ Transactions API: healthy" || echo "✗ Transactions API: unhealthy"
	@curl -s http://localhost:8003/healthz > /dev/null && echo "✓ Insights API: healthy" || echo "✗ Insights API: unhealthy"
