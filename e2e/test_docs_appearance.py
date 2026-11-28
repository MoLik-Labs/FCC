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
@pytest.mark.parametrize(
    "path", ["", "vi/", "vi.html", "molik.html", "reference/", "reference/vi/"]
)
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
    regions = page.locator('.nb-table-scroll[role="region"][tabindex="0"][aria-label]')
    assert regions.count() == page.locator("main table").count() > 0
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


def test_static_agent_guides_at_project_base(page: Page, docs_url: str) -> None:
    index = page.request.get(docs_url + "llms.txt")
    assert index.ok
    text = index.text()
    assert "https://molik-labs.github.io/FCC/guide.md" in text
    assert "https://molik-labs.github.io/FCC/vi/guide.md" in text
    assert "https://molik-labs.github.io/FCC/reference/index.md" in text
    assert "https://molik-labs.github.io/FCC/reference/vi/index.md" in text
    for path in ("guide.md", "vi/guide.md"):
        response = page.request.get(docs_url + path)
        assert response.ok
        assert "PROVIDER_MAX_CONCURRENCY" in response.text()
        assert "fcc-update" in response.text()
    for path in ("reference/index.md", "reference/vi/index.md"):
        response = page.request.get(docs_url + path)
        assert response.ok
        assert "/v1/responses" in response.text()
    page.goto(docs_url + "vi/")
    page.get_by_role("link", name="Markdown cho agent").click()
    expect(page).to_have_url(docs_url + "vi/guide.md")


def test_nimbus_table_scrolls_with_keyboard(page: Page, docs_url: str) -> None:
    page.set_viewport_size({"width": 320, "height": 900})
    page.goto(docs_url)
    region = page.locator('.nb-table-scroll[role="region"]').first
    assert region.evaluate("el => el.scrollWidth > el.clientWidth")
    assert region.get_attribute("aria-label")
    region.focus()
    page.keyboard.press("ArrowRight")
    page.wait_for_function("document.querySelector('.nb-table-scroll').scrollLeft > 0")


def test_reference_navigation_and_search(page: Page, docs_url: str) -> None:
    page.goto(docs_url)
    page.get_by_role("link", name="API reference", exact=True).click()
    expect(page).to_have_url(docs_url + "reference/")
    expect(page.get_by_role("heading", name="FCC API overview")).to_be_visible()
    page.get_by_role("button", name="Search").click()
    page.locator("#localsearch-input").fill("responses")
    result = page.locator("#localsearch-list a.result").first
    expect(result).to_be_visible()
    assert result.get_attribute("href") == "/FCC/reference/#api-endpoints"
    page.keyboard.press("Escape")
    page.get_by_role("link", name="Security and remote access").click()
    expect(page).to_have_url(docs_url + "#security")
    page.goto(docs_url + "vi/")
    page.get_by_role("link", name="Tham khảo API", exact=True).click()
    expect(page).to_have_url(docs_url + "reference/vi/")
    expect(page.get_by_role("heading", name="Tổng quan API FCC")).to_be_visible()
    page.get_by_role("button", name="Tìm kiếm").click()
    page.locator("#localsearch-input").fill("responses")
    expect(page.locator("#localsearch-list a.result").first).to_be_visible()
    page.keyboard.press("Escape")
    page.get_by_role("link", name="Bảo mật và truy cập từ xa").click()
    expect(page).to_have_url(docs_url + "vi/#security")


def test_reference_works_without_javascript(browser: Browser, docs_url: str) -> None:
    context = browser.new_context(
        java_script_enabled=False, viewport={"width": 320, "height": 900}
    )
    try:
        page = context.new_page()
        page.goto(docs_url + "reference/")
        expect(page.get_by_role("heading", name="FCC API overview")).to_be_visible()
        page.get_by_role("link", name="Security and remote access").click()
        expect(page).to_have_url(docs_url + "#security")
        page.goto(docs_url + "reference/vi/")
        expect(page.get_by_role("heading", name="Tổng quan API FCC")).to_be_visible()
        assert page.evaluate("document.documentElement.scrollWidth <= innerWidth")
    finally:
        context.close()


def test_guide_search_indexes_vitepress_reference(page: Page, docs_url: str) -> None:
    page.goto(docs_url)
    page.get_by_role("button", name="Search").first.click()
    page.get_by_placeholder("Search", exact=True).fill("count_tokens")
    result = page.get_by_role("link", name="FCC API overview", exact=False).first
    expect(result).to_be_visible(timeout=10000)
    assert "/FCC/reference/" in (result.get_attribute("href") or "")
    page.goto(docs_url + "vi/")
    page.locator("site-search button[data-open-modal]").first.click()
    page.locator(".pagefind-ui__search-input").fill("count_tokens")
    vietnamese = page.get_by_role("link", name="Tổng quan API FCC", exact=False).first
    expect(vietnamese).to_be_visible(timeout=10000)
    assert "/FCC/reference/vi/" in (vietnamese.get_attribute("href") or "")
