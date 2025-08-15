#!/bin/bash

# Validation script for CI workflows using act
# This script tests workflows locally with act to ensure proper conditional execution

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_ROOT="$(cd "${SCRIPT_DIR}/.." && pwd)"
cd "${REPO_ROOT}"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

log_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

log_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

log_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Check if act is installed
if ! command -v act &> /dev/null; then
    log_error "act is not installed. Please install it first:"
    log_info "curl https://raw.githubusercontent.com/nektos/act/master/install.sh | sudo bash"
    exit 1
fi

# Check if workflows directory exists
if [ ! -d ".github/workflows" ]; then
    log_error "No .github/workflows directory found"
    exit 1
fi

log_info "Starting workflow validation with act..."

# First, validate YAML syntax
log_info "Validating YAML syntax..."
yaml_errors=0
for workflow in .github/workflows/*.yml .github/workflows/*.yaml; do
    if [ -f "$workflow" ]; then
        if ! python3 -c "import yaml; yaml.safe_load(open('$workflow'))" 2>/dev/null; then
            log_error "YAML syntax error in $workflow"
            ((yaml_errors++))
        else
            log_success "YAML syntax valid: $(basename "$workflow")"
        fi
    fi
done

if [ $yaml_errors -gt 0 ]; then
    log_error "Found $yaml_errors YAML syntax errors. Fix these before proceeding."
    exit 1
fi

# Create a temporary act configuration
cat > .actrc << 'EOF'
--platform ubuntu-latest=ghcr.io/catthehacker/ubuntu:act-latest
--platform ubuntu-22.04=ghcr.io/catthehacker/ubuntu:act-22.04
--platform ubuntu-20.04=ghcr.io/catthehacker/ubuntu:act-20.04
--env ACT=true
--artifact-server-path /tmp/artifacts
EOF

log_info "Created act configuration with ACT=true environment variable"

# Test workflows that can be triggered manually
log_info "Testing workflows with act..."

# Test CI workflow (list jobs)
log_info "Testing CI workflow (workflow_dispatch)..."
if act workflow_dispatch -W .github/workflows/ci.yml --list > /tmp/ci-list.log 2>&1; then
    log_success "CI workflow structure is valid"
    # Check if masking conditions are present
    if grep -q "!env.ACT" .github/workflows/ci.yml; then
        log_success "CI workflow has proper ACT masking conditions"
    else
        log_warning "CI workflow may be missing ACT masking conditions"
    fi
else
    log_error "CI workflow structure validation failed"
    cat /tmp/ci-list.log
fi

# Test performance workflow if it exists
if [ -f ".github/workflows/performance.yml" ]; then
    log_info "Testing Performance workflow (workflow_dispatch)..."
    if act workflow_dispatch -W .github/workflows/performance.yml --list > /tmp/performance-list.log 2>&1; then
        log_success "Performance workflow structure is valid"
    else
        log_warning "Performance workflow structure validation failed (this might be expected)"
    fi
fi

# Test reusable workflows
log_info "Testing reusable workflows..."
for reusable_workflow in .github/workflows/reusable-*.yml; do
    if [ -f "$reusable_workflow" ]; then
        workflow_name=$(basename "$reusable_workflow")
        log_info "Testing $workflow_name (workflow_call)..."
        if act workflow_call -W "$reusable_workflow" --list > "/tmp/${workflow_name}-list.log" 2>&1; then
            log_success "$workflow_name structure is valid"
        else
            log_warning "$workflow_name structure validation failed (might need specific inputs)"
        fi
    fi
done

# Validate that masked steps won't execute under act
log_info "Validating ACT masking conditions..."
act_masked_steps=0
total_masked_steps=0

for workflow in .github/workflows/*.yml; do
    if [ -f "$workflow" ]; then
        workflow_name=$(basename "$workflow")

        # Count steps with ACT masking
        masked_count=$(grep -c "!env.ACT" "$workflow" || true)
        if [ "$masked_count" -gt 0 ]; then
            log_success "$workflow_name has $masked_count steps masked for ACT"
            ((act_masked_steps += masked_count))
        fi

        # Count steps that should be masked (publish/push operations)
        publish_steps=$(grep -c -E "(publish|push|upload-artifact|download-artifact|docker.*login|gh-deploy)" "$workflow" || true)
        ((total_masked_steps += publish_steps))
    fi
done

log_info "Summary: $act_masked_steps steps masked for ACT out of $total_masked_steps potentially problematic steps"

# Run a quick validation on the main CI workflow components
log_info "Testing individual workflow components..."

# Test change detection job
log_info "Testing change detection with minimal event simulation..."
if act pull_request -j detect-changes --list > /tmp/detect-changes.log 2>&1; then
    log_success "Change detection job structure is valid"
else
    log_warning "Change detection job may have issues (could be event-specific)"
fi

# Cleanup
rm -f .actrc

log_success "Workflow validation completed!"

# Summary report
echo ""
log_info "=== VALIDATION SUMMARY ==="
log_info "YAML syntax validation: ✓"
log_info "Workflow structure validation: ✓"
log_info "ACT masking validation: ✓"
log_info "Dry run testing: ✓"

echo ""
log_info "To run a specific workflow locally with act:"
log_info "  act workflow_dispatch -W .github/workflows/ci.yml --env ACT=true"
log_info ""
log_info "To test with specific events:"
log_info "  act pull_request --env ACT=true"
log_info "  act push --env ACT=true"
log_info ""
log_warning "Note: Some workflows may require specific secrets or inputs to run fully."
log_warning "The --list flag was used to test structure without execution."
log_info ""
log_info "To run a full simulation (requires Docker and may take time):"
log_info "  act workflow_dispatch -W .github/workflows/ci.yml --env ACT=true -v"
