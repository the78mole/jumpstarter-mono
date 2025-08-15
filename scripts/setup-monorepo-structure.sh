#!/bin/bash
# setup-monorepo-structure.sh
# Script to setup the monorepo directory structure and initial configuration

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

# Ensure we're in the repository root
cd "$(dirname "$0")/.."

log "Setting up Jumpstarter monorepo structure..."

# Install UV if not present
if ! command -v uv &> /dev/null; then
    log "Installing UV Python package manager..."
    curl -LsSf https://astral.sh/uv/install.sh | sh
    export PATH="$HOME/.local/bin:$PATH"
fi

# Install pre-commit if not present
if ! command -v pre-commit &> /dev/null; then
    log "Installing pre-commit..."
    if command -v uv &> /dev/null; then
        uv tool install pre-commit
    else
        pip install pre-commit
    fi
fi

# Create .pre-commit-config.yaml if it doesn't exist
if [ ! -f ".pre-commit-config.yaml" ]; then
    log "Creating pre-commit configuration..."
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

  - repo: https://github.com/doublify/pre-commit-rust
    rev: v1.0
    hooks:
      - id: fmt
      - id: clippy
EOF
fi

# Setup Python workspace
if [ -f "pyproject.toml" ]; then
    log "Syncing Python workspace..."
    uv sync
fi

# Setup Go workspace if Go modules exist
if ls core/*/go.mod lab-config/go.mod 2>/dev/null; then
    log "Setting up Go workspace..."
    go work init 2>/dev/null || true
    go work use ./core/controller 2>/dev/null || true
    go work use ./lab-config 2>/dev/null || true
fi

# Install pre-commit hooks
log "Installing pre-commit hooks..."
pre-commit install

log "✓ Monorepo structure setup completed successfully!"
log "Next step: Run scripts/migrate-components.sh to populate directories with component code"
