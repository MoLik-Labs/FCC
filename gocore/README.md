# gocore

`gocore` is an **in-progress** Go port of selected FCC streaming and provider
components. The Python application remains the executable reference and the
production entry point; this module is not a complete replacement server. It
uses only the Go standard library.

## Layout

| Package | Responsibility |
| --- | --- |
| `sse`, `wirejson`, `wireuuid` | SSE framing and wire-compatible JSON and UUID helpers |
| `retention`, `recovery`, `delivery` | Retention accounting, recovery decisions and delivered-response state |
| `toolbuf`, `continuation`, `pipeline` | Tool-call holdback, continuation and streaming composition |
| `providerpolicy`, `admission` | Failure policy, retry timing, capacity and execution lifecycle |
| `openaichat` | OpenAI-compatible Chat leaf utilities, tool-name handling, stream lifecycle and Anthropic SSE output |
| `cmd/bench` | Benchmark using the shared streaming fixture |

Coverage is incremental, not a claim of complete behavioral or byte parity.
In particular, provider adapters and their runners, additional wire writers,
service and Admin APIs, storage, launchers and bridges remain unfinished. The
[master plan](../docs/rewrite/00-master-plan.md) tracks the phases and current
limitations. Python's isolated surrogate handling, Go integer bounds, future
Unicode versions, and metadata key ordering have explicit parity limitations;
ordered JSON objects are required when insertion order affects wire bytes.

## Parity workflow

Golden data under `testdata/` is generated from the Python implementation by
scripts in `../bench_local/`. Run the applicable generator from the repository
root with the project virtualenv active, then review fixture changes and run
the Go tests. For example:

```bash
.venv/bin/python bench_local/generate_golden.py
cd gocore
go test ./...
go vet ./...
gofmt -l .
```

Other phases have their own `generate_golden_*` scripts. Many still target the
older Python reference and cannot currently run against the newer packet-based
implementation. Committed fixtures capture the reference revision used to
generate them; passing Go tests alone cannot prove compatibility with a newer
Python reference. See the [integration review](../docs/rewrite/integration-review.md)
for the generator inventory and known divergences.

## Benchmark

From the repository root:

```bash
.venv/bin/python bench_local/bench_public_delivery.py
cd gocore
go run ./cmd/bench
go test -bench=HotPath -benchmem ./sse/
```

The historical `bench_golden.py` targets the removed API-level tool buffer and
is not directly comparable to the new packet-based public delivery benchmark.
These are component benchmarks, not measurements of the complete FCC server.
Results depend on the machine, runtime and implementation revision. Historical
measurements and the language decision are documented in
[the language verdict](../docs/rewrite/01-language-verdict.md).

## Porting invariants

- Only CR, LF and CRLF delimit SSE frames; Unicode line and paragraph
  separators inside JSON are ordinary characters.
- Numbers stay on the wire as `json.Number`; no float round-trip.
- Cached payloads are read-only; copy before mutation.
- Retention uses an approximation, not an exact serialized byte count. Its
  behavior is covered by reference fixtures but should not be treated as a
  strict bound on JSON wire size.
- Preserve reference behavior with differential tests, including fallback paths
  and malformed inputs within the supported contract.
