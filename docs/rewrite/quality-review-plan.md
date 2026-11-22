# FCC quality review plan

This is an incremental work plan, not an assertion that the review or rewrite
is complete. The active branch retains the Python implementation as the
executable specification and the Go work as a partial port.

## Skill inventory (2026-10-27)

The following source repositories were cloned at these revisions, and their
skill directories and source licenses were copied to `/home/user/skill-installed/`
for manual use in this environment. Compressed copies are kept under the
ignored local `agent_workspace/skills/` directory; extract each archive outside
the repository to restore the files after an environment reset. Keeping loose
third-party Python files under the repository interferes with its recursive CI
suppression check. These are local working files, not vendored project
dependencies; they are not committed or automatically installed for other
contributors.

| Library | Revision | Local manual-use directory | Contents |
| --- | --- | --- | --- |
| [Taste](https://github.com/Leonxlnx/taste-skill) | `f6b76575f949707b4e3c37cd86cb44b0b98643e8` | `taste-skill` | One frontend skill with assets and scripts |
| [Impeccable](https://github.com/pbakaus/impeccable) | `b49177d32788f2dc585c42831d845dacf1ebed2c` | `impeccable` | One design skill with references and launcher |
| [Superpowers](https://github.com/obra/superpowers) | `22f0b463e7a5e579904dd6211a520321035526f7` | `superpowers` | 15 development-method skills |
| [ECC](https://github.com/affaan-m/ecc) | `4adeee6732a12d1540f49a4d50a38ce8c97f9a53` | `ecc` | 328 skill directories |

Arena does not advertise a native plugin loader for these libraries. Copying
the files allows their instructions and resources to be consulted manually;
it does not activate slash commands, automatic skill invocation, hooks,
subagents, telemetry or the full ECC agent framework. Native Claude/Codex
integrations require their respective hosts, which were not present on PATH.
Impeccable's executable engine is a separate dependency; its presence or
functionality is not established by copying the skill. Some Taste and
Impeccable checks require a real browser; Node 22 is present. A temporary
headless Chromium 153 from npm was made to launch using the package's bundled
NSS/NSPR libraries. The full existing browser suite passed (312 cases),
including new reduced-motion navigation coverage. This temporary browser
uses `--no-sandbox` and is not a substitute for normal cross-browser,
physical-device, or security testing. Optional video checks also require
ffmpeg, which is absent.
Do not run remote installers or enable hooks without reviewing their effects.

Use Superpowers' planning, debugging and test-first practices; ECC's applicable
language, security and coding guidance; Impeccable's operate-mode guidance for
the Admin and Code interfaces; and Taste's product-specific, anti-template
principles where appropriate. These are references, not substitutes for FCC's
requirements, actual tests, or accessibility judgment. In particular, a
working admin tool does not need marketing-site animations, imagery or arbitrary
visual effects.

## Ordered work and evidence gates

1. **Integrate current main.** Fetch and inspect changes, resolve conflicts
   semantically, retain branch performance work and the intentional remote
   Admin/Code behavior, then run tests against the combined tree. As of this
   review, `origin/main` is `0d7e83e`, six commits past the prior integration
   base. A trial merge produced conflicts in three runtime files and twelve
   test files, including the remote-access tests and the Gemini SDK test bridge;
   it was aborted rather than leaving a conflicted or unverified tree. Main
   also moves CI from `ty` to strict Pyrefly, affecting hundreds of files.
   Neither old green results nor a textual merge establish correctness.
2. **Establish a reproducible baseline.** Rebuild the pinned toolchains and
   dependencies, run current CI formatting, lint, type and unit gates, Go
   tests and race checks, and remote/tunnel regression tests. Record environment
   and unavailable platforms or services. Compare packet-delivery performance
   against the same main revision with matching output hashes.
3. **Review the product surface.** Audit README/install provenance, deployment
   risks, configuration and error copy, then the Admin/Code workflows. Check
   keyboard and assistive technology semantics, focus, reduced motion, form
   errors, contrast, responsive layout, touch behavior and loading/failure
   states. Test desktop and narrow viewports in a real browser; do not infer
   visual or accessibility compliance solely from source. Avoid cosmetic
   redesigns unsupported by product needs or actual evidence.
4. **Review Python and Go subsystem by subsystem.** Start with public APIs,
   retries, stream delivery, request outcomes, tool calls, security boundaries,
   storage and provider integrations. Preserve the Python reference and test
   new Go behavior against freshly generated current-reference goldens. Repair
   stale generators individually; do not overwrite historical fixtures blindly.
   Measure resource contracts and benchmark representative workloads rather
   than extrapolating from isolated microbenchmarks.
5. **Release gate.** Complete required parity and user-facing features, validate
   supported operating systems and browsers, rerun all relevant checks, review
   residual risks, then merge only when justified. Inspect branches before
   removing any; preserve valuable work. The present review and partial Go
   port do not satisfy this gate.

Remote deployment through Cloudflare Tunnel, ngrok and Tailscale remains an
explicit requirement: do not reintroduce localhost-only checks or add a new
authentication requirement. Existing exposure risks must remain clearly
documented; do not misrepresent the proxy API token as Admin authentication.
