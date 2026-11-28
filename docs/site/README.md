# FCC documentation site

This directory contains the static source for the FCC guide. The site covers
installation, providers, clients, APIs, messaging, security, maintenance and
contribution in English and Vietnamese. It links to the project README for the
full provider catalog and platform-specific details; update both translations
and the README when behavior changes. The site intentionally uses dependency-free
static HTML, CSS and JavaScript: no external fonts, framework runtime or package
install is needed to build. The bilingual guides offer Light, Forest and
Catppuccin Mocha reading themes, stored locally in the browser; with JavaScript
disabled they remain readable in Light. This is separate from the Admin theme.
The Catppuccin palette is MIT licensed; see the [license](../../src/free_claude_code/api/admin_static/CATPPUCCIN-LICENSE.txt).

## Build and preview

From the repository root, with Python 3.9 or newer:

```sh
python3 -m unittest tests.scripts.test_docs_site -v
python3 scripts/build_docs_site.py
python3 -m http.server 8000 --directory site-build
```

Open <http://localhost:8000/> on the machine running the server. To avoid
accidentally overwriting files, the builder requires a new or empty output
directory. Remove the generated, ignored `site-build/` directory before a
subsequent build, or use `--output` with another empty directory. It copies the
site source and the Molik vector mark, creates `.nojekyll`, and checks
local links, fragments, duplicate IDs and image alt attributes. It does not
verify external links, rendered appearance or product behavior. No additional
packages are needed. To exercise the reading theme, language navigation,
restricted storage and the no-JavaScript layout with a locally installed
Playwright Chromium browser, run `uv run pytest -n 0 e2e/test_docs_appearance.py`.
These browser tests do not replace manual accessibility or cross-browser review.

Admin's Nord and Dracula themes are palette adaptations; Nord foregrounds are
adjusted for contrast. The Admin and documentation theme choices are independent.
See [Nord](https://www.nordtheme.com/) and [Dracula](https://draculatheme.com/)
for their original palettes.

## Deploy

The `Documentation site` workflow checks every pull request. Pushes to `main`
and manually dispatched runs on `main` build a Pages artifact and deploy it
using GitHub Actions. A repository administrator must select **Settings → Pages →
Build and deployment → Source: GitHub Actions**, and any required environment
approval must be granted. See the workflow's deploy job for the actual page
URL. This workflow never publishes a pull request or working branch, and adding
it to a branch does not itself publish a site or change repository settings.

Links to project source and reference material point to the `main` branch of
`MoLik-Labs/FCC`; they represent the published documentation after these
changes reach that branch. The quick start clones the current development branch
of this fork. The inherited installers target upstream, not this checkout; the
new packaged updater is disabled, while older installed copies may still run
upstream's installer. The inherited PyPI publish job is disabled on this fork.

The Oh My Molik [theme collection](molik.html) is a static design preview;
it does not demonstrate a working agent UI or theme persistence. The builder
copies the vector mark from `ohmymolik/brand/` into the published site.
