# Makefile for ping-uptime project

.PHONY: all dev dev-backend dev-frontend build build-frontend build-backend clean test help

# Default target
all: help

help: ## Show this help message
	@echo "Usage: make [target]"
	@echo ""
	@echo "Targets:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-16s\033[0m %s\n", $$1, $$2}'

dev: ## Run backend (using air) and frontend (vite dev) concurrently
	@echo "Starting development servers..."
	@if command -v npm >/dev/null 2>&1; then \
		(cd web && npm run dev) & \
	fi; \
	if command -v air >/dev/null 2>&1; then \
		air; \
	else \
		go run main.go; \
	fi

dev-backend: ## Run backend with live-reloads (air)
	@echo "Starting backend with Air..."
	@if command -v air >/dev/null 2>&1; then \
		air; \
	else \
		go run main.go; \
	fi

dev-frontend: ## Run frontend dev server (vite)
	@echo "Starting frontend dev server..."
	@cd web && npm run dev

build: ## Build both frontend and backend for production
	@$(MAKE) build-frontend
	@$(MAKE) build-backend

build-frontend: ## Compile Vue 3 frontend assets into static folder
	@echo "Building frontend..."
	@cd web && npm install && npm run build

build-backend: ## Compile Go backend binary with embedded assets
	@echo "Building backend..."
	@mkdir -p bin
	@go build -o bin/ping-uptime main.go

clean: ## Clean build outputs and temporary directories
	@echo "Cleaning build artifacts..."
	@rm -rf tmp bin static/assets static/index.html static/favicon.ico static/logo.png

test: ## Run backend and frontend tests
	@echo "Running tests..."
	@go test -v ./...
	@if [ -d "web" ] && [ -f "web/package.json" ]; then \
		cd web && npm run test:unit || true; \
	fi
