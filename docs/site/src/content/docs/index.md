---
title: FCC field guide
description: Install, configure and operate the FCC gateway.
---

[Hướng dẫn tiếng Việt](vi/) · [Previous guide](legacy.html)

<a id="overview"></a>

FCC connects coding tools and API clients to your chosen cloud or local model providers. Python is the supported runtime; the Go port and Oh My Molik are experimental.

<a id="quickstart"></a>

## From zero to ready.

Start locally, configure one provider and verify that the model
is reachable.

**This installs the fork from source.** The commands
use the current development branch of MoLik-Labs/FCC, not the
upstream release. The Go port is experimental; Python remains the
runtime.

### Get the fork

Install [Git](https://git-scm.com/downloads) and
[uv](https://docs.astral.sh/uv/getting-started/installation/). In a terminal or Windows PowerShell, run:

```text
git clone --branch arena/5fe6a43f-fcc https://github.com/MoLik-Labs/FCC.git
cd FCC
uv sync --locked
```

The branch is under development. After integration into
main, choose an available branch as needed. uv can provision
Python 3.14.10. Install coding agents separately; this
source setup does not create desktop shortcuts or install
optional tools.

### Start the server

From the FCC checkout on Windows, macOS or Linux:

```text
uv run fcc-server
```

Keep the terminal open and use the Admin URL from the log.
In a second terminal, return to this checkout to run a
client.

### Connect a model provider

In **Providers**, open a provider, add its key
and save. In **Model Config**, select a model
and click **Apply**. For example, NVIDIA NIM
uses `NVIDIA_NIM_API_KEY`; obtain a key from
[NVIDIA](https://build.nvidia.com/settings/api-keys). The provider controls access, pricing and limits.

If your provider cannot list models, enter its exact ID as
`provider-id/model-id`. For subscription
authentication, use
**Providers → OAuth providers**.

### Launch your agent

Install the agent separately, then launch Claude Code or
Codex from a terminal in the checkout. The launcher can
start FCC automatically.

```text
uv run fcc-claude
# or
uv run fcc-codex
```

Use the agent’s model picker when available. See
[all launchers and integrations](#clients).

Choose **Appearance** in the Admin sidebar for
Forest, Catppuccin Mocha, Nord, Dracula, Ocean, Dusk or Dawn. The
selection is stored in this browser when storage is available, not
on the server or in the separate Molik CLI.

<a id="first-request"></a>

## Make your first request.

Once the server and provider are ready, list models and send a
request locally.

The examples use the default local address. If your server log
shows a different address or port, use that instead. Select an
available model ID from the model list; the example model requires
the corresponding provider to be configured.

LIST AVAILABLE MODELS

```text
curl "http://127.0.0.1:8082/v1/models?view=messages"
```

SEND AN ANTHROPIC MESSAGES REQUEST

```text
curl http://127.0.0.1:8082/v1/messages \
  -H "Content-Type: application/json" \
  -d '{
    "model": "nvidia_nim/nvidia/nemotron-3-ultra-550b-a55b",
    "max_tokens": 1024,
    "messages": [{"role": "user", "content": "Hello!"}]
  }'
```

If **Proxy Authentication** is enabled, also send
`Authorization: Bearer <token>`. Keep tokens out
of shared logs and shell history. A successful response depends on
the provider, its credentials and model availability.

<a id="providers"></a>

## Providers & models.

Use a hosted service, your own compatible endpoint, or a model
running locally.

### Hosted providers

Open **Providers**, add credentials or connect
OAuth, then choose a model. Numbered keys can be added and
reordered; key 1 is the default. Some providers need extra
region, endpoint or account settings.

### Custom providers

Choose
**Providers → Custom providers → Add provider**.
Enter the base URL, API format and optional keys. Leave model
IDs empty for discovery or enter exact IDs, then save and
select a model.

### Local models

Start the model server first. Use `lmstudio/`,
`llamacpp/` or `ollama/` with its exact
model ID. After starting it, refresh models or use
**Test** in its provider card.

### Routing & fallback

`MODEL` is the default. Configure optional tier
overrides, reasoning effort and ordered fallback models in
**Model Config**. A failed request may consume
usage on multiple providers.

Local defaults: LM Studio `http://localhost:1234/v1`,
llama.cpp `http://localhost:8080/v1` and Ollama
`http://localhost:11434`. Prefer tool-capable models
with sufficient context for coding agents.

[Explore the full provider catalog and settings
↗](https://github.com/MoLik-Labs/FCC/blob/main/README.md#choose-a-provider)

<a id="clients"></a>

## Bring your tools.

FCC supports terminal launchers, editor and desktop
integrations, browser sessions, and direct API clients.

For this source checkout, install the agent separately and prefix
the launcher commands below with `uv run` from the FCC
directory.

| TERMINAL CLIENT | COMMAND | TERMINAL CLIENT | COMMAND |
| --- | --- | --- | --- |
| Claude Code | `fcc-claude` | Codex | `fcc-codex` |
| Pi | `fcc-pi` | OpenCode | `fcc-opencode` |
| Cline | `fcc-cline` | Hermes | `fcc-hermes` |
| DeepSeek Harness | `fcc-dsh` | Grok Build | `fcc-grok` |
| Muse Code | `fcc-muse` | Aider | `fcc-aider` |
| Antigravity CLI | `fcc-agy` |  |  |

### Editors and desktop apps

1. Install and start the client on the same computer and user
   environment as FCC.
2. Start FCC; open **Admin → Integrations** and select
   **Connect** on the relevant card.
3. Reload VS Code or restart the app or IDE, then select an FCC
   model where applicable. Use **Disconnect** on the
   same card to undo the integration.

Claude Desktop must be fully quit before connecting or
disconnecting. JetBrains ACP has additional agent setup
requirements. Consult the
[integration-specific instructions](https://github.com/MoLik-Labs/FCC/blob/main/README.md#integrations).

<a id="api"></a>

## Use the API directly.

Point a compatible SDK or HTTP client at FCC. The examples
assume the default local port.

| PROTOCOL | ENDPOINT | STREAMING |
| --- | --- | --- |
| Anthropic Messages | `POST /v1/messages` | Add `"stream": true` |
| OpenAI Chat Completions | `POST /v1/chat/completions` | Add `"stream": true` |
| OpenAI Responses | `POST /v1/responses` | Always streamed |
| Gemini | `POST /v1/models/{model}:generateContent` | Use `:streamGenerateContent?alt=sse` |

OpenAI-compatible SDK base URL:
`http://127.0.0.1:8082/v1`. Anthropic and Gemini base
URL: `http://127.0.0.1:8082`. When proxy authentication
is enabled, pass its token as the API key or bearer token. Model
capabilities and provider-specific behavior still apply.

[See request examples for all four APIs
↗](https://github.com/MoLik-Labs/FCC/blob/main/README.md#use-the-api-directly)

<a id="messaging"></a>

## Work from your phone.

Set up Discord or Telegram in
**Admin → Messaging**, then click Apply.

DISCORD

### Connect a bot

Create a bot in the Discord Developer Portal, enable Message
Content Intent and invite it with read, send, history and
Manage Messages permissions. Configure the token, allowed
channel IDs and an absolute allowed directory.

TELEGRAM

### Connect a bot

Create a bot with BotFather, obtain your numeric user ID and
configure the token, allowed user ID and an absolute allowed
directory. In groups, grant the bot permission to delete
messages.

Restart the server if prompted. `/stats` shows session
state; standalone `/stop` cancels work and standalone
`/clear` clears FCC conversations in the chat. Replying
to a message with these commands targets the selected request or
conversation chain. Voice notes support NVIDIA NIM or optional
local Whisper; see the
[full setup](https://github.com/MoLik-Labs/FCC/blob/main/README.md#use-from-your-phone).

<a id="workloads"></a>

OPERATE / CAPACITY

## Match the workload.

Start conservatively, observe the result, and adjust.

In **Admin → Providers → Runtime configuration**,
workload buttons stage `PROVIDER_MAX_CONCURRENCY`:
Light 1, Balanced 2 (default), High 4. Review the value and click
**Apply**. A process environment override locks the
field and disables the suggestions.

### Light · 1

Start here on constrained machines, particularly with a local
model. Reduce model size or context and avoid optional local
voice when resources are tight.

### Balanced · 2

The default starting point for a typical single-user setup.
Watch memory, response latency and errors.

### High · 4

Try only when the provider quota and local inference server
have headroom. Reduce it if memory or latency increases.

### Measure & adjust

The limit is per provider, not a machine-wide resource cap.
Rate Limit and Rate Window remain independent. Cloud inference
runs on the provider, not your computer.

These are suggestions, not hardware detection or speed guarantees.
Local inference depends on model, context, quantization and the
inference server; concurrency does not change those settings.

<a id="security"></a>

## Operate with care.

Before opening access beyond your machine, understand the trust
boundary.

**Admin and Code endpoints have no separate Admin
authentication.**
Remote clients that can reach them may change configuration or
operate code sessions. Proxy Authentication protects the proxy
API; do not assume it protects Admin or Code. Do not expose them
to untrusted networks or the public internet without independent
access controls.

### Remote access and tunnels

Cloudflare Tunnel, ngrok and Tailscale can provide remote
connectivity when configured appropriately. Restrict access at the
network or tunnel layer with controls you trust; a tunnel URL by
itself is not an authorization policy. Review its provider’s
access configuration and test that unauthorized users cannot reach
Admin or Code before use. Remote HTTP transmits tokens and request
contents without HTTP transport encryption; use HTTPS or a
properly protected private tunnel, including for Claude Desktop.
Preserve the server and credentials on a trusted host.

### Credential hygiene

- Use the least privilege available for provider keys, bots and
  tunnel access. Restrict messaging users, channels and allowed
  directories.
- Never commit tokens or paste them into public issues. Redact
  logs and diagnostic output before sharing.
- Review provider terms, model pricing, and the directory
  accessible to code sessions and messaging agents.

<a id="maintenance"></a>

## Keep things healthy.

Check versions, update deliberately, and know what an uninstall
removes.

| TASK | COMMAND OR ACTION |
| --- | --- |
| Check version | `uv run fcc-server --version` |
| Diagnose | `uv run fcc-doctor` |
| Update | Stop FCC, then run `git pull --ff-only` and `uv sync --locked` in the checkout. |
| Uninstall | Stop FCC and connected apps, back up data, then remove the source checkout if no longer needed. |

The source checkout contains its own `.venv`; FCC data
such as `~/.fcc/` is separate and is not removed with
the checkout. In this revision, `fcc-update` is
disabled and makes no changes; older installed copies may still
invoke the upstream installer. The inherited install scripts also
target upstream releases, not this checkout. Automatic fork
release publication is disabled pending a separate
package-identity review.

[View platform-specific maintenance instructions
↗](https://github.com/MoLik-Labs/FCC/blob/main/README.md#manage-your-installation)

<a id="troubleshooting"></a>

## When things go wrong.

Start with the smallest failing connection and work outward.

Admin does not open or the API cannot connect

Confirm FCC is running and use the address and port printed in
its log. Check local firewall or tunnel access controls. Do
not expose Admin publicly as a workaround.


The model is missing or a request fails

Check provider credentials, model ID, account access and
provider availability. For local providers, start the model
server and refresh models or use Test. Verify that the model
supports the tools and context your agent needs.


An editor or desktop client does not see FCC

Ensure FCC and the client run in the same user environment,
reconnect in Admin, and reload or restart the client. Some
integrations require a complete quit or additional setup.


I need to report a bug

Run `uv run fcc-doctor` from the checkout,
reproduce the problem and record the FCC version, active
model, model mappings, steps and complete error. Redact
credentials and private data before sharing. Issues are
currently disabled on this fork; use the upstream tracker only
for upstream issues and releases.

<a id="development"></a>

## Build and verify.

Python is the supported application. Go is an experimental,
incomplete port, not a replacement server.

For this checkout, install
[uv](https://docs.astral.sh/uv/) and use its pinned
Python and locked dependencies:

```text
git clone --branch arena/5fe6a43f-fcc https://github.com/MoLik-Labs/FCC.git
cd FCC
uv sync --locked
uv run fcc-server
```

Run `./scripts/ci.sh` on macOS/Linux or
`.\scripts\ci.ps1` on Windows. Individual checks
include Ruff, Pyrefly and pytest. Review
[CONTRIBUTING.md](https://github.com/MoLik-Labs/FCC/blob/main/CONTRIBUTING.md)
before opening a pull request. Repository-wide type checks and
some historical benchmark generators have known failures on the
development branch; see the checkpoints rather than assuming all
gates pass.

For Go parity work, read
[gocore/README.md](https://github.com/MoLik-Labs/FCC/blob/main/gocore/README.md)
and the
[rewrite plan](https://github.com/MoLik-Labs/FCC/blob/main/docs/rewrite/00-master-plan.md). Fixture success alone does not establish current-provider or
full-service parity.

IN DEVELOPMENT / SEPARATE PROJECT

### Oh My Molik

The Rust and Go prototypes connect to Chat endpoints and
inspect files read-only. Rust also offers a limited read-only
MCP stdio server. Neither is a complete coding agent.

[Explore the project →](https://github.com/MoLik-Labs/FCC/blob/main/ohmymolik/README.md)
  
[View theme collection →](molik.html)   
[MCP setup and security →](molik.html#mcp-stdio)

<a id="pages"></a>

## Publish this site

This site uses Astro Starlight. Run `npm ci` and `npm run build` in `docs/site`. GitHub Actions builds pull requests and publishes `main` when GitHub Pages is configured to use Actions. A live deployment requires repository permissions and a successful workflow.

<a id="resources"></a>

## Go further.

Detailed references and project information, all in one place.

[**Full project README** — Provider catalog, API examples and integrations ↗](https://github.com/MoLik-Labs/FCC/blob/main/README.md)

[**Contributing**Development standards and quality checks↗](https://github.com/MoLik-Labs/FCC/blob/main/CONTRIBUTING.md)

[**Upstream issues**For the upstream project and its releases↗](https://github.com/Alishahryar1/free-claude-code/issues)

[**License**GNU AGPL version 3 only↗](https://github.com/MoLik-Labs/FCC/blob/main/LICENSE)
