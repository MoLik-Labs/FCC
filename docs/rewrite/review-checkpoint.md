# Quality review checkpoint

This checkpoint covers branch `ce3fc6b` plus the reduced-motion change in this
commit, not current `main` (`0d7e83e`). It is not a release approval or a
complete A–Z audit. See [the plan](quality-review-plan.md) and
[integration review](integration-review.md) for outstanding gates.

## UI

The Admin navigation forced smooth scrolling on every view change, including
when the user requested reduced motion. The navigation now uses instant
scrolling for that preference and retains smooth scrolling otherwise. A
browser regression exercises both settings and verifies that the destination
view is shown. The reduced-motion case was observed to fail before the fix.

The existing Admin asset and provider tests and the complete current browser
suite passed against a temporary headless Chromium 153 using two workers and
`--no-sandbox`: **312 passed**. The suite includes desktop and narrow viewport
cases, not physical-device, screen-reader or cross-browser validation. A
rendered desktop Model Config screenshot showed a restrained system font,
consistent control shapes, selected-navigation state, and explicit field
labels; this is one state, not evidence of universal visual or accessibility
quality. More tasks and states need manual review before redesign or release.

## Baseline verification

- Python tests: **11,528 passed, 13 skipped**, two workers, editable pip install;
  this is not a `uv --locked` run. Focused Admin, Code routes and Gemini SDK
  recovery: **328 passed** after the UI change.
- Go 1.22.12: `go vet`, `go test ./...`, and `go test -race ./...` passed.
  The race run required `CGO_ENABLED=1` in this environment.
- Ruff format and lint, `ty check src tests`, the suppression check and
  JavaScript syntax check passed. Repository-wide `ty check` still reports
  **74 diagnostics**, primarily in historical fixture generators; the
  current main branch replaces this gate with strict Pyrefly, which remains
  to be integrated and evaluated.
- The first Python run failed one CI-script test because loose third-party
  skills under the ignored `agent_workspace` were nevertheless scanned by
  the suppression check. They were moved outside the repository and archived
  locally; the check and full suite then passed.
- Initial browser setup failed because Chromium could not find NSS/NSPR and
  pytest used its default executable. The bundled libraries and a temporary
  pytest launch override resolved these environment issues. The override and
  browser are not committed and are not a production deployment configuration.

The browser evidence and green tests do not close the current-main merge,
current-reference Go fixture drift, remaining Go implementation, performance
benchmarking, or deployment risk reviews. Keep remote hosting functional;
exposed Admin and Code routes do not have separate Admin authentication.
