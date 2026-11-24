# Oh My Molik (formerly Oh My Fox): research and design checkpoint

## Source inventory

The following selected public repositories were cloned as shallow partial clones with
`--no-checkout --filter=blob:none --depth=1` outside the FCC tree, at the
revisions below. Revisions are pinned at or before 2026-10-27 UTC; earlier
research notes used newer checkout tips and are superseded by this snapshot.
Selected README, license and source files were read with
`git show`; the entire trees were not checked out or exhaustively audited.
The copies live at `/home/user/oh-my-fox-research` in this environment and
are temporary, not part of the repository or guaranteed to persist. No code
from them is vendored here. Package/subdirectory licenses and notices require
individual review before reuse. To inspect one of the listed revisions after a
reset, for example:

```sh
git clone --no-checkout --filter=blob:none --depth=1 https://github.com/can1357/oh-my-pi.git oh-my-pi-research
git -C oh-my-pi-research fetch --depth=1 origin 79a26141c905f713903b0784a74337d7f4e27225
git -C oh-my-pi-research show FETCH_HEAD:packages/agent/src/agent-loop.ts
```

Use the corresponding URL and SHA in the table for other projects. Fetching
blobs on demand requires network access; respect each project's license.

| Project | Revision | Root license |
| --- | --- | --- |
| [Oh My Pi](https://github.com/can1357/oh-my-pi) | `79a26141c905f713903b0784a74337d7f4e27225` | MIT |
| [Pi](https://github.com/earendil-works/pi) | `a92283056d915c758799fcd904948ff0cb42628c` | MIT |
| [Codex](https://github.com/openai/codex) | `e8e6a922c0008bc08bacc69dd2ce9b6c3c12944c` | Apache-2.0 |
| [Claude Code public repository](https://github.com/anthropics/claude-code) | `58d6da1e34d6cd1c6379dcbf568bd66f5f9b6102` | Commercial terms; see LICENSE.md |
| [OpenCode](https://github.com/anomalyco/opencode) | `891a020bcc935457a34585a6d6a0ea6fb5b92eec` | MIT |
| [DeepSeek Harness](https://github.com/deepseek-ai/deepseek-harness) | `f3803c96974ad4a94195fb97e5c637e9c2a3afe3` | MIT |
| [Gemini CLI](https://github.com/google-gemini/gemini-cli) | `fa5c17365754a4936c3b97429df798966305267d` | Apache-2.0 |
| [Goose](https://github.com/aaif-goose/goose) | `8a16788137c8af25f9d6d17ee87bbedb53ea9325` | Apache-2.0 |
| [Aider](https://github.com/Aider-AI/aider) | `5dc9490bb35f9729ef2c95d00a19ccd30c26339c` | Apache-2.0 |
| [OpenHands / Agent Canvas](https://github.com/OpenHands/OpenHands) | `6354932e87f99e1583f482f2cea49d6d573443e0` | MIT |
| [Cline](https://github.com/cline/cline) | `292a9e14cca3f6d03f5a5277842b291c07d5aae4` | Apache-2.0 |
| [Hermes Agent](https://github.com/NousResearch/hermes-agent) | `392c6f7d9ec64db759baf81c6321c34a78f85898` | MIT |
| [Kilo Code](https://github.com/Kilo-Org/kilocode) | `307ff45f27ec796309ffb4e2f00f8dab85db3838` | MIT |
| [Qwen Code](https://github.com/QwenLM/qwen-code) | `127c2ab9048839c64b806e58e6d7329cbb2d9a8a` | Apache-2.0 |
| [SWE-agent](https://github.com/SWE-agent/SWE-agent) | `3ea751c087f32b16e039a2233dd6eefecef325d5` | MIT |
| [Continue](https://github.com/continuedev/continue) | `5522c6f44ca0ac3528b37244818fbfa39b5af470` | Apache-2.0 |
| [Roo Code](https://github.com/RooVetGit/Roo-Code) | `b867ec9145750d0ae1ff7f02d35406e9bf2a0b16` | Apache-2.0 |

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
  There is no evidence for a claim that Oh My Molik is the strongest agent.

## Focused source study: decisions and gaps

The following files were inspected at the revisions above. These are selected
entry points, not exhaustive audits of the projects or their full histories.
Links include the pinned SHA to keep the observations reproducible.

| Source examined | Observation | Oh My Molik decision / gap |
| --- | --- | --- |
| [Pi agent loop](https://github.com/earendil-works/pi/blob/a92283056d915c758799fcd904948ff0cb42628c/packages/agent/src/agent-loop.ts) | Events and streaming are separated from the model-specific stream function. | Keep state machine independent of transport; actual event bus and model adapter remain to be built. |
| [Oh My Pi loop](https://github.com/can1357/oh-my-pi/blob/79a26141c905f713903b0784a74337d7f4e27225/packages/agent/src/agent-loop.ts), [approval](https://github.com/can1357/oh-my-pi/blob/79a26141c905f713903b0784a74337d7f4e27225/packages/coding-agent/src/tools/approval.ts) | Explicit interruption handling and read/write/exec tiers matter; live settings and approval caches add complexity. | Assemble bounded fragments before dispatch and require per-call mutation consent. No cached grants or real interruption handling yet. |
| [Codex sandbox adapter](https://github.com/openai/codex/blob/e8e6a922c0008bc08bacc69dd2ce9b6c3c12944c/codex-rs/core/src/sandboxing/mod.rs), [approvals](https://github.com/openai/codex/blob/e8e6a922c0008bc08bacc69dd2ce9b6c3c12944c/codex-rs/core/src/tools/approvals.rs) | Execution context carries sandbox, network, expiration and approval metadata. | Never equate approval with isolation. Process execution stays unimplemented until OS-specific enforcement and tests exist. |
| [OpenCode permissions](https://github.com/anomalyco/opencode/blob/891a020bcc935457a34585a6d6a0ea6fb5b92eec/packages/core/src/permission.ts), [Kilo permissions](https://github.com/Kilo-Org/kilocode/blob/307ff45f27ec796309ffb4e2f00f8dab85db3838/packages/core/src/permission.ts) | Permissions address action and resource; missing child-agent permissions can be denied. | The current tool-name/access-only model is intentionally narrower; resource scopes and child-agent policy are future work, not claimed features. |
| [DeepSeek session stream](https://github.com/deepseek-ai/deepseek-harness/blob/f3803c96974ad4a94195fb97e5c637e9c2a3afe3/packages/api/session-controller/src/assistant-stream.ts) | Stream attempts have separate accumulators and snapshots. | Never merge retries or replay partial chunks automatically. SSE and delta accumulation exist as separate non-network components only. |
| [Qwen agent loop](https://github.com/QwenLM/qwen-code/blob/127c2ab9048839c64b806e58e6d7329cbb2d9a8a/packages/core/src/core/agent-loop.ts) | Stop reasons distinguish cancellation, limits and duplicate provider calls; retry attempts need explicit boundaries. | Reject repeated IDs and incomplete tool batches; real retries, cancellation and provider adapters remain open. |
| [Roo command approvals](https://github.com/RooVetGit/Roo-Code/blob/b867ec9145750d0ae1ff7f02d35406e9bf2a0b16/src/core/auto-approval/commands.ts) | Shell syntax and substitution make naive command allowlists dangerous. | Do not add shell execution or treat regex/prefix matching as a sandbox. |

No implementation code was transplanted from these sources. Features should
only be called compatible after protocol, security and real-world tests. There
is no objective finite list of all “best” harnesses, and proprietary runtimes
cannot be fully cloned or reviewed; extend this inventory as evidence warrants.

## Language choice and migration

**Rust is the selected implementation language for Oh My Molik**, replacing the
previous provisional Go choice. Neither language alone guarantees performance
or security; no comparative Oh My Molik benchmark has been run. FCC's existing
Python runtime, remote access and streaming are unaffected by this decision.
The Go parity work for FCC is separate from Oh My Molik.

`ohmymolik/rust/` includes a state machine, stream framing, theme tokens and an
experimental non-streaming Chat adapter using system curl. It can list and read
files in a trusted workspace, but has no editing, shell, sandbox, TUI or user
approval UI. The Go implementation at `ohmyfox/` remains available; complete
parity, live providers and platform security remain unverified. Neither is a
complete coding agent. Go 1.22 was the local bootstrap/test
compatibility target, not a release recommendation.

## Current prototype and gates

The Go prototype has an explicit remote-or-local OpenAI-compatible Chat
endpoint and two workspace inspection tools. It has no editing, shell, approval
UI, sandbox, streaming, session persistence, TUI, MCP, plugins, direct
Anthropic/Gemini protocol, model routing, parallel subagents or FCC launcher
integration. Prompts and requested file contents are sent to the selected
endpoint. See the [Rust architecture and migration plan](../../ohmymolik/rust/README.md)
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
