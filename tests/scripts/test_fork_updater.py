"""The fork's update wrapper must never run the upstream installer."""

import os
import subprocess
from pathlib import Path

import pytest

ROOT = Path(__file__).resolve().parents[2]


@pytest.mark.skipif(os.name == "nt", reason="POSIX shell wrapper")
@pytest.mark.parametrize(
    "arguments", [(), ("--voice-local",), ("$(touch unexpected)",)]
)
def test_posix_update_fails_closed_without_executing_download(
    tmp_path: Path, arguments: tuple[str, ...]
) -> None:
    stub = tmp_path / "curl"
    stub.write_text("#!/bin/sh\nprintf invoked > \"$MARKER\"\nprintf 'exit 0\\n'\n")
    stub.chmod(0o755)
    marker = tmp_path / "network-called"
    result = subprocess.run(
        ["/bin/sh", str(ROOT / "scripts/update/fcc-update"), *arguments],
        cwd=tmp_path,
        env={**os.environ, "PATH": str(tmp_path), "MARKER": str(marker)},
        capture_output=True,
        text=True,
        timeout=10,
        check=False,
    )
    assert result.returncode != 0
    assert not marker.exists()
    assert not (tmp_path / "unexpected").exists()
    assert "git pull --ff-only" in result.stderr
    assert "uv sync --locked" in result.stderr
    assert "upstream" in result.stderr.lower()


def test_windows_wrapper_has_no_network_or_dynamic_execution() -> None:
    script = (ROOT / "scripts/update/fcc-update.cmd").read_text().lower()
    assert "git pull --ff-only" in script
    assert "uv sync --locked" in script
    assert "exit /b 2" in script
    for forbidden in (
        "powershell",
        "invoke-restmethod",
        "curl",
        "raw.githubusercontent",
        "%*",
    ):
        assert forbidden not in script
