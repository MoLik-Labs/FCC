"""Tests for the standalone documentation build and local-link validator."""

import tempfile
import unittest
from pathlib import Path

from scripts.build_docs_site import build, validate_site


class DocumentationSiteTests(unittest.TestCase):
    def test_build_contains_site_without_stale_screenshots(self) -> None:
        with tempfile.TemporaryDirectory() as temp:
            output = Path(temp) / "site"
            build(output)
            self.assertTrue((output / "index.html").is_file())
            self.assertTrue((output / "vi.html").is_file())
            self.assertTrue((output / "molik.html").is_file())
            self.assertTrue((output / ".nojekyll").is_file())
            self.assertFalse((output / "assets" / "admin-page.png").exists())
            self.assertFalse((output / "assets" / "browser-code-session.png").exists())
            self.assertTrue((output / "assets" / "molik-mark.svg").is_file())
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

    def test_workload_guidance_and_fork_only_commands(self) -> None:
        root = Path(__file__).resolve().parents[2]
        readme = (root / "README.md").read_text(encoding="utf-8")
        self.assertNotIn("raw.githubusercontent.com/Alishahryar1", readme)
        for name in ("index.html", "vi.html"):
            with self.subTest(page=name):
                page = (root / "docs" / "site" / name).read_text(encoding="utf-8")
                self.assertIn('id="workloads"', page)
                self.assertIn("PROVIDER_MAX_CONCURRENCY", page)
                self.assertIn("Catppuccin Mocha", page)
                self.assertIn("Nord", page)
                self.assertIn("Dracula", page)

    def test_bilingual_guides_offer_reading_themes(self) -> None:
        root = Path(__file__).resolve().parents[2] / "docs" / "site"
        for name in ("index.html", "vi.html"):
            with self.subTest(page=name):
                page = (root / name).read_text(encoding="utf-8")
                self.assertIn('id="reading-theme"', page)
                for value in ("light", "forest", "catppuccin-mocha"):
                    self.assertIn(f'value="{value}"', page)

    def test_molik_gallery_tracks_the_rust_theme_tokens(self) -> None:
        import re

        root = Path(__file__).resolve().parents[2]
        rust = (root / "ohmymolik/rust/src/theme.rs").read_text(encoding="utf-8")
        gallery = (root / "docs/site/molik.html").read_text(encoding="utf-8")
        blocks = list(
            re.finditer(r'Theme \{\s*name: "([a-z-]+)",(.*?)\n    \}', rust, re.S)
        )
        self.assertEqual(len(blocks), 5)
        for block in blocks:
            name = block.group(1)
            with self.subTest(theme=name):
                self.assertIn(f'id="{name}"', gallery)
                for match in re.finditer(r'([a-z]+): "(#[A-F0-9]{6})"', block.group(2)):
                    self.assertIn(
                        f"{match.group(1)} <code>{match.group(2)}</code>", gallery
                    )

    def test_mcp_setup_is_linked_and_describes_its_boundary(self) -> None:
        root = Path(__file__).resolve().parents[2] / "docs" / "site"
        for name in ("index.html", "vi.html"):
            with self.subTest(page=name):
                self.assertIn(
                    'href="molik.html#mcp-stdio"',
                    (root / name).read_text(encoding="utf-8"),
                )
        guide = (root / "molik.html").read_text(encoding="utf-8")
        self.assertIn('id="mcp-stdio"', guide)
        self.assertIn("2025-06-18", guide)
        self.assertIn("not a sandbox", guide)
        self.assertIn('"mcp", "serve", "--workspace"', guide)

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
