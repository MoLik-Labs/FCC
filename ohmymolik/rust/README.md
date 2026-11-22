# Oh My Molik: Rust core

This is an experimental, dependency-free foundation, not a coding agent yet.
The CLI supports help, version and inspection of theme tokens. The model-connected Go prototype is
retained separately until Rust is verified to replace it. FCC remains Python;
this project does not change its server, streaming, or remote/tunnel support.

## Build and verify

Rust 1.88 or newer, including rustfmt and clippy:

```sh
cargo fmt --all -- --check
cargo clippy --locked --offline --all-targets -- -D warnings
cargo test --locked --offline
cargo run --locked --offline -- --help
cargo run --locked --offline -- theme list
cargo run --locked --offline -- theme show forest
```

The lockfile is committed. There are no third-party crates today. CI checks
Linux, macOS and Windows using the official Rust toolchain; a workflow definition
is not evidence that remote checks have run or passed.

## Themes and stream framing

`theme list` and `theme show NAME` expose the Forest, Catppuccin Mocha,
Ocean, Dusk and Dawn design tokens. They do not configure or render an agent
TUI. Color contrast tests cover text and semantic tokens on both base and
surface colors; these tests do not guarantee accessibility for every UI state.
See the [brand guide](../brand/README.md) for attribution and preview.

The SSE decoder processes UTF-8 bytes incrementally, including split multibyte
characters, BOM, CR/LF/CRLF and multiline `data` fields. It bounds individual
lines and accumulated event fields, rejects malformed UTF-8, and discards
incomplete events on EOF. It deliberately ignores `id` and `retry`; it is not
an HTTP client, JSON parser, reconnect implementation, provider adapter, or
proof of end-to-end streaming. If any frame in a `push` call is invalid, that
call returns an error and its events are not returned; discard the decoder.
Callers must bound each input chunk and enforce overall response/time limits.

`stream::Assembler` accepts already-decoded deltas for exactly one provider
choice. It bounds bytes and call indexes, assembles interleaved fragments and
rejects gaps, reused IDs, mixed assistant content/tool calls and incompatible
finish reasons. A failed append poisons the assembler. Only a finished reply
can be submitted to the session; no partial tool call can be dispatched. The
adapter must validate JSON, tool type, one choice and protocol ordering, and
must discard partial output on interruption or retry. It must not replay chunks
or silently treat an unsupported finish reason as success.

## Boundaries

`Session` is an in-memory, synchronous protocol state machine. It accepts a
prompt and registered tool descriptors, validates complete model replies, and
exposes one next action at a time. A full tool batch is validated before any
call can be claimed; unknown tool names and reused call IDs are rejected.
Registered read-only tools may be claimed directly. Each call declared mutating
requires its own explicit decision; denial is recorded as a tool result.
`begin_execution` claims a call once and `record_result` accepts only the
claimed ID. Calls execute sequentially and every call result is returned in
model order. There are per-item, transcript, batch and turn limits; lengths
count UTF-8 bytes, not tokens or serialized protocol overhead. A call batch is
rejected on the last turn so that it cannot execute without a final model turn.
Errors leave the current operation uncommitted; if a denial/result cannot fit
in the budget, the caller must abort the session, not execute it anyway.

The caller owns model transport, parsing, tool schemas, policy classification,
identity of the approver, timeouts, cancellation and execution. There is no
sandbox or secure filesystem here. In particular, a falsely labelled read-only
tool can write or leak secrets, and a malicious caller can bypass this library.
A claim is an in-memory state transition, not a distributed exactly-once
execution guarantee. Serialized requests, provider response size, credentials,
network redirects and prompt injection also require independent controls.
Never use this library alone as permission to run shell commands or write files.

## Migration gates

1. Add a bounded model adapter with safe credential handling, local/remote
   HTTP(S), streaming and cancellation; test it against real protocol fixtures.
2. Port read-only workspace tools with reviewed OS-level path safety, then test
   parity against the existing Go prototype and real FCC endpoints.
3. Design a real user approval surface, capability-scoped write/process tools,
   OS isolation, process-tree cleanup, adversarial testing and secrets policy.
4. Build persistence, interactive CLI, plugins/MCP and provider integrations in
   independently tested increments; do not claim compatibility in advance.
5. Benchmark startup, latency, memory and reproducible coding tasks on equivalent
   models and permissions. Measure before asserting performance or quality.

No third-party source was copied into this core; see the
[research inventory](../../docs/oh-my-fox/research.md) for inspected projects and
license notes. The root repository license applies.
