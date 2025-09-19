.PHONY: build docs deps dev run fmt sec test migrate migrate-status migrate-new setup clean install help

help:
	@echo "Available commands:"
	@echo "  build           - Build the application"
	@echo "                    Usage examples:"
	@echo "                      make build                    # Auto-detect: valid semver tag or commit hash"
	@echo "                      make build VERSION=v1.2.3     # Override with explicit version"
	@echo "                      make build COMMIT=abc1234     # Override with explicit commit"
	@echo "                    Note: Only valid semver tags (v1.2.3, v2.0.0-beta.1) are used;"
	@echo "                          invalid tags fall back to commit hash"
	@echo "  docs            - Generate API documentation (Swagger)"
	@echo "  deps            - Download dependencies"
	@echo "  dev             - Run with hot reload (requires air)"
	@echo "  run             - Run the application"
	@echo "  fmt             - Format code"
	@echo "  sec             - Check for security vulnerabilities"
	@echo "  test            - Run tests"
	@echo "  migrate         - Apply database migrations"
	@echo "  migrate-status  - Show migration status"
	@echo "  migrate-new     - Create new migration (usage: make migrate-new name=migration_name)"
	@echo "  setup           - Setup development environment"
	@echo "  clean           - Clean build artefacts"
	@echo "  install         - Install application (requires root)"
	@echo "  help            - Show this help message"

# Auto-detect version from git or use provided override
# Validates semver format for tags, falls back to commit hash if invalid
# Regex pattern validates: v?MAJOR.MINOR.PATCH(-prerelease)?(+build)?
VERSION ?= $(shell \
	git describe --tags --exact-match 2>/dev/null | \
	grep -E '^v?[0-9]+\.[0-9]+\.[0-9]+(-[0-9A-Za-z-]+(\.[0-9A-Za-z-]+)*)?(\+[0-9A-Za-z-]+(\.[0-9A-Za-z-]+)*)?$$' >/dev/null && \
	git describe --tags --exact-match 2>/dev/null || \
	git rev-parse --short HEAD 2>/dev/null || \
	echo "dev" \
)

build:
	@echo "Building the application..."
	@echo "Using version: $(VERSION)"
	go build -ldflags="-X 'main.Version=$(VERSION)'" -o target/api ./cmd/api
	@echo "Build complete! Binary available at target/api"

run: build
	@echo "Starting the API..."
	./target/api

fmt:
	@echo "Running fmt..."
	go fmt ./...

sec:
	@echo "Running security checks..."
	gosec ./...

lint:
	@echo "Running lint..."
	golangci-lint run

test:
	@echo "Running tests..."
	go test ./...

clean:
	@echo "Cleaning build artefacts..."
	rm -rf target/
	go clean

migrate:
	@echo "Applying database migrations..."
	atlas migrate apply --env local

migrate-status:
	@echo "Migration status:"
	atlas migrate status --env local

migrate-new:
	@if [ -z "$(name)" ]; then \
		echo "Usage: make migrate-new name=migration_name"; \
		exit 1; \
	fi
	@echo "Creating new migration: $(name)"
	atlas migrate diff --env local $(name)

deps:
	@echo "Installing dependencies..."
	go mod tidy
	go mod download

dev:
	@echo "Starting development server with hot reload..."
	@if ! command -v air &> /dev/null; then \
		echo "Installing air for hot reload..."; \
		go install github.com/cosmtrek/air@latest; \
	fi
	air

docs:
	@echo "Generating API documentation..."
	@if ! command -v swag &> /dev/null; then \
		echo "Installing swag..."; \
		go install github.com/swaggo/swag/cmd/swag@latest; \
	fi
	swag init -d ./cmd/api,internal/*/http,pkg/middleware -o api/docs
	@echo "Documentation generated in api/docs/"

setup: deps docs
	@echo "Development environment setup complete"

install: build
	@echo "Installing dunhayat..."
	@if [ "$$(uname)" != "FreeBSD" ]; then \
		echo "Error: Install target only supports FreeBSD at the moment"; \
		exit 1; \
	fi
	@if [ "$$(id -u)" -eq 0 ]; then \
		./scripts/platform/FreeBSD/install.sh; \
	else \
		echo "Installation requires root privileges."; \
		exit 1; \
	fi
