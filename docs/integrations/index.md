# Integrations

This section covers integrations with CI/CD systems and development tools.

## Integration Architecture Overview

```mermaid
graph TB
    subgraph "CI/CD Platforms"
        TEKTON[Tekton Pipelines<br/>K8s Native CI/CD]
        GITHUB[GitHub Actions<br/>Cloud CI/CD]
        JENKINS[Jenkins<br/>Self-hosted CI/CD]
    end

    subgraph "Development Tools"
        VSCODE[VS Code Extension<br/>Editor Integration]
        DEVSPACE[DevSpace<br/>Development Environment]
        DOCKER[Docker Images<br/>Containerized Deployment]
    end

    subgraph "Jumpstarter Core"
        CONTROLLER[Kubernetes Controller<br/>Go Service]
        CLI[CLI Tools<br/>Python Library]
        API[REST API<br/>Integration Interface]
    end

    subgraph "Monitoring & Observability"
        PROMETHEUS[Prometheus<br/>Metrics Collection]
        GRAFANA[Grafana<br/>Visualization]
        LOGGING[Centralized Logging<br/>ELK Stack]
    end

    subgraph "Hardware Layer"
        HARDWARE[DUT Link Boards<br/>Physical Hardware]
        DEVICES[Test Devices<br/>Under Test]
    end

    TEKTON --> CONTROLLER
    GITHUB --> CLI
    JENKINS --> API

    VSCODE --> CLI
    DEVSPACE --> CONTROLLER
    DOCKER --> CONTROLLER

    CONTROLLER --> API
    CLI --> API

    API --> PROMETHEUS
    CONTROLLER --> GRAFANA
    CLI --> LOGGING

    CONTROLLER --> HARDWARE
    HARDWARE --> DEVICES

    style CONTROLLER fill:#f3e5f5
    style CLI fill:#e1f5fe
    style API fill:#fff3e0
    style HARDWARE fill:#ffebee
```

## CI/CD Integrations

### Tekton Pipelines

Tekton tasks and pipelines for cloud-native CI/CD.

```mermaid
graph LR
    subgraph "Tekton Pipeline Flow"
        TRIGGER[Pipeline Trigger<br/>Git Push/PR]
        PROVISION[Device Provision<br/>Task]
        BUILD[Build & Test<br/>Task]
        HARDWARE[Hardware Test<br/>Task]
        REPORT[Results Report<br/>Task]
        CLEANUP[Cleanup<br/>Task]
    end

    TRIGGER --> PROVISION
    PROVISION --> BUILD
    BUILD --> HARDWARE
    HARDWARE --> REPORT
    REPORT --> CLEANUP

    style PROVISION fill:#e8f5e8
    style HARDWARE fill:#fff3e0
    style REPORT fill:#e1f5fe
```

#### Available Tasks

- **jumpstarter-test**: Run hardware tests in Tekton pipelines
- **device-provision**: Provision test devices
- **results-collect**: Collect and process test results

#### Example Pipeline

```yaml
apiVersion: tekton.dev/v1beta1
kind: Pipeline
metadata:
  name: hardware-test-pipeline
spec:
  params:
    - name: device-config
      type: string
  tasks:
    - name: provision-device
      taskRef:
        name: device-provision
      params:
        - name: config
          value: $(params.device-config)
    - name: run-tests
      runAfter: [provision-device]
      taskRef:
        name: jumpstarter-test
      params:
        - name: test-suite
          value: integration
```

### GitHub Actions Integration Flow

```mermaid
sequenceDiagram
    participant Dev as Developer
    participant GH as GitHub
    participant Runner as Self-hosted Runner
    participant JS as Jumpstarter
    participant HW as Hardware

    Dev->>GH: Push code/Create PR
    GH->>Runner: Trigger workflow
    Runner->>JS: Setup Jumpstarter
    Runner->>JS: Configure devices
    JS->>HW: Provision hardware
    JS->>HW: Execute tests
    HW-->>JS: Test results
    JS-->>Runner: Aggregate results
    Runner-->>GH: Upload artifacts
    GH-->>Dev: Show test results
```

#### Available Actions

- `jumpstarter-dev/setup-action`: Setup Jumpstarter environment
- `jumpstarter-dev/test-action`: Run hardware tests
- `jumpstarter-dev/report-action`: Generate test reports

#### Example Workflow

```yaml
name: Hardware CI
on: [push, pull_request]

jobs:
  hardware-tests:
    runs-on: self-hosted
    steps:
      - uses: actions/checkout@v4
      - uses: jumpstarter-dev/setup-action@v1
        with:
          version: latest
      - uses: jumpstarter-dev/test-action@v1
        with:
          config: tests/ci-config.yaml
          devices: raspberry-pi,arduino
      - uses: jumpstarter-dev/report-action@v1
        if: always()
        with:
          format: junit
```

## Development Environment Integrations

### VS Code Extension Architecture

