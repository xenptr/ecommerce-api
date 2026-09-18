default: help

include .env
export

# Application
APP_NAME := $(shell basename "$(PWD)")
APP_ENTRYPOINT := ./cmd/api

BIN_DIR := out
APP_BINARY := ./$(BIN_DIR)/$(APP_NAME)

AIR_DIR := tmp
AIR_BINARY := ./$(AIR_DIR)/main

# Database
MIGRATION_PATH := ./db/migrations
DATABASE_URL := postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable
MIGRATE := migrate -path $(MIGRATION_PATH) -database "$(DATABASE_URL)"

# Colors
MAUVE    := $(shell tput setaf 5)
BLUE     := $(shell tput setaf 4)
LAVENDER := $(shell tput setaf 13)
SUBTEXT  := $(shell tput setaf 7)
GREEN    := $(shell tput setaf 2)
RESET    := $(shell tput sgr0)


## Database

migrate-create: ## Create a new migration. Usage: make migrate-create name=create_users
	@echo "  > Creating migration: $(name)"
	@migrate create -ext sql -dir $(MIGRATION_PATH) -seq $(name)

migrate-up: ## Apply all pending migrations
	@echo "  > Applying database migrations..."
	@$(MIGRATE) up

migrate-down: ## Roll back the last migration
	@echo "  > Rolling back last migration..."
	@$(MIGRATE) down 1

migrate-force: ## Force migration version. Usage: make migrate-force version=1
	@echo "  > Forcing migration version: $(version)"
	@$(MIGRATE) force $(version)


## Application

build: ## Build the application
	@echo "  > Building $(APP_NAME)..."
	@mkdir -p $(BIN_DIR)
	@go build -o $(APP_BINARY) $(APP_ENTRYPOINT)
	@echo "  > $(GREEN)Build successful$(RESET)"

run: build ## Build and run the application
	@echo "  > Starting $(APP_NAME)..."
	@$(APP_BINARY)

clean: ## Remove build artifacts
	@echo "  > Cleaning build artifacts..."
	@go clean
	@rm -rf $(BIN_DIR) $(AIR_DIR)


## Development

air-build: ## Build application binary for Air
	@mkdir -p $(AIR_DIR)
	@go build -o $(AIR_BINARY) $(APP_ENTRYPOINT)

dev: ## Start development server with Air
	@echo "  > Starting development server..."
	@air

fmt: ## Format Go source files
	@echo "  > Formatting Go source files..."
	@go fmt ./...

vet: ## Run go vet
	@echo "  > Running go vet..."
	@go vet ./...

tidy: ## Tidy Go dependencies
	@echo "  > Tidying Go dependencies..."
	@go mod tidy

test: ## Run tests
	@echo "  > Running tests..."
	@go test ./...

check: fmt vet test ## Run formatting, vet, and tests

compile: check build ## Check and build the application


## Help

help: ## Show this help
	@echo ""
	@echo "  $(MAUVE)make$(RESET) $(SUBTEXT)<target>$(RESET)"
	@echo ""
	@echo "  $(BLUE)Targets$(RESET)"
	@awk 'BEGIN {FS = ":.*?## "} \
		/^[a-zA-Z_-]+:.*?##/ {printf "    $(LAVENDER)%-20s$(RESET) $(SUBTEXT)%s$(RESET)\n", $$1, $$2}' \
		$(MAKEFILE_LIST)


.PHONY: \
	migrate-create \
	migrate-up \
	migrate-down \
	migrate-force \
	build \
	run \
	clean \
	air-build \
	dev \
	fmt \
	vet \
	tidy \
	test \
	check \
	compile \
	help
