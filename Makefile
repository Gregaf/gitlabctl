SOURCE := $(shell find . -name "*.go")

TEST_DIRS := $(shell find . -name "*_test.go" | xargs -I {} dirname {} | sort -u)
COVERAGE_FILE := coverage.out

ENTRY_DIR := cmd
ENTRY_POINTS := $(shell find $(ENTRY_DIR) -name "*.go")

BIN_DIR := bin
BIN_NAME := main

VERSION := $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")

LDFLAGS := -ldflags "-X main.version=$(VERSION) -s -w"

.DEFAULT_GOAL := help

lint: ## Lint the project
	golangci-lint run -v

test: ## Run all tests
	go test $(TEST_DIRS) -v

test-cov: ## Run all tests and output coverage file
	go test $(TEST_DIRS) -coverprofile=$(COVERAGE_FILE) -v
	go tool cover -func=$(COVERAGE_FILE)	
	go tool cover -html=$(COVERAGE_FILE)

build: $(BIN_DIR) ## Build all entrypoints for current platform
	@for entry in $(ENTRY_POINTS); do \
		entry_name=$$(basename $$(dirname $$entry)); \
		output="$(BIN_DIR)/$(BIN_NAME)"; \
		echo "Building $$output"; \
		go build $(LDFLAGS) -o $$output $$entry; \
	done

build-all: build-linux build-windows build-darwin ## Build for all platforms

$(BIN_DIR):
	mkdir -p $@ 

build-linux: $(BIN_DIR) ## Build for Linux
	@for entry in $(ENTRY_POINTS); do \
		entry_name=$$(basename $$(dirname $$entry)); \
		echo "Building $(BIN_NAME) for Linux..."; \
		GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o $(BIN_DIR)/$(BIN_NAME)-linux-amd64 $$entry; \
		GOOS=linux GOARCH=arm64 go build $(LDFLAGS) -o $(BIN_DIR)/$(BIN_NAME)-linux-arm64 $$entry; \
	done

build-windows: $(BIN_DIR) ## Build for Windows
	@for entry in $(ENTRY_POINTS); do \
		entry_name=$$(basename $$(dirname $$entry)); \
		echo "Building $(BIN_NAME) for Windows..."; \
		GOOS=windows GOARCH=amd64 go build $(LDFLAGS) -o $(BIN_DIR)/$(BIN_NAME)-windows-amd64.exe $$entry; \
	done

build-darwin: $(BIN_DIR) ## Build for macOS
	@for entry in $(ENTRY_POINTS); do \
		entry_name=$$(basename $$(dirname $$entry)); \
		echo "Building $(BIN_NAME) for macOS..."; \
		GOOS=darwin GOARCH=amd64 go build $(LDFLAGS) -o $(BIN_DIR)/$(BIN_NAME)-darwin-amd64 $$entry; \
		GOOS=darwin GOARCH=arm64 go build $(LDFLAGS) -o $(BIN_DIR)/$(BIN_NAME)-darwin-arm64 $$entry; \
	done

version: ## Show the version that would be embedded
	@echo "Version: $(VERSION)"

tidy: ## Tidy go modules
	go mod tidy

clean: ## Remove built artifacts
	@echo "Cleaning up..."
	rm -rf $(BIN_DIR)

help: ## Display this help
	@$(info APPLICATION NAME)
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z0-9_-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

.PHONY: help clean lint test test-cov build build-all build-linux build-windows build-darwin version tidy