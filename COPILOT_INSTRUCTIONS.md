# Copilot Instructions for Jumpstarter Monorepo Integration

## Overview

This document provides comprehensive instructions for integrating all Jumpstarter components into a single monorepo. The goal is to consolidate the distributed Jumpstarter ecosystem into `jumpstarter-mono` while maintaining functionality, improving developer experience, and enabling easier management.

## Current State Analysis

### Existing Jumpstarter Ecosystem (jumpstarter-dev organization)

The Jumpstarter project currently consists of 13 separate repositories:

#### Core Components
1. **jumpstarter** - Main Python library and CLI tools
   - Language: Python
   - Purpose: Core Jumpstarter functionality, drivers, CLI
   - Size: ~10MB, 64 open issues
   - Key files: pyproject.toml, packages/, examples/, docs/

2. **jumpstarter-controller** - Kubernetes controller
   - Language: Go
   - Purpose: Kubernetes operator for managing Jumpstarter resources
   - Size: ~1MB, 14 open issues
   - Key files: go.mod, Makefile, api/, cmd/, internal/

3. **jumpstarter-protocol** - Protocol definitions
   - Language: Makefile/Protocol buffers
   - Purpose: Shared protocol definitions between components
   - Size: ~173KB, 2 open issues

#### Hardware and Firmware
4. **dutlink-firmware** - Device firmware
   - Language: Rust
   - Purpose: Firmware for the DUTLink test harness board
   - Size: ~107KB, 1 open issue

5. **dutlink-board** - Hardware design
   - Purpose: Open hardware design files for test harness
   - Size: ~12MB

#### Supporting Tools and Templates
6. **jumpstarter-lab-config** - Lab configuration management
   - Language: Go
   - Purpose: Lab configuration and management tools
   - Size: ~207KB, 3 open issues

7. **jumpstarter-tekton-tasks** - CI/CD integration
   - Purpose: Tekton tasks and pipeline examples
   - Size: ~31KB, 3 open issues

8. **jumpstarter-driver-template** - Driver template
   - Language: Python
   - Purpose: Template for creating new drivers
   - Size: ~69KB
   - Type: Template repository

9. **jumpstarter-devspace** - Development environment
   - Purpose: DevSpaces configuration for development
   - Size: ~87KB
   - Type: Template repository

10. **jumpstarter-e2e** - End-to-end testing
    - Language: Shell
    - Purpose: End-to-end test suite
    - Size: ~43KB, 1 open issue

11. **vscode-jumpstarter** - VS Code extension
    - Language: TypeScript
    - Purpose: VS Code plugin for Jumpstarter development
    - Size: ~53KB

12. **packages** - Package repository
    - Language: Shell
    - Purpose: Python package repository generator for pkg.jumpstarter.dev
    - Size: ~16KB

13. **.github** - Organization templates
    - Purpose: GitHub organization configuration and templates
    - Size: ~13KB

## Target Monorepo Structure

```
jumpstarter-mono/
├── README.md
├── LICENSE
├── .gitignore
├── Makefile                    # Root build orchestration
├── pyproject.toml             # Python workspace configuration
├── go.work                    # Go workspace configuration
├── .github/                   # CI/CD workflows
│   └── workflows/
├── docs/                      # Consolidated documentation
├── tools/                     # Build and development tools
├── scripts/                   # Utility scripts
├── 
├── core/                      # Core Jumpstarter components
│   ├── jumpstarter/           # Main Python library (from jumpstarter)
│   ├── controller/            # Kubernetes controller (from jumpstarter-controller)
│   └── protocol/              # Protocol definitions (from jumpstarter-protocol)
│
├── hardware/                  # Hardware-related components
│   ├── dutlink-board/         # Hardware design files
│   └── dutlink-firmware/      # Rust firmware
│
├── packages/                  # Package management and distribution
│   ├── python/                # Python packages and wheels
│   ├── debian/                # Debian packages
│   ├── rpm/                   # RPM packages
│   └── container/             # Container images
│
├── integrations/              # CI/CD and tooling integrations
│   ├── tekton/                # Tekton tasks (from jumpstarter-tekton-tasks)
│   ├── vscode/                # VS Code extension (from vscode-jumpstarter)
│   └── devspace/              # Development environments (from jumpstarter-devspace)
│
├── templates/                 # Templates and scaffolding
│   └── driver/                # Driver template (from jumpstarter-driver-template)
│
├── testing/                   # Testing infrastructure
│   ├── e2e/                   # End-to-end tests (from jumpstarter-e2e)
│   ├── integration/           # Integration tests
│   └── fixtures/              # Test fixtures and data
│
└── lab-config/                # Lab configuration tools
    └── src/                   # Lab config source (from jumpstarter-lab-config)
```

