# Makefile for ping-uptime project

.PHONY: all dev dev-backend dev-frontend build build-frontend build-backend clean test help

# Default target
all: help

## help: Show this help message
help:
	@echo "Usage: make [target]"
	@echo ""
	@echo "Targets:"
	@fgrep -h "##" $(MAKEFILE_LIST) | fgrep -v fgrep | sed -e 's/\\$$//' | sed -e 's/##//'

## dev: Run backend (using air) and frontend (vite dev) concurrently
dev:
	@echo "Starting development servers..."
	@if command -v npm >/dev/null 2>&1; then \
		(cd web && npm run dev) & \
	fi
	@if command -v air >/dev/null 2>&1; then \
		air; \
	else \
		go run main.go; \
	fi

## dev-backend: Run backend with live-reloads (air)
dev-backend:
	@echo "Starting backend with Air..."
	@if command -v air >/dev/null 2>&1; then \
		air; \
	else \
		go run main.go; \
	fi

## dev-frontend: Run frontend dev server (vite)
dev-frontend:
	@echo "Starting frontend dev server..."
	@cd web && npm run dev

## build: Build both frontend and backend for production
build: build-frontend build-backend

## build-frontend: Compile Vue 3 frontend assets into static folder
build-frontend:
	@echo "Building frontend..."
	@cd web && npm install && npm run build

## build-backend: Compile Go backend binary with embedded assets
build-backend:
	@echo "Building backend..."
	@go build -o bin/ping-uptime main.go

## clean: Clean build outputs and temporary directories
clean:
	@echo "Cleaning build artifacts..."
	@rm -rf tmp bin static/assets static/index.html static/favicon.ico static/logo.png

## test: Run backend and frontend tests
test:
	@echo "Running tests..."
	@go test -v ./...
	@cd web && npm run test:unit
