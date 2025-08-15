# Jumpstarter Monorepo Documentation

Welcome to the Jumpstarter monorepo documentation. This unified repository contains all Jumpstarter components in a single, manageable structure.

## Overview

The Jumpstarter monorepo brings together:

- **Core Components**: Main library, controller, and protocol definitions
- **Hardware**: DUT link board designs and firmware
- **Packages**: Distribution packages for various platforms
- **Integrations**: CI/CD tooling and development environment integrations
- **Templates**: Scaffolding for driver development
- **Testing**: Comprehensive testing infrastructure

## Quick Navigation

- [Installation Guide](installation/index.md) - Get started with Jumpstarter
- [Development Guide](development/index.md) - Contribute to the project
- [Architecture](architecture/index.md) - Understand the system design
- [User Guide](user-guide/index.md) - Learn how to use Jumpstarter
- [Hardware](hardware/index.md) - Hardware components and specifications
- [Integrations](integrations/index.md) - CI/CD and tooling integrations

## Repository Structure

The monorepo is organized into logical components:

```mermaid
graph TB
    subgraph "Jumpstarter Monorepo"
        ROOT[jumpstarter-mono/]

        subgraph "Core Platform"
            CORE[core/]
            CORE_JS[├── jumpstarter/<br/>│   Python Library & CLI]
            CORE_CTL[├── controller/<br/>│   Kubernetes Controller Go]
            CORE_PROTO[└── protocol/<br/>    Protocol Definitions]
        end

        subgraph "Hardware"
            HW[hardware/]
            HW_BOARD[├── dutlink-board/<br/>│   PCB Design Files]
            HW_FW[└── dutlink-firmware/<br/>    Rust Firmware]
        end

        subgraph "Distribution"
            PKG[packages/]
            PKG_PY[├── python/]
            PKG_DEB[├── debian/]
            PKG_RPM[├── rpm/]
            PKG_CONT[└── container/]
        end

        subgraph "Integration & Tools"
            INT[integrations/]
            INT_TEK[├── tekton/]
            INT_VS[├── vscode/]
            INT_DEV[└── devspace/]

            TMPL[templates/]
            TMPL_DRV[└── driver/]

            TEST[testing/]
            TEST_E2E[├── e2e/]
            TEST_INT[├── integration/]
            TEST_FIX[└── fixtures/]
        end

        subgraph "Configuration & Docs"
            LAB[lab-config/]
            DOCS[docs/]
            TOOLS[tools/]
            SCRIPTS[scripts/]
        end
    end

    ROOT --> CORE
    CORE --> CORE_JS
    CORE --> CORE_CTL
    CORE --> CORE_PROTO

    ROOT --> HW
    HW --> HW_BOARD
    HW --> HW_FW

    ROOT --> PKG
    PKG --> PKG_PY
    PKG --> PKG_DEB
    PKG --> PKG_RPM
    PKG --> PKG_CONT

    ROOT --> INT
    INT --> INT_TEK
    INT --> INT_VS
    INT --> INT_DEV

    ROOT --> TMPL
    TMPL --> TMPL_DRV

    ROOT --> TEST
    TEST --> TEST_E2E
    TEST --> TEST_INT
    TEST --> TEST_FIX

    ROOT --> LAB
    ROOT --> DOCS
    ROOT --> TOOLS
    ROOT --> SCRIPTS

    style ROOT fill:#e8f5e8
    style CORE fill:#e1f5fe
    style HW fill:#fff3e0
    style PKG fill:#f3e5f5
    style INT fill:#ffebee
```

## Getting Started

To get started with development:

```bash
# Clone the repository
git clone https://github.com/the78mole/jumpstarter-mono.git
cd jumpstarter-mono

# Setup development environment
make setup

# View available commands
make help
```

For detailed instructions, see the [Installation Guide](installation/index.md).
