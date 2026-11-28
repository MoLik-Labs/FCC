"""Build and validate the Astro Starlight GitHub Pages documentation site."""

import argparse
import os
import shutil
import subprocess
from collections.abc import Callable
from html.parser import HTMLParser
from pathlib import Path
from typing import TYPE_CHECKING, cast
from urllib.parse import unquote, urlsplit

if TYPE_CHECKING:
    from typing import override
else:
    override = cast(Callable[[object], object], lambda method: method)


ROOT = Path(__file__).resolve().parent.parent
SITE = ROOT / "docs" / "site"


class PageParser(HTMLParser):
    def __init__(self) -> None:
        super().__init__(convert_charrefs=True)
        self.ids: set[str] = set()
        self.links: list[str] = []
        self.errors: list[str] = []

    @override
    def handle_starttag(self, tag: str, attrs: list[tuple[str, str | None]]) -> None:
        values = dict(attrs)
        identifier = values.get("id")
        if identifier:
            if identifier in self.ids:
                self.errors.append(f"duplicate id: {identifier}")
            self.ids.add(identifier)
        if tag == "img" and "alt" not in values:
            self.errors.append("image missing alt attribute")
        for attr in ("href", "src"):
            target = values.get(attr)
            if target:
                self.links.append(target)


def site_base() -> str:
    base = os.environ.get("DOCS_BASE")
    if not base and os.environ.get("GITHUB_ACTIONS"):
        repository = os.environ.get("GITHUB_REPOSITORY", "").split("/")[-1]
        base = f"/{repository}/" if repository else "/"
    return base or "/"


def validate_site(output: Path, base: str | None = None) -> None:
    base = base or site_base()
    pages = list(output.rglob("*.html"))
    if not pages:
        raise ValueError("No HTML pages in built site")
    parsed: dict[Path, PageParser] = {}
    for page in pages:
        parser = PageParser()
        parser.feed(page.read_text(encoding="utf-8"))
        if parser.errors:
            raise ValueError(f"{page}: {', '.join(parser.errors)}")
        parsed[page.resolve()] = parser
    for page, parser in parsed.items():
        for link in parser.links:
            parts = urlsplit(link)
            if parts.scheme or parts.netloc:
                continue
            if parts.path.startswith("/"):
                prefix = "/" + base.strip("/") if base.strip("/") else ""
                if prefix and not (
                    parts.path == prefix or parts.path.startswith(prefix + "/")
                ):
                    raise ValueError(f"{page}: link outside site base: {link}")
                target = (
                    output / unquote(parts.path[len(prefix) :]).lstrip("/")
                ).resolve()
            else:
                target = (
                    (page.parent / unquote(parts.path)).resolve()
                    if parts.path
                    else page
                )
            if not target.is_relative_to(output.resolve()):
                raise ValueError(f"{page}: link escapes site: {link}")
            if target.is_dir():
                target = target / "index.html"
            if not target.is_file():
                raise ValueError(f"{page}: missing local target: {link}")
            if parts.fragment and (
                target not in parsed
                or unquote(parts.fragment) not in parsed[target].ids
            ):
                raise ValueError(f"{page}: missing fragment: {link}")


def build(output: Path) -> None:
    output = output.resolve()
    if output.exists() and (not output.is_dir() or any(output.iterdir())):
        raise ValueError(f"Output must be a new or empty directory: {output}")
    subprocess.run(["npm", "run", "build"], cwd=SITE, check=True)
    shutil.copytree(SITE / "dist", output, dirs_exist_ok=True)
    (output / ".nojekyll").touch()
    validate_site(output)


class BuildArgs(argparse.Namespace):
    output: Path


def main() -> None:
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument("--output", type=Path, default=ROOT / "site-build")
    args = parser.parse_args(namespace=BuildArgs())
    build(args.output)
    print(f"Documentation built and validated: {args.output}")


if __name__ == "__main__":
    main()
