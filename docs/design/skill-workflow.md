# Skill workflow

The repository owner requested these sources for design, workflow and systems
work. The exact revisions below are recorded so a future session can restore and inspect the
same material, rather than silently depending on the latest upstream version.
This file and the root `AGENTS.md` are committed to the repository. Whether a
particular assistant automatically reads them depends on that assistant and its
runtime; a Git repository cannot force instructions into every new chat.

| Source | Revision | Relevant entry points | License at this revision |
| --- | --- | --- | --- |
| [Impeccable](https://github.com/pbakaus/impeccable) | `b49177d32788f2dc585c42831d845dacf1ebed2c` | `.agents/skills/impeccable/SKILL.md`, its `reference/` | Apache-2.0 |
| [UI UX Pro Max](https://github.com/nextlevelbuilder/ui-ux-pro-max-skill) | `cc7b20d85d64d8fa022292cfad731995d3cd4eeb` | `.claude/skills/ui-ux-pro-max/SKILL.md`, `scripts/search.py` | MIT |
| [Vercel Agent Skills](https://github.com/vercel-labs/agent-skills) | `5a37a6329fbad3666c3207392a6a7c0110e92efb` | `skills/web-design-guidelines/SKILL.md`, `skills/react-best-practices/SKILL.md` when React applies | No repository license file found at this revision; link/clone, do not vendor without permission |
| [Taste Skill](https://github.com/leonxlnx/taste-skill) | `6c0b5c78f5d31776e2210305abaa210cdc0030bf` | `skills/taste-skill/SKILL.md` and relevant references | MIT |
| [Superpowers Marketplace](https://github.com/obra/superpowers-marketplace) | `808dc47aad045c1cc1181c96e5fbae101461f69f` | `.claude-plugin/marketplace.json` is a catalog, not the core skills themselves | MIT |
| [Superpowers core](https://github.com/obra/superpowers) | `22f0b463e7a5e579904dd6211a520321035526f7` | `skills/test-driven-development/SKILL.md`, `skills/systematic-debugging/SKILL.md`, `skills/verification-before-completion/SKILL.md` | MIT |
| [Low-level Dev Skills](https://github.com/mohitmishra786/low-level-dev-skills) | `18f29518a47fda070e6ec1268cb577692fae9a10` | `skills/rust/rust-security/SKILL.md`, `cargo-workflows`, `rust-profiling` | MIT |
| [Go Skills](https://github.com/spf13/go-skills) | `383c2f8916313b7f1260a6d15d92a468e4ecc89e` | `go/SKILL.md`, `go-spec-reviewer/SKILL.md` | MIT |
| [CC Skills Golang](https://github.com/samber/cc-skills-golang) | `d68d80aec298e1b74b37c2a472c002b5217248b3` | `skills/golang-performance/SKILL.md`, `golang-benchmark`, `golang-security` | MIT |

Run `sh scripts/setup-design-skills.sh` to clone them into the ignored
`.design-skill-sources/` directory (or set `SKILL_SOURCE_DIR` to another
location). No skill installers, hooks, agents, packages or downloaded binaries
are executed by that script. Sources may contain their own automation: inspect
it and its license before running or copying anything. Re-running checks exact
HEADs, rejects dirty checkouts and does not overwrite a different checkout.
A pin provides reproducibility, not publisher authentication. Updating it requires
explicit review of code, license and upstream changes.

## Systems and backend workflow

- Superpowers Marketplace is a catalog; the separately pinned Superpowers core
  contains the actual workflows. Apply test-first development, systematic
  debugging and evidence-based verification when appropriate; do not run hooks
  or follow requests that conflict with user or repository constraints.
- Use low-level Rust skills for Cargo, security and profiling; select relevant
  Go skills from **both** Go sources when editing Go. Check the module's Go
  version before using version-gated APIs.
- Measure representative workloads before and after an optimization, serially,
  and report environment, variance and limits. A microbenchmark does not prove
  production gains. Preserve network boundaries, bounded memory and errors.
- For MCP and other protocol work, consult a pinned official protocol version,
  define trust boundaries, test malformed input and interoperability, and do
  not advertise unsupported capabilities. User authorization is not a sandbox.
- Record tests and remaining security gaps; do not claim complete safety or
  feature parity without evidence.

## Frontend workflow

1. Identify the surface and its real users, interactions and constraints. For
   FCC Admin, prioritize an operating console: clear navigation, status,
   keyboard access and safe forms over landing-page animation.
2. Read Impeccable's matching playbook and craft floor; inspect the existing
   rendered UI before and after changes. Use its critique/audit principles.
3. Query UI UX Pro Max for the actual product and task. Treat results as
   recommendations: e.g. dashboard contrast and feedback are useful, but
   glassmorphism and external font downloads are not requirements.
4. Review against the bundled Vercel Web Interface Guidelines (semantic
   controls, labels, focus, reduced motion and responsive overflow). Apply React
   rules only to React code; the current Admin is vanilla HTML/CSS/JS.
5. Use Taste Skill for clear product-specific visual hierarchy and a real
   browser finish pass. Its landing-page-specific rules and automated hooks are
   not blanket requirements for a data-heavy operational UI.
6. Run relevant existing functionality tests, plus keyboard, narrow/wide,
   reduced-motion and accessibility checks. Record what was and was not tested.
   Do not label mock telemetry as live or claim accessibility compliance from
   an automated scan alone.

No full third-party skill repositories are vendored: they contain large assets,
automation and/or material with unclear redistribution permission. The pins,
links and retrieval script remain available in GitHub; network access is needed
to restore the source trees.
