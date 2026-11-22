"""Tests for the standalone documentation build and local-link validator."""

import tempfile
import unittest
from pathlib import Path

from scripts.build_docs_site import build, validate_site


class DocumentationSiteTests(unittest.TestCase):
    def test_build_contains_site_and_screenshots(self) -> None:
        with tempfile.TemporaryDirectory() as temp:
            output = Path(temp) / "site"
            build(output)
            self.assertTrue((output / "index.html").is_file())
            self.assertTrue((output / "vi.html").is_file())
            self.assertTrue((output / ".nojekyll").is_file())
            self.assertTrue((output / "assets" / "admin-page.png").is_file())
            self.assertTrue((output / "assets" / "browser-code-session.png").is_file())
            validate_site(output)

    def test_quick_starts_use_the_fork_and_warn_about_upstream_updater(self) -> None:
        root = Path(__file__).resolve().parents[2]
        readme = (root / "README.md").read_text(encoding="utf-8")
        quick_start = readme.split("## Quick Start", 1)[1].split(
            "## Choose A Provider", 1
        )[0]
        self.assertIn("https://github.com/MoLik-Labs/FCC.git", quick_start)
        self.assertIn("uv sync --locked", quick_start)
        self.assertIn("uv run fcc-server", quick_start)
        self.assertNotIn("raw.githubusercontent.com/Alishahryar1", quick_start)
        for name in ("index.html", "vi.html"):
            with self.subTest(page=name):
                page = (root / "docs" / "site" / name).read_text(encoding="utf-8")
                quick_start = page.split('id="quickstart"', 1)[1].split(
                    'id="first-request"', 1
                )[0]
                self.assertIn("https://github.com/MoLik-Labs/FCC.git", quick_start)
                self.assertIn("uv sync --locked", quick_start)
                self.assertIn("uv run fcc-server", quick_start)
                self.assertNotIn("raw.githubusercontent.com/Alishahryar1", quick_start)
                maintenance = page.split('id="maintenance"', 1)[1].split(
                    'id="troubleshooting"', 1
                )[0]
                self.assertIn("fcc-update", maintenance)
                self.assertIn("upstream", maintenance)

    def test_missing_target_fragment_duplicate_id_and_alt_are_rejected(self) -> None:
        examples = (
            '<a href="missing.html">missing</a>',
            '<a href="#unknown">missing fragment</a>',
            '<p id="same"></p><p id="same"></p>',
            '<img src="img.png">',
            '<a href="../outside">escape</a>',
        )
        for html in examples:
            with self.subTest(html=html), tempfile.TemporaryDirectory() as temp:
                output = Path(temp)
                (output / "index.html").write_text(html, encoding="utf-8")
                with self.assertRaises(ValueError):
                    validate_site(output)

    def test_nonempty_output_is_never_overwritten(self) -> None:
        with tempfile.TemporaryDirectory() as temp:
            output = Path(temp)
            existing = output / "keep.txt"
            existing.write_text("keep", encoding="utf-8")
            with self.assertRaises(ValueError):
                build(output)
            self.assertEqual(existing.read_text(encoding="utf-8"), "keep")


if __name__ == "__main__":
    unittest.main()
