"""The documentation source must be an actual bilingual Starlight project."""

import json
from pathlib import Path
from typing import TypedDict, cast


class Package(TypedDict):
    dependencies: dict[str, str]


ROOT = Path(__file__).resolve().parents[2]
SITE = ROOT / "docs/site"


def test_starlight_project_pins_dependencies_and_pages_base() -> None:
    package = cast(
        Package, json.loads((SITE / "package.json").read_text(encoding="utf-8"))
    )
    assert "@astrojs/starlight" in package["dependencies"]
    assert "@cloudflare/nimbus-docs" in package["dependencies"]
    assert "astro" in package["dependencies"]
    assert (SITE / "package-lock.json").is_file()
    config = (SITE / "astro.config.mjs").read_text(encoding="utf-8")
    assert "@astrojs/starlight" in config
    assert "@cloudflare/nimbus-docs/markdown" in config
    assert "base" in config
    assert "vi" in config


def test_bilingual_guides_retain_core_operational_details() -> None:
    for name in ("index.md", "vi/index.md"):
        page = (SITE / "src/content/docs" / name).read_text(encoding="utf-8")
        for value in (
            "MoLik-Labs/FCC.git",
            "uv sync --locked",
            "uv run fcc-server",
            "PROVIDER_MAX_CONCURRENCY",
            "fcc-update",
            "/v1/messages",
            "/v1/chat/completions",
            "Catppuccin Mocha",
        ):
            assert value in page, (name, value)


def test_all_legacy_section_anchors_are_present_in_bilingual_guides() -> None:
    from scripts.build_docs_site import PageParser

    for name in ("index.html", "vi.html"):
        old = PageParser()
        old.feed((SITE / name).read_text(encoding="utf-8"))
        new_name = "index.md" if name == "index.html" else "vi/index.md"
        new = (SITE / "src/content/docs" / new_name).read_text(encoding="utf-8")
        for anchor in (
            "overview",
            "quickstart",
            "first-request",
            "providers",
            "clients",
            "api",
            "messaging",
            "workloads",
            "security",
            "maintenance",
            "troubleshooting",
            "development",
            "pages",
            "resources",
        ):
            assert anchor in old.ids
            assert f'<a id="{anchor}"></a>' in new
