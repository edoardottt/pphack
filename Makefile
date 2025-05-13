# Makefile for pphack

REPO        := github.com/edoardottt/pphack
BINARY_NAME := pphack
CMD_DIR     := ./cmd/$(BINARY_NAME)
BIN_PATH    := /usr/local/bin/$(BINARY_NAME)

.PHONY: all remod update lint build clean test install uninstall

all: build

remod:
	@echo "Reinitializing Go module..."
	@rm -f go.mod go.sum
	@go mod init $(REPO)
	@go get ./...
	@go mod tidy -v
	@echo "Go module reinitialized."

update:
	@echo "Updating dependencies..."
	@go get -u ./...
	@go mod tidy -v
	@echo "Dependencies updated."

lint:
	@echo "Running linter..."
	@golangci-lint run
	@echo "Linting complete."

build:
	@echo "Building $(BINARY_NAME)..."
	@go build -o $(BINARY_NAME) $(CMD_DIR)
	@echo "Build complete."

install: build
	@echo "Installing $(BINARY_NAME) to $(BIN_PATH)..."
	@sudo mv $(BINARY_NAME) $(BIN_PATH)
	@echo "Installed successfully."

uninstall:
	@echo "Uninstalling $(BINARY_NAME) from $(BIN_PATH)..."
	@sudo rm -f $(BIN_PATH)
	@echo "Uninstalled."

test:
	@echo "Running tests with race detector..."
	@go test -race ./...
	@echo "Tests complete."

clean:
	@echo "Cleaning build artifacts..."
	@rm -f $(BINARY_NAME)
	@echo "Clean complete."
