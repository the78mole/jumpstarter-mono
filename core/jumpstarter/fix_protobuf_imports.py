#!/usr/bin/env python3
"""
Fix protobuf imports in generated files to use correct module paths.

This script fixes imports like 'from jumpstarter.v1 import' to use the
correct module path 'from jumpstarter_protocol.jumpstarter.v1 import'.
"""

import re
from pathlib import Path


def fix_protobuf_imports():
    """Fix imports in generated protobuf files."""
    protocol_dir = Path("packages/jumpstarter-protocol/jumpstarter_protocol")

    # Pattern to match absolute imports that need to be fixed
    import_pattern = re.compile(r"^from jumpstarter\.([a-zA-Z0-9_.]+) import (.+)$", re.MULTILINE)

    # Find all Python files in the protocol directory
    for py_file in protocol_dir.rglob("*.py"):
        if py_file.name in ("__init__.py", "py.typed"):
            continue

        print(f"Processing {py_file}")

        # Read the file
        content = py_file.read_text()

        # Check if there are any problematic imports
        if "from jumpstarter." in content:
            # Fix the imports
            def replace_import(match):
                module_path = match.group(1)
                imports = match.group(2)
                return f"from jumpstarter_protocol.jumpstarter.{module_path} import {imports}"

            # Apply the fix
            fixed_content = import_pattern.sub(replace_import, content)

            # Write back if changes were made
            if fixed_content != content:
                py_file.write_text(fixed_content)
                print(f"  Fixed imports in {py_file}")
            else:
                print(f"  No changes needed in {py_file}")


if __name__ == "__main__":
    fix_protobuf_imports()