```mermaid
graph TB
    subgraph "VS Code Extension"
        UI[Extension UI<br/>Panels & Views]
        LANG[Language Support<br/>YAML, Python]
        DEBUG[Debug Interface<br/>Breakpoints]
        TERM[Integrated Terminal<br/>CLI Integration]
    end

    subgraph "Language Server"
        LSP[Language Server<br/>Protocol]
        VALIDATE[Config Validation<br/>Real-time]
        INTELLISENSE[IntelliSense<br/>Completions]
    end

    subgraph "Jumpstarter Integration"
        CLI_INT[CLI Integration<br/>Command Execution]
        DEVICE_MGR[Device Manager<br/>Hardware Control]
        TEST_RUNNER[Test Runner<br/>Execution Engine]
    end

    UI --> LSP
    LANG --> VALIDATE
    DEBUG --> CLI_INT
    TERM --> CLI_INT

    LSP --> INTELLISENSE
    VALIDATE --> INTELLISENSE

    CLI_INT --> DEVICE_MGR
    CLI_INT --> TEST_RUNNER

    style UI fill:#e1f5fe
    style LSP fill:#fff3e0
    style CLI_INT fill:#e8f5e8
```

The Jumpstarter VS Code extension provides:

- Syntax highlighting for configuration files
- IntelliSense for test definitions
- Integrated test runner
- Device management interface
- Real-time test monitoring

#### Installation

```bash
code --install-extension jumpstarter.jumpstarter-vscode
```

#### Features

1. **Configuration Validation**: Real-time validation of YAML configs
2. **Test Runner**: Run tests directly from the editor
3. **Device Explorer**: Browse and manage connected devices
4. **Log Viewer**: View test logs with syntax highlighting
5. **Debugging**: Set breakpoints in test scripts

### DevSpace Development Environment

```mermaid
graph TB
    subgraph "Local Development"
        DEV[Developer Workstation]
        DEVSPACE[DevSpace CLI]
        CONFIG[devspace.yaml]
    end

    subgraph "Kubernetes Cluster"
        NAMESPACE[Dev Namespace]
        CONTROLLER[Controller Pod]
        RUNNER[Test Runner Pod]
        STORAGE[Persistent Storage]
    end

    subgraph "Development Features"
        SYNC[File Sync<br/>Real-time]
        PORT[Port Forwarding<br/>Local Access]
        LOGS[Log Streaming<br/>Real-time]
        SHELL[Remote Shell<br/>Debug Access]
    end

    DEV --> DEVSPACE
    DEVSPACE --> CONFIG
    CONFIG --> NAMESPACE

    NAMESPACE --> CONTROLLER
    NAMESPACE --> RUNNER
    NAMESPACE --> STORAGE

    DEVSPACE --> SYNC
    DEVSPACE --> PORT
    DEVSPACE --> LOGS
    DEVSPACE --> SHELL

    SYNC --> CONTROLLER
    PORT --> CONTROLLER
    LOGS --> RUNNER
    SHELL --> RUNNER

    style DEV fill:#e8f5e8
    style DEVSPACE fill:#e1f5fe
    style CONTROLLER fill:#fff3e0
```

DevSpace configuration for development environments.

#### Setup

```yaml
# devspace.yaml
version: v2beta1
name: jumpstarter-dev

pipelines:
  dev:
    run: |
      start_dev hardware-controller
      start_dev test-runner

deployments:
  hardware-controller:
    helm:
      chart:
        name: jumpstarter-controller
  test-runner:
    kubectl:
      manifests:
        - k8s/test-runner.yaml

dev:
  hardware-controller:
    imageSelector: jumpstarter/controller
    workingDir: /app
    ports:
      - port: "8080:8080"
```

## Container Integrations

### Docker Images

Pre-built Docker images for easy deployment:

- `jumpstarter/controller`: Kubernetes controller
- `jumpstarter/cli`: Command-line interface
- `jumpstarter/test-runner`: Test execution environment

#### Example Usage

```bash
# Run CLI in container
docker run --rm -v $(pwd):/workspace jumpstarter/cli run tests/

# Start controller
docker run -d --name controller jumpstarter/controller

# Run test runner
docker run --rm --device /dev/ttyUSB0 jumpstarter/test-runner
```

### Kubernetes Deployment

Deploy Jumpstarter in Kubernetes:

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: jumpstarter-controller
spec:
  replicas: 1
  selector:
    matchLabels:
      app: jumpstarter-controller
  template:
    metadata:
      labels:
        app: jumpstarter-controller
    spec:
      containers:
        - name: controller
          image: jumpstarter/controller:latest
          ports:
            - containerPort: 8080
          env:
            - name: CONFIG_PATH
              value: /etc/jumpstarter/config.yaml
```

## Monitoring and Observability

### Prometheus Metrics

Jumpstarter exports metrics for monitoring:

- Test execution duration
- Device availability
- Error rates
- Queue depth

### Grafana Dashboards

Pre-built dashboards for visualization:

- Test execution overview
- Device health monitoring
- Performance metrics
- Error analysis

### Logging Integration

Integration with logging systems:

- Structured JSON logging
- Correlation IDs for tracing
- Configurable log levels
- Log forwarding to external systems

## Custom Integrations

### Plugin Architecture

Create custom integrations using the plugin system:

```python
from jumpstarter.plugin import BasePlugin

class CustomIntegration(BasePlugin):
    def __init__(self, config):
        self.config = config

    def on_test_start(self, test_info):
        # Custom logic for test start
        pass

    def on_test_complete(self, test_results):
        # Custom logic for test completion
        pass
```

### API Integration

REST API for external integrations:

```bash
# Start a test
curl -X POST http://localhost:8080/api/v1/tests \
  -H "Content-Type: application/json" \
  -d '{"config": "test-config.yaml"}'

# Get test status
curl http://localhost:8080/api/v1/tests/12345/status

# Get test results
curl http://localhost:8080/api/v1/tests/12345/results
```
