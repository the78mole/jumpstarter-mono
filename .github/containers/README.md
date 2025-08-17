# CI Container Images

This directory contains container definitions for optimized CI environments. These images provide pre-installed tools to speed up CI builds and ensure consistent environments.

## Available Images

### `ci-base`
Basic CI environment with common system dependencies:
- Ubuntu 24.04 base
- Git, make, build-essential
- USB development libraries (libudev-dev, libusb-1.0-0-dev)
- Non-root `ci` user

### `ci-python`
Python development environment:
- Everything from `ci-base`
- Python 3.12 with venv support
- UV package manager
- Pre-configured virtual environment

### `ci-go`
Go development environment:
- Everything from `ci-base`
- Go 1.22.9
- Pre-configured GOPATH and build cache directories

### `ci-rust`
Rust development environment:
- Everything from `ci-base`
- Rust stable toolchain
- thumbv7em-none-eabihf target (for embedded development)
- Cargo and rustc

### `ci-node`
Node.js/TypeScript development environment:
- Everything from `ci-base`
- Node.js 18 LTS
- npm and pnpm package managers

### `ci-multi`
All-in-one environment with all tools:
- Everything from all above images
- Suitable for workflows that need multiple toolchains
- Larger image size but maximum compatibility

## Usage

These images are automatically built and published to `ghcr.io/the78mole/jumpstarter-mono/ci-*:latest` when:
- Container definitions change
- CI configuration changes
- Manually triggered via workflow_dispatch

The reusable workflow templates automatically use these containers when available, falling back to tool installation if needed (e.g., when running with `act` locally).

## Building Locally

To build and test a container locally:

```bash
# Build a specific container
docker build -f .github/containers/Containerfile.ci-python -t ci-python .

# Test the container
docker run --rm ci-python python3.12 --version
```

## ACT Integration

When testing workflows locally with `act`, the containers are not used (to allow for faster local testing). The workflows automatically detect the `ACT` environment variable and fall back to installing tools via GitHub Actions.

## Cache Strategy

The containers include:
- Pre-installed development tools
- Optimized layer caching
- Build cache directories pre-created
- Package registries/indexes cached where possible

This reduces CI build time by eliminating repetitive tool installation and setup steps.