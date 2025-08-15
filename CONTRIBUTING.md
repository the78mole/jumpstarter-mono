# Contributing to Jumpstarter

Thank you for your interest in contributing to Jumpstarter! We are an open community and welcome contributions of all kinds.

## Monorepo Structure

This is a monorepo containing all Jumpstarter components:

- `core/jumpstarter/` - Core Python packages and drivers
- `core/controller/` - Kubernetes controller (Go)
- `core/protocol/` - Protocol definitions
- `hardware/` - Hardware designs and firmware
- `integrations/` - Third-party integrations (Tekton, VSCode, DevSpace)
- `lab-config/` - Configuration management (Go)
- `packages/` - Repository tools and utilities
- `templates/` - Component templates
- `testing/` - End-to-end testing

## Getting Started

### Prerequisites

- Python 3.12+ with [uv](https://docs.astral.sh/uv/) package manager
- Go 1.21+ for controller and lab-config components
- Rust for firmware development (optional)
- Node.js and pnpm for VSCode extension development (optional)
- qemu & qemu-user

### Development Setup

> **Note**: These manual setup steps will be automated with a devcontainer configuration in the near future to streamline the development experience.

1. Clone the repository:

   ```bash
   git clone https://github.com/jumpstarter-dev/jumpstarter-mono.git
   cd jumpstarter-mono
   ```

2. Install pre-commit hooks:

   ```bash
   pip install pre-commit
   pre-commit install
   ```

3. Set up Python environment:

   ```bash
   uv sync
   ```

4. Set up Go workspace:
   ```bash
   go work sync
   ```

### Building Components

Use the unified Makefile for cross-component builds:

```bash
# Build all components
make build

# Build specific languages
make build-go
make build-python

# Run tests
make test

# Run linting
make lint
```

## Development Workflow

1. **Create a feature branch** from `main`
2. **Make your changes** following the coding standards
3. **Run tests and linting** locally
4. **Commit your changes** with conventional commit messages
5. **Submit a pull request** with a clear description

### Coding Standards

- **Python**: Follow PEP 8, use Ruff for linting and formatting
- **Go**: Follow Go conventions, use golangci-lint
- **Rust**: Follow Rust conventions, use clippy and rustfmt
- **Commit messages**: Use [Conventional Commits](https://www.conventionalcommits.org/) format

### Testing

- **Python**: Write pytest tests in `tests/` directories
- **Go**: Write Go tests alongside source code
- **Integration**: Use the `testing/e2e/` framework for end-to-end tests

## Component-Specific Guidelines

### Core Python Packages (core/jumpstarter/)

For detailed Python development guidelines, see [core/jumpstarter/docs/source/contributing/](core/jumpstarter/docs/source/contributing/).

### Controller (core/controller/)

Follow Kubernetes controller development best practices. The controller uses controller-runtime framework.

### Hardware (hardware/)

- **Firmware**: Written in Rust for embedded targets
- **PCB designs**: Use KiCad for hardware designs

## Documentation

- Use clear, concise language
- Include code examples where appropriate
- Update relevant documentation with your changes
- Documentation is built with MkDocs

## Getting Help

- Open an [issue](https://github.com/jumpstarter-dev/jumpstarter-mono/issues) for bugs or feature requests
- Join our community discussions
- Check existing documentation at [jumpstarter.dev](https://jumpstarter.dev)

## License

By contributing to Jumpstarter, you agree that your contributions will be licensed under the [Apache License 2.0](LICENSE).
