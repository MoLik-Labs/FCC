# Oh My Fox: research and design checkpoint

## Source inventory

The following public repositories were cloned as shallow partial clones with
`--no-checkout --filter=blob:none --depth=1` outside the FCC tree, at the
revisions below. Selected README, license and source files were read with
`git show`; the entire trees were not checked out or exhaustively audited.
The copies live at `/home/user/oh-my-fox-research` in this environment and
are temporary, not part of the repository or guaranteed to persist. No code
from them is vendored here. Package/subdirectory licenses and notices require
individual review before reuse. To inspect one of the listed revisions after a
reset, for example:

```sh
git clone --no-checkout --filter=blob:none --depth=1 https://github.com/can1357/oh-my-pi.git oh-my-pi-research
git -C oh-my-pi-research fetch --depth=1 origin b505770284deff6b2888a4eb26a7cc95497cc111
git -C oh-my-pi-research show FETCH_HEAD:packages/agent/src/agent-loop.ts
```

Use the corresponding URL and SHA in the table for other projects. Fetching
blobs on demand requires network access; respect each project's license.

| Project | Revision | Root license |
| --- | --- | --- |
| [Oh My Pi](https://github.com/can1357/oh-my-pi) | `b505770284deff6b2888a4eb26a7cc95497cc111` | MIT |
| [Pi](https://github.com/earendil-works/pi) | `9c4652688590abb836b40ebd1e4725ecf9a7c7d1` | MIT |
| [Codex](https://github.com/openai/codex) | `cf3ff1530d821ec2d26452da4971d8e3e7d18d25` | Apache-2.0 |
| [Claude Code public repository](https://github.com/anthropics/claude-code) | `1006c2e4f783f016963fb290d3eb4afaa2d448b1` | Commercial terms; see LICENSE.md |
| [OpenCode](https://github.com/anomalyco/opencode) | `1a215a4b2082ab66e016e130be37be097d7234f4` | MIT |
| [DeepSeek Harness](https://github.com/deepseek-ai/deepseek-harness) | `4b399c39a6d1479cd84ca591b4d7c88e24d38a96` | MIT |
| [Gemini CLI](https://github.com/google-gemini/gemini-cli) | `eecccfcd5aff25de5e78135ca20f7abe3f566955` | Apache-2.0 |
| [Goose](https://github.com/aaif-goose/goose) | `36b29113aff691e1a46a21a5c25900824caf9949` | Apache-2.0 |
| [Aider](https://github.com/Aider-AI/aider) | `5dc9490bb35f9729ef2c95d00a19ccd30c26339c` | Apache-2.0 |
| [OpenHands / Agent Canvas](https://github.com/OpenHands/OpenHands) | `3e9cf433b0669661833c2025f8d25ac9ba417d31` | MIT |
| [Cline](https://github.com/cline/cline) | `8075d2e35cf5cad67c24f180c80a05d311b74a3e` | Apache-2.0 |
| [Hermes Agent](https://github.com/NousResearch/hermes-agent) | `0f602c5f6f18615632535251d62a2db74ce93b69` | MIT |

The public Claude Code repository declares commercial terms and does not
establish rights to copy or reimplement its closed runtime from code. Treat it
as documentation only; do not copy proprietary materials. License labels are
inventory notes, not legal advice or blanket permission to transplant files.

## Observations, not feature-completeness claims

- Pi exposes a small agent loop, tool contracts and extension points. Oh My Pi
  builds on a larger tool stack and has a Rust native core. Their public READMEs
  and `packages/agent/src/agent-loop.ts` were inspected; this prototype is not
  a port or extension-compatible with either.
- Codex exposes explicit tool approval and sandbox machinery; OpenCode separates
  plan and build agents. These illustrate the importance of consent and isolation
  before introducing writes and process execution.
- DeepSeek Harness describes a plugin-first architecture; its session controller
  has explicit stream and session state. Compatibility promises require actual
  protocol and lifecycle tests, not only a plugin interface.
- The other listed repositories were cloned and their top-level descriptions
  inspected; they have not received a comprehensive feature or security review.
  There is no evidence for a claim that Oh My Fox is the strongest agent.

## Language choice and migration

**Rust is the selected implementation language for Oh My Fox**, replacing the
previous provisional Go choice. Neither language alone guarantees performance
or security; no comparative Oh My Fox benchmark has been run. FCC's existing
Python runtime, remote access and streaming are unaffected by this decision.
The Go parity work for FCC is separate from Oh My Fox.

`ohmyfox/rust/` now contains a dependency-free state machine and a minimal
help/version CLI. It validates bounded batches and tracks per-call approval and
execution claims, but does not connect to a model, perform I/O, sandbox, or
implement a user approval UI. The Go implementation at `ohmyfox/` remains the
only runnable read-only model-connected prototype until Rust reaches tested
parity; it is not a complete coding agent. Go 1.22 was the local bootstrap/test
compatibility target, not a release recommendation.

## Current prototype and gates

The Go prototype has an explicit remote-or-local OpenAI-compatible Chat
endpoint and two workspace inspection tools. It has no editing, shell, approval
UI, sandbox, streaming, session persistence, TUI, MCP, plugins, direct
Anthropic/Gemini protocol, model routing, parallel subagents or FCC launcher
integration. Prompts and requested file contents are sent to the selected
endpoint. See the [Rust architecture and migration plan](../../ohmyfox/rust/README.md)
for the staged replacement and verification requirements.

Before adding mutation or command execution: implement reviewed capability-based
approval, a real OS/container isolation strategy, cancellation and process-tree
cleanup, secrets policy, adversarial tests and cross-platform verification.
Path checks in the prototype are best effort and cannot prevent symlink races
or replace a sandbox. Evaluate prompt injection, malicious files, network
redirects, huge responses and token leakage. Avoid running unknown plugin hooks.

Establish representative cold/warm startup, latency, throughput, memory and
real coding-task benchmarks with comparable models, prompts and tool permissions
before making performance or quality claims. Continue FCC Go parity separately;
the existing Python runtime and its tunnel functionality remain authoritative.
