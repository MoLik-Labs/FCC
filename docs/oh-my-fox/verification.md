# Oh My Molik (formerly Oh My Fox) and remote access: verification checkpoint

This is an incremental checkpoint, not an approval to release the Go rewrite or
claim state-of-the-art agent performance.

- The complete Python suite passed on Linux with two workers: **11,724 passed,
  19 skipped, 7 subtests passed**. The preceding run had one timeout in
  `test_real_responses_post_start_disconnect_has_no_failed_event`; that test
  passed when run alone and the full suite passed on the subsequent run. The
  timeout is recorded rather than treated as proof that the test cannot flake.
- A focused Admin, Code, Desktop integration, startup and documentation run
  passed **318 tests, 7 subtests**. Remote tests use simulated client addresses
  and tunnel hostnames in-process; no real Cloudflare, ngrok or Tailscale
  deployment was exercised.
- `go vet ./...`, `go test ./...` and `go test -race ./...` passed for both
  `ohmyfox/` and `gocore/` using the temporary Go 1.22.12 bootstrap on Linux.
  The Oh My Fox tests include a local fake Chat endpoint, tool-call replay,
  redirect handling, workspace bounds and a CLI invocation. No real model,
  external provider, Windows/macOS or currently supported Go release was verified.
- Ruff, Pyrefly on `src tests e2e`, the suppression check and the local Pages
  build passed. Repository-wide Pyrefly remains an open historical benchmark
  gate; the earlier baseline had 330 diagnostics and was not rerun here.
- The new Go GitHub Actions workflow has been added, but its hosted execution
  has not yet been observed. No Oh My Fox Go-versus-Rust or coding-task
  benchmark was performed, and no live deployment was made.

See `research.md` for the source inventory and remaining feature, performance,
license and security gates, and `remote-access.md` for the trust boundary.

## Rust core checkpoint

On Linux, Rust 1.88.0 ran `cargo fmt --check`, Clippy with `-D warnings`,
`cargo test --locked --offline` (seven tests), and the help command. The Go
prototype's `go test ./...` and the documentation builder/tests also passed.
These are local component tests, not evidence of model, filesystem, FCC,
tunnel, sandbox, performance or cross-platform parity. The Rust CI matrix is
configured for Linux, macOS and Windows but has not been observed running.

The sandbox could not reach `static.rust-lang.org` or crates.io. For this local
check only, compiler components were unpacked from third-party npm `@rustbin`
packages outside the repository, without running their install scripts. Their
contents were not independently verified against official Rust release checksums;
therefore local results carry a toolchain provenance limitation. The CI workflow
is configured to install Rust 1.88.0 using the usual Rust distribution instead.
No compiler binaries are included in this repository.

The updated bilingual site was tested locally in headless Chromium at 1440,
390 and 320 pixels: no overflow, page/HTTP errors or axe WCAG 2/2.1 A/AA
violations were detected, and mobile menu/Escape navigation worked. The fox
mark loaded at desktop and narrow widths. This is not cross-browser or physical
device testing and does not mean GitHub Pages is publicly deployed.

## Oh My Molik rename and theme / streaming checkpoint

- Rust 1.88.0: fmt, Clippy (`-D warnings`) and **16 offline tests** passed on
  Linux. The new tests cover SSE chunk boundaries and invalid input; delta
  assembly, size limits and finish reasons; and palette contrast and CLI output.
  This does not verify HTTP, a provider or real streamed model output.
- The toolchain provenance limitation above still applies. The CI matrix has not
  been observed; Linux local tests do not prove Windows/macOS behavior.
- The retained Go prototype passed `go vet`, `go test` and `go test -race`
  using Go 1.22.12; race testing required `CGO_ENABLED=1` in this environment.
- Documentation tests passed (**5 tests, 12 subtests**), Ruff and Pyrefly passed
  for the changed Python files, and the site builder validated local links.
  Headless Chromium and axe checks found no overflow, page/HTTP errors or
  WCAG 2/2.1 A/AA violations for the gallery and EN/VI pages at 1440, 390 and
  320 pixels. This is not cross-browser or physical-device validation.
- The selected reference image is user-provided. The vector mark is a manual
  interpretation; no legal trademark clearance is claimed. Themes are token
  sets with illustrative previews, not a functioning theme-switching TUI.

## Rust read-only Chat adapter checkpoint

The Rust executable now uses the system curl process for explicitly selected
HTTP(S) Chat endpoints and supports two read-only workspace tools. Local Rust
unit/integration tests exercise protocol validation, stdin prompting, a mocked
HTTP tool-call loop, size/path checks, ANSI-control filtering and rejection of
redirects. A separate test starts the actual FCC ASGI application on a local
TCP socket with proxy authentication and a **mock upstream provider**: it
observed a `read_file` call, delivered its result through FCC and received the
final reply. This is not a live model, HTTPS/tunnel, streaming, sandbox or
cross-platform test.

Reproduce the FCC test after building the Rust binary and syncing FCC's Python
environment:

```sh
MOLIK_BINARY="$PWD/ohmymolik/rust/target/debug/oh-my-molik" \
  uv run --locked pytest -n 0 tests/ohmymolik/test_fcc_integration.py -q
```

With `MOLIK_BINARY` unset, the test skips. A Linux GitHub Actions job is
configured to build the binary and run it; hosted execution has not yet been
observed. The Rust CLI remains non-streaming; FCC's own streaming is unchanged.

