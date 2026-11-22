# Oh My Fox and remote access: verification checkpoint

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
