#!/bin/bash
#
# cleanup-orphaned-releases.sh - Clean up releases that point to non-existent commits
#
# This script identifies GitHub releases that reference commits that no longer exist
# in the repository (e.g., after squashing commits) and optionally deletes them.
#
# Usage:
#   ./cleanup-orphaned-releases.sh [--dry-run] [--force]
#
# Options:
#   --dry-run    Show what would be deleted without actually deleting
#   --force      Delete without confirmation prompts
#

set -euo pipefail

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Options
DRY_RUN=false
FORCE=false

# Parse command line arguments
while [[ $# -gt 0 ]]; do
    case $1 in
        -n|--dry-run)
            DRY_RUN=true
            shift
            ;;
        -f|--force)
            FORCE=true
            shift
            ;;
        -h|--help)
            echo "Usage: $0 [-n|--dry-run] [-f|--force]"
            echo ""
            echo "Clean up GitHub releases that point to non-existent commits"
            echo ""
            echo "Options:"
            echo "  -n, --dry-run   Show what would be deleted without actually deleting"
            echo "  -f, --force     Delete without confirmation prompts"
            echo "  -h, --help      Show this help message"
            exit 0
            ;;
        *)
            echo "Unknown option $1"
            exit 1
            ;;
    esac
done

# Check if gh CLI is available
if ! command -v gh &> /dev/null; then
    echo -e "${RED}Error: GitHub CLI (gh) is not installed or not in PATH${NC}"
    exit 1
fi

# Check if we're authenticated with GitHub
if ! gh auth status &> /dev/null; then
    echo -e "${RED}Error: Not authenticated with GitHub. Run 'gh auth login' first.${NC}"
    exit 1
fi

# Check if we're in a git repository
if ! git rev-parse --git-dir &> /dev/null; then
    echo -e "${RED}Error: Not in a git repository${NC}"
    exit 1
fi

echo -e "${BLUE}🔍 Checking for orphaned releases...${NC}"
echo

# Get all releases and their target commits
releases_json=$(gh release list --limit 1000 --json tagName,name,createdAt)

if [[ "$releases_json" == "[]" ]]; then
    echo -e "${GREEN}✅ No releases found${NC}"
    exit 0
fi

orphaned_releases_file=$(mktemp)
valid_releases_file=$(mktemp)

# Cleanup temp files on exit
trap 'rm -f "$orphaned_releases_file" "$valid_releases_file"' EXIT

# Process each release
while IFS= read -r tag_name; do
    # Get release details
    release_details=$(echo "$releases_json" | jq -r --arg tag "$tag_name" '.[] | select(.tagName == $tag)')
    release_name=$(echo "$release_details" | jq -r '.name')
    created_at=$(echo "$release_details" | jq -r '.createdAt')

    # Get the commit that the tag points to
    if target_commit=$(git rev-list -n 1 "$tag_name" 2>/dev/null); then
        # Check if the target commit is reachable from HEAD (main branch)
        if git merge-base --is-ancestor "$target_commit" HEAD 2>/dev/null; then
            echo -e "${GREEN}✅ Valid:${NC} $tag_name (commit: ${target_commit:0:8}) - $release_name"
            echo -e "   ${GREEN}Reason:${NC} Commit reachable from main branch"
            echo "$tag_name" >> "$valid_releases_file"
        else
            # Check if the commit is reachable from any remote branch
            reachable_branches=$(git branch -r --contains "$target_commit" 2>/dev/null | head -5 | tr '\n' ' ' | sed 's/^ *//')
            if [ -n "$reachable_branches" ]; then
                echo -e "${YELLOW}⚠️  Valid (on branch):${NC} $tag_name (commit: ${target_commit:0:8}) - $release_name"
                echo -e "   ${YELLOW}Created:${NC} $created_at"
                echo -e "   ${YELLOW}Reason:${NC} Commit reachable from branches: $reachable_branches"
                echo "$tag_name" >> "$valid_releases_file"
            else
                # Check if the commit object exists but is not reachable from any branch
                if git cat-file -e "$target_commit" 2>/dev/null; then
                    echo -e "${RED}❌ Orphaned:${NC} $tag_name (commit: ${target_commit:0:8}) - $release_name"
                    echo -e "   ${YELLOW}Created:${NC} $created_at"
                    echo -e "   ${YELLOW}Reason:${NC} Commit exists but not reachable from any branch (likely squashed/rebased)"
                    echo "$tag_name" >> "$orphaned_releases_file"
                else
                    echo -e "${RED}❌ Orphaned:${NC} $tag_name (commit: ${target_commit:0:8}) - $release_name"
                    echo -e "   ${YELLOW}Created:${NC} $created_at"
                    echo -e "   ${YELLOW}Reason:${NC} Commit object no longer exists"
                    echo "$tag_name" >> "$orphaned_releases_file"
                fi
            fi
        fi
    else
        # Tag doesn't exist at all
        echo -e "${RED}❌ Orphaned:${NC} $tag_name (tag missing) - $release_name"
        echo -e "   ${YELLOW}Created:${NC} $created_at"
        echo -e "   ${YELLOW}Reason:${NC} Tag points to non-existent commit"
        echo "$tag_name" >> "$orphaned_releases_file"
    fi
