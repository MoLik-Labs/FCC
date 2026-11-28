# Current-main integration checkpoint

The working branch integrates `main` at `0d7e83e` (six commits beyond the
previous `62db82e` integration) while retaining the intentional remote Admin
and Code behavior. This is a branch integration checkpoint, not a release
approval; repeat the comparison if main moves.

## Resolution

- Retained packet-delivery caching, SSE fast paths and failure observation,
  alongside main's JSON vocabulary and stricter type annotations. The original
  event observer remains available for its contract tests.
- Preserved tunnel and remote-client tests rather than restoring the old
  loopback guard. The Gemini test bridge continues to return `httpx2` responses
  to the SDK. Main's failed-background-delivery behavior and diagnostics
  changes are included.
- Corrected type errors in the branch's integration tests and the decoder
  without using type suppressions. The renderer still honors reduced motion.

## Verification

Python 3.14.10 and Go 1.22.12 were built from the repository bootstrap scripts.
Python dependencies were synchronized with `uv sync --locked`; the test runs
below used that environment on Linux.

- Python: **11,707 passed, 19 skipped**, two workers. The focused stream,
  Admin, Code and Gemini regression run passed **462 tests**.
- Headless Chromium: **312 passed**, two workers, using the temporary npm
  Chromium 153 build and `--no-sandbox`. This does not verify Firefox, WebKit,
  native screen readers, physical touch devices or other operating systems.
- Go `vet`, tests and race tests passed. The current identity fixture
  regenerated unchanged with SHA-256
  `e118055efbdc986ddb69a89408b465a6405b516b43dbb37607fe1eb758b6e24c`.
- Ruff format and lint, the suppression check, and strict Pyrefly over
  `src tests e2e` passed. The current identity generator and the packet
  benchmark also pass Pyrefly. **Repository-wide Pyrefly still fails with 334
  diagnostics**, all located in 19 `bench_local` scripts. Many scripts also
  import APIs removed from current Python, so this cannot be fixed safely by
  suppressing the diagnostics or blindly regenerating the fixtures. The
  repository-wide CI type gate remains red.

For a 2,004-frame packet-based Messages delivery component benchmark, three
interleaved runs gave medians of **222.39–225.85 ms** on `0d7e83e` versus
**188.78–198.62 ms** on the integrated branch; each run produced the same
output SHA-256 `5cde2bedb002bf497ba4f28c158340f2e596613d1c9fb5307f9ecf205d1ab6b8`.
These are local component measurements, not an end-to-end latency or
throughput guarantee and not proof that the earlier 56% gain generalizes.

## Remaining gates

The repository-wide type check and migration of historical golden generators
remain open. Go tests against historical fixtures do not establish parity with
the current Python executable specification. The Go service, adapters,
storage, transport, admin UI integration and many other rewrite phases remain
incomplete. Comprehensive cross-platform, accessibility, security and
end-to-end performance review also remains. Do not merge this working branch
into main or delete valuable work until those gates are resolved. Remote
hosting remains permitted; exposed Admin and Code routes have no separate
Admin authentication and must be deployed accordingly.
