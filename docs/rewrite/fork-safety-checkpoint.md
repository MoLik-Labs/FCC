# Fork safety and verification checkpoint

## Update and release boundaries

The package and the retained installer scripts still carry the upstream
`free-claude-code` identity. The inherited POSIX and Windows `fcc-update`
wrappers previously downloaded and executed upstream's installer. The wrappers
in this revision now stop with exit code 2 and print manual source-checkout
instructions; they do not download, interpret arguments or update anything.
This is a safety stop, **not an automatic fork updater**. Old installed copies
are not replaced by this change and can still execute the upstream installer.
The inherited `scripts/install.sh` and `scripts/install.ps1` also still target
upstream releases; do not use them to install this checkout.

The inherited PyPI publication job is disabled in this fork. It must not be
re-enabled until the fork has an agreed package identity, distribution source,
versioning policy and credentials, and the entire release flow has been
reviewed and tested. Manually invoking historical release scripts remains a
separate action and is not made safe by disabling the workflow. GitHub Pages is
a different workflow and remains subject to repository settings and deployment
permissions.

## Test collection

When live smoke tests were disabled, the smoke plugin marked **every** collected
item skipped, including unrelated tests in mixed invocations. It now applies the
skip only to items in the `smoke/` tree. A subprocess regression runs unit and
smoke paths together with live smoke off and verifies the unit tests run. Live
smoke tests remain opt-in and may contact configured services or change data;
this checkpoint did not run live-provider tests.

## Local evidence and limits

- Python test suite: **11,740 passed, 21 skipped**, two workers, Linux.
- The installed-updater smoke test passed on Linux in an isolated uv tool
  environment with the network-command path restricted. A separately built wheel
  contained both disabled wrappers; its POSIX wrapper was installed without
  dependencies in an isolated virtual environment and returned exit code 2.
  The wheel build warned that the repository is shallow. Windows wrapper source
  is checked and the native test is written for Windows, but Windows execution
  has not been performed here.
- Documentation builder and its eight tests passed; the four Chromium reading
  preference tests passed. This does not establish live GitHub Pages deployment.
- Go 1.22.12 `go vet`, `go test` and `go test -race` passed; some Go results
  came from cache and the partial Go port is not a replacement server.
- Rust fmt, Clippy and all 30 tests passed with the existing third-party local
  Rust toolchain and temporary source replacement. Those local sources have
  not been verified byte-for-byte against crates.io archives. The opt-in local
  MCP SDK and FCC mock integration tests both passed; they do not prove general
  MCP conformance, provider compatibility or workspace isolation.
- Repository Ruff and scoped `pyrefly check src tests e2e` passed. The actual
  repository-wide `pyrefly check` still reports **330 diagnostics**, largely in
  historical benchmark and generator scripts; the existing CI gate is not
  green on that evidence. No suppression or narrowing of the CI gate was added.
- Remote tunnel behavior, live provider accounts, all OS/browser combinations,
  real package publication and deployed Pages were not tested. Admin and Code
  do not have a separate authentication boundary; protect remote access at the
  network/tunnel layer. See [remote access](../oh-my-fox/remote-access.md).

## Starlight documentation migration

The documentation now builds with pinned Astro and Starlight dependencies. A
fresh `npm ci --engine-strict` on Node 22.23.4 succeeded, and `npm audit`
reported no known vulnerabilities at the time of the check. The documentation
unit and Chromium tests reported **25 passed and 18 subtests passed**. They
exercise the `/FCC/` base, local search, both languages, legacy fragments,
mobile layouts and navigation with JavaScript disabled. Selected axe WCAG
2.0/2.1 A/AA automated checks reported no violations on the two main language
pages at 320 and 1440 pixels. These checks do not constitute a full manual
accessibility review or cross-browser testing.

The workflow was exercised locally with GitHub Actions repository variables;
actual GitHub Actions and Pages deployment remain unverified. The repository
administrator must configure Pages and the publishing changes must reach the
publishing branch. The full backend suite above was not rerun for this docs-only
migration. The old reading-theme controls remain on archived pages; the new
Starlight pages have light/dark appearance, while FCC Admin independently
retains its additional themes.

## Nimbus and static agent-facing docs

Nimbus 0.20.2's MIT-licensed `tableScroll` helper is used with the Starlight
Markdown pipeline for named, keyboard-scrollable tables. Testing both complete
Astro integrations together showed duplicated search/sitemap generation and
Nimbus did not discover the Starlight routes, so this site deliberately does
not claim to run a full Nimbus site. Static, Nimbus-inspired `llms.txt` and
English/Vietnamese Markdown guides are built from the same source files. No
hosted AI search or remote MCP server is provided by GitHub Pages.

With the pinned dependencies, a clean `npm ci --engine-strict` succeeded and
`npm audit` reported no known vulnerabilities. The documentation and Chromium
checks reported **27 passed and 18 subtests passed**; the tested cases include
`/FCC/` search and agent links, no-JavaScript navigation and keyboard table
scrolling. Selected axe WCAG 2.0/2.1 A/AA checks on both languages at 320 and
1440 pixels in light and dark mode reported no violations. These automated
checks are not full accessibility or cross-browser certification.

The GitHub Pages API reports a workflow-based Pages configuration at
`https://molik-labs.github.io/FCC/`, with the `github-pages` environment
restricted to `main`. The documentation workflow is not yet on the default
branch. The current working branch also contains broader application changes;
deploying via `main` requires integration and review, not simply pushing this
branch. No successful hosted deployment was verified in this checkpoint.
