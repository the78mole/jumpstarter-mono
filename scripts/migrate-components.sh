#!/bin/bash
# migrate-components.sh
# Script to migrate components from individual repositories into the monorepo

set -e

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

log() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

warn() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Configuration
TEMP_DIR="./temp-migration"
GITHUB_ORG="jumpstarter-dev"

# Component mappings: "source_repo:target_path"
declare -A COMPONENTS=(
    ["jumpstarter"]="core/jumpstarter"
    ["jumpstarter-controller"]="core/controller"
    ["jumpstarter-protocol"]="core/protocol"
    ["dutlink-firmware"]="hardware/dutlink-firmware"
    ["dutlink-board"]="hardware/dutlink-board"
    ["jumpstarter-tekton-tasks"]="integrations/tekton"
    ["vscode-jumpstarter"]="integrations/vscode"
    ["jumpstarter-devspace"]="integrations/devspace"
    ["jumpstarter-driver-template"]="templates/driver"
    ["jumpstarter-e2e"]="testing/e2e"
    ["jumpstarter-lab-config"]="lab-config"
    ["packages"]="packages/repository-tools"
)

migrate_component() {
    local repo_name=$1
    local target_path=$2
    
    log "Migrating ${repo_name} to ${target_path}..."
    
    # Create temporary directory
    mkdir -p "${TEMP_DIR}"
    
    # Clone repository
    log "Cloning ${GITHUB_ORG}/${repo_name}..."
    if git clone "https://github.com/${GITHUB_ORG}/${repo_name}.git" "${TEMP_DIR}/${repo_name}"; then
        # Create target directory
        mkdir -p "${target_path}"
        
        # Copy files (excluding .git)
        log "Copying files to ${target_path}..."
        rsync -av --exclude='.git' "${TEMP_DIR}/${repo_name}/" "${target_path}/"
        
        # Clean up temporary directory
        rm -rf "${TEMP_DIR}/${repo_name}"
        
        log "✓ Successfully migrated ${repo_name}"
    else
        warn "Failed to clone ${repo_name}, may not exist or be accessible"
        rm -rf "${TEMP_DIR}/${repo_name}" 2>/dev/null || true
    fi
}

update_go_modules() {
    log "Updating Go module paths..."
    
    # Update controller module
    if [ -f "core/controller/go.mod" ]; then
        cd core/controller
        go mod edit -module github.com/the78mole/jumpstarter-mono/core/controller
        go mod tidy || true
        cd ../..
    fi
    
    # Update lab-config module
    if [ -f "lab-config/go.mod" ]; then
        cd lab-config
        go mod edit -module github.com/the78mole/jumpstarter-mono/lab-config
        go mod tidy || true
        cd ..
    fi
    
    log "✓ Go modules updated"
}

update_python_configs() {
    log "Updating Python configurations..."
    
    # Update jumpstarter pyproject.toml
    if [ -f "core/jumpstarter/pyproject.toml" ]; then
        sed -i 's|name = "jumpstarter"|name = "jumpstarter-core"|g' core/jumpstarter/pyproject.toml
    fi
    
    # Update driver template pyproject.toml
    if [ -f "templates/driver/pyproject.toml" ]; then
        sed -i 's|name = .*|name = "jumpstarter-driver-template"|g' templates/driver/pyproject.toml
    fi
    
    log "✓ Python configurations updated"
}

update_imports() {
    log "Updating import paths..."
    
    # Update Python imports (basic pattern matching)
    find core/jumpstarter -name "*.py" -type f -exec sed -i 's|from jumpstarter|from jumpstarter_core|g' {} \; 2>/dev/null || true
    find templates/driver -name "*.py" -type f -exec sed -i 's|from jumpstarter|from jumpstarter_core|g' {} \; 2>/dev/null || true
    
    # Update Go imports
    find core/controller -name "*.go" -type f -exec sed -i 's|github.com/jumpstarter-dev/jumpstarter-controller|github.com/the78mole/jumpstarter-mono/core/controller|g' {} \; 2>/dev/null || true
    find lab-config -name "*.go" -type f -exec sed -i 's|github.com/jumpstarter-dev/jumpstarter-lab-config|github.com/the78mole/jumpstarter-mono/lab-config|g' {} \; 2>/dev/null || true
    
    log "✓ Import paths updated"
}

# Main migration process
main() {
    log "Starting Jumpstarter monorepo migration..."
    
    # Migrate all components
    for repo in "${!COMPONENTS[@]}"; do
        migrate_component "$repo" "${COMPONENTS[$repo]}"
    done
    
    # Update configurations
    update_go_modules
    update_python_configs
    update_imports
    
    # Clean up
    rm -rf "${TEMP_DIR}"
    
    log "✓ Migration completed successfully!"
    warn "Please review the migrated files and update any remaining references manually."
    warn "Don't forget to test builds: make build"
}

# Run migration
main "$@"