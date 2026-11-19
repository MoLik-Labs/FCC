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

The trial merge produced conflicts in `api/admin_routes.py`,
`api/stream_delivery.py`, `api/tool_call_buffer.py` (modified/deleted),
`core/delivered_response.py`, and `tests/api/test_antigravity_routes.py`.
The trial was aborted without changing either branch. In particular, `main`
removed Antigravity connection endpoints; preserve that deliberate removal rather
than resurrecting obsolete routes while retaining remote access for supported
Admin routes.

Before merging: resolve and review all conflicts; test the combined Python and
Go code, regenerate reference fixtures, check formatting, lint, type checks,
race tests, security-sensitive routes and browser tests; then benchmark the
combined path against an appropriate baseline. Record unavailable checks and
failures rather than treating historical branch-only results as a merge gate.

## Branch-only verification

On the reviewed branch, Go formatting, `go vet`, `go test ./...` and
`go test -race ./...` passed. The Python test suite passed with 9,794 tests
and 10 skips using two workers. These results do not test the trial merge,
regenerated fixtures against `main`, Playwright browsers, or performance of
the integrated pipeline.

At review time the remote had only `main` and the active working branch. No
obsolete remote branch was identified for deletion. Preserve the working
branch until integration is verified.
