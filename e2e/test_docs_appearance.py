"""Documentation reading preferences work across languages and without storage."""

from collections.abc import Iterator
from functools import partial
from http.server import SimpleHTTPRequestHandler, ThreadingHTTPServer
from pathlib import Path
from threading import Thread

import pytest
from playwright.sync_api import Browser, Page, expect

from scripts.build_docs_site import build


@pytest.fixture(scope="module")
def docs_url(tmp_path_factory: pytest.TempPathFactory) -> Iterator[str]:
    output: Path = tmp_path_factory.mktemp("fcc-docs") / "site"
    build(output)
    server = ThreadingHTTPServer(
        ("127.0.0.1", 0), partial(SimpleHTTPRequestHandler, directory=str(output))
    )
    thread = Thread(target=server.serve_forever, daemon=True)
    thread.start()
    try:
        yield f"http://127.0.0.1:{server.server_port}"
    finally:
        server.shutdown()
        thread.join(timeout=5)
        server.server_close()


@pytest.mark.parametrize("width", [320, 1440])
def test_docs_theme_persists_across_languages(
    page: Page, docs_url: str, width: int
) -> None:
    page.set_viewport_size({"width": width, "height": 900})
    page.goto(docs_url)
    if width == 320:
        page.locator(".menu-toggle").click()
    page.get_by_label("Reading theme").select_option("catppuccin-mocha")
    expect(page.locator("html")).to_have_attribute(
        "data-reading-theme", "catppuccin-mocha"
    )
    page.get_by_role("link", name="VI — read the Vietnamese guide").click()
    if width == 320:
        page.locator(".menu-toggle").click()
    expect(page.get_by_label("Giao diện đọc")).to_have_value("catppuccin-mocha")
    assert page.evaluate("document.documentElement.scrollWidth <= innerWidth")
    page.get_by_label("Giao diện đọc").select_option("forest")
    page.reload()
    expect(page.locator("html")).to_have_attribute("data-reading-theme", "forest")


def test_docs_theme_storage_errors_and_invalid_values(
    page: Page, docs_url: str
) -> None:
    page.goto(docs_url)
    page.evaluate("localStorage.setItem('fcc-docs-theme', 'invalid')")
    page.reload()
    expect(page.get_by_label("Reading theme")).to_have_value("light")
    page.add_init_script(
        "Object.defineProperty(window, 'localStorage', { get() { throw new Error('blocked'); } })"
    )
    page.reload()
    page.get_by_label("Reading theme").select_option("forest")
    expect(page.locator("html")).to_have_attribute("data-reading-theme", "forest")


def test_docs_navigation_remains_available_without_javascript(
    browser: Browser, docs_url: str
) -> None:
    context = browser.new_context(
        java_script_enabled=False, viewport={"width": 320, "height": 900}
    )
    try:
        page = context.new_page()
        page.goto(docs_url)
        expect(page.get_by_role("link", name="Workload settings")).to_be_visible()
        expect(page.locator(".reading-control")).to_be_hidden()
        assert page.evaluate("document.documentElement.scrollWidth <= innerWidth")
    finally:
        context.close()
