# Jumpstarter Monorepo Implementation Guide

## Quick Start Commands

This document provides the specific commands and scripts to implement the monorepo migration outlined in `COPILOT_INSTRUCTIONS.md`.

## Phase 1: Repository Setup

### 1.1 Create Directory Structure

```bash
#!/bin/bash
# setup-monorepo-structure.sh

# Create main directory structure
mkdir -p core/{jumpstarter,controller,protocol}
mkdir -p hardware/{dutlink-board,dutlink-firmware}
mkdir -p packages/{python,debian,rpm,container}
mkdir -p integrations/{tekton,vscode,devspace}
mkdir -p templates/driver
mkdir -p testing/{e2e,integration,fixtures}
mkdir -p lab-config/src
mkdir -p docs/{installation,development,architecture,user-guide,hardware,integrations}
mkdir -p tools
mkdir -p scripts

echo "Monorepo directory structure created successfully!"
```

### 1.2 Create Workspace Configuration Files

**pyproject.toml (Root)**
```toml
[build-system]
requires = ["hatchling"]
build-backend = "hatchling.build"

[project]
name = "jumpstarter-mono"
version = "0.1.0"
description = "Jumpstarter monorepo containing all components"
readme = "README.md"
license = "Apache-2.0"
authors = [
    {name = "Jumpstarter Contributors"},
]
classifiers = [
    "Development Status :: 4 - Beta",
    "Intended Audience :: Developers",
    "License :: OSI Approved :: Apache Software License",
    "Programming Language :: Python :: 3",
    "Programming Language :: Python :: 3.12",
]
dependencies = []

[tool.uv.workspace]
members = [
    "core/jumpstarter",
    "templates/driver",
]

[tool.uv.sources]

[tool.ruff]
target-version = "py312"
line-length = 88

[tool.ruff.lint]
select = ["E", "F", "I", "N", "W"]
ignore = []

[tool.black]
line-length = 88
target-version = ['py312']

[tool.pytest.ini_options]
testpaths = ["core/jumpstarter/tests", "templates/driver/tests"]
python_files = "test_*.py"
python_classes = "Test*"
python_functions = "test_*"
```

**go.work (Root)**
```go
go 1.22

use (
	./core/controller
	./lab-config
)
```

**Makefile (Root)**
```makefile
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
```

## Phase 2: Migration Scripts

### 2.1 Component Migration Script

