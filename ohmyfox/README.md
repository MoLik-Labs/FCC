# Oh My Fox — experimental read-only prototype

Oh My Fox is a separate Go module exploring a coding-assistant tool loop. It is
**not** a complete coding agent or an Oh My Pi fork, and has not been benchmarked
against other products. It currently sends a prompt to an explicitly selected
OpenAI-compatible Chat Completions endpoint and can list and read text files
inside the selected workspace. It does not edit files, execute commands, stream,
persist sessions or implement a terminal UI, plugins, MCP, or subagents.

## Run

With a supported Go toolchain, from `ohmyfox/`:

```sh
go test ./...
go run ./cmd/ohmyfox \
  --endpoint http://127.0.0.1:8082/v1/chat/completions \
  --model 'provider/model-id' --workspace /path/to/project \
  'Explain the project structure'
```

Replace the example URL and model with an active FCC or compatible endpoint
and model. HTTPS remote/tunnel endpoints also work; the endpoint is explicitly
chosen by the user and is **not limited to localhost**. If authentication is
enabled, supply the token in the `OMF_API_KEY` environment variable. Do not
put it in the URL or command line. HTTP transmits the token and content without
transport encryption; prefer HTTPS or an appropriately protected private tunnel.
The request path is a full Chat Completions URL, not an automatically inferred
base URL. The current API is non-streaming and may take time to respond.

A command-line prompt may be visible in shell history or process listings. To
avoid putting it in the command line, pass `--prompt-stdin` and provide the
prompt on standard input instead. Never include credentials in a prompt.

The model provider receives your prompt and any files it requests. Only use a
trusted endpoint and workspace; source files can contain secrets. Git metadata
and symlink paths are rejected, and files, directory listings, model responses,
number of tool calls and turns have limits. These restrictions are **not an OS
sandbox**: concurrent filesystem changes can race path checks, and allowed
files (including configuration and secrets) can still be read and transmitted.
Do not use on an untrusted or concurrently modified workspace. Ctrl+C cancels
the request; no persistent session or automatic retry is provided.

The prototype was developed and tested with the repository's Go 1.22 bootstrap;
that is not a recommendation to deploy an unsupported Go version. See
[research and remaining gates](../docs/oh-my-fox/research.md) before considering
release or expanding its capabilities. Python remains FCC's supported runtime.