## Integration Strategy

### Phase 1: Repository Setup and Core Integration

#### 1.1 Initialize Monorepo Structure
```bash
# Create directory structure
mkdir -p {core,hardware,packages,integrations,templates,testing,lab-config}
mkdir -p {docs,tools,scripts}
mkdir -p packages/{python,debian,rpm,container}
mkdir -p integrations/{tekton,vscode,devspace}
mkdir -p testing/{e2e,integration,fixtures}
mkdir -p templates/driver
```

#### 1.2 Setup Multi-language Workspace Configuration

**Python Workspace (pyproject.toml)**
```toml
[build-system]
requires = ["hatchling"]
build-backend = "hatchling.build"

[project]
name = "jumpstarter-mono"
description = "Jumpstarter monorepo containing all components"
readme = "README.md"
license = "Apache-2.0"
authors = [
    {name = "Jumpstarter Contributors"},
]

[tool.hatch.build.targets.wheel]
packages = ["core/jumpstarter/packages"]

[tool.uv.workspace]
members = [
    "core/jumpstarter",
    "templates/driver",
    "testing/e2e",
]
```

**Go Workspace (go.work)**
```go
go 1.22

use (
    ./core/controller
    ./lab-config
)
```

**Root Makefile**
```makefile
.PHONY: help build test clean lint fmt install

help: ## Show this help
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

build: build-python build-go build-rust build-web ## Build all components

build-python: ## Build Python components
	cd core/jumpstarter && uv build
	cd templates/driver && uv build

build-go: ## Build Go components
	cd core/controller && make build
	cd lab-config && go build ./...

build-rust: ## Build Rust components
	cd hardware/dutlink-firmware && cargo build --release

build-web: ## Build web components
	cd integrations/vscode && npm install && npm run compile

test: test-python test-go test-rust test-e2e ## Run all tests

test-python: ## Run Python tests
	cd core/jumpstarter && uv run pytest
	cd templates/driver && uv run pytest

test-go: ## Run Go tests
	cd core/controller && make test
	cd lab-config && go test ./...

test-rust: ## Run Rust tests
	cd hardware/dutlink-firmware && cargo test

test-e2e: ## Run end-to-end tests
	cd testing/e2e && ./run-tests.sh

lint: lint-python lint-go lint-rust ## Run all linters

clean: ## Clean build artifacts
	find . -name "*.pyc" -delete
	find . -name "__pycache__" -type d -exec rm -rf {} +
	cd core/controller && make clean
	cd hardware/dutlink-firmware && cargo clean
	cd integrations/vscode && rm -rf node_modules

install: ## Install all components
	pip install -e core/jumpstarter
	cd core/controller && make install
```

### Phase 2: Component Migration

#### 2.1 Core Components Migration

**Migrate jumpstarter (Main Python Library)**
```bash
# Clone and move main jumpstarter repository
git clone https://github.com/jumpstarter-dev/jumpstarter.git temp-jumpstarter
mv temp-jumpstarter/* core/jumpstarter/
rm -rf temp-jumpstarter

# Update imports and paths in core/jumpstarter/
# Preserve git history using git subtree or git filter-branch if needed
```

**Migrate jumpstarter-controller (Go Kubernetes Controller)**
```bash
# Clone and move controller
git clone https://github.com/jumpstarter-dev/jumpstarter-controller.git temp-controller
mv temp-controller/* core/controller/
rm -rf temp-controller

# Update go.mod paths and module names
cd core/controller
go mod edit -module github.com/the78mole/jumpstarter-mono/core/controller
```

