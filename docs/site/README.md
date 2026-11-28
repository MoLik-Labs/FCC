# FCC documentation site

The bilingual FCC guide is built with [Astro Starlight](https://starlight.astro.build/).
English and Vietnamese source is in `src/content/docs/`. Keep both translations
and the repository README aligned with actual behavior and safety boundaries.
Starlight provides responsive navigation, local Pagefind search, and light/dark
appearance. The bilingual API reference at `reference/` is built separately with
[VitePress](https://vitepress.dev/) into the same static Pages artifact. Its
source is in `reference/`, and its overview links back to the Starlight guide. Forest accents are used on this site. FCC Admin has separate
browser-local themes including Forest and Catppuccin Mocha; those are not the
Starlight theme selector.

The site also uses the MIT-licensed [Nimbus](https://github.com/cloudflare/nimbus)
`tableScroll` helper to make wide tables named, keyboard-scrollable regions.
Nimbus and Starlight are both Astro documentation frameworks rather than
composable themes. We do not run the full Nimbus integration alongside
Starlight: both would generate search and sitemap output, while Nimbus's content
routing does not discover the Starlight pages. The pinned Nimbus package is
used for the table helper, not for its full CLI or its server features.

The two generators share one base-path and canonical-origin configuration. The
final build re-indexes both with Pagefind, so search in the main guide finds the
API reference too. VitePress also has its own reference-only local search; its
appearance preference is independent of Starlight's. Each framework's sitemap
is included in the published sitemap index. VitePress is pinned to a 2.0 alpha:
the current 1.6 stable dependency tree reports Vite advisories. Alpha APIs can
change; review upgrades and repeat the full build and browser checks.

For agent readers, this static site publishes a Nimbus-inspired `llms.txt` index
and the same source content at `guide.md`, `vi/guide.md`,
`reference/index.md` and `reference/vi/index.md`. These are generated at build
time from the bilingual guides and need no external service. They are
public just like the HTML. GitHub Pages does not run an MCP server, AI search,
or server-side Nimbus features; do not put secrets in any published content.

## Build and preview

Install Node.js 22.23.4 or later, then from the repository root:

```sh
npm ci --prefix docs/site
npm run dev --prefix docs/site
```

The dev command starts Astro and a local VitePress server on port 4322, proxied
through Astro under `/reference/`; use the Astro address it prints. The port
must be free. To work on VitePress alone, use `npm run reference:dev --prefix
docs/site`. Search uses production indexes and may differ in development.

For a production build and local-link validation, run:

```sh
python3 -m unittest tests.scripts.test_docs_site -v
python3 scripts/build_docs_site.py
python3 -m http.server 8000 --directory site-build
```

The Python wrapper builds both generators, merges their sitemaps, indexes the
combined pages, and needs Node/npm and installed dependencies. It requires a new
or empty output directory and checks local links, fragments, duplicate IDs and
image alt attributes. It does not verify external sites or product behavior.
Use `--output` for another destination. For the project-site base path, set
`DOCS_BASE=/FCC/`; serve its contents under `/FCC/` rather than the server root.
`DOCS_SITE` can override the canonical site origin for custom domains. GitHub
Actions derives the repository base automatically. Use `DOCS_BASE=/` for a
user/organization site or custom domain served at the root.

The old English and Vietnamese HTML sources are retained as compatibility pages
at `legacy.html` and `vi.html`; historic `/#quickstart`, `/vi.html#...`, and
`/molik.html#...` links continue to work under the configured base. The old
reading themes are available on those archived pages only. The build syncs
those files and the Molik mark into `public/`; edit the original files, not the
generated copies. The Molik gallery is a static design preview, not a working
agent UI. The MCP guide documents the experimental read-only server and its
security boundaries.

## Deploy

The `Documentation site` workflow builds pull requests and publishes pushes or
manual runs on `main` with GitHub Actions. The intended project-site address is
<https://molik-labs.github.io/FCC/> after a successful deployment. A repository administrator must
select **Settings → Pages → Build and deployment → Source: GitHub Actions** and
approve the environment if required. A deployment URL exists only after the
change reaches the publishing branch and a deployment succeeds. This workflow
does not publish a working branch or change repository settings.

The quick start uses the current development branch of the fork. The inherited
installers target upstream, the new packaged updater is disabled, and older
installed copies may still execute upstream's installer. The inherited PyPI
publish job is disabled on this fork. Do not treat generated docs as evidence
that deployment, provider connectivity, or platform-specific behavior has been
verified.