The sandbox cannot access crates.io. To run local tests, upstream Git tags were
copied into a temporary source replacement outside the repository; two serde
manifests were adapted for standalone directory use. The lockfile checksums were
read from the public crates.io index on GitHub, but the local test sources were
not verified byte-for-byte against crates.io archives. The compiler provenance
limitation above also still applies. CI is configured to fetch the locked crates
normally, but cannot be considered verified until it actually runs.

Local follow-up: **22 Rust tests** passed with fmt and Clippy, including two
TCP mock-endpoint CLI tests; **55 Python tests and 12 subtests** passed across
the opt-in FCC roundtrip, documentation and Chat protocol matrix. Ruff for the
repository and Pyrefly for `src tests e2e` passed. A local HTTPS fixture also
succeeded when curl was given its test CA and failed for the same untrusted
certificate without it. This checks TLS validation locally, not remote provider
or tunnel operation. The same provenance and hosted-CI limitations apply.

## Admin appearance and SSE prefix checkpoint

FCC Admin now provides browser-local Forest, Catppuccin Mocha, Ocean, Dusk and
Dawn themes; this is separate from the Rust agent, which still has no themed
TUI. Existing configuration and Code-session behavior remains the supported
Python implementation. A storage failure or an unrecognized theme falls back
to Forest without blocking navigation. The Catppuccin MIT notice is included
with the Python package's Admin assets.

- Full Python `tests`: **11,726 passed, 19 skipped, 12 subtests passed** with
  two workers and the opt-in Rust/FCC mock roundtrip enabled.
- Chromium `e2e`: **322 passed**. Four subsequent keyboard/Code appearance tests
  were added and the appearance suite was run separately: **14 passed**. This
  is Chromium only, on Linux, using the local browser fixture and no real
  provider accounts.
- An axe-core scan of Providers, Model Config, Messaging and Integrations for
  each of the five themes at 1440px reported no WCAG 2/2.1 A/AA violations;
  a separate 390px scan of the NVIDIA provider dialog in all themes and the
  Code library in Dawn also reported none. No page errors were observed in the
  first scan. These are selected automated checks, not a complete accessibility
  audit, and do not validate all forms, states or browsers.
- The site build, Prettier on the changed site pages, repository Ruff checks,
  and Pyrefly for `src tests e2e` passed. A wheel was built locally and its Admin
  CSS, JS and Catppuccin license were present; the shallow checkout produced a
  build-version warning.
- The Rust SSE decoder gained `push_into` so callers can retain complete events
  preceding an error in the same input chunk. It still poisons the decoder on
  error, and the existing `push` API intentionally discards that call's prefix.
  This is not an end-to-end streaming implementation or a performance claim.
  Rust fmt, Clippy and **23 tests** passed with the local source replacement and
  the previously noted unverified toolchain provenance.
- All four design skill sources were restored at recorded pins into a temporary
  directory, then verified on a second run. A dirty checkout was rejected.
  Sources remain linked rather than vendored; fresh restoration needs GitHub
  access. Future agents must actually read the repository guidance for it to
  take effect.

## Read-only MCP stdio checkpoint

The experimental Rust binary includes a deliberately limited MCP server over
stdio. Its protocol reference is the official [2025-06-18 specification](https://github.com/modelcontextprotocol/specification/tree/b0be8f0f2481be8d5ea5294cf3fa57cd03cedc48/docs/specification/2025-06-18).
It advertises only the two existing read-only workspace tools; it is not an MCP
client or an HTTP MCP server. Rust tests exercise initialization, negotiation,
notification handling, malformed frames and params, bounded input, rate
limiting, directory traversal and Git metadata denial. These tests are part of
the Rust CI matrix, whose hosted results have not yet been observed.

A local opt-in interoperability test used the official Python MCP SDK **1.30.0**
to initialize, list tools, read a fixture and reject an escape. The SDK and its
transitive dependencies were installed into a temporary environment outside the
repository; this was not a locked or audited production dependency installation.
To repeat on a Unix-like system after building Molik (this installs SDK 1.30.0
and its currently resolved transitive dependencies into an isolated environment):

```sh
uv venv /tmp/molik-mcp-sdk
uv pip install --python /tmp/molik-mcp-sdk/bin/python mcp==1.30.0
MOLIK_BINARY="$PWD/ohmymolik/rust/target/debug/oh-my-molik" \
MCP_SDK_PYTHON=/tmp/molik-mcp-sdk/bin/python \
  uv run --locked pytest -n 0 tests/ohmymolik/test_mcp_sdk.py -q
```

The test skips when either environment variable is absent. This is one SDK and
one local OS, not a full compatibility or security certification. Access is
controlled only by launching the process and OS permissions; the workspace path
checks are not a sandbox, secret filter or defense against hostile local races.
See the MCP setup and security notes in `ohmymolik/README.md`.

Local final gates: Rust fmt/Clippy and **30 tests** passed using the same
unverified local compiler and temporary dependency replacement described above;
**164 focused Python tests and 14 subtests** passed including the opt-in SDK and
FCC mock roundtrips; Ruff, Pyrefly (`src tests e2e`) and the site build passed.
Go 1.22.12 `vet`, `test` and race tests passed for `gocore`; this does not imply
full Go/Python parity. Chromium/axe checks of the English, Vietnamese and Molik
pages at 320, 390 and 1440 pixels found no overflow, browser errors or WCAG
2/2.1 A/AA violations on those pages. This is not a cross-browser audit. All
nine pinned source trees (the eight requested repositories plus Superpowers
core) were restored/checked; third-party hooks were not enabled. No performance
benchmark or live-provider/tunnel test was run for this checkpoint.
