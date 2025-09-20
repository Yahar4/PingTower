include .env
export

# Project configuration
PROJECT_NAME := ping-tower
CMD_PATH := cmd/pingtower/main.go
BINARY_NAME := $(PROJECT_NAME)

# Go variables
GO := go
GOBUILD := $(GO) build
GOCLEAN := $(GO) clean
GOTEST := $(GO) test
GOMOD := $(GO) mod

# Migration configuration
MIGRATION_PATH := ./database/migrations
DSN := $(DSN)

# Build flags
LDFLAGS := -ldflags "-X main.Version=$(VERSION)"

.PHONY: all build clean test run install

all: build

build:
	@echo "Building $(PROJECT_NAME)..."
	$(GOBUILD) $(LDFLAGS) -o bin/$(BINARY_NAME) $(CMD_PATH)

run:
	@echo "Running $(PROJECT_NAME)..."
	$(GO) run $(CMD_PATH)

clean:
	@echo "Cleaning..."
	$(GOCLEAN)
	rm -f bin/$(BINARY_NAME)

test:
	@echo "Testing..."
	$(GOTEST) -v ./...

install:
	@echo "Installing dependencies..."
	$(GOMOD) download

dev:
	@echo "Starting in development mode..."
	$(GO) run $(CMD_PATH)

watch:
	@echo "Watching for changes..."
	reflex -r '\.go$$' -s -- sh -c 'make run'

docker-build:
	@echo "Building Docker image..."
	docker build -t $(PROJECT_NAME):$(VERSION) .

docker-run:
	@echo "Running in Docker..."
	docker run -p 8080:8080 $(PROJECT_NAME):$(VERSION)

docker-compose:
	@echo "Starting Docker conatiners..."
	docker-compose up -d

# Database migration commands
db-status:
	@GOOSE_DRIVER=postgres GOOSE_DBSTRING=$(DSN) goose -dir=$(MIGRATION_PATH) status

db-up:
	@GOOSE_DRIVER=postgres GOOSE_DBSTRING=$(DSN) goose -dir=$(MIGRATION_PATH) up

db-down:
	@GOOSE_DRIVER=postgres GOOSE_DBSTRING=$(DSN) goose -dir=$(MIGRATION_PATH) down

db-reset:
	@GOOSE_DRIVER=postgres GOOSE_DBSTRING=$(DSN) goose -dir=$(MIGRATION_PATH) reset

db-create:
	@echo "Enter migration name: " && read name && \
	GOOSE_DRIVER=postgres GOOSE_DBSTRING=$(DB_URL) goose -dir=$(MIGRATION_PATH) create "$$name" sql

# ---- Help ----
help:
	@echo "Available targets:"
	@echo "  build          - Build the application"
	@echo "  run            - Run the application"
	@echo "  clean          - Clean build artifacts"
	@echo "  test           - Run tests"
	@echo "  install        - Install dependencies"
	@echo "  dev            - Run in development mode"
	@echo "  watch          - Watch for changes and restart"
	@echo ""
	@echo "Database migration targets:"
	@echo "  db-status      - Show migration status"
	@echo "  db-up          - Apply all up migrations"
	@echo "  db-down        - Apply one down migration"
	@echo "  db-reset       - Reset database (down then up)"
	@echo "  db-create      - Create new migration file"
	@echo ""
	@echo "Docker targets:"
	@echo "  docker-build   - Build Docker image"
	@echo "  docker-run     - Run in Docker container"
	@echo "  docker-compose - Start containers with docker-compose"
