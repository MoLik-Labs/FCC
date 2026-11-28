"""Exercise the Starlight site at the GitHub Pages project subpath."""

from collections.abc import Iterator
from functools import partial
from http.server import SimpleHTTPRequestHandler, ThreadingHTTPServer
from threading import Thread

import pytest
from playwright.sync_api import Browser, Page, expect

from scripts.build_docs_site import build


@pytest.fixture(scope="module")
def docs_url(tmp_path_factory: pytest.TempPathFactory) -> Iterator[str]:
    root = tmp_path_factory.mktemp("fcc-docs")
    import os

    previous = os.environ.get("DOCS_BASE")
    os.environ["DOCS_BASE"] = "/FCC/"
    try:
        build(root / "FCC")
    finally:
        if previous is None:
            os.environ.pop("DOCS_BASE", None)
        else:
            os.environ["DOCS_BASE"] = previous
    server = ThreadingHTTPServer(
        ("127.0.0.1", 0), partial(SimpleHTTPRequestHandler, directory=str(root))
    )
    thread = Thread(target=server.serve_forever, daemon=True)
    thread.start()
    try:
        yield f"http://127.0.0.1:{server.server_port}/FCC/"
    finally:
        server.shutdown()
        thread.join(timeout=5)
        server.server_close()


@pytest.mark.parametrize("width", [320, 390, 1440])
@pytest.mark.parametrize("path", ["", "vi/", "vi.html", "molik.html"])
def test_docs_render_at_pages_base(
    page: Page, docs_url: str, width: int, path: str
) -> None:
    errors: list[str] = []
    page.on("pageerror", lambda error: errors.append(str(error)))
    page.set_viewport_size({"width": width, "height": 900})
    response = page.goto(docs_url + path)
    assert response and response.status == 200
    expect(page.locator("h1").first).to_be_visible()
    assert page.evaluate("document.documentElement.scrollWidth <= innerWidth")
    assert not errors


def test_search_and_language_navigation(page: Page, docs_url: str) -> None:
    page.goto(docs_url)
    page.get_by_role("link", name="Hướng dẫn tiếng Việt").click()
    expect(page).to_have_url(docs_url + "vi/")
    expect(page.locator("#security")).to_have_count(1)
    page.get_by_role("link", name="Cài đặt nhanh", exact=True).first.click()
    expect(page).to_have_url(docs_url + "vi/#quickstart")
    page.goto(docs_url)
    page.get_by_role("button", name="Search").first.click()
    search = page.get_by_placeholder("Search", exact=True)
    assert page.locator("pre:not([tabindex])").count() == 0
    assert page.locator("table:not([tabindex])").count() == 0
    search.fill("PROVIDER_MAX_CONCURRENCY")
    expect(
        page.get_by_role("link", name="FCC field guide", exact=False).first
    ).to_be_visible(timeout=10000)


def test_navigation_and_legacy_fragments_without_javascript(
    browser: Browser, docs_url: str
) -> None:
    context = browser.new_context(
        java_script_enabled=False, viewport={"width": 320, "height": 900}
    )
    try:
        page = context.new_page()
        for path, anchor in [
            ("", "quickstart"),
            ("vi/", "security"),
            ("vi.html", "quickstart"),
            ("molik.html", "mcp-stdio"),
        ]:
            page.goto(docs_url + path + "#" + anchor)
            expect(page.locator("#" + anchor)).to_have_count(1)
            assert page.evaluate("document.documentElement.scrollWidth <= innerWidth")
        page.goto(docs_url)
        expect(page.get_by_role("link", name="Hướng dẫn tiếng Việt")).to_be_visible()
        page.get_by_text("On this page", exact=True).first.click()
        expect(
            page.get_by_role("link", name="From zero to ready.").first
        ).to_be_visible()
    finally:
        context.close()