done < <(echo "$releases_json" | jq -r '.[].tagName')

# Read results from temp files
mapfile -t orphaned_releases < "$orphaned_releases_file"
mapfile -t valid_releases < "$valid_releases_file"

# Count results
orphaned_count=${#orphaned_releases[@]}
valid_count=${#valid_releases[@]}

echo
echo -e "${BLUE}📊 Summary:${NC}"
echo -e "  Valid releases: ${GREEN}$valid_count${NC}"
echo -e "  Orphaned releases: ${RED}$orphaned_count${NC}"

if [[ $orphaned_count -eq 0 ]]; then
    echo -e "${GREEN}✅ No orphaned releases found!${NC}"
    exit 0
fi

if [[ "$DRY_RUN" == "true" ]]; then
    echo
    echo -e "${YELLOW}🧪 DRY RUN MODE - Would delete the following releases:${NC}"
    for tag in "${orphaned_releases[@]}"; do
        echo -e "  ${RED}•${NC} $tag"
    done
    echo
    echo -e "${YELLOW}Run without --dry-run to actually delete these releases${NC}"
    exit 0
fi

echo
echo -e "${RED}⚠️  WARNING: This will permanently delete $orphaned_count orphaned releases!${NC}"
echo

if [[ "$FORCE" != "true" ]]; then
    echo -e "${YELLOW}Releases to be deleted:${NC}"
    for tag in "${orphaned_releases[@]}"; do
        echo -e "  ${RED}•${NC} $tag"
    done
    echo

    read -r -p "Do you want to proceed with deletion? (yes/no): " confirm
    if [[ "$confirm" != "yes" ]]; then
        echo -e "${YELLOW}❌ Aborted by user${NC}"
        exit 0
    fi
fi

echo
echo -e "${BLUE}🗑️  Deleting orphaned releases...${NC}"

# Temporarily disable the trap to prevent premature temp file cleanup during deletion loop
trap '' EXIT

deleted_count=0
failed_count=0

for tag in "${orphaned_releases[@]}"; do
    echo -n -e "  Deleting $tag... "

    # Temporarily disable exit on error for this command
    set +e
    gh release delete "$tag" --yes 2>/dev/null
    delete_result=$?
    set -e

    if [[ $delete_result -eq 0 ]]; then
        echo -e "${GREEN}✅ Done${NC}"
        ((deleted_count++))
    else
        echo -e "${RED}❌ Failed${NC}"
        ((failed_count++))
    fi
done

# Re-enable cleanup trap
trap 'rm -f "$orphaned_releases_file" "$valid_releases_file"' EXIT

echo
echo -e "${BLUE}📊 Deletion Summary:${NC}"
echo -e "  Successfully deleted: ${GREEN}$deleted_count${NC}"
echo -e "  Failed to delete: ${RED}$failed_count${NC}"

if [[ $failed_count -eq 0 ]]; then
    echo -e "${GREEN}✅ All orphaned releases have been successfully cleaned up!${NC}"
else
    echo -e "${YELLOW}⚠️  Some releases could not be deleted. Check permissions and try again.${NC}"
    exit 1
fi
