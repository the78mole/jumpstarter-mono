# Architecture

This document describes the overall architecture of the Jumpstarter monorepo and how components interact.

## System Overview

The Jumpstarter system consists of multiple components working together to provide a comprehensive testing and automation platform.

```mermaid
graph TB
    subgraph "User Interface"
        CLI[CLI Tools]
        API[REST API]
        WEB[Web Interface]
    end
    
    subgraph "Core Platform"
        LIB[Jumpstarter Library<br/>Python]
        CTL[Kubernetes Controller<br/>Go]
        PROTO[Protocol Definitions<br/>Protocol Buffers]
    end
    
    subgraph "Hardware Layer"
        FW[DUTLink Firmware<br/>Rust]
        HW[DUTLink Board<br/>Hardware]
        DUT[Device Under Test]
    end
    
    subgraph "Integration & Tooling"
        TEKTON[Tekton CI/CD]
        VSCODE[VS Code Extension]
        DEV[DevSpace Templates]
    end
    
    subgraph "Testing Infrastructure"
        E2E[End-to-End Tests]
        INT[Integration Tests]
        UNIT[Unit Tests]
    end
    
    CLI --> LIB
    API --> LIB
    WEB --> LIB
    
    LIB <--> CTL
    LIB <--> PROTO
    CTL <--> PROTO
    
    LIB --> FW
    FW --> HW
    HW --> DUT
    
    TEKTON --> CTL
    VSCODE --> LIB
    DEV --> LIB
    
    E2E --> LIB
    INT --> LIB
    UNIT --> LIB
    
    style LIB fill:#e1f5fe
    style CTL fill:#f3e5f5
    style FW fill:#fff3e0
    style HW fill:#ffebee
```

## Component Architecture

### Core Components

#### Jumpstarter Library (`core/jumpstarter/`)
- Main Python library and CLI
- Provides core functionality and APIs
- Plugin system for extensibility

#### Controller (`core/controller/`)
- Kubernetes controller written in Go
- Manages test environments and resources
- Handles orchestration and scheduling

#### Protocol (`core/protocol/`)
- Communication protocol definitions
- Shared data structures
- API specifications

### Hardware Components

#### DUT Link Firmware (`hardware/dutlink-firmware/`)
- Rust-based firmware for hardware control
- Low-level device interaction
- Real-time communication protocols

#### DUT Link Board (`hardware/dutlink-board/`)
- Hardware design files
- PCB layouts and schematics
- Component specifications

## Monorepo Build Architecture

The monorepo uses a unified build system that coordinates between different technologies:

```mermaid
graph TD
    subgraph "Language Ecosystems"
        PY[Python Components<br/>core/jumpstarter<br/>templates/driver]
        GO[Go Components<br/>core/controller<br/>lab-config]
        RUST[Rust Components<br/>hardware/dutlink-firmware]
        TS[TypeScript Components<br/>integrations/vscode]
    end
    
    subgraph "Build Tools"
        UV[UV Package Manager<br/>Python Workspace]
        GOWORK[Go Workspace<br/>go.work]
        CARGO[Cargo<br/>Rust Build]
        NPM[NPM<br/>Node.js Build]
    end
    
    subgraph "Unified Orchestration"
        MAKE[Root Makefile<br/>40+ Build Targets]
        CI[GitHub Actions<br/>Multi-language CI]
    end
    
    subgraph "Output Artifacts"
        WHEELS[Python Wheels]
        BINS[Go Binaries]
        FIRMWARE[Rust Firmware]
        EXTENSION[VS Code Extension]
        CONTAINERS[Container Images]
        PACKAGES[Distribution Packages]
    end
    
    PY --> UV
    GO --> GOWORK
    RUST --> CARGO
    TS --> NPM
    
    UV --> MAKE
    GOWORK --> MAKE
    CARGO --> MAKE
    NPM --> MAKE
    
    MAKE --> CI
    
    MAKE --> WHEELS
    MAKE --> BINS
    MAKE --> FIRMWARE
    MAKE --> EXTENSION
    MAKE --> CONTAINERS
    MAKE --> PACKAGES
    
    style MAKE fill:#e8f5e8
    style CI fill:#fff9c4
```

## Design Principles

1. **Modularity**: Each component can be developed and tested independently
2. **Consistency**: Unified build, test, and deployment processes
3. **Scalability**: Components can scale independently
4. **Maintainability**: Clear separation of concerns
5. **Extensibility**: Plugin architecture for custom functionality

## Data Flow Architecture

```mermaid
sequenceDiagram
    participant User
    participant CLI
    participant Library
    participant Controller
    participant Firmware
    participant Hardware
    participant DUT

    User->>CLI: Execute test command
    CLI->>Library: Parse and validate
    Library->>Controller: Schedule test job
    Controller->>Library: Allocate resources
    Library->>Firmware: Send control commands
    Firmware->>Hardware: Control hardware pins
    Hardware->>DUT: Physical interaction
    DUT-->>Hardware: Response signals
    Hardware-->>Firmware: Read status
    Firmware-->>Library: Report results
    Library-->>Controller: Update job status
    Controller-->>CLI: Return test results
    CLI-->>User: Display results
```

## Deployment Architecture

The system supports multiple deployment models:

```mermaid
graph TB
    subgraph "Standalone Deployment"
        S_CLI[CLI Tools]
        S_LIB[Local Library]
        S_HW[Direct Hardware]
    end
    
    subgraph "Kubernetes Deployment"
        K_API[API Gateway]
        K_CTL[Controller Pods]
        K_LIB[Library Services]
        K_HW[Hardware Nodes]
    end
    
    subgraph "Hybrid Deployment"
        H_CLI[Local CLI]
        H_K8S[Remote K8s Cluster]
        H_HW[Local Hardware]
    end
    
    S_CLI --> S_LIB
    S_LIB --> S_HW
    
    K_API --> K_CTL
    K_CTL --> K_LIB
    K_LIB --> K_HW
    
    H_CLI --> H_K8S
    H_CLI --> H_HW
    H_K8S --> H_HW
    
    style S_LIB fill:#e1f5fe
    style K_CTL fill:#f3e5f5
    style H_K8S fill:#fff3e0
```

## Package Distribution Flow

```mermaid
graph LR
    subgraph "Source Code"
        SRC[Monorepo Source]
    end
    
    subgraph "Build Process"
        BUILD[Unified Build System]
    end
    
    subgraph "Package Types"
        PY_PKG[Python Wheels<br/>PyPI]
        DEB_PKG[Debian Packages<br/>APT Repo]
        RPM_PKG[RPM Packages<br/>YUM Repo]
        CONT[Container Images<br/>Registry]
        FW_PKG[Firmware Binaries<br/>Releases]
    end
    
    subgraph "Distribution"
        PYPI[PyPI Repository]
        DEB_REPO[Debian Repository]
        RPM_REPO[RPM Repository]
        DOCKER[Container Registry]
        GITHUB[GitHub Releases]
    end
    
    SRC --> BUILD
    BUILD --> PY_PKG
    BUILD --> DEB_PKG
    BUILD --> RPM_PKG
    BUILD --> CONT
    BUILD --> FW_PKG
    
    PY_PKG --> PYPI
    DEB_PKG --> DEB_REPO
    RPM_PKG --> RPM_REPO
    CONT --> DOCKER
    FW_PKG --> GITHUB
    
    style BUILD fill:#e8f5e8
```

## Security Considerations

- Component isolation
- Secure communication protocols
- Access control and authentication
- Audit logging