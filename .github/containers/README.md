# CI Container Images

This directory contains container definitions for optimized CI environments. These images provide pre-installed tools to speed up CI builds and ensure consistent environments.

## Available Images

### CI Container Images (for development and testing)

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
- Go 1.24.6 (latest from golang.org)
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
- Go 1.24.6, Python 3.12, Rust stable, Node.js 18 LTS
- Suitable for workflows that need multiple toolchains
- Larger image size but maximum compatibility

### Production Container Images (Red Hat UBI)

The production applications also use containerized environments based on Red Hat Universal Base Images (UBI):

#### Core Controller & Lab Config
- **Base Image**: `registry.access.redhat.com/ubi9/go-toolset:1.24`
- **Runtime**: `registry.access.redhat.com/ubi9/ubi-micro:9.5`
- **Go Version**: 1.24.4 (Red Hat enterprise-ready distribution)
- **Usage**: Production deployment of Jumpstarter controller and lab configuration services
- **Security**: Red Hat UBI provides enterprise-grade security and support

The production containers use Red Hat UBI for enhanced security, enterprise support, and compliance requirements. These images are automatically tested and validated against Red Hat's security policies.

## Usage

These images are automatically built and published to `ghcr.io/the78mole/jumpstarter-mono/ci-*:latest` when:
- Container definitions change
- CI configuration changes
- Manually triggered via workflow_dispatch
- Weekly on schedule (Sundays at 2 AM UTC)

### Enabling Container Usage

Currently, containers are **disabled by default** in the CI workflows to ensure compatibility with local testing via `act`. To enable container usage in production:

1. Edit `.github/workflows/ci.yml`
2. Change `use-container: false` to `use-container: true` for each workflow
3. Commit and push the changes

Example:
```yaml
python-lint:
  uses: ./.github/workflows/reusable-python-lint.yml
  with:
    working-directory: core/jumpstarter
    use-container: true  # Enable container usage
```

### Performance Benefits

When enabled, containers provide:
- **2-5 minutes saved per job** (eliminates tool installation)
- Consistent environments across all CI runs
- Better caching through pre-built layers
- Reduced network usage during CI runs

## Building Locally

To build and test a container locally:

```bash
# Build a specific container
docker build -f .github/containers/Containerfile.ci-python -t ci-python .

# Test the container
docker run --rm ci-python python3.12 --version
```

### Validating Red Hat Container Access

To validate that Red Hat UBI containers are accessible and working correctly:

```bash
# Run the validation script
./scripts/validate-redhat-containers.sh
```

This script verifies:
- Red Hat container registry connectivity
- Go toolset and runtime image availability
- Container functionality with current environment

## ACT Integration

When testing workflows locally with `act`, the workflows automatically fall back to installing tools via GitHub Actions (regardless of the `use-container` setting). This ensures local testing remains fast and doesn't require pulling large container images.

## Cache Strategy

The containers include:
- Pre-installed development tools
- Optimized layer caching
- Build cache directories pre-created
- Package registries/indexes cached where possible

This reduces CI build time by eliminating repetitive tool installation and setup steps.