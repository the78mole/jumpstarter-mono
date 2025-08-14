# Development Guide

This guide covers the development workflow for the Jumpstarter monorepo.

## Development Environment

The monorepo uses a unified development environment with multi-language support.

```mermaid
graph TB
    subgraph "Development Setup"
        CLONE[Clone Repository]
        SETUP[make setup]
        ENV[Development Environment]
    end

    subgraph "Multi-language Tools"
        PY_TOOLS[Python Tools<br/>UV, Black, Ruff]
        GO_TOOLS[Go Tools<br/>gofmt, golangci-lint]
        RUST_TOOLS[Rust Tools<br/>cargo, clippy, rustfmt]
        TS_TOOLS[TypeScript Tools<br/>npm, eslint, prettier]
    end

    subgraph "Unified Commands"
        BUILD[make build]
        TEST[make test]
        LINT[make lint]
        FMT[make fmt]
    end

    CLONE --> SETUP
    SETUP --> ENV
    ENV --> PY_TOOLS
    ENV --> GO_TOOLS
    ENV --> RUST_TOOLS
    ENV --> TS_TOOLS

    PY_TOOLS --> BUILD
    GO_TOOLS --> BUILD
    RUST_TOOLS --> BUILD
    TS_TOOLS --> BUILD

    BUILD --> TEST
    TEST --> LINT
    LINT --> FMT

    style ENV fill:#e8f5e8
    style BUILD fill:#e1f5fe
```

### Available Commands

Use `make help` to see all available commands:

- `make setup` - Initial setup of development environment
- `make build` - Build all components
- `make test` - Run all tests
- `make lint` - Run all linters
- `make fmt` - Format all code
- `make clean` - Clean all build artifacts

### Language-specific Development

#### Python

```bash
make dev-python    # Start Python development environment
make test-python   # Run Python tests
make lint-python   # Lint Python code
```

#### Go

```bash
make dev-go       # Start Go development environment
make test-go      # Run Go tests
make lint-go      # Lint Go code
```

#### Rust

```bash
make dev-rust     # Start Rust development environment
make test-rust    # Run Rust tests
make lint-rust    # Lint Rust code
```

## Contributing

### Workflow

```mermaid
flowchart TD
    START[Start Development]
    BRANCH[Create Feature Branch]
    CODE[Make Changes]
    LOCAL[Test Locally with Act]
    BUILD[Run make build]
    TEST[Run make test]
    LINT[Run make lint]
    COMMIT[Commit Changes]
    PR[Submit Pull Request]
    REVIEW[Code Review]
    MERGE[Merge to Main]

    START --> BRANCH
    BRANCH --> CODE
    CODE --> LOCAL
    LOCAL --> BUILD
    BUILD --> TEST
    TEST --> LINT
    LINT --> COMMIT
    COMMIT --> PR
    PR --> REVIEW
    REVIEW -->|Approved| MERGE
    REVIEW -->|Changes Requested| CODE

    BUILD -->|Failed| CODE
    TEST -->|Failed| CODE
    LINT -->|Failed| CODE
    LOCAL -->|Failed| CODE

    style START fill:#e8f5e8
    style MERGE fill:#e1f5fe
    style CODE fill:#fff3e0
    style LOCAL fill:#ffebee
```

1. Create a feature branch
2. Make your changes
3. **Test workflows locally**: See [Local Workflow Testing](act-integration.md) for act integration
4. Run tests: `make test`
5. Run linting: `make lint`
6. Submit a pull request

### Code Standards

- Follow language-specific conventions
- Ensure all tests pass
- Maintain documentation
- Use pre-commit hooks: `make pre-commit-install`

## Project Structure

Each component follows its own conventions while integrating with the monorepo build system.

```mermaid
graph TB
    subgraph "Monorepo Structure"
        ROOT[jumpstarter-mono/]

        subgraph "Core Components"
            CORE_JS[core/jumpstarter/]
            CORE_CTL[core/controller/]
            CORE_PROTO[core/protocol/]
        end

        subgraph "Hardware Components"
            HW_FW[hardware/dutlink-firmware/]
            HW_BOARD[hardware/dutlink-board/]
        end

        subgraph "Support Components"
            PKG[packages/]
            INT[integrations/]
            TMPL[templates/]
            TEST[testing/]
            LAB[lab-config/]
        end
    end

    ROOT --> CORE_JS
    ROOT --> CORE_CTL
    ROOT --> CORE_PROTO
    ROOT --> HW_FW
    ROOT --> HW_BOARD
    ROOT --> PKG
    ROOT --> INT
    ROOT --> TMPL
    ROOT --> TEST
    ROOT --> LAB

    style ROOT fill:#e8f5e8
    style CORE_JS fill:#e1f5fe
    style CORE_CTL fill:#f3e5f5
    style HW_FW fill:#fff3e0
```

### Core Components

- `core/jumpstarter/` - Main Python library
- `core/controller/` - Kubernetes controller (Go)
- `core/protocol/` - Protocol definitions

### Hardware Components

- `hardware/dutlink-firmware/` - Rust firmware
- `hardware/dutlink-board/` - Hardware design files

## Testing

The monorepo includes comprehensive testing:

```mermaid
graph TB
    subgraph "Testing Pyramid"
        E2E[End-to-End Tests<br/>Full system validation]
        INT[Integration Tests<br/>Component interactions]
        UNIT[Unit Tests<br/>Individual functions]
    end

    subgraph "Test Execution"
        ALL[make test<br/>All Tests]
        LANG_PY[make test-python]
        LANG_GO[make test-go]
        LANG_RUST[make test-rust]
        E2E_ONLY[make test-e2e]
    end

    subgraph "Workflow Testing"
        ACT[act - Local Workflow Testing]
        CI_LOCAL[Test CI/CD Locally]
        VALIDATE[Validation Scripts]
    end

    subgraph "Test Types"
        PERF[Performance Tests]
        SEC[Security Tests]
        COMPAT[Compatibility Tests]
    end

    UNIT --> INT
    INT --> E2E

    ALL --> LANG_PY
    ALL --> LANG_GO
    ALL --> LANG_RUST
    ALL --> E2E_ONLY

    E2E --> PERF
    E2E --> SEC
    E2E --> COMPAT

    ACT --> CI_LOCAL
    CI_LOCAL --> VALIDATE

    style E2E fill:#ffebee
    style INT fill:#fff3e0
    style UNIT fill:#e8f5e8
    style ACT fill:#e1f5fe
```

### Local Workflow Testing

Test GitHub Actions workflows locally using [act](https://github.com/nektos/act):

```bash
# Test the full CI pipeline locally
act workflow_dispatch -W .github/workflows/ci.yml --env ACT=true

# Validate all workflows
./scripts/validate-ci-with-act.sh
```

For detailed instructions, see [Local Workflow Testing](act-integration.md).

### Code Testing

- Unit tests for each component
- Integration tests
- End-to-end tests
- Performance tests

Run all tests with:

```bash
make test
```
