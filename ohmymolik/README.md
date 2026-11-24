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
