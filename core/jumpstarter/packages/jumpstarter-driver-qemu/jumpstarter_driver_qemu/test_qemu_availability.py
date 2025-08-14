"""Simple smoke test to verify QEMU is available in CI."""

import subprocess
import pytest


def test_qemu_commands_available():
    """Test that essential QEMU commands are available."""
    qemu_commands = ["qemu-img", "qemu-system-x86_64", "qemu-system-aarch64"]

    for cmd in qemu_commands:
        try:
            result = subprocess.run([cmd, "--version"], capture_output=True, text=True, timeout=5)
            assert result.returncode == 0, f"{cmd} is not working properly"
            assert "QEMU" in result.stdout or "qemu" in result.stdout, f"{cmd} doesn't appear to be QEMU"
        except (subprocess.TimeoutExpired, FileNotFoundError) as e:
            pytest.fail(f"{cmd} is not available or not working: {e}")


def test_qemu_version_info():
    """Test that we can get QEMU version information."""
    try:
        result = subprocess.run(["qemu-system-x86_64", "--version"], capture_output=True, text=True, timeout=5)
        assert result.returncode == 0
        version_output = result.stdout.strip()
        print(f"QEMU version: {version_output}")
        # Just verify we get some version output
        assert len(version_output) > 0
    except Exception as e:
        pytest.fail(f"Failed to get QEMU version: {e}")
