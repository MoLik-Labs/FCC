"""Appearance preferences remain local and do not interrupt Admin navigation."""

import pytest
from playwright.sync_api import Page, expect


@pytest.mark.parametrize(
    "theme", ["forest", "catppuccin-mocha", "nord", "dracula", "ocean", "dusk", "dawn"]
)
def test_theme_persists_and_navigation_still_works(
    page: Page, admin_base_url: str, theme: str
):
    page.goto(f"{admin_base_url}/admin")
    select = page.get_by_label("Appearance theme")
    select.select_option(theme)
    expect(page.locator("html")).to_have_attribute("data-theme", theme)
    page.get_by_role("button", name="Model Config", exact=True).click()
    expect(page.locator("#pageTitle")).to_have_text("Model Config")
    page.reload()
    expect(select).to_have_value(theme)
    expect(page.locator("html")).to_have_attribute("data-theme", theme)
    expect(page.locator("#pageTitle")).to_have_text("Model Config")


def test_invalid_or_unavailable_storage_falls_back_without_breaking_ui(
    page: Page, admin_base_url: str
):
    page.goto(f"{admin_base_url}/admin")
    page.evaluate("localStorage.setItem('fcc-appearance-theme', 'unknown')")
    page.reload()
    expect(page.get_by_label("Appearance theme")).to_have_value("forest")
    page.add_init_script(
        "Object.defineProperty(window, 'localStorage', { get() { throw new Error('blocked'); } })"
    )
    page.reload()
    page.get_by_label("Appearance theme").select_option("ocean")
    expect(page.locator("html")).to_have_attribute("data-theme", "ocean")
    page.get_by_role("button", name="Messaging", exact=True).click()
    expect(page.locator("#pageTitle")).to_have_text("Messaging")


@pytest.mark.parametrize("width", [320, 390, 900, 1440])
def test_appearance_and_navigation_fit_viewport(
    page: Page, admin_base_url: str, width: int
):
    page.set_viewport_size({"width": width, "height": 800})
    page.goto(f"{admin_base_url}/admin")
    expect(page.get_by_label("Appearance theme")).to_be_visible()
    assert page.evaluate("document.documentElement.scrollWidth <= window.innerWidth")
    page.get_by_role("button", name="Integrations", exact=True).click()
    expect(page.locator("#pageTitle")).to_have_text("Integrations")
    assert page.evaluate("document.documentElement.scrollWidth <= window.innerWidth")


def test_skip_link_reaches_main_content_by_keyboard(page: Page, admin_base_url: str):
    page.goto(f"{admin_base_url}/admin")
    page.keyboard.press("Tab")
    link = page.get_by_role("link", name="Skip to main content")
    expect(link).to_be_focused()
    page.keyboard.press("Enter")
    expect(page.locator("#main-content")).to_be_focused()


@pytest.mark.parametrize("width", [320, 390, 1440])
def test_code_session_navigation_keeps_theme_control_usable(
    page: Page, admin_base_url: str, width: int
):
    page.set_viewport_size({"width": width, "height": 800})
    page.goto(f"{admin_base_url}/admin/code")
    expect(page.locator("#codeRoot")).to_be_visible()
    page.get_by_label("Appearance theme").select_option("catppuccin-mocha")
    expect(page.locator("html")).to_have_attribute("data-theme", "catppuccin-mocha")
    assert page.evaluate("document.documentElement.scrollWidth <= window.innerWidth")


@pytest.mark.parametrize(
    "theme,background",
    [
        ("forest", "rgb(11, 33, 27)"),
        ("catppuccin-mocha", "rgb(30, 30, 46)"),
        ("nord", "rgb(46, 52, 64)"),
        ("dracula", "rgb(40, 42, 54)"),
    ],
)
def test_theme_uses_expected_base_palette(
    page: Page, admin_base_url: str, theme: str, background: str
):
    page.goto(f"{admin_base_url}/admin")
    page.get_by_label("Appearance theme").select_option(theme)
    assert (
        page.evaluate("getComputedStyle(document.body).backgroundColor") == background
    )
