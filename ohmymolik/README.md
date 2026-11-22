# Oh My Molik

<img src="brand/molik-mark.svg" alt="" width="80" height="80">

**An independent developer-agent project by MoLik Labs, in development.**
Previously called Oh My Fox, it is being built in Rust. FCC remains the separate
Python gateway; this rename does not change FCC's streaming or remote support.

## Current status

The Rust core includes a bounded sequential session state machine, an incremental
SSE framing parser and five accessible color-token sets. The CLI can display its
help, version and theme tokens. It **does not yet connect to a model, render a
TUI, persist themes, execute tools, or implement a sandbox or an approval UI**.
It is not yet a complete coding agent. The earlier model-connected Go read-only
prototype is retained at [`../ohmyfox/`](../ohmyfox/README.md); its paths and
`OMF_API_KEY` stay unchanged for compatibility, and it has its own limitations.

## Try the Rust core

From `ohmymolik/rust/` with Rust 1.88 or newer:

```sh
cargo test --locked --offline
cargo run --locked --offline -- --help
cargo run --locked --offline -- theme list
cargo run --locked --offline -- theme show forest
cargo run --locked --offline -- theme show catppuccin-mocha
```

The five available token sets are `forest`, `catppuccin-mocha`, `ocean`, `dusk`
and `dawn`. These commands inspect tokens only; they do not configure a TUI.
Open the [theme collection](../docs/site/molik.html) through the built
[documentation site](../docs/site/README.md) for visual mockups. The
[brand guide](brand/README.md), [Rust architecture](rust/README.md) and
[research and roadmap](../docs/oh-my-fox/research.md) document provenance,
safety boundaries and remaining work. No top-ranking or performance claim has
been independently established.
