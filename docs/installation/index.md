# Installation Guide

This guide will help you install and set up the Jumpstarter monorepo for development or production use.

## Prerequisites

### System Requirements

- Python 3.12 or later
- Go 1.22 or later
- Rust (latest stable)
- Node.js 18 or later
- Git

### Tools

- [uv](https://github.com/astral-sh/uv) - Python package manager
- Make - Build orchestration
- Docker (optional) - For containerized development

## Quick Setup

1. Clone the repository:
   ```bash
   git clone https://github.com/the78mole/jumpstarter-mono.git
   cd jumpstarter-mono
   ```

2. Run the setup command:
   ```bash
   make setup
   ```

This will install all dependencies for all components in the monorepo.

## Component-specific Setup

### Python Components

Python components use `uv` for dependency management:

```bash
make build-python
```

### Go Components

Go components are managed with Go modules:

```bash
make build-go
```

### Rust Components

Rust components use Cargo:

```bash
make build-rust
```

### Web Components

TypeScript/Node.js components use npm:

```bash
make build-web
```

## Verification

To verify your installation:

```bash
make test
```

## Next Steps

- [Development Guide](../development/index.md) - Learn about the development workflow
- [Architecture](../architecture/index.md) - Understand the system design