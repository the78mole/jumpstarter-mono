# Local Workflow Testing with Act

This repository includes comprehensive integration with [act](https://github.com/nektos/act) for local testing of GitHub Actions workflows. Act allows you to run workflows locally in Docker containers, enabling faster iteration and debugging of CI/CD changes.

## Quick Start

```bash
# Install act (see installation options below)
curl https://raw.githubusercontent.com/nektos/act/master/install.sh | sudo bash

# Test the main CI workflow
act workflow_dispatch -W .github/workflows/ci.yml --env ACT=true

# Run comprehensive validation on all workflows
./scripts/validate-ci-with-act.sh

# Test specific masking conditions
./scripts/test-act-masking.sh
```

> **💡 Tip:** Use the validation scripts first to ensure your environment is properly configured before running individual workflows.

## Installation

### Option 1: Official Installer (Recommended)

Install act using the official installer:

```bash
curl https://raw.githubusercontent.com/nektos/act/master/install.sh | sudo bash
```

### Option 2: Package Managers

**macOS with Homebrew:**

```bash
brew install act
```

**Arch Linux:**

```bash
sudo pacman -S act
```

**Manual Download:**
If you experience network issues with the installer, download directly from GitHub releases:

```bash
# Download the latest release for your platform
LATEST_VERSION=$(curl -s https://api.github.com/repos/nektos/act/releases/latest | grep -oP '"tag_name": "\K(.*)(?=")')
wget "https://github.com/nektos/act/releases/download/${LATEST_VERSION}/act_Linux_x86_64.tar.gz"
tar xzf act_Linux_x86_64.tar.gz
sudo mv act /usr/local/bin/
```

### Option 3: Using Go

If you have Go installed:

```bash
go install github.com/nektos/act@latest
```

### Option 4: Docker (No Installation Required)

Run act directly with Docker:

```bash
# Create an alias for easier use
alias act='docker run --rm -it -v /var/run/docker.sock:/var/run/docker.sock -v "$PWD":/workspace -w /workspace nektos/act'

# Use normally
act workflow_dispatch -W .github/workflows/ci.yml --env ACT=true
```

## Configuration

The workflows are designed to work with act out of the box. When running under act, the `ACT` environment variable is automatically set to `true`, which conditionally disables steps that:

- Push artifacts to external services
- Make changes to the repository (tags, releases)
- Require authentication tokens
- Perform Docker registry operations

### Repository .actrc Configuration

The repository includes a `.actrc` configuration file with optimal defaults:

```ini
# Use medium-sized Ubuntu container for better compatibility
-P ubuntu-latest=catthehacker/ubuntu:act-latest

# Set ACT environment variable automatically
--env ACT=true

# Enable artifact server for local artifact handling
--artifact-server-path /tmp/artifacts

# Default event for testing
--eventpath .github/workflows/events/workflow_dispatch.json
```

This configuration ensures consistent behavior across all local testing scenarios and eliminates the need to specify `--env ACT=true` manually.

## Usage Examples

### Basic Workflow Testing

Test the CI workflow structure:

```bash
act workflow_dispatch -W .github/workflows/ci.yml --list
```

Run the CI workflow (using .actrc configuration):

```bash
act workflow_dispatch -W .github/workflows/ci.yml
```

### Event-Specific Testing

Test pull request events:

```bash
act pull_request
```

Test push events:

```bash
act push
```

### Job-Specific Testing

Run only the change detection job:

```bash
act workflow_dispatch -j detect-changes
```

Run specific language builds:

```bash
act workflow_dispatch -j python-test
act workflow_dispatch -j controller-build
```

### Reusable Workflow Testing

Test reusable workflows:

```bash
act workflow_call -W .github/workflows/reusable-rust-build.yml \
  --input working-directory=hardware/dutlink-firmware \
  --input targets=thumbv7em-none-eabihf
```

## Automated Validation

The repository includes comprehensive validation tools for testing workflows with act:

### Validation Script

Use the main validation script to test all workflows:

```bash
./scripts/validate-ci-with-act.sh
```

This script performs:

- **YAML syntax validation** for all workflow files
- **Workflow structure validation** using act's `--list` mode
- **ACT masking condition verification** to ensure proper conditional execution
- **Basic workflow execution testing** with dry-run mode
- **Summary reporting** of validation results

### ACT Masking Test

Test specific masking conditions:

```bash
./scripts/test-act-masking.sh
```

This script validates:

- All steps that should be masked under act have proper `if: ${{ !env.ACT }}` conditions
- Steps that push artifacts, publish packages, or modify repositories are properly masked
- Local-only debug steps work correctly under act

### Manual Validation

You can also run individual validations:

```bash
# Validate workflow structure
act workflow_dispatch -W .github/workflows/ci.yml --list

# Test workflow with act environment
act workflow_dispatch -W .github/workflows/ci.yml --env ACT=true --dry-run

# Validate specific jobs
act workflow_dispatch -j detect-changes --list
```

## Container Images

Act supports different runner images:

- **Micro** (~200MB): Basic NodeJS environment, fastest but limited
- **Medium** (~500MB): Includes common tools, good balance
- **Large** (~17GB): Full GitHub runner environment, most compatible

Choose based on your workflow requirements:

```bash
# Use medium image (recommended)
act --platform ubuntu-latest=catthehacker/ubuntu:act-latest

# Use micro image for simple workflows
act --platform ubuntu-latest=node:16-alpine

# Use large image for complex workflows
act --platform ubuntu-latest=catthehacker/ubuntu:full-latest
```

## Conditional Masking

Steps that should not run under act are masked with:

```yaml
- name: Step that pushes artifacts
  if: ${{ !env.ACT }}
  uses: actions/upload-artifact@v4
  # ... rest of step
```

Steps that only run under act:

```yaml
- name: Debug step for local testing
  if: ${{ env.ACT }}
  run: echo "Running under act for local testing"
```

## Limitations

When running with act:

1. **No external authentication**: Steps requiring GitHub tokens, PyPI tokens, etc. are masked
2. **No artifact persistence**: Artifacts are not uploaded to GitHub
3. **Limited container images**: Some tools may not be available in minimal images
4. **No real external services**: External API calls may fail
5. **File system differences**: Container filesystem differs from GitHub runners

## Troubleshooting

### Common Issues and Solutions

#### 1. Paths-filter Action Errors

**Error**: "This action requires 'base' input to be configured"

**Solution**: The CI workflow now includes automatic base branch configuration for act:

```yaml
- name: Check for changes
  uses: dorny/paths-filter@v3
  with:
    base: ${{ env.ACT && 'main' || '' }}
    filters: |
      # ... filters
```

This ensures the paths-filter action works correctly in both GitHub Actions and act environments.

#### 2. Certificate/Network Issues

**Error**: "self-signed certificate in certificate chain"

This is common when using actions that download tools. Solutions:

```bash
# Option 1: Use pre-built containers with tools included
act -P ubuntu-latest=ghcr.io/catthehacker/ubuntu:full-latest

# Option 2: Skip problematic jobs during testing
act workflow_dispatch -j detect-changes  # Test only specific jobs

# Option 3: Use dryrun mode for structure validation
act workflow_dispatch --dryrun
```

#### 3. Container Image Issues

**Error**: Command not found or missing tools

```bash
# Use full-featured container for complete tool compatibility
act --platform ubuntu-latest=ghcr.io/catthehacker/ubuntu:full-latest

# Or specify different container images per job
act -P ubuntu-latest=ubuntu:22.04  # Official Ubuntu with apt packages
```

### Network and Firewall Issues

**DNS Resolution Problems:**
If you encounter DNS blocks for `api.nektosact.com` or other act-related domains:

```bash
# Option 1: Use Docker-based act (bypasses local network issues)
alias act='docker run --rm -it -v /var/run/docker.sock:/var/run/docker.sock -v "$PWD":/workspace -w /workspace nektos/act'

# Option 2: Download manually and install
# See installation options above for manual download instructions

# Option 3: Use alternative container registries
act --platform ubuntu-latest=ubuntu:20.04  # Use official Ubuntu instead
```

**Corporate Firewall/Proxy Issues:**

```bash
# Configure Docker to use proxy
export HTTPS_PROXY=your-proxy:port
export HTTP_PROXY=your-proxy:port

# Or configure act to use different registries
act --platform ubuntu-latest=registry-1.docker.io/library/ubuntu:20.04
```

### Common Issues

**Git not found:**

```bash
# Use a larger image with git installed
act --platform ubuntu-latest=catthehacker/ubuntu:act-latest

# Or specify in ~/.actrc
echo "--platform ubuntu-latest=catthehacker/ubuntu:act-latest" >> ~/.actrc
```

**Bash not found:**

```bash
# Use an image with bash
act --platform ubuntu-latest=ubuntu:20.04

# For minimal setups, install bash in container
act --platform ubuntu-latest=node:16-alpine --env SETUP_BASH=true
```

**Action compatibility issues:**

```bash
# Use the full image for maximum compatibility
act --platform ubuntu-latest=catthehacker/ubuntu:full-latest

# Check specific action requirements
act workflow_dispatch --list  # Shows action compatibility
```

**Docker daemon not running:**

```bash
# Start Docker service
sudo systemctl start docker

# Or for Docker Desktop
open -a Docker  # macOS
# Windows: Start Docker Desktop from Start Menu
```

**Permission denied errors:**

```bash
# Add user to docker group (Linux)
sudo usermod -aG docker $USER
newgrp docker

# Or run with sudo (not recommended for security)
sudo act workflow_dispatch --env ACT=true
```

**Container image pull failures:**

```bash
# Use smaller images that are more likely to be cached
act --platform ubuntu-latest=node:16-alpine

# Pre-pull images manually
docker pull catthehacker/ubuntu:act-latest
docker pull ghcr.io/catthehacker/ubuntu:act-latest

# Use local registry if available
act --platform ubuntu-latest=localhost:5000/ubuntu:act-latest
```

### Debug Mode

Run with verbose output for detailed troubleshooting:

```bash
# Maximum verbosity
act workflow_dispatch --verbose --env ACT=true

# With Docker debug info
act workflow_dispatch --env ACT=true --container-daemon-socket /var/run/docker.sock -v

# Test specific job with debug
act workflow_dispatch -j python-test --env ACT=true --verbose
```

### Environment-Specific Issues

**Resource Constraints:**

```bash
# Use micro image for resource-constrained environments
act --platform ubuntu-latest=node:16-alpine

# Limit container resources
act --container-options "--memory=1g --cpus=1"

# Clean up containers after runs
act --rm
```

**File System Permissions:**

```bash
# Fix workspace permissions
act --userns host

# Mount with proper ownership
act --env RUNNER_UID=$(id -u) --env RUNNER_GID=$(id -g)
```

**Multiple Docker Contexts:**

```bash
# List available contexts
docker context ls

# Use specific context
act --container-daemon-socket /var/run/docker.sock

# Or set specific context
docker context use desktop-linux  # or your preferred context
```

### Local Configuration

Create `~/.actrc` for persistent configuration:

```
--platform ubuntu-latest=catthehacker/ubuntu:act-latest
--env ACT=true
--artifact-server-path /tmp/artifacts
```

## Examples

### Test Python Components

```bash
act workflow_dispatch -j python-lint --env ACT=true
act workflow_dispatch -j python-test --env ACT=true
```

### Test Go Components

```bash
act workflow_dispatch -j controller-build --env ACT=true
act workflow_dispatch -j lab-config-build --env ACT=true
```

### Test Rust Components

```bash
act workflow_dispatch -j rust-build --env ACT=true
```

### Test Documentation

```bash
act workflow_dispatch -j docs-build --env ACT=true
```

## Integration with Development Workflow

1. **Before committing**: Run act to validate workflow changes
2. **When debugging CI failures**: Reproduce issues locally with act
3. **When adding new workflows**: Test locally before pushing
4. **When modifying reusable workflows**: Validate with different inputs

## Performance Considerations

### Image Selection Strategy

Choose container images based on your testing needs:

```bash
# For quick syntax/structure validation (fastest)
act --platform ubuntu-latest=node:16-alpine --list

# For most development workflows (recommended balance)
act --platform ubuntu-latest=catthehacker/ubuntu:act-latest

# For complex workflows with many dependencies (slowest but most compatible)
act --platform ubuntu-latest=catthehacker/ubuntu:full-latest
```

### Optimization Tips

**Cache Docker Images:**

```bash
# Pre-pull commonly used images
docker pull catthehacker/ubuntu:act-latest
docker pull ghcr.io/catthehacker/ubuntu:act-latest
docker pull node:16-alpine

# Use local images without version tags to avoid pulls
act --platform ubuntu-latest=catthehacker/ubuntu:act-latest
```

**Selective Job Execution:**

```bash
# Test only changed components
act workflow_dispatch -j detect-changes --env ACT=true
act workflow_dispatch -j python-test --env ACT=true

# Skip jobs that aren't relevant to your changes
act workflow_dispatch --job controller-build --env ACT=true
```

**Resource Management:**

```bash
# Limit container resources for faster startup
act --container-options "--memory=1g --cpus=2"

# Clean up containers automatically
act --rm

# Use bind mounts instead of volumes for faster I/O
act --bind
```

**Parallel Execution:**
Act supports running multiple jobs in parallel by default, but you can control this:

```bash
# Limit parallelism for resource-constrained systems
act --job python-test --job controller-build  # Runs in parallel by default

# Run jobs sequentially if needed
act --job python-test --env ACT=true && act --job controller-build --env ACT=true
```

- **Use appropriate image sizes**: Micro for simple jobs, Medium for most cases
- **Cache Docker images**: Reuse images across runs to avoid repeated downloads
- **Parallel execution**: Act supports running multiple jobs in parallel when dependencies allow
- **Selective job execution**: Use `-j` flag to run only specific jobs you're working on
- **Local artifact storage**: Configure `--artifact-server-path` for faster artifact handling

This integration ensures that workflows can be thoroughly tested locally, reducing CI/CD feedback cycles and improving development velocity.

## Quick Reference

### Essential Commands

| Command                                                            | Purpose                           |
| ------------------------------------------------------------------ | --------------------------------- |
| `act workflow_dispatch -W .github/workflows/ci.yml --env ACT=true` | Test CI workflow                  |
| `act pull_request --env ACT=true`                                  | Test pull request workflow        |
| `act push --env ACT=true`                                          | Test push workflow                |
| `act workflow_dispatch -j python-test --env ACT=true`              | Test specific job                 |
| `./scripts/validate-ci-with-act.sh`                                | Validate all workflows            |
| `act --list`                                                       | List available workflows and jobs |

### Common Flags

| Flag                             | Description                                        |
| -------------------------------- | -------------------------------------------------- |
| `--env ACT=true`                 | Set ACT environment variable (masks publish steps) |
| `--list`                         | List workflows/jobs without execution              |
| `--dry-run`                      | Show what would run without execution              |
| `--verbose`                      | Verbose output for debugging                       |
| `-j JOB_NAME`                    | Run specific job only                              |
| `-W WORKFLOW_FILE`               | Specify workflow file                              |
| `--platform ubuntu-latest=IMAGE` | Use specific container image                       |

### Container Images

| Image                             | Size   | Use Case                      |
| --------------------------------- | ------ | ----------------------------- |
| `node:16-alpine`                  | ~200MB | Simple Node.js workflows      |
| `catthehacker/ubuntu:act-latest`  | ~500MB | General purpose (recommended) |
| `catthehacker/ubuntu:full-latest` | ~17GB  | Maximum compatibility         |
| `ubuntu:20.04`                    | ~200MB | Basic Ubuntu environment      |

### Masked Operations

The following operations are automatically skipped when `ACT=true`:

- ✅ Artifact uploads/downloads (`actions/upload-artifact`, `actions/download-artifact`)
- ✅ Package publishing (PyPI, npm, container registries)
- ✅ Repository modifications (creating tags, releases)
- ✅ External service deployments
- ✅ Marketplace publishing (VSCode extensions)
- ✅ Docker registry authentication and pushes

### Configuration Files

Create `~/.actrc` for persistent configuration:

```
--platform ubuntu-latest=catthehacker/ubuntu:act-latest
--env ACT=true
--artifact-server-path /tmp/artifacts
--rm
```

Create `.actrc` in repository root for project-specific settings:

```
--platform ubuntu-latest=catthehacker/ubuntu:act-latest
--env ACT=true
```
