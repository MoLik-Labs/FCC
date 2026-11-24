# Oh My Molik: Rust core

This is an experimental read-only assistant, not a complete coding agent.
The CLI supports non-streaming Chat Completions via system curl and bounded
workspace inspection, as well as help, version and theme-token inspection. The
Go prototype remains available during migration. FCC remains Python;
this project does not change its server, streaming, or remote/tunnel support.

## Build and verify

Rust 1.88 or newer, including rustfmt and clippy:

```sh
cargo fmt --all -- --check
cargo clippy --locked --all-targets -- -D warnings
cargo test --locked
cargo run --locked -- --help
cargo run --locked -- theme list
cargo run --locked -- theme show forest
```

The committed lockfile pins `serde_json` and transitive dependencies from
crates.io (MIT / Apache-2.0 and other listed upstream licenses). Fresh builds
need registry access, after which Cargo can work offline with cached crates.
System curl 7.63+ is required at runtime. CI is configured for Linux, macOS and
Windows; a workflow definition is not evidence that remote checks have passed.
See the [CLI guide](../README.md) for setup and the security limitations.

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

The library caller remains responsible for its transport, parsing, tool schemas,
policy classification, approver identity, timeouts and cancellation. The
experimental CLI supplies a non-streaming adapter using a curl subprocess; it
is not a general transport or a user approval implementation. There is no
sandbox or secure filesystem here. In particular, a falsely labelled read-only
tool can write or leak secrets, and a malicious caller can bypass this library.
A claim is an in-memory state transition, not a distributed exactly-once
execution guarantee. Serialized requests, provider response size, credentials,
network redirects and prompt injection also require independent controls.
Never use this library alone as permission to run shell commands or write files.

## Migration gates

1. Replace the curl subprocess with a reviewed native HTTPS transport, add
   streaming, cancellation and protocol replay tests across providers.
2. Replace best-effort workspace checks with reviewed OS-level path safety;
   test parity against the Go prototype and real FCC/provider/tunnel deployments.
3. Design a real user approval surface, capability-scoped write/process tools,
   OS isolation, process-tree cleanup, adversarial testing and secrets policy.
4. Build persistence, interactive CLI, plugins/MCP and provider integrations in
   independently tested increments; do not claim compatibility in advance.
5. Benchmark startup, latency, memory and reproducible coding tasks on equivalent
   models and permissions. Measure before asserting performance or quality.

No third-party source is vendored into this core; see the
[research inventory](../../docs/oh-my-fox/research.md) for inspected projects and
license notes. The root repository license applies.
