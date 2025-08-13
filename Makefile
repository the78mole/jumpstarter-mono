.DEFAULT_GOAL := help
.PHONY: help build test clean lint fmt install setup

# Colors
RED    := \033[31m
GREEN  := \033[32m
YELLOW := \033[33m
BLUE   := \033[34m
RESET  := \033[0m

help: ## Show this help message
	@echo "$(BLUE)Jumpstarter Monorepo Commands$(RESET)"
	@echo ""
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "$(GREEN)%-20s$(RESET) %s\n", $$1, $$2}' $(MAKEFILE_LIST)

# Setup
setup: ## Initial setup of development environment
	@echo "$(YELLOW)Setting up development environment...$(RESET)"
	uv sync
	cd core/controller && go mod download
	cd lab-config && go mod download
	cd hardware/dutlink-firmware && cargo fetch
	cd integrations/vscode && npm install
	@echo "$(GREEN)Setup completed!$(RESET)"

# Build targets
build: build-python build-go build-rust build-web ## Build all components
	@echo "$(GREEN)All components built successfully!$(RESET)"

build-python: ## Build Python components
	@echo "$(YELLOW)Building Python components...$(RESET)"
	uv build

build-go: ## Build Go components
	@echo "$(YELLOW)Building Go components...$(RESET)"
	cd core/controller && make build
	cd lab-config && go build ./...

build-rust: ## Build Rust components
	@echo "$(YELLOW)Building Rust components...$(RESET)"
	cd hardware/dutlink-firmware && cargo build --release

build-web: ## Build TypeScript/Node.js components
	@echo "$(YELLOW)Building web components...$(RESET)"
	cd integrations/vscode && npm run compile

# Test targets
test: test-python test-go test-rust test-e2e ## Run all tests
	@echo "$(GREEN)All tests completed!$(RESET)"

test-python: ## Run Python tests
	@echo "$(YELLOW)Running Python tests...$(RESET)"
	uv run pytest

test-go: ## Run Go tests
	@echo "$(YELLOW)Running Go tests...$(RESET)"
	cd core/controller && make test
	cd lab-config && go test ./...

test-rust: ## Run Rust tests
	@echo "$(YELLOW)Running Rust tests...$(RESET)"
	cd hardware/dutlink-firmware && cargo test

test-e2e: ## Run end-to-end tests
	@echo "$(YELLOW)Running end-to-end tests...$(RESET)"
	cd testing/e2e && bash run-tests.sh

# Lint targets
lint: lint-python lint-go lint-rust lint-web ## Run all linters
	@echo "$(GREEN)All linting completed!$(RESET)"

lint-python: ## Run Python linters
	@echo "$(YELLOW)Linting Python code...$(RESET)"
	uv run ruff check .
	uv run black --check .

lint-go: ## Run Go linters
	@echo "$(YELLOW)Linting Go code...$(RESET)"
	cd core/controller && golangci-lint run
	cd lab-config && golangci-lint run

lint-rust: ## Run Rust linters
	@echo "$(YELLOW)Linting Rust code...$(RESET)"
	cd hardware/dutlink-firmware && cargo clippy -- -D warnings

lint-web: ## Run web linters
	@echo "$(YELLOW)Linting TypeScript code...$(RESET)"
	cd integrations/vscode && npm run lint

# Format targets
fmt: fmt-python fmt-go fmt-rust ## Format all code
	@echo "$(GREEN)All code formatted!$(RESET)"

fmt-python: ## Format Python code
	uv run black .
	uv run ruff check --fix .

fmt-go: ## Format Go code
	cd core/controller && go fmt ./...
	cd lab-config && go fmt ./...

fmt-rust: ## Format Rust code
	cd hardware/dutlink-firmware && cargo fmt

# Clean targets
clean: clean-python clean-go clean-rust clean-web ## Clean all build artifacts
	@echo "$(GREEN)All build artifacts cleaned!$(RESET)"

clean-python: ## Clean Python build artifacts
	find . -name "*.pyc" -delete
	find . -name "__pycache__" -type d -exec rm -rf {} +
	find . -name "*.egg-info" -type d -exec rm -rf {} +
	rm -rf build/ dist/

clean-go: ## Clean Go build artifacts
	cd core/controller && make clean
	cd lab-config && go clean -cache

clean-rust: ## Clean Rust build artifacts
	cd hardware/dutlink-firmware && cargo clean

clean-web: ## Clean web build artifacts
	cd integrations/vscode && rm -rf node_modules out

# Install targets
install: ## Install all components for development
	@echo "$(YELLOW)Installing all components...$(RESET)"
	uv pip install -e .
	cd core/controller && make install
	@echo "$(GREEN)Installation completed!$(RESET)"

# Package targets
package: package-python package-debian package-rpm package-container ## Build all packages

package-python: ## Build Python packages
	uv build

package-debian: ## Build Debian packages
	@echo "$(YELLOW)Building Debian packages...$(RESET)"
	cd packages/debian && ./build.sh

package-rpm: ## Build RPM packages
	@echo "$(YELLOW)Building RPM packages...$(RESET)"
	cd packages/rpm && ./build.sh

package-container: ## Build container images
	@echo "$(YELLOW)Building container images...$(RESET)"
	cd packages/container && ./build.sh

# Documentation targets
docs: ## Build documentation
	@echo "$(YELLOW)Building documentation...$(RESET)"
	cd docs && mkdocs build

docs-serve: ## Serve documentation locally
	cd docs && mkdocs serve

# Development targets
dev: ## Start development environment
	@echo "$(YELLOW)Starting development environment...$(RESET)"
	@echo "Run 'make dev-python', 'make dev-go', or 'make dev-rust' for specific environments"

dev-python: ## Start Python development environment
	uv run python

dev-go: ## Start Go development environment
	cd core/controller && go run ./cmd/manager/main.go

dev-rust: ## Start Rust development environment
	cd hardware/dutlink-firmware && cargo run

# Pre-commit hooks
pre-commit: lint test ## Run pre-commit checks
	@echo "$(GREEN)Pre-commit checks passed!$(RESET)"

pre-commit-install: ## Install pre-commit hooks
	pre-commit install

# Git hooks
hooks: ## Setup git hooks
	cp tools/hooks/pre-commit .git/hooks/
	chmod +x .git/hooks/pre-commit