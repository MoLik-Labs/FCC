# Current source-order buffer checkpoint

The standalone Go `toolbuf.OrderedBuffer` implements the current Python
`OrderedToolCallBuffer` state transitions: reserve, append, grow, lookup,
complete, discard, discard-pending, clear and lazy, ready-prefix publication.
An optional shared `retention.PendingBudget` charges and releases bytes and
records. A rejected charge leaves the buffer unchanged; clear releases only
its own slots, not another owner's charges. Draining clears backing-array
references, and the return boolean distinguishes a yielded nil from no value.

`bench_local/generate_golden_ordered_buffer_current.py` executes the Python
reference, inspecting its private queue and pending keys for per-step state.
Its 24 scenarios include eight designed cases and sixteen seeded sequences of
35 operations, with shared budget operations, bounds, errors, key reuse,
structured values and partial lazy release. Go replays all steps and compares
results, error categories, pending keys, queue contents and budget counters.
The fixture regenerates identically (SHA-256
`84f30e96b8019ea90c4dd5a3ac9e331686553409938add4a6f83105099605ca9`).
Additional Go tests cover default limits, reference release and interaction with
the Chat identity resolver using the same budget. Go vet, all Go tests and race
tests pass; the focused Python buffer/limit tests pass (10 tests). Ruff,
suppression checks and Pyrefly on `src tests e2e` and both touched generators
also pass.

The ASCII JSON fixture generator now uses the typed wrapper around Python's
standard `json.loads`; regeneration of its 1,005 existing cases is byte-for-byte
unchanged. Repository-wide Pyrefly still reports 330 diagnostics in historical
benchmark scripts and is not a passing release gate.

This is a tested standalone primitive, **not** a completed Go provider rewrite:
it is not wired into the current Chat adapter or stream publication path, and
it does not reproduce full Python execution-failure metadata. Its generic Go
keys are comparable, not all Python hashables; integer accounting is limited to
representable machine integers. Concurrent use is unsupported, matching the
request-local reference. Broader provider parity, end-to-end performance,
cross-platform/browser/accessibility and security review remain open. See
`quality-review-plan.md` and the master plan before treating this as releasable.
