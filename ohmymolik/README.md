# Oh My Molik

<img src="brand/molik-mark.svg" alt="" width="80" height="80">

**An independent developer-agent project by MoLik Labs, in development.**
Previously called Oh My Fox, it is being built in Rust. FCC remains the separate
Python gateway; this rename does not change FCC's streaming or remote support.

## Current status

The Rust CLI can use an explicitly selected local or remote OpenAI-compatible
Chat Completions endpoint to answer a prompt and inspect a trusted workspace
with `read_file` and `list_files`. It uses system `curl` and is non-streaming.
The endpoint receives the prompt and any file contents the model requests.
It has **no editing, shell, sandbox, approval UI, TUI, theme persistence or
session persistence**, and is not yet a complete coding agent. The earlier Go
prototype is retained at [`../ohmyfox/`](../ohmyfox/README.md), with its
historical path and `OMF_API_KEY` unchanged for compatibility.

## Try the Rust read-only preview

Install Rust 1.88+ and curl (version 7.63+); from `ohmymolik/rust/`:

```sh
cargo test --locked
# If the endpoint requires authentication, set OMM_API_KEY in your shell.
# Do not put a token in the URL or command line.
printf '%s' 'Explain the project structure' | cargo run --locked -- run \
  --endpoint http://127.0.0.1:8082/v1/chat/completions \
  --model 'provider/model-id' --workspace /path/to/project --prompt-stdin
```

Replace the URL, model and workspace with your own. HTTPS remote and tunnel
endpoints are supported, not limited to localhost. HTTP transmits prompts,
files and credentials without transport encryption: prefer HTTPS or a protected
private tunnel. The endpoint is a complete Chat Completions URL; credentials,
query strings and fragments in the URL are rejected. Curl uses normal TLS
certificate and proxy configuration; prompts can also be visible to a proxy.
The experimental implementation requires the system `curl` executable in PATH;
it ignores curlrc and does not follow redirects. Its timeout is 120 seconds per
model request. Ctrl+C handling and process-tree cancellation have not been
verified across platforms. The model endpoint must be trusted.

The workspace checks reject parent paths, symlinks and Git metadata, but are
best effort, not OS isolation: concurrent filesystem changes, hardlinks,
platform-specific filesystem behavior and sensitive files can defeat the
intended boundary. Only use on a trusted, stable workspace and inspect what you
send. No automatic secret filtering is provided. A command-line prompt can appear
in shell history or process listings; `--prompt-stdin` avoids that exposure.

## Theme tokens

```sh
cargo run --locked -- theme list
cargo run --locked -- theme show forest
cargo run --locked -- theme show catppuccin-mocha
```

The five available token sets are `forest`, `catppuccin-mocha`, `ocean`, `dusk`
and `dawn`. These commands inspect tokens only; they do not configure a TUI.
Open the [theme collection](../docs/site/molik.html) through the built
[documentation site](../docs/site/README.md) for visual mockups. The
[brand guide](brand/README.md), [Rust architecture](rust/README.md) and
[research and roadmap](../docs/oh-my-fox/research.md) document provenance,
safety boundaries and remaining work. No top-ranking or performance claim has
been independently established.

## Experimental read-only MCP server

A trusted MCP client can launch the Molik binary over **stdio**. Build first
with `cargo build --locked --release`, then configure the client with the
absolute binary and workspace paths, for example:

```json
{
  "mcpServers": {
    "oh-my-molik": {
      "command": "/absolute/path/to/oh-my-molik",
      "args": ["mcp", "serve", "--workspace", "/absolute/path/to/trusted/project"]
    }
  }
}
```

The client-specific location and format of this configuration can vary; these
paths are examples and must be replaced. This server implements a limited
`2025-06-18` MCP stdio lifecycle, `ping`, `tools/list` and `tools/call` for
`list_files` and `read_file` only. It does **not** make Molik an MCP client and
does not provide HTTP transport, resources, prompts, editing or shell tools.
An unsupported requested version negotiates down to `2025-06-18`; a client that
does not support it should disconnect.

**Trust boundary:** starting it grants the client read access to files in the
selected workspace, and the client may send results to its model or elsewhere.
There is no separate authentication or per-call approval in stdio mode: rely on
OS process permissions, use a trusted client and workspace, and do not expose
stdio through an unprotected network service. The same best-effort filesystem
checks and race/hardlink limitations described above apply; this is not a
sandbox or secret filter. File content is untrusted input, not instructions.
Individual input frames are capped at 256 KiB, files at 64 KiB, listings at 200
entries and tool invocations at 120 per rolling minute per server process. The
rate limit is local to that process and is not a distributed protection.
See the [verification notes](../docs/oh-my-fox/verification.md) for the exact
SDK and test coverage, which are not a broad MCP interoperability guarantee.
