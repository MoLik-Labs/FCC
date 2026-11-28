"""Installed fork wrappers fail closed rather than running an upstream installer."""

import os
import shutil
import subprocess
import sys
import zipfile
from pathlib import Path

import pytest

ROOT = Path(__file__).resolve().parents[2]
SCRIPTS = ("fcc-update", "fcc-update.cmd")


def fixture_wheel(area: Path) -> Path:
    dist = "free_claude_code-1.0.0"
    wheel = area / f"{dist}-py3-none-any.whl"
    with zipfile.ZipFile(wheel, "w") as archive:
        archive.writestr(
            f"{dist}.dist-info/METADATA",
            "Metadata-Version: 2.1\nName: free-claude-code\nVersion: 1.0.0\n",
        )
        archive.writestr(
            f"{dist}.dist-info/WHEEL",
            "Wheel-Version: 1.0\nRoot-Is-Purelib: true\nTag: py3-none-any\n",
        )
        for name in SCRIPTS:
            source = (ROOT / "scripts/update" / name).read_bytes()
            archive.writestr(f"{dist}.data/scripts/{name}", source)
        archive.writestr("free_claude_code/__init__.py", "")
        archive.writestr(
            f"{dist}.dist-info/RECORD",
            "".join(f"{name},,\n" for name in archive.namelist())
            + f"{dist}.dist-info/RECORD,,\n",
        )
    return wheel


def test_installed_updater_stops_without_contacting_installer(tmp_path: Path) -> None:
    uv = shutil.which("uv")
    if uv is None:
        pytest.skip("uv is required for native script lifecycle coverage")
    area = tmp_path / "José 测试 & (spaces) ! apostrophe'"
    area.mkdir()
    env = {
        **os.environ,
        "UV_TOOL_DIR": str(area / "tools"),
        "UV_TOOL_BIN_DIR": str(area / "bin"),
        "UV_CACHE_DIR": str(area / "cache"),
        "UV_NO_CONFIG": "1",
        "UV_PYTHON_DOWNLOADS": "never",
    }
    marker = area / "network-command-was-run"
    env["FCC_TEST_NETWORK_MARKER"] = str(marker)
    if os.name == "nt":
        stub = area / "powershell.cmd"
        stub.write_text('@echo off\r\necho called > "%FCC_TEST_NETWORK_MARKER%"\r\n')
    else:
        stub = area / "curl"
        stub.write_text('#!/bin/sh\nprintf called > "$FCC_TEST_NETWORK_MARKER"\n')
        stub.chmod(0o755)
    env["PATH"] = str(area) + os.pathsep + env.get("PATH", "")
    wheel = fixture_wheel(area)
    installed = subprocess.run(
        [
            uv,
            "tool",
            "install",
            "--python",
            getattr(sys, "_base_executable", sys.executable),
            str(wheel),
        ],
        env=env,
        capture_output=True,
        text=True,
        timeout=90,
        check=False,
    )
    assert installed.returncode == 0, installed.stderr
    try:
        name = "fcc-update.cmd" if os.name == "nt" else "fcc-update"
        launcher = area / "bin" / name
        command = (
            f'"{os.environ.get("COMSPEC", "cmd.exe")}" /d /s /c ""{launcher}""'
            if os.name == "nt"
            else [str(launcher), "--voice-local"]
        )
        # A regression must not reach a real network command. The installed
        # updater itself needs no PATH tools; only the harmless stub is visible.
        run_env = {**env, "PATH": str(area)}
        result = subprocess.run(
            command,
            env=run_env,
            cwd=area,
            capture_output=True,
            text=True,
            timeout=30,
            check=False,
        )
        assert result.returncode == 2, result.stdout + result.stderr
        assert "git pull --ff-only" in result.stderr
        assert "uv sync --locked" in result.stderr
        assert not marker.exists()
    finally:
        removed = subprocess.run(
            [uv, "tool", "uninstall", "free-claude-code"],
            env=env,
            capture_output=True,
            text=True,
            timeout=30,
            check=False,
        )
        assert removed.returncode == 0, removed.stderr
    assert not any((area / "bin" / name).exists() for name in SCRIPTS)