**Migrate jumpstarter-protocol**
```bash
# Clone and move protocol definitions
git clone https://github.com/jumpstarter-dev/jumpstarter-protocol.git temp-protocol
mv temp-protocol/* core/protocol/
rm -rf temp-protocol

# Update references in other components
```

#### 2.2 Hardware Components Migration

**Migrate dutlink-firmware (Rust)**
```bash
# Clone and move firmware
git clone https://github.com/jumpstarter-dev/dutlink-firmware.git temp-firmware
mv temp-firmware/* hardware/dutlink-firmware/
rm -rf temp-firmware

# Update Cargo.toml if needed
```

**Migrate dutlink-board (Hardware)**
```bash
# Clone and move hardware design files
git clone https://github.com/jumpstarter-dev/dutlink-board.git temp-board
mv temp-board/* hardware/dutlink-board/
rm -rf temp-board
```

#### 2.3 Supporting Components Migration

**Migrate remaining components following similar patterns:**
- jumpstarter-tekton-tasks → integrations/tekton/
- vscode-jumpstarter → integrations/vscode/
- jumpstarter-devspace → integrations/devspace/
- jumpstarter-driver-template → templates/driver/
- jumpstarter-e2e → testing/e2e/
- jumpstarter-lab-config → lab-config/
- packages → packages/

### Phase 3: Build System Integration

#### 3.1 Python Build Integration
- Consolidate all Python packages under UV workspace
- Update pyproject.toml files to reference monorepo structure
- Maintain separate package builds but unified development workflow
- Update import paths and dependencies

#### 3.2 Go Build Integration
- Setup Go workspace with go.work
- Update module paths to reference monorepo
- Consolidate Go tooling and linting configuration
- Maintain separate Go modules for different components

#### 3.3 Multi-language Tooling
- Setup pre-commit hooks for all languages
- Integrate formatters: black (Python), gofmt (Go), rustfmt (Rust)
- Setup linters: ruff (Python), golangci-lint (Go), clippy (Rust)
- Unified testing orchestration via root Makefile

### Phase 4: CI/CD Integration

#### 4.1 GitHub Actions Workflows

**.github/workflows/ci.yml**
```yaml
name: CI

on:
  push:
    branches: [ main ]
  pull_request:
    branches: [ main ]

jobs:
  python:
    runs-on: ubuntu-latest
    steps:
    - uses: actions/checkout@v4
    - uses: astral-sh/setup-uv@v3
    - name: Install dependencies
      run: uv sync
    - name: Run tests
      run: make test-python
    - name: Run linting
      run: make lint-python

  go:
    runs-on: ubuntu-latest
    steps:
    - uses: actions/checkout@v4
    - uses: actions/setup-go@v4
      with:
        go-version: '1.22'
    - name: Run tests
      run: make test-go
    - name: Run linting
      run: make lint-go

  rust:
    runs-on: ubuntu-latest
    steps:
    - uses: actions/checkout@v4
    - uses: actions-rs/toolchain@v1
      with:
        toolchain: stable
    - name: Run tests
      run: make test-rust

  e2e:
    runs-on: ubuntu-latest
    needs: [python, go]
    steps:
    - uses: actions/checkout@v4
    - name: Setup test environment
      run: |
        # Setup kind cluster or other test infrastructure
    - name: Run E2E tests
      run: make test-e2e
```

#### 4.2 Release Strategy
- Unified versioning strategy across all components
- Automated release workflows
- Multi-platform package building
- Container image building and publishing

### Phase 5: Documentation Consolidation

#### 5.1 Documentation Structure
```
docs/
├── README.md                  # Main documentation
├── installation/
│   ├── from-source.md
│   ├── containers.md
│   └── packages.md
├── development/
│   ├── getting-started.md
│   ├── building.md
│   ├── testing.md
│   └── contributing.md
├── architecture/
│   ├── overview.md
│   ├── components.md
│   └── protocols.md
├── user-guide/
│   ├── cli.md
│   ├── drivers.md
│   └── kubernetes.md
├── hardware/
│   ├── dutlink-board.md
│   └── firmware.md
└── integrations/
    ├── tekton.md
    ├── vscode.md
    └── devspace.md
```

#### 5.2 Documentation Generation
- Setup unified documentation build with MkDocs or Sphinx
- Auto-generate API documentation
- Include hardware documentation and schematics
- Consolidate examples and tutorials

