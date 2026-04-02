##############################################################################
# worker-service Makefile
##############################################################################

BINARY        := worker
DLQ_READER    := dlq-reader
CMD_WORKER    := ./cmd/worker
CMD_DLQ       := ./cmd/dlq-reader
BUILD_DIR     := ./bin
GOFLAGS       := -ldflags="-s -w"

.PHONY: all build test lint vet fmt clean run run-dlq \
        docker-build docker-up docker-down migrate

## ── Build ──────────────────────────────────────────────────────────────────

all: build

build: ## Build all binaries
	@mkdir -p $(BUILD_DIR)
	go build $(GOFLAGS) -o $(BUILD_DIR)/$(BINARY)    $(CMD_WORKER)
	go build $(GOFLAGS) -o $(BUILD_DIR)/$(DLQ_READER) $(CMD_DLQ)
	@echo "✓ Binaries in $(BUILD_DIR)/"

## ── Run ────────────────────────────────────────────────────────────────────

run: ## Run the worker locally (requires .env)
	go run $(CMD_WORKER)


run-dlq: ## Tail the dead-letter queue
	go run $(CMD_DLQ)

## ── Test ───────────────────────────────────────────────────────────────────

test: ## Run all unit tests
	go test ./... -v -race -count=1

test-short: ## Run tests without the race detector (fast)
	go test ./... -count=1

cover: ## Generate HTML coverage report
	go test ./... -coverprofile=coverage.out -race
	go tool cover -html=coverage.out -o coverage.html
	@echo "✓ Coverage report: coverage.html"

## ── Code quality ───────────────────────────────────────────────────────────

vet: ## Run go vet
	go vet ./...

fmt: ## Format all Go files
	gofmt -s -w .

lint: ## Run golangci-lint (must be installed)
	golangci-lint run ./...

## ── Docker ─────────────────────────────────────────────────────────────────

docker-build: ## Build the Docker image
	docker build -t worker-service:latest .

docker-up: ## Start all services with docker-compose
	docker compose up --build -d

docker-down: ## Stop all docker-compose services
	docker compose down

docker-logs: ## Follow worker logs
	docker compose logs -f worker

## ── Database ───────────────────────────────────────────────────────────────

migrate: ## Apply DB migrations (runs at startup automatically; use for manual apply)
	@echo "Migrations run automatically on startup. Start the worker to apply."

## ── Helpers ────────────────────────────────────────────────────────────────

clean: ## Remove build artefacts
	rm -rf $(BUILD_DIR) coverage.out coverage.html

tidy: ## Tidy and verify go modules
	go mod tidy
	go mod verify

help: ## Show this help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-20s\033[0m %s\n", $$1, $$2}'
