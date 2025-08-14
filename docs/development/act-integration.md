# Act Integration for Local Workflow Testing

This repository includes comprehensive integration with [act](https://github.com/nektos/act) for local testing of GitHub Actions workflows. Act allows you to run workflows locally in Docker containers, enabling faster iteration and debugging.

## Installation

Install act using the official installer:

```bash
curl https://raw.githubusercontent.com/nektos/act/master/install.sh | sudo bash
```

Or on macOS with Homebrew:

```bash
brew install act
```

## Configuration

The workflows are designed to work with act out of the box. When running under act, the `ACT` environment variable is automatically set to `true`, which conditionally disables steps that:

- Push artifacts to external services
- Make changes to the repository (tags, releases)
- Require authentication tokens
- Perform Docker registry operations

## Usage Examples

### Basic Workflow Testing

Test the CI workflow structure:
```bash
act workflow_dispatch -W .github/workflows/ci.yml --list
```

Run the CI workflow with act environment:
```bash
act workflow_dispatch -W .github/workflows/ci.yml --env ACT=true
```

### Event-Specific Testing

Test pull request events:
```bash
act pull_request --env ACT=true
```

Test push events:
```bash
act push --env ACT=true
```

### Job-Specific Testing

Run only the change detection job:
```bash
act workflow_dispatch -j detect-changes --env ACT=true
```

Run specific language builds:
```bash
act workflow_dispatch -j python-test --env ACT=true
act workflow_dispatch -j controller-build --env ACT=true
```

### Reusable Workflow Testing

Test reusable workflows:
```bash
act workflow_call -W .github/workflows/reusable-rust-build.yml \
  --input working-directory=hardware/dutlink-firmware \
  --input targets=thumbv7em-none-eabihf
```

## Automated Validation

Use the provided validation script to test all workflows:

```bash
./scripts/validate-ci-with-act.sh
```

This script performs:
- YAML syntax validation
- Workflow structure validation
- ACT masking condition verification
- Basic workflow execution testing

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

### Common Issues

**Git not found**:
```bash
# Use a larger image with git installed
act --platform ubuntu-latest=catthehacker/ubuntu:act-latest
```

**Bash not found**:
```bash
# Use an image with bash
act --platform ubuntu-latest=ubuntu:20.04
```

**Action compatibility**:
```bash
# Use the full image for maximum compatibility
act --platform ubuntu-latest=catthehacker/ubuntu:full-latest
```

### Debug Mode

Run with verbose output for debugging:
```bash
act workflow_dispatch --verbose --env ACT=true
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

- **Use appropriate image sizes**: Micro for simple jobs, Medium for most cases
- **Cache Docker images**: Reuse images across runs
- **Parallel execution**: Act supports running multiple jobs in parallel
- **Selective job execution**: Use `-j` flag to run specific jobs

This integration ensures that workflows can be thoroughly tested locally, reducing CI/CD feedback cycles and improving development velocity.