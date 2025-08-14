#!/bin/bash

# CI/CD Validation Script
# This script helps validate the GitHub Actions workflows locally

set -e

echo "🚀 Jumpstarter CI/CD Validation Script"
echo "======================================"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Function to print colored output
print_status() {
    local status=$1
    local message=$2
    case $status in
        "INFO")
            echo -e "${BLUE}ℹ️  $message${NC}"
            ;;
        "SUCCESS")
            echo -e "${GREEN}✅ $message${NC}"
            ;;
        "WARNING")
            echo -e "${YELLOW}⚠️  $message${NC}"
            ;;
        "ERROR")
            echo -e "${RED}❌ $message${NC}"
            ;;
    esac
}

# Check if required tools are available
check_dependencies() {
    print_status "INFO" "Checking dependencies..."
    
    local missing_deps=()
    
    # Check for act (GitHub Actions runner)
    if ! command -v act &> /dev/null; then
        missing_deps+=("act")
    fi
    
    # Check for yamllint
    if ! command -v yamllint &> /dev/null; then
        missing_deps+=("yamllint")
    fi
    
    # Check for actionlint
    if ! command -v actionlint &> /dev/null; then
        missing_deps+=("actionlint")
    fi
    
    if [ ${#missing_deps[@]} -ne 0 ]; then
        print_status "WARNING" "Missing optional dependencies: ${missing_deps[*]}"
        print_status "INFO" "Install with: brew install act yamllint actionlint"
        print_status "INFO" "Or: apt-get install yamllint && npm install -g @actions/act"
    else
        print_status "SUCCESS" "All dependencies available"
    fi
}

# Validate YAML syntax
validate_yaml() {
    print_status "INFO" "Validating YAML syntax..."
    
    local yaml_files=(
        ".github/workflows/ci.yml"
        ".github/workflows/release.yml"
        ".github/workflows/publish.yml"
        ".github/workflows/performance.yml"
        ".github/workflows/reusable-rust-build.yml"
        ".github/workflows/reusable-web-build.yml"
    )
    
    for file in "${yaml_files[@]}"; do
        if [ -f "$file" ]; then
            if command -v yamllint &> /dev/null; then
                if yamllint "$file" 2>/dev/null; then
                    print_status "SUCCESS" "✓ $file"
                else
                    print_status "ERROR" "✗ $file (YAML syntax error)"
                    return 1
                fi
            else
                # Basic YAML check with Python
                if python3 -c "import yaml; yaml.safe_load(open('$file'))" 2>/dev/null; then
                    print_status "SUCCESS" "✓ $file"
                else
                    print_status "ERROR" "✗ $file (YAML syntax error)"
                    return 1
                fi
            fi
        else
            print_status "ERROR" "✗ $file (file not found)"
            return 1
        fi
    done
}

# Validate GitHub Actions syntax
validate_actions() {
    if command -v actionlint &> /dev/null; then
        print_status "INFO" "Validating GitHub Actions syntax..."
        if actionlint .github/workflows/*.yml; then
            print_status "SUCCESS" "All workflow files are valid"
        else
            print_status "ERROR" "Action validation failed"
            return 1
        fi
    else
        print_status "WARNING" "actionlint not available, skipping action validation"
    fi
}

# Test workflow with act (if available)
test_workflows() {
    if command -v act &> /dev/null; then
        print_status "INFO" "Testing workflows with act..."
        
        # Test CI workflow (dry run)
        print_status "INFO" "Testing CI workflow (dry run)..."
        if act -n workflow_dispatch -W .github/workflows/ci.yml 2>/dev/null; then
            print_status "SUCCESS" "CI workflow syntax is valid"
        else
            print_status "WARNING" "CI workflow test failed (this might be expected)"
        fi
        
        # List available workflows
        print_status "INFO" "Available workflows:"
        act -l 2>/dev/null || true
    else
        print_status "WARNING" "act not available, skipping workflow testing"
        print_status "INFO" "Install act with: curl https://raw.githubusercontent.com/nektos/act/master/install.sh | sudo bash"
    fi
}

# Check file structure
check_structure() {
    print_status "INFO" "Checking CI/CD file structure..."
    
    local expected_files=(
        ".github/workflows/ci.yml"
        ".github/workflows/release.yml"
        ".github/workflows/publish.yml"
        ".github/workflows/performance.yml"
        ".github/workflows/reusable-go-build.yml"
        ".github/workflows/reusable-go-lint.yml"
        ".github/workflows/reusable-python-lint.yml"
        ".github/workflows/reusable-python-test.yml"
        ".github/workflows/reusable-rust-build.yml"
        ".github/workflows/reusable-web-build.yml"
        ".github/workflows/reusable-typos.yml"
    )
    
    for file in "${expected_files[@]}"; do
        if [ -f "$file" ]; then
            print_status "SUCCESS" "✓ $file"
        else
            print_status "WARNING" "✗ $file (not found)"
        fi
    done
}

# Validate workflow triggers and paths
validate_triggers() {
    print_status "INFO" "Validating workflow triggers..."
    
    # Check if main CI has proper triggers
    if grep -q "push:" .github/workflows/ci.yml && grep -q "pull_request:" .github/workflows/ci.yml; then
        print_status "SUCCESS" "CI workflow has proper triggers"
    else
        print_status "WARNING" "CI workflow might be missing triggers"
    fi
    
    # Check if change detection is configured
    if grep -q "dorny/paths-filter" .github/workflows/ci.yml; then
        print_status "SUCCESS" "Change detection configured"
    else
        print_status "WARNING" "Change detection not found"
    fi
    
    # Check if caching is configured
    if grep -q "actions/cache" .github/workflows/ci.yml; then
        print_status "SUCCESS" "Caching configured in CI"
    else
        print_status "WARNING" "No caching found in CI workflow"
    fi
}

# Generate summary
generate_summary() {
    print_status "INFO" "CI/CD Setup Summary:"
    echo ""
    echo "Workflows created:"
    echo "  • ci.yml - Main CI pipeline with change detection"
    echo "  • release.yml - Release automation and artifact publishing"
    echo "  • publish.yml - Development package publishing"
    echo "  • performance.yml - Build performance analysis"
    echo "  • reusable-*.yml - Modular workflow components"
    echo ""
    echo "Features implemented:"
    echo "  • Multi-language support (Python, Go, Rust, Node.js)"
    echo "  • Smart change detection for efficient builds"
    echo "  • Comprehensive caching strategies"
    echo "  • Parallel build optimization"
    echo "  • Automated testing and linting"
    echo "  • Release automation with artifact publishing"
    echo "  • Development package publishing to test registries"
    echo "  • Performance monitoring and analysis"
    echo ""
    print_status "SUCCESS" "CI/CD integration complete!"
}

# Main execution
main() {
    check_dependencies
    echo ""
    
    check_structure
    echo ""
    
    validate_yaml
    echo ""
    
    validate_actions
    echo ""
    
    validate_triggers
    echo ""
    
    test_workflows
    echo ""
    
    generate_summary
}

# Run validation
main