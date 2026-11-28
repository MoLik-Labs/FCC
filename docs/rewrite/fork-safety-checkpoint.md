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