```bash
#!/bin/bash
# migrate-components.sh

set -e

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

log() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

warn() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Configuration
TEMP_DIR="./temp-migration"
GITHUB_ORG="jumpstarter-dev"

# Component mappings: "source_repo:target_path"
declare -A COMPONENTS=(
    ["jumpstarter"]="core/jumpstarter"
    ["jumpstarter-controller"]="core/controller"
    ["jumpstarter-protocol"]="core/protocol"
    ["dutlink-firmware"]="hardware/dutlink-firmware"
    ["dutlink-board"]="hardware/dutlink-board"
    ["jumpstarter-tekton-tasks"]="integrations/tekton"
    ["vscode-jumpstarter"]="integrations/vscode"
    ["jumpstarter-devspace"]="integrations/devspace"
    ["jumpstarter-driver-template"]="templates/driver"
    ["jumpstarter-e2e"]="testing/e2e"
    ["jumpstarter-lab-config"]="lab-config"
    ["packages"]="packages/repository-tools"
)

migrate_component() {
    local repo_name=$1
    local target_path=$2
    
    log "Migrating ${repo_name} to ${target_path}..."
    
    # Create temporary directory
    mkdir -p "${TEMP_DIR}"
    
    # Clone repository
    log "Cloning ${GITHUB_ORG}/${repo_name}..."
    git clone "https://github.com/${GITHUB_ORG}/${repo_name}.git" "${TEMP_DIR}/${repo_name}"
    
    # Create target directory
    mkdir -p "${target_path}"
    
    # Copy files (excluding .git)
    log "Copying files to ${target_path}..."
    rsync -av --exclude='.git' "${TEMP_DIR}/${repo_name}/" "${target_path}/"
    
    # Clean up temporary directory
    rm -rf "${TEMP_DIR}/${repo_name}"
    
    log "✓ Successfully migrated ${repo_name}"
}

update_go_modules() {
    log "Updating Go module paths..."
    
    # Update controller module
    if [ -f "core/controller/go.mod" ]; then
        cd core/controller
        go mod edit -module github.com/the78mole/jumpstarter-mono/core/controller
        go mod tidy
        cd ../..
    fi
    
    # Update lab-config module
    if [ -f "lab-config/go.mod" ]; then
        cd lab-config
        go mod edit -module github.com/the78mole/jumpstarter-mono/lab-config
        go mod tidy
        cd ..
    fi
    
    log "✓ Go modules updated"
}

update_python_configs() {
    log "Updating Python configurations..."
    
    # Update jumpstarter pyproject.toml
    if [ -f "core/jumpstarter/pyproject.toml" ]; then
        sed -i 's|name = "jumpstarter"|name = "jumpstarter-core"|g' core/jumpstarter/pyproject.toml
    fi
    
    # Update driver template pyproject.toml
    if [ -f "templates/driver/pyproject.toml" ]; then
        sed -i 's|name = .*|name = "jumpstarter-driver-template"|g' templates/driver/pyproject.toml
    fi
    
    log "✓ Python configurations updated"
}

update_imports() {
    log "Updating import paths..."
    
    # Update Python imports (basic pattern matching)
    find core/jumpstarter -name "*.py" -type f -exec sed -i 's|from jumpstarter|from jumpstarter_core|g' {} \;
    find templates/driver -name "*.py" -type f -exec sed -i 's|from jumpstarter|from jumpstarter_core|g' {} \;
    
    # Update Go imports
    find core/controller -name "*.go" -type f -exec sed -i 's|github.com/jumpstarter-dev/jumpstarter-controller|github.com/the78mole/jumpstarter-mono/core/controller|g' {} \;
    find lab-config -name "*.go" -type f -exec sed -i 's|github.com/jumpstarter-dev/jumpstarter-lab-config|github.com/the78mole/jumpstarter-mono/lab-config|g' {} \;
    
    log "✓ Import paths updated"
}

# Main migration process
main() {
    log "Starting Jumpstarter monorepo migration..."
    
    # Migrate all components
    for repo in "${!COMPONENTS[@]}"; do
        migrate_component "$repo" "${COMPONENTS[$repo]}"
    done
    
    # Update configurations
    update_go_modules
    update_python_configs
    update_imports
    
    # Clean up
    rm -rf "${TEMP_DIR}"
    
    log "✓ Migration completed successfully!"
    warn "Please review the migrated files and update any remaining references manually."
    warn "Don't forget to test builds: make build"
}

# Run migration
main "$@"
```

### 2.2 Build System Setup

