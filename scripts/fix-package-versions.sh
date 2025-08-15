#!/bin/bash
# fix-package-versions.sh
# Fix VCS version issues by setting static versions

set -e

log() {
    echo -e "\033[32m[INFO]\033[0m $1"
}

cd /home/runner/work/jumpstarter-mono/jumpstarter-mono

log "Fixing package versions to use static versioning..."

# Find all pyproject.toml files with dynamic version
find . -name "pyproject.toml" -exec grep -l 'dynamic.*version' {} \; | while read -r file; do
    log "Fixing version in $file"

    # Remove dynamic version and add static version
    sed -i 's/dynamic = \["version", "urls"\]/version = "0.1.0"/' "$file"
    sed -i 's/dynamic = \["version"\]/version = "0.1.0"/' "$file"

    # Remove [tool.hatch.version] sections that depend on VCS
    sed -i '/\[tool\.hatch\.version\]/,/^$/d' "$file"

    # Remove [tool.hatch.metadata] sections that depend on VCS
    sed -i '/\[tool\.hatch\.metadata\]/,/^$/d' "$file"
done

log "✓ Package versions fixed"
