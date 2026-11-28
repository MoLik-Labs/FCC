"""Disabling live smoke must not silently skip unrelated unit tests."""

import os
import subprocess
import sys
from pathlib import Path

ROOT = Path(__file__).resolve().parents[2]


def test_mixed_smoke_and_unit_collection_only_skips_smoke(tmp_path: Path) -> None:
    # The smoke plugin loads managed settings during collection. Keep the
    # subprocess away from the developer's real home and optional config file.
    env = dict(os.environ)
    env.pop("FCC_LIVE_SMOKE", None)
    env.pop("FCC_ENV_FILE", None)
    env.pop("PYTEST_ADDOPTS", None)
    env["HOME"] = env["USERPROFILE"] = str(tmp_path)
    result = subprocess.run(
        [
            sys.executable,
            "-m",
            "pytest",
            "-n",
            "0",
            "-q",
            "tests/scripts/test_fork_updater.py",
            "smoke/installers/test_update_wrapper.py",
        ],
        cwd=ROOT,
        env=env,
        capture_output=True,
        text=True,
        timeout=60,
        check=False,
    )
    assert result.returncode == 0, result.stdout + result.stderr
    if os.name == "nt":
        assert "1 passed, 4 skipped" in result.stdout
    else:
        assert "4 passed, 1 skipped" in result.stdout