```bash
#!/bin/bash
# setup-build-system.sh

set -e

log() {
    echo -e "\033[32m[INFO]\033[0m $1"
}

setup_python_workspace() {
    log "Setting up Python workspace..."
    
    # Install UV if not present
    if ! command -v uv &> /dev/null; then
        curl -LsSf https://astral.sh/uv/install.sh | sh
        source $HOME/.local/bin/env
    fi
    
    # Sync workspace
    uv sync
    
    log "✓ Python workspace configured"
}

setup_go_workspace() {
    log "Setting up Go workspace..."
    
    # Initialize Go workspace
    go work init
    go work use ./core/controller
    go work use ./lab-config
    
    # Download dependencies
    cd core/controller && go mod download && cd ../..
    cd lab-config && go mod download && cd ..
    
    log "✓ Go workspace configured"
}

setup_rust_environment() {
    log "Setting up Rust environment..."
    
    # Check if Rust is installed
    if ! command -v cargo &> /dev/null; then
        curl --proto '=https' --tlsv1.2 -sSf https://sh.rustup.rs | sh -s -- -y
        source $HOME/.cargo/env
    fi
    
    # Fetch dependencies
    cd hardware/dutlink-firmware && cargo fetch && cd ../..
    
    log "✓ Rust environment configured"
}

setup_node_environment() {
    log "Setting up Node.js environment..."
    
    # Setup VS Code extension
    if [ -d "integrations/vscode" ]; then
        cd integrations/vscode
        npm install
        cd ../..
    fi
    
    log "✓ Node.js environment configured"
}

setup_pre_commit() {
    log "Setting up pre-commit hooks..."
    
    # Install pre-commit if not present
    if ! command -v pre-commit &> /dev/null; then
        uv tool install pre-commit
    fi
    
    # Create pre-commit config
    cat > .pre-commit-config.yaml << 'EOF'
repos:
  - repo: https://github.com/pre-commit/pre-commit-hooks
    rev: v4.5.0
    hooks:
      - id: trailing-whitespace
      - id: end-of-file-fixer
      - id: check-yaml
      - id: check-added-large-files

  - repo: https://github.com/psf/black
    rev: 23.9.1
    hooks:
      - id: black
        language_version: python3.12

  - repo: https://github.com/astral-sh/ruff-pre-commit
    rev: v0.1.3
    hooks:
      - id: ruff
        args: [--fix]

  - repo: https://github.com/golangci/golangci-lint
    rev: v1.54.2
    hooks:
      - id: golangci-lint
        args: [--timeout=5m]
        files: \.go$

  - repo: https://github.com/domlysz/BlenderGIS
    rev: master
    hooks:
      - id: rustfmt
        files: \.rs$
EOF
    
    # Install hooks
    pre-commit install
    
    log "✓ Pre-commit hooks configured"
}

# Main setup
main() {
    log "Setting up build system for Jumpstarter monorepo..."
    
    setup_python_workspace
    setup_go_workspace
    setup_rust_environment
    setup_node_environment
    setup_pre_commit
    
    log "✓ Build system setup completed!"
    log "Run 'make help' to see available commands"
}

main "$@"
```

## Phase 3: CI/CD Setup

### 3.1 GitHub Actions Workflow

