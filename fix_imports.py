#!/usr/bin/env python3

import re
import sys
from pathlib import Path


def fix_imports_in_file(filepath):
    """Fix jumpstarter_core* imports in a Python file."""
    try:
        with open(filepath, "r", encoding="utf-8") as f:
            content = f.read()

        original_content = content

        # Define replacement patterns
        patterns = [
            # Basic jumpstarter_core -> jumpstarter
            (r"from jumpstarter_core\.", "from jumpstarter."),
            (r"import jumpstarter_core\.", "import jumpstarter."),
            (r"import jumpstarter_core\b", "import jumpstarter"),
            # Protocol package
            (r"from jumpstarter_core_protocol\b", "from jumpstarter_protocol"),
            (r"import jumpstarter_core_protocol\b", "import jumpstarter_protocol"),
            # Driver packages
            (r"from jumpstarter_core_driver_", "from jumpstarter_driver_"),
            (r"import jumpstarter_core_driver_", "import jumpstarter_driver_"),
            # CLI packages
            (r"from jumpstarter_core_cli_", "from jumpstarter_cli_"),
            (r"import jumpstarter_core_cli_", "import jumpstarter_cli_"),
            # Imagehash package
            (r"from jumpstarter_core_imagehash\b", "from jumpstarter_imagehash"),
            (r"import jumpstarter_core_imagehash\b", "import jumpstarter_imagehash"),
            # Testing package
            (r"from jumpstarter_core_testing\b", "from jumpstarter_testing"),
            (r"import jumpstarter_core_testing\b", "import jumpstarter_testing"),
            # Kubernetes package
            (r"from jumpstarter_core_kubernetes\b", "from jumpstarter_kubernetes"),
            (r"import jumpstarter_core_kubernetes\b", "import jumpstarter_kubernetes"),
        ]

        # Apply replacements
        for pattern, replacement in patterns:
            content = re.sub(pattern, replacement, content)

        # Only write if content changed
        if content != original_content:
            with open(filepath, "w", encoding="utf-8") as f:
                f.write(content)
            print(f"Fixed: {filepath}")
            return True
        return False

    except Exception as e:
        print(f"Error processing {filepath}: {e}")
        return False


def main():
    if len(sys.argv) > 1:
        root_dir = sys.argv[1]
    else:
        root_dir = "core/jumpstarter"

    root_path = Path(root_dir)
    if not root_path.exists():
        print(f"Directory {root_dir} does not exist")
        return

    fixed_count = 0
    total_count = 0

    # Find all Python files
    for py_file in root_path.rglob("*.py"):
        total_count += 1
        if fix_imports_in_file(py_file):
            fixed_count += 1

    print(f"\nProcessed {total_count} Python files")
    print(f"Fixed imports in {fixed_count} files")


if __name__ == "__main__":
    main()
