.PHONY: all build run clean

all: check test-small build

build: ## Build the project
	go build -v -o bin/server .

run: ## Run the server
	go run .

clean: ## Clean build artifacts
	rm -rf bin/
	rm -f coverage.txt

.PHONY: test-small test-medium test-large test-all test-cover

test-small: ## Run small size tests
	go test -v ./...

test-medium: ## Run medium size tests
	go test -tags=medium -v ./...

test-large: ## Run large size tests
	go test -tags=large -v ./...

test-all: test-small test-medium test-large ## Run all tests

test-cover: ## Run tests with coverage
	go test -race -coverprofile=coverage.txt -covermode=atomic ./...
	go tool cover -func=coverage.txt

.PHONY: generate generate-api generate-db

generate: generate-api generate-db ## Run all code generation

generate-api: ## Generate OpenAPI code
	oapi-codegen -config oapi-codegen.yaml ../spec/openapi.yml

generate-db: ## Generate sqlc code
	cd ../database && sqlc generate
	@for dir in internal/*/; do \
		[ -f "$${dir}db.go" ] && mv "$${dir}db.go" "$${dir}db_gen.go" || true; \
		[ -f "$${dir}models.go" ] && mv "$${dir}models.go" "$${dir}models_gen.go" || true; \
	done

.PHONY: check lint fmt

check: lint fmt ## Run lint and format checks

lint: ## Run lint
	golangci-lint run

fmt: ## Run go fmt
	go fmt ./...

.PHONY: help

help: ## Display this help screen
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'
