# Integration review

This is a working review, not a release approval. It compares branch tip
`e4d0051` with `main` at `62db82e`; repeat the comparison if either moves.
The merge base is `56ee72e`.

## Findings and merge gates

1. **Remote administration is a trust boundary change.** This branch removes
   loopback checks from Admin and Code session routes so tunnels and remote
   hosts work. Admin routes do not require the proxy API's bearer token, and
   the status route reflects the request Origin. Exposure can allow untrusted
   clients to change configuration or operate code sessions. The README now
   warns about this; remote access remains intentional. Review deployment and
   browser-origin risks, including newly added routes on `main`, before release.
2. **The streaming conflict is architectural.** `main` replaces the API-level
   `ToolCallBuffer` with a `StreamOutput` packet and completed-call proof path,
   and changes replay finalization and Chat recovery. This branch optimizes
   parsing and retention in the older delivery path. A textual resolution alone
   cannot establish atomic tool release, continuation correctness, or preserved
   performance. Port the applicable optimizations to the new architecture, test
   failure and replay paths, and measure the integrated implementation.
3. **Go fixtures reference an older Python specification.** For example,
   `main` removes `providers/openai_chat/tool_calls.py` and
   `tool_call_identity.py`; branch generators import them. Adapt the generators
   to the new reference, regenerate all affected fixtures and investigate
   differences before claiming Go/Python parity. Existing Go tests only verify
   the committed fixtures.
4. **Retention is estimated, not a strict byte limit.** `_json_byte_size` and
   its Go port do not count JSON escaping. For a string of 1,000 NUL characters,
   the estimate is 1,002 bytes while JSON serialization takes 6,002 bytes.
   Neither serialized size nor live memory is guaranteed by this estimate.
   Decide the intended resource contract and test adversarial inputs before
   changing it; any implementation change must be ported to both languages and
   its performance and behavior assessed.
5. **Documentation and installation provenance.** The root README's install
   commands reference the upstream repository, not this fork. The Go README
   previously described the partial port as byte-for-byte complete and omitted
   later packages and limitations; both documents now clarify their scope.

The integration produced conflicts in `api/admin_routes.py`,
`api/stream_delivery.py`, `api/tool_call_buffer.py` (modified/deleted),
`core/delivered_response.py`, and `tests/api/test_antigravity_routes.py`.
The working branch integrates the new packet release path and preserves remote
Admin access. `main` removed Antigravity connection endpoints; the integration
retains that removal rather than resurrecting obsolete routes. This is not yet
an approval to merge the working branch into `main`.

Before merging: resolve and review all conflicts; test the combined Python and
Go code, regenerate reference fixtures, check formatting, lint, type checks,
race tests, security-sensitive routes and browser tests; then benchmark the
combined path against an appropriate baseline. Record unavailable checks and
failures rather than treating historical branch-only results as a merge gate.

## Fixture migration inventory

A trial regeneration against the integrated Python reference was run in an
isolated directory, without overwriting historical fixtures. The base SSE/JSON
fixture was regenerated after fixing malformed named-error handling. The ASCII
JSON and output-cap edge fixtures regenerated unchanged.

- The phase 2, 4a, 4b, 4c-2, 4c-3c and streamed-name generators import
  interfaces removed or relocated by the new reference.
- Phase 3b still reads the removed delivery sequence field; 4c-1, 4c-3a,
  4c-4a, tool-name and stream-output generators depend on changed signatures,
  state or lifecycle operations. Anthropic Chat output now emits complete
  tool calls rather than the prior start/delta/stop interface.
- Phase 3 overlap and stream generators run but produce different events,
  sizes or sequence fields. The 4c-3b recovery generator also produces different
  execution-state/timing results. These changes need semantic review and Go
  reimplementation, not a blind fixture refresh.

The historical fixtures remain useful as records of earlier behavior, but a
green Go test run against them is not a current-reference parity gate. Update
and test each subsystem against the new Python behavior before release. As a
first migration slice, `golden_chat_identity_current.json` checks the current
unbudgeted resolver's returned slot, learned ID and error text over 27
scenarios, including Python-specific control-character whitespace. The separate
`golden_chat_identity_budget.json` now checks ten current-reference budget
scenarios, including exact limits, UTF-8 ID sizes, repeated IDs, partial state
on a rejected ID charge, and rejection without mutation for a failed slot
charge. The Go pending-budget primitive accounts for bytes and records and
preserves failed-add atomicity. It uses machine integers, assumes representable
counters and valid UTF-8 inputs, and only models the error message, not the
full Python execution-failure metadata. The surrounding current source
tool-call buffer is still not ported.

## Integrated verification

The integrated Python test suite passed 11,528 tests with 13 skips on Linux
using two workers and a pip-installed environment; the SDK test bridge now
uses the same `httpx2` response type as the Gemini SDK. Go formatting, `go vet`,
`go test ./...` and `go test -race ./...` passed against the committed fixtures.
Repository-wide Ruff and `ty check src tests` passed. Repository-wide `ty check`
still reports diagnostics in historical benchmark/generator scripts that need
migration; it is not a green CI gate. Playwright browsers and other operating
systems have not been tested here.

`bench_local/bench_public_delivery.py` exercises the new packet-based Messages
path. On this sandbox, repeated 2,004-frame runs gave approximately 334–343 ms
for `main` and 273 ms for the integrated code after a warm-up; output SHA-256
was identical (`5cde2bedb002bf497ba4f28c158340f2e596613d1c9fb5307f9ecf205d1ab6b8`).
These are noisy component measurements, not an end-to-end throughput claim or
proof that the earlier 56% improvement survives every new path.

At review time the remote had only `main` and the active working branch. No
obsolete remote branch was identified for deletion. Preserve the working
branch until integration is verified.