### Phase 6: Migration Validation

#### 6.1 Functionality Testing
- Verify all components build successfully
- Run existing test suites
- Validate cross-component integrations
- Test package generation and installation

#### 6.2 Performance Validation
- Compare build times before/after migration
- Validate container image sizes
- Test development workflow efficiency
- Measure CI/CD pipeline performance

## Implementation Checklist

### Pre-Migration Setup
- [ ] Create monorepo directory structure
- [ ] Setup multi-language workspace configuration
- [ ] Create root Makefile for build orchestration
- [ ] Setup initial CI/CD workflows
- [ ] Create documentation framework

### Core Component Migration
- [ ] Migrate jumpstarter main library
- [ ] Migrate jumpstarter-controller
- [ ] Migrate jumpstarter-protocol
- [ ] Update cross-component dependencies
- [ ] Validate core functionality

### Hardware Component Migration
- [ ] Migrate dutlink-firmware
- [ ] Migrate dutlink-board
- [ ] Validate firmware build process
- [ ] Update hardware documentation

### Supporting Component Migration
- [ ] Migrate Tekton tasks and examples
- [ ] Migrate VS Code extension
- [ ] Migrate DevSpace configuration
- [ ] Migrate driver template
- [ ] Migrate E2E test suite
- [ ] Migrate lab configuration tools
- [ ] Migrate package repository tools

### Build System Integration
- [ ] Consolidate Python packages under UV workspace
- [ ] Setup Go workspace configuration
- [ ] Integrate Rust build into main workflow
- [ ] Setup TypeScript/Node.js build for VS Code extension
- [ ] Create unified linting and formatting
- [ ] Setup pre-commit hooks

### CI/CD Integration
- [ ] Create multi-language CI pipeline
- [ ] Setup automated testing for all components
- [ ] Create release automation
- [ ] Setup package publishing
- [ ] Setup container image building

### Documentation and Cleanup
- [ ] Consolidate all documentation
- [ ] Update README files
- [ ] Create migration guide
- [ ] Archive old repositories
- [ ] Update external references

### Validation and Testing
- [ ] Full integration test suite
- [ ] Performance benchmarking
- [ ] Developer workflow testing
- [ ] Documentation review
- [ ] Community feedback collection

## Risk Mitigation

### Identified Risks
1. **Build complexity**: Managing multiple languages and build systems
2. **Repository size**: Large repository with binary hardware files
3. **Development workflow**: Potential slowdown for developers
4. **CI/CD performance**: Longer build and test times
5. **Dependency conflicts**: Cross-component dependency management

### Mitigation Strategies
1. **Modular builds**: Only build/test changed components
2. **Git LFS**: Use Git LFS for large binary files
3. **Sparse checkout**: Enable developers to work on subsets
4. **Caching**: Aggressive build and dependency caching
5. **Dependency isolation**: Clear dependency boundaries between components

## Success Criteria

### Technical Success Criteria
- [ ] All components build successfully in monorepo
- [ ] All existing tests pass
- [ ] CI/CD pipeline completes in reasonable time (<30 minutes)
- [ ] Developer setup time reduced compared to multi-repo setup
- [ ] Documentation is consolidated and accessible

### Process Success Criteria
- [ ] Migration completed without data loss
- [ ] Existing workflows continue to function
- [ ] Community can continue contributing
- [ ] Release process is simplified
- [ ] Maintenance overhead is reduced

## Timeline Estimate

- **Week 1-2**: Monorepo setup and core component migration
- **Week 3**: Hardware component migration and build integration
- **Week 4**: Supporting component migration
- **Week 5**: CI/CD integration and documentation
- **Week 6**: Testing, validation, and refinement

## Conclusion

This migration will consolidate the Jumpstarter ecosystem into a manageable monorepo while preserving functionality and improving developer experience. The phased approach allows for validation at each step and reduces migration risks.

The resulting monorepo will provide:
- Unified development workflow
- Simplified dependency management
- Consolidated documentation
- Streamlined release process
- Better cross-component integration
- Reduced maintenance overhead

Next steps involve executing this plan systematically, starting with the repository structure setup and core component migration.