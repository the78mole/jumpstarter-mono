# GitHub Actions Cleanup Scripts

This directory contains utility scripts for managing GitHub Actions workflows and runs.

## cleanup-runs.sh

A script to clean up old GitHub Actions workflow runs, keeping only the most recent ones.

### Features

- 🧹 **Smart cleanup**: Keeps the N most recent runs, deletes the rest
- 📊 **Configurable**: Default keeps 20 runs, but accepts any number as argument
- 🔒 **Safe**: Asks for confirmation before deleting anything (unless `-y` flag is used)
- ⚡ **Automated**: Support for `-y` flag to skip confirmation for automated usage
- ⚡ **Efficient**: Deletes runs in batches with API-friendly delays
- 📈 **Informative**: Shows detailed progress and summary

### Usage

```bash
# Keep default 20 most recent runs (with confirmation)
./.github/cleanup-runs.sh

# Keep default 20 most recent runs (skip confirmation)
./.github/cleanup-runs.sh -y

# Keep 50 most recent runs (with confirmation)
./.github/cleanup-runs.sh 50

# Keep 50 most recent runs (skip confirmation)
./.github/cleanup-runs.sh -y 50

# Keep only 10 most recent runs (skip confirmation)
./.github/cleanup-runs.sh -y 10
```

### Requirements

- [GitHub CLI (gh)](https://cli.github.com/) installed and authenticated
- Appropriate repository permissions to delete workflow runs

### Examples

```bash
# Dry run - see what would be deleted without actually deleting
./.github/cleanup-runs.sh 20

# Automated cleanup without confirmation
./.github/cleanup-runs.sh -y 20

# Example output:
🧹 GitHub Actions Run Cleanup
📊 Keeping the 20 most recent runs...
📈 Analyzing workflow runs...
   Found 45 total runs
🗑️  Will delete 25 old runs (keeping newest 20)
⚡ Auto-confirming deletion (--yes flag used)
```

### Safety Features

- ✅ Validates input parameters
- ✅ Checks for required tools (gh CLI)
- ✅ Verifies GitHub authentication
- ✅ Asks for confirmation before deletion
- ✅ Handles API rate limits with delays
- ✅ Reports success/failure for each deletion

### Tips

- Run this script regularly as part of repository maintenance
- Consider adding it to a scheduled workflow for automatic cleanup
- Use a higher number (e.g., 50-100) for repositories with important CI history
- Use a lower number (e.g., 10-15) for repositories with frequent runs

### Integration with GitHub Actions

You can also run this script automatically using a scheduled GitHub Action:

```yaml
name: Cleanup Old Runs
on:
  schedule:
    - cron: "0 2 * * 0" # Weekly on Sunday at 2 AM
  workflow_dispatch:

jobs:
  cleanup:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - name: Cleanup old runs
        run: ./.github/cleanup-runs.sh -y 30
        env:
          GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}
```