```yaml
# .github/workflows/ci.yml
name: CI

on:
  push:
    branches: [ main, develop ]
  pull_request:
    branches: [ main ]

env:
  CARGO_TERM_COLOR: always

jobs:
  detect-changes:
    runs-on: ubuntu-latest
    outputs:
      python: ${{ steps.changes.outputs.python }}
      go: ${{ steps.changes.outputs.go }}
      rust: ${{ steps.changes.outputs.rust }}
      web: ${{ steps.changes.outputs.web }}
      docs: ${{ steps.changes.outputs.docs }}
    steps:
      - uses: actions/checkout@v4
      - uses: dorny/paths-filter@v2
        id: changes
        with:
          filters: |
            python:
              - 'core/jumpstarter/**'
              - 'templates/driver/**'
              - 'pyproject.toml'
              - 'uv.lock'
            go:
              - 'core/controller/**'
              - 'lab-config/**'
              - 'go.work'
              - 'go.work.sum'
            rust:
              - 'hardware/dutlink-firmware/**'
            web:
              - 'integrations/vscode/**'
            docs:
              - 'docs/**'
              - '*.md'

  python:
    needs: detect-changes
    if: needs.detect-changes.outputs.python == 'true'
    runs-on: ubuntu-latest
    strategy:
      matrix:
        python-version: [3.12]
    steps:
      - uses: actions/checkout@v4
      - name: Set up Python
        uses: actions/setup-python@v4
        with:
          python-version: ${{ matrix.python-version }}
      - name: Install UV
        uses: astral-sh/setup-uv@v3
      - name: Install dependencies
        run: uv sync
      - name: Run linting
        run: |
          uv run ruff check .
          uv run black --check .
      - name: Run tests
        run: uv run pytest
      - name: Build packages
        run: uv build

  go:
    needs: detect-changes
    if: needs.detect-changes.outputs.go == 'true'
    runs-on: ubuntu-latest
    strategy:
      matrix:
        go-version: [1.22]
    steps:
      - uses: actions/checkout@v4
      - name: Set up Go
        uses: actions/setup-go@v4
        with:
          go-version: ${{ matrix.go-version }}
      - name: Install dependencies
        run: |
          cd core/controller && go mod download
          cd ../../lab-config && go mod download
      - name: Run linting
        run: |
          go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
          cd core/controller && golangci-lint run
          cd ../../lab-config && golangci-lint run
      - name: Run tests
        run: |
          cd core/controller && make test
          cd ../../lab-config && go test ./...
      - name: Build
        run: |
          cd core/controller && make build
          cd ../../lab-config && go build ./...

  rust:
    needs: detect-changes
    if: needs.detect-changes.outputs.rust == 'true'
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - name: Set up Rust
        uses: actions-rs/toolchain@v1
        with:
          toolchain: stable
          components: rustfmt, clippy
      - name: Cache cargo registry
        uses: actions/cache@v3
        with:
          path: ~/.cargo/registry
          key: ${{ runner.os }}-cargo-registry-${{ hashFiles('**/Cargo.lock') }}
      - name: Cache cargo index
        uses: actions/cache@v3
        with:
          path: ~/.cargo/git
          key: ${{ runner.os }}-cargo-index-${{ hashFiles('**/Cargo.lock') }}
      - name: Cache cargo build
        uses: actions/cache@v3
        with:
          path: hardware/dutlink-firmware/target
          key: ${{ runner.os }}-cargo-build-target-${{ hashFiles('**/Cargo.lock') }}
      - name: Check formatting
        run: cd hardware/dutlink-firmware && cargo fmt --all -- --check
      - name: Run clippy
        run: cd hardware/dutlink-firmware && cargo clippy -- -D warnings
      - name: Run tests
        run: cd hardware/dutlink-firmware && cargo test
      - name: Build release
        run: cd hardware/dutlink-firmware && cargo build --release

  web:
    needs: detect-changes
    if: needs.detect-changes.outputs.web == 'true'
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - name: Set up Node.js
        uses: actions/setup-node@v3
        with:
          node-version: '18'
          cache: 'npm'
          cache-dependency-path: integrations/vscode/package-lock.json
      - name: Install dependencies
        run: cd integrations/vscode && npm ci
      - name: Run linting
        run: cd integrations/vscode && npm run lint
      - name: Run tests
        run: cd integrations/vscode && npm test
      - name: Build
        run: cd integrations/vscode && npm run compile

  e2e:
    needs: [python, go]
    if: always() && (needs.python.result == 'success' || needs.python.result == 'skipped') && (needs.go.result == 'success' || needs.go.result == 'skipped')
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - name: Set up Python
        uses: actions/setup-python@v4
        with:
          python-version: 3.12
      - name: Set up Go
        uses: actions/setup-go@v4
        with:
          go-version: 1.22
      - name: Install UV
        uses: astral-sh/setup-uv@v3
      - name: Setup test environment
        run: |
          # Install kind for Kubernetes testing
          curl -Lo ./kind https://kind.sigs.k8s.io/dl/v0.20.0/kind-linux-amd64
          chmod +x ./kind
          sudo mv ./kind /usr/local/bin/kind
          
          # Create test cluster
          kind create cluster --name jumpstarter-test
      - name: Install dependencies
        run: |
          uv sync
          cd core/controller && go mod download
      - name: Build components
        run: |
          cd core/controller && make build
          uv build
      - name: Run E2E tests
        run: cd testing/e2e && bash run-tests.sh
      - name: Cleanup
        if: always()
        run: kind delete cluster --name jumpstarter-test

  docs:
    needs: detect-changes
    if: needs.detect-changes.outputs.docs == 'true'
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - name: Set up Python
        uses: actions/setup-python@v4
        with:
          python-version: 3.12
      - name: Install dependencies
        run: |
          pip install mkdocs mkdocs-material
      - name: Build docs
        run: cd docs && mkdocs build --strict
      - name: Deploy docs (main branch only)
        if: github.ref == 'refs/heads/main'
        run: cd docs && mkdocs gh-deploy --force
```

## Usage Instructions

1. **Run the setup script**:
   ```bash
   chmod +x setup-monorepo-structure.sh
   ./setup-monorepo-structure.sh
   ```

2. **Migrate components**:
   ```bash
   chmod +x migrate-components.sh
   ./migrate-components.sh
   ```

3. **Setup build system**:
   ```bash
   chmod +x setup-build-system.sh
   ./setup-build-system.sh
   ```

4. **Test the setup**:
   ```bash
   make help
   make build
   make test
   ```

This implementation guide provides the concrete steps and scripts needed to execute the monorepo migration plan outlined in the main instructions document.