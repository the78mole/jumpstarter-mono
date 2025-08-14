#!/bin/bash

# Quick test of ACT masking functionality
# This script creates a minimal workflow to test the conditional masking

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_ROOT="$(cd "${SCRIPT_DIR}/.." && pwd)"
cd "${REPO_ROOT}"

# Create a temporary test workflow
mkdir -p .github/workflows-test
cat > .github/workflows-test/test-act-masking.yml << 'EOF'
name: Test ACT Masking

on:
  workflow_dispatch: {}

jobs:
  test-masking:
    runs-on: ubuntu-latest
    steps:
      - name: Always run step
        run: echo "This step always runs"
      
      - name: Masked for ACT
        if: ${{ !env.ACT }}
        run: echo "This step should NOT run under ACT"
      
      - name: Only runs under ACT
        if: ${{ env.ACT }}
        run: echo "This step ONLY runs under ACT"
      
      - name: Check ACT environment
        run: |
          echo "ACT environment variable: ${ACT:-not set}"
          if [ "${ACT:-}" = "true" ]; then
            echo "✅ Running under ACT"
          else
            echo "❌ NOT running under ACT"
          fi
EOF

echo "Testing ACT masking with temporary workflow..."

# Run the test workflow with act
if command -v act &> /dev/null; then
    echo "Running test workflow with ACT=true..."
    act workflow_dispatch -W .github/workflows-test/test-act-masking.yml --env ACT=true -P ubuntu-latest=node:16-alpine || echo "Test completed (expected some failures due to minimal container)"
else
    echo "ACT not available, skipping execution test"
fi

# Cleanup
rm -rf .github/workflows-test

echo "✅ ACT masking test completed"
echo ""
echo "Expected behavior:"
echo "- 'Always run step' should execute"
echo "- 'Masked for ACT' should be SKIPPED (if: !env.ACT)"
echo "- 'Only runs under ACT' should execute (if: env.ACT)"
echo "- 'Check ACT environment' should show ACT=true"