<div align="center">

<h1>
  <picture>
    <source media="(prefers-color-scheme: light)" srcset="assets/free-claude-code-wordmark-light.svg">
    <img src="assets/free-claude-code-wordmark-dark.svg" alt="Free Claude Code" width="610">
  </picture>
</h1>

[![License: AGPL v3](https://img.shields.io/badge/License-AGPL_v3-2e7d32.svg?style=for-the-badge)](LICENSE)
[![Python 3.14](https://img.shields.io/badge/python-3.14-3776ab.svg?style=for-the-badge&logo=python&logoColor=white)](https://www.python.org/downloads/)
[![Package Manager: uv](assets/package-manager-uv.svg)](https://github.com/astral-sh/uv)
[![Testing: Pytest](https://img.shields.io/badge/Testing-Pytest-ad1457.svg?style=for-the-badge)](https://github.com/pytest-dev/pytest)
[![Type Checker: Pyrefly](https://img.shields.io/badge/Type%20Checker-pyrefly-fdd835.svg?style=for-the-badge)](https://github.com/pyrefly-org/pyrefly)
[![Formatter: Ruff](https://img.shields.io/badge/Formatter-ruff-bf4b00.svg?style=for-the-badge)](https://github.com/astral-sh/ruff)
[![Logging: Loguru](https://img.shields.io/badge/logging-loguru-00695c.svg?style=for-the-badge)](https://github.com/Delgan/loguru)

[Documentation](#documentation) · [Quick Start](#quick-start) · [Providers](#choose-a-provider) · [Integrations](#integrations) · [API](#use-the-api-directly) · [Phone](#use-from-your-phone) · [Manage](#manage-your-installation)

</div>

<p align="center">
  <em>Independent open-source project. Not affiliated with or endorsed by Anthropic. Claude and Claude Code are trademarks of Anthropic.</em>
</p>

Free Claude Code (FCC) is an open-source gateway that connects coding agents,
editors and API clients to cloud or local model providers. Configure routing in
the Admin UI, then use a supported launcher or one of the compatible APIs.

This repository, [MoLik-Labs/FCC](https://github.com/MoLik-Labs/FCC), is a fork
of [Alishahryar1/free-claude-code](https://github.com/Alishahryar1/free-claude-code).
The Python application is the supported runtime. The Go work in `gocore/` is
experimental and is not yet a replacement for the Python server.

## Documentation

- [English guide](docs/site/src/content/docs/index.md) · [Hướng dẫn tiếng Việt](docs/site/src/content/docs/vi/index.md)
  — source pages for the bilingual installation, API and operations guide.
- [Build and publish the documentation](docs/site/README.md) — local preview and
  GitHub Pages setup. The public URL is available from the deployment workflow
  after Pages has been enabled and a deployment succeeds.
- [Development and contribution](CONTRIBUTING.md) · [Go rewrite status](gocore/README.md)
  · [Rewrite plan and checkpoints](docs/rewrite/00-master-plan.md)
  · [Fork safety and verification](docs/rewrite/fork-safety-checkpoint.md).
- [Skill workflow](docs/design/skill-workflow.md) — pinned sources and
  restoration instructions for future development sessions.

The site uses Astro Starlight, Nimbus's accessible table helper and static
Markdown for agent readers, with a VitePress API reference. To preview it from a checkout with Node.js
22.23.4 or newer:

```sh
npm ci --prefix docs/site
npm run dev --prefix docs/site
```

See the build guide for production validation and GitHub Pages deployment. The
site's light/dark preference is separate from the Admin's Forest, Catppuccin
Mocha and other appearance choices.

## Oh My Molik (experimental)

[Oh My Molik](ohmymolik/README.md) is a separate experimental agent project. Its
Rust and historical Go prototypes can inspect files read-only via an explicit
local or remote Chat Completions endpoint. The Rust CLI uses system curl and
is non-streaming. It also has an experimental read-only MCP stdio server;
neither prototype is a complete coding agent. See the [research and roadmap](docs/oh-my-fox/research.md)
for sources, migration status and safety gates.

## Features

- Route any supported harness to models from any supported provider.
- 12 supported harnesses: [Claude Code](https://code.claude.com/docs/en/overview), [Codex](https://github.com/openai/codex), [Pi](https://github.com/earendil-works/pi), [OpenCode](https://github.com/anomalyco/opencode), [Cline](https://github.com/cline/cline), [Hermes](https://github.com/NousResearch/hermes-agent), [DeepSeek Harness](https://github.com/deepseek-ai/deepseek-harness), [Grok Build](https://github.com/xai-org/grok-build), [Muse Code](https://research.meta.ai/blog/introducing-muse-code-and-muse-spark-1-2/), [Aider](https://aider.chat/), [Antigravity CLI](https://antigravity.google/docs/getting-started/), or [VS Code Chat](https://code.visualstudio.com/).
  - Supports streaming, tool use, and interleaved thinking where supported by the client and provider.
  - For Claude Code, route [Fable](https://www.anthropic.com/claude/fable), [Opus](https://www.anthropic.com/claude/opus), [Sonnet](https://www.anthropic.com/claude/sonnet), and [Haiku](https://www.anthropic.com/claude/haiku) to supported models.
- Use any model through any of four APIs: [Anthropic Messages](https://platform.claude.com/docs/en/api/messages) (`/v1/messages`), [OpenAI Chat Completions](https://platform.openai.com/docs/api-reference/chat/create) (`/v1/chat/completions`), [OpenAI Responses](https://platform.openai.com/docs/api-reference/responses/create) (`/v1/responses`), or [Gemini](https://ai.google.dev/api/generate-content) (`generateContent`). See [Use The API Directly](#use-the-api-directly).
- Supports free-tier, paid, subscription, and local providers. Model access,
  quotas, pricing, and permitted use depend on each provider and can change.
- Use through the terminal, desktop apps ([Codex App](https://learn.chatgpt.com/docs/app), [Claude Desktop](https://claude.ai/download), and [DeepSeek Desktop](https://github.com/deepseek-ai/deepseek-harness)), IDE extensions ([VS Code](https://code.visualstudio.com/) and [JetBrains](https://www.jetbrains.com/)), [your phone](#use-from-your-phone) ([Discord](https://discord.com/) and [Telegram](https://telegram.org/)), or your browser through native Code sessions.
  - Voice notes supported for Discord and Telegram using local [Whisper](https://github.com/openai/whisper) or [NVIDIA NIM](https://docs.nvidia.com/nim/speech/latest/asr/deploy-asr-models/whisper.html).
- Optional token optimization using [RTK](https://github.com/rtk-ai/rtk) and [Headroom](https://github.com/headroomlabs-ai/headroom); savings depend on the workload.
- Reliability features for provider failures:
  - Automatic retries.
  - Configurable rate limits.
  - Configurable concurrency control.
  - Automatic model/provider fallback.
  - Automatic multi-API key fallback.

Free-tier availability and limits are controlled by each provider and may change.

> **Security:** Admin and Code session endpoints are reachable by remote clients
> when the server is exposed; they do not have a separate Admin authentication
> boundary. They can change configuration and operate code sessions. Do not
> expose them to untrusted networks or the public internet without protection
> provided by your network or tunnel. Proxy Authentication protects the proxy
> API; do not assume it also protects Admin. Remote HTTP sends credentials
> and request contents without HTTP transport encryption; use HTTPS or a
> properly protected private tunnel, including for Claude Desktop gateways.
>
> The Go port in [`gocore/`](gocore/README.md) is experimental and incomplete;
> the Python application remains the supported runtime.

## Quick Start

<a id="install"></a>

### 1. Install this fork from source

Install [Git](https://git-scm.com/downloads) and [uv](https://docs.astral.sh/uv/getting-started/installation/) first. In a terminal (or Windows PowerShell), run:

```bash
git clone --branch arena/5fe6a43f-fcc https://github.com/MoLik-Labs/FCC.git
cd FCC
uv sync --locked
```

This checks out the current development branch of **this fork**, not an upstream release. The branch is still under development; after these changes are integrated into `main`, use the desired available branch instead. Python 3.14.10 is required; uv can provision it. This source setup does not install coding agents, create desktop shortcuts or provide the upstream installer's optional setup. Install the client you intend to use separately.

### 2. Start FCC

From the `FCC` checkout on Windows, macOS or Linux:

```bash
uv run fcc-server
```

Keep the terminal open and open the Admin URL shown in the server log. In a
second terminal, change into the same checkout before running a launcher with
`uv run`.

<a id="nvidia-nim-provider"></a>

### 3. Configure NVIDIA NIM

1. Create an API key at [build.nvidia.com/settings/api-keys](https://build.nvidia.com/settings/api-keys).
2. Open the Admin UI URL from the server log.
3. Open NVIDIA NIM in **Providers**, click **Add key**, paste the key, and click **Save**.
4. Leave `MODEL` on the default `nvidia_nim/nvidia/nemotron-3-ultra-550b-a55b`, or search the model dropdown and select another model.
5. Click **Apply** to save your model selection.

To protect the local proxy with a bearer token, enable **Proxy Authentication**
in Admin. This token does not protect Admin or Code sessions; see the [remote
access guidance](docs/oh-my-fox/remote-access.md) before exposing them.

Use **Appearance** in the sidebar to choose Forest, Catppuccin Mocha, Nord, Dracula, Ocean,
Dusk or Dawn. The choice is saved in this browser when storage is available;
it does not change server configuration or the separate Molik CLI.

### 4. Run Your Coding Agent

Install your chosen agent separately, then run one of these commands from the checkout:

Claude Code:

```bash
uv run fcc-claude
```

Codex:

```bash
uv run fcc-codex
```

Pi:

```bash
uv run fcc-pi
```

OpenCode 2:

```bash
uv run fcc-opencode
```

Cline:

```bash
uv run fcc-cline
```

Hermes:

```bash
uv run fcc-hermes
```

DeepSeek Harness Web:

```bash
uv run fcc-dsh
```

Grok Build:

```bash
uv run fcc-grok
```

Muse Code:

```bash
uv run fcc-muse
```

Aider:

```bash
uv run fcc-aider
```

Antigravity CLI:

```bash
uv run fcc-agy
```

Other clients and SDKs can call FCC directly. See [Use The API Directly](#use-the-api-directly).

<a id="model-picker"></a>

## Workload and hardware guidance

In **Admin → Providers → Runtime configuration**, the workload buttons stage
only **Provider Max Concurrency**: Light = 1, Balanced = 2 (the default), High
= 4. Review the field and click **Apply** to save. You can enter a different
positive value manually. If a process environment variable locks the field,
change that variable instead; the buttons are disabled.

The limit is per provider, not a global CPU/RAM limit. Provider Rate Limit and
Rate Window remain separate and can still delay attempts. These are starting
points, not automatic hardware detection or speed guarantees:

- **Low-resource machine / local model:** start at 1; use a smaller model and
  context, and avoid optional local Whisper if memory or disk is tight.
- **Typical single-user setup:** start at 2; observe memory, latency and errors.
- **More capacity or multiple users:** try 4 only when the provider quota and
  local inference server can sustain it; higher concurrency can increase memory
  use and latency. Lower it if requests fail or the machine becomes unresponsive.

Cloud inference runs on the provider's hardware; local inference requirements
are determined primarily by the model, context, quantization and inference
server. These controls do not change model settings or limit the entire machine.

## Choose A Provider

1. Open a provider link below for its key, models, or setup instructions.
2. In the Admin UI, configure the listed setting. For OpenAI / ChatGPT
   subscription access, use **Providers → OAuth providers** instead.
3. Search the `MODEL` dropdown and select a model. If the provider cannot list
   models, enter `<provider-id>/<exact-provider-model-id>` manually.
4. Click **Apply**.

Optional: add an ordered **Fallback Models** list under **Model Config**. It
applies to every connected client. A failed request may reach and consume usage
from more than one provider before succeeding.

Provider modals support numbered API keys. Use **Add key**, remove rows, or drag their handles to reorder them. FCC uses key 1.

For another API endpoint, open **Providers → Custom providers → Add provider**. Enter its API base URL, optional keys, and API format. Leave **Model IDs** empty to discover models automatically, or enter one upstream model ID per line. Save, then select a model in **Model Config**.

<details>
<summary><strong>Provider catalog</strong></summary>

| Provider | Admin UI setting | Example `MODEL` |
| --- | --- | --- |
| [NVIDIA NIM](https://build.nvidia.com/settings/api-keys) | `NVIDIA_NIM_API_KEY` | `nvidia_nim/nvidia/nemotron-3-ultra-550b-a55b` |
| [OpenRouter](https://openrouter.ai/keys) | `OPENROUTER_API_KEY` | `open_router/openrouter/free` |
| [Groq](https://console.groq.com/keys) | `GROQ_API_KEY` | `groq/llama-3.3-70b-versatile` |
| [ClinePass](https://docs.cline.bot/getting-started/clinepass) | `CLINE_API_KEY` | `cline_pass/cline-pass/kimi-k3` |
| [OpenAI / ChatGPT](https://learn.chatgpt.com/docs/auth) | Connect ChatGPT in the Admin UI | `openai/<model-id>` |
| [OpenAI API](https://platform.openai.com/api-keys) | `OPENAI_API_KEY` | `openai_api/gpt-5.6-sol` |
| [GitHub Copilot](https://docs.github.com/en/copilot/how-tos/copilot-sdk/auth/authenticate) | Connect GitHub Copilot in the Admin UI | `github_copilot/<model-id>` |
| [xAI (Grok)](https://console.x.ai/team/default/api-keys) | `XAI_API_KEY` | `xai/grok-4.5` |
| [Alibaba Cloud](https://www.alibabacloud.com/help/en/model-studio/get-api-key) | `ALIBABA_CLOUD_API_KEY` | `alibaba_cloud/<model-id>` |
| [QwenCloud Token Plan](https://home.qwencloud.com/api-keys) | `QWENCLOUD_API_KEY` | `qwencloud/qwen3.7-plus` |
| [QwenCloud Coding Plan](https://home.qwencloud.com/api-keys) | `QWENCLOUD_CODING_API_KEY` | `qwencloud_coding/qwen3.7-plus` |
| [Together AI](https://api.together.ai/settings/api-keys) | `TOGETHER_API_KEY` | `together/zai-org/GLM-5.2` |
| [DeepInfra](https://deepinfra.com/dash/api_keys) | `DEEPINFRA_API_KEY` | `deepinfra/deepseek-ai/DeepSeek-V4-Flash` |
| [SiliconFlow](https://cloud.siliconflow.com/account/ak) | `SILICONFLOW_API_KEY` | `siliconflow/Qwen/Qwen3-32B` |
| [Nebius Token Factory](https://tokenfactory.nebius.com/project/api-keys) | `NEBIUS_API_KEY` | `nebius/Qwen/Qwen3-30B-A3B` |
| [Chutes](https://chutes.ai/docs/getting-started/authentication) | `CHUTES_API_KEY` | `chutes/Qwen/Qwen3-32B-TEE` |
| [Featherless AI](https://featherless.ai/account/api-keys) | `FEATHERLESS_API_KEY` | `featherless/Qwen/Qwen3-32B` |
| [Agnes AI](https://agnes-ai.com/) | `AGNES_API_KEY` | `agnes/agnes-2.0-flash` |
| [ZenMux](https://zenmux.ai/platform/pay-as-you-go) | `ZENMUX_API_KEY` | `zenmux/deepseek/deepseek-v4-flash-free` |
| [W&B Inference](https://wandb.ai/settings) | `WANDB_API_KEY` | `wandb/openai/gpt-oss-20b` |
| [Anthropic](https://platform.claude.com/settings/keys) | `ANTHROPIC_API_KEY` | `anthropic/<model-id>` |
| [Azure OpenAI](https://learn.microsoft.com/azure/foundry/openai/how-to/chatgpt) | `AZURE_OPENAI_API_KEY` and `AZURE_OPENAI_BASE_URL` | `azure_openai/<deployment-name>` |
| [Google AI Studio (Gemini)](https://aistudio.google.com/apikey) | `GEMINI_API_KEY` | `gemini/models/gemini-3.1-flash-lite` |
| [Google Vertex AI](https://cloud.google.com/vertex-ai/generative-ai/docs/start/openai) | `VERTEX_PROJECT_ID` + ADC | `vertex/google/gemini-3.5-flash` |
| [DeepSeek](https://platform.deepseek.com/api_keys) | `DEEPSEEK_API_KEY` | `deepseek/deepseek-chat` |
| [Mistral La Plateforme](https://console.mistral.ai/) | `MISTRAL_API_KEY` | `mistral/devstral-small-latest` |
| [Mistral Codestral](https://console.mistral.ai/) | `CODESTRAL_API_KEY` | `mistral_codestral/codestral-latest` |
| [OpenCode Zen](https://opencode.ai/auth) | `OPENCODE_API_KEY` | `opencode_zen/gpt-5.3-codex` |
| [OpenCode Go](https://opencode.ai/auth) | `OPENCODE_API_KEY` | `opencode_go/minimax-m2.7` |
| [Vercel AI Gateway](https://vercel.com/docs/ai-gateway/models-and-providers) | `AI_GATEWAY_API_KEY` | `vercel/openai/gpt-5.5` |
| [Amazon Bedrock](https://console.aws.amazon.com/bedrock/) | `AWS_BEARER_TOKEN_BEDROCK` | `bedrock/openai.gpt-oss-120b` |
| [Hugging Face Inference Providers](https://huggingface.co/settings/tokens) | `HUGGINGFACE_API_KEY` | `huggingface/Qwen/Qwen3-Coder-480B-A35B-Instruct:fastest` |
| [Cohere](https://dashboard.cohere.com/api-keys) | `COHERE_API_KEY` | `cohere/command-a-plus-05-2026` |
| [Wafer](https://wafer.ai/) | `WAFER_API_KEY` | `wafer/DeepSeek-V4-Pro` |
| [Kimi API](https://platform.moonshot.ai/console/api-keys) | `KIMI_API_KEY` | `kimi/kimi-k2.5` |
| [Kimi Code](https://www.kimi.com/code/console) | `KIMI_CODE_API_KEY` | `kimi_code/k3` |
| [MiniMax](https://platform.minimax.io/user-center/basic-information/interface-key) | `MINIMAX_API_KEY` | `minimax/MiniMax-M3` |
| [Cerebras Inference](https://cloud.cerebras.ai/) | `CEREBRAS_API_KEY` | `cerebras/gpt-oss-120b` |
| [SambaNova](https://cloud.sambanova.ai/apis) | `SAMBANOVA_API_KEY` | `sambanova/Meta-Llama-3.3-70B-Instruct` |
| [Kilo.ai](https://kilo.ai) | `KILO_API_KEY` | `kilo/kilo-auto/free` |
| [Fireworks AI](https://fireworks.ai/account/api-keys) | `FIREWORKS_API_KEY` | `fireworks/accounts/fireworks/models/llama-v3p3-70b-instruct` |
| [Novita AI](https://novita.ai/settings/key-management) | `NOVITA_API_KEY` | `novita/deepseek/deepseek-v4-flash-0731` |
| [Cloudflare Workers AI](https://developers.cloudflare.com/workers-ai/) | `CLOUDFLARE_API_TOKEN` and `CLOUDFLARE_ACCOUNT_ID` | `cloudflare/@cf/moonshotai/kimi-k2.6` |
| [Z.ai Coding Plan](https://z.ai/manage-apikey/apikey-list) | `ZAI_API_KEY` | `zai/glm-5.2` |
| [Z.ai API (pay as you go)](https://z.ai/manage-apikey/apikey-list) | `ZAI_API_KEY` | `zai_api/glm-4.7-flash` |
| [TokenRouter](https://www.tokenrouter.com/) | `TOKENROUTER_API_KEY` | `tokenrouter/moonshotai/kimi-k3-free` |
| [NaraRoute](https://router.bynara.id/) | `NARAROUTE_API_KEY` | `nararoute/kimi-k3-free` |
| [Poolside AI](https://platform.poolside.ai/) | `POOLSIDE_API_KEY` | `poolside/poolside/laguna-s-2.1` |
| [LLM7.io](https://dash.llm7.io/) | `LLM7_API_KEY` | `llm7/default` |
| [Scaleway](https://console.scaleway.com/iam/api-keys) | `SCW_SECRET_KEY` | `scaleway/deepseek/deepseek-v4-flash` |
| [Lightning AI](https://lightning.ai/) | `LIGHTNING_API_KEY` | `lightning/lightning-ai/Qwen3.8-27B` |
| [Experiential Labs](https://platform.experientiallabs.ai/) | `EXPLABS_API_KEY` | `experiential/union-alpha` |
| [Cheaper Inference](https://cheaperinference.com/signup) | `CHEAPER_INFERENCE_API_KEY` | `cheaperinference/gpt-5.4-mini` |
| [OrcaRouter](https://www.orcarouter.ai/) | `ORCAROUTER_API_KEY` | `orcarouter/deepseek/deepseek-v4-flash-free` |
| [Atlas Cloud](https://console.atlascloud.ai/api-keys) | `ATLASCLOUD_API_KEY` | `atlascloud/deepseek-ai/deepseek-v4-flash` |
| [xKiro](https://xkiro.com/) | `XKIRO_API_KEY` | `xkiro/qwen/qwen3.7-flash:free` |
| [Opper](https://platform.opper.ai/) | `OPPER_API_KEY` | `opper/claude-sonnet-4-6` |
| [Ollama Cloud](https://ollama.com/settings/keys) | `OLLAMA_API_KEY` | `ollama_cloud/qwen3-coder:480b` |
| [LM Studio](https://lmstudio.ai/) | `LM_STUDIO_BASE_URL` | `lmstudio/<model-id>` |
| [llama.cpp](https://github.com/ggml-org/llama.cpp) | `LLAMACPP_BASE_URL` | `llamacpp/<model-id>` |
| [Ollama](https://ollama.com/) | `OLLAMA_BASE_URL` | `ollama/<model-tag>` |

</details>

<details>
<summary><strong>Provider-specific setup</strong></summary>

- OpenAI / ChatGPT uses your ChatGPT subscription rather than an API key. Connect from
  **Providers → OAuth providers → OpenAI / ChatGPT → Connect** in the Admin UI
  and finish signing in through your browser. Restart an already-running agent after connecting.
- Anthropic uses a separately billed API key. Add it under
  **Providers → Cloud providers → Anthropic → Configure**, then choose a discovered
  model. Multi-workspace keys also require `ANTHROPIC_WORKSPACE_ID`.
- OpenAI API uses a separate Platform API key. Enter it under
  **Providers → Cloud providers → OpenAI API → Configure**. The model list may
  include IDs that cannot handle coding requests. Choose a text-generation model.
- GitHub Copilot uses your signed-in GitHub account and subscription. Install
  [Copilot CLI 1.0.83](https://github.com/github/copilot-cli/releases/tag/v1.0.83)
  on PATH, then choose **Providers → OAuth providers → GitHub Copilot → Connect**.
  Finish signing in, then select a `github_copilot/<model-id>` from the model picker.
  Available models and quotas depend on your subscription and organization policies.
  Restart an already-running agent after connecting.
- Azure OpenAI uses the deployment names from your resource. Set
  `AZURE_OPENAI_BASE_URL` to its complete v1 endpoint, such as
  `https://YOUR-RESOURCE-NAME.openai.azure.com/openai/v1/`, and select a
  deployment that supports Chat Completions. Enter the deployment name as a
  custom model slug if it does not appear in the model dropdown.
- Mistral Codestral uses a separate key from Mistral La Plateforme.
- Kimi Code subscription keys use `kimi_code/`. Kimi API credit keys use
  `kimi/`. Kimi Code plans are for personal interactive coding-agent use under
  [Kimi's community guidelines](https://www.kimi.com/code/docs/en/kimi-code/community-guidelines.html).
- Alibaba Cloud uses a Model Studio pay-as-you-go key and defaults to Singapore.
  For another region or workspace, set `ALIBABA_CLOUD_BASE_URL` to its
  [OpenAI-compatible endpoint](https://www.alibabacloud.com/help/en/model-studio/base-url)
  ending in `/compatible-mode/v1`. The key and endpoint must use the same region.
- QwenCloud Coding Plan keys use `qwencloud_coding/`. QwenCloud Token Plan keys
  use `qwencloud/`. The keys and endpoints are not interchangeable. Coding Plan
  is for local, personal, interactive coding-agent use under the
  [Coding Plan terms](https://www.alibabacloud.com/help/en/model-studio/coding-plan).
- OpenCode Zen and OpenCode Go share `OPENCODE_API_KEY` but use the explicit
  `opencode_zen/` and `opencode_go/` model prefixes.
- For Amazon Bedrock, set `BEDROCK_BASE_URL` to the URL for the same region as
  the API key and select one of the listed models.
- Vertex AI uses Google Application Default Credentials instead of an API key.
  Locally, run `gcloud auth application-default login` once. Service-account
  files and attached service accounts also work. Set `VERTEX_PROJECT_ID`, and
  optionally change `VERTEX_LOCATION` from its `global` default.
- Cloudflare requires both its API token and account ID.
- For Ollama Cloud, use the exact model IDs shown in the model picker. Local
  Ollama uses the separate `ollama/` prefix.
- Prefer tool-capable models for coding agents. Local models also need enough context for the agent's system prompt and tool definitions.

</details>

<details>
<summary><strong>Local provider setup</strong></summary>

FCC automatically discovers models from running local servers at startup. If you start a server later, click **Refresh models** in Admin to add its models to the picker.

### LM Studio

Start LM Studio's local server, load a tool-capable model, and use the model identifier shown by LM Studio with the `lmstudio/` prefix. The default URL is `http://localhost:1234/v1`.

### llama.cpp

Start `llama-server` with its OpenAI-compatible Chat Completions API and enough context for the model. Use the local model ID with the `llamacpp/` prefix. `LLAMACPP_BASE_URL` defaults to `http://localhost:8080/v1`.

### Ollama

```bash
ollama pull llama3.1
ollama serve
```

Use the tag shown by `ollama list` with the `ollama/` prefix. `OLLAMA_BASE_URL` defaults to `http://localhost:11434`.

Provider cards show the checked endpoint and response time, with a useful message when a local server is unavailable. After starting the server or correcting its settings, click **Test** to refresh its models.

</details>

<details>
<summary><strong>Optional model-tier routing</strong></summary>

`MODEL` is the fallback for every request. Select a model for `MODEL_FABLE`, `MODEL_OPUS`, `MODEL_SONNET`, or `MODEL_HAIKU` to override an individual Claude Code tier. Select **None** to use `MODEL`.

In **Admin UI → Model Config**, add **Fallback Models** in the order FCC should try them if the selected model fails. Each fallback has its own reasoning effort: **Inherit**, **Off**, **From client**, **Low**, **Medium**, **High**, **X-High**, or **Max**.

Choose a model, effort, and API key number in the bottom row, then click **+**. **Key 1**, **Key 2**, and so on refer to the current numbered keys in that provider's settings. Edit the selections in place, drag the handle to change the order, or click the red trash button to remove a fallback.

**Inherit** keeps the original request's resolved reasoning setting. **From client** uses the effort or token budget sent by your coding agent, bypassing the global or tier setting. You can select the same model with another key number.

</details>

<details>
<summary><strong>Reasoning control</strong></summary>

Open **Admin UI → Model Config → Model Routing** and choose an effort beside each model. **Default Effort** starts at **From client**. Fable, Opus, Sonnet, and Haiku start at **Inherit**.

| Selection | Behavior |
| --- | --- |
| **From client** | Use the effort sent by your coding agent. If none is sent, keep the provider default. |
| **Off** | Request reasoning to be disabled. |
| **Low**, **Medium**, **High**, **X-High**, or **Max** | Override the client with the selected reasoning level. |
| **Inherit** | Tier overrides use Default Effort. Fallbacks keep the original request's resolved effort. |

</details>

<a id="integrations"></a>
<a id="connect-your-client"></a>

## Integrations

FCC launchers automatically start the server.

For editor and app integrations, install the client, start FCC, then open
**Admin UI → Integrations** and click **Connect** on its card.

- **Claude Code in VS Code**: Install the [Claude Code extension](https://marketplace.visualstudio.com/items?itemName=anthropic.claude-code).
- **VS Code Chat**: Install the latest [VS Code](https://code.visualstudio.com/), connect, and select an FCC model in its built-in Chat model picker.
- **Claude Desktop**: Install [Claude Desktop](https://claude.ai/download). Fully quit it before connecting or disconnecting, then reopen it. Disconnect returns to normal Claude sign-in.
- **Codex in VS Code and App**: Install the [Codex extension](https://marketplace.visualstudio.com/items?itemName=openai.chatgpt) or Codex App.
- **DeepSeek Harness Desktop**: Install and open the app once before connecting. Start a new session to use FCC models.
- **Claude Code in JetBrains ACP**: Install Claude Agent in JetBrains AI Assistant and start it once, then click **Connect** in FCC. Reopen the IDE, select **Claude Code (FCC)**, and start a new chat. After JetBrains updates the agent, restart FCC before starting a new chat. Requires a local IDE in its standard installation locations.

Reload VS Code or restart the app/IDE after connecting. In Codex and Claude Desktop, select an FCC
model from the model picker. Reload or restart the client when FCC reports updated settings. Use
**Disconnect** on the same card to remove the integration.

To complete Claude Code onboarding separately, open **Admin UI → Integrations →
Claude Code onboarding** and click **Complete onboarding**. Connecting Claude Code
in VS Code also completes onboarding.

Run FCC on the same computer and in the same user environment as the client you
are configuring.

<a id="use-the-api-directly"></a>

## Use The API Directly

Any client or SDK for these APIs can call FCC. Set the base URL to
`http://127.0.0.1:8082/v1` for OpenAI clients and `http://127.0.0.1:8082` for
Anthropic and Gemini clients.
Every FCC model works on every endpoint. If Proxy
Authentication is enabled, use its token as the client's API key, or add
`-H "Authorization: Bearer <token>"` to each `curl` request.
List the available models with `curl "http://127.0.0.1:8082/v1/models?view=messages"`.

[`POST /v1/messages`](https://platform.claude.com/docs/en/api/messages):

```bash
curl http://127.0.0.1:8082/v1/messages \
  -H "Content-Type: application/json" \
  -d '{
    "model": "nvidia_nim/nvidia/nemotron-3-ultra-550b-a55b",
    "max_tokens": 1024,
    "messages": [{"role": "user", "content": "Hello!"}]
  }'
```

[`POST /v1/chat/completions`](https://platform.openai.com/docs/api-reference/chat/create):

```bash
curl http://127.0.0.1:8082/v1/chat/completions \
  -H "Content-Type: application/json" \
  -d '{
    "model": "nvidia_nim/nvidia/nemotron-3-ultra-550b-a55b",
    "messages": [{"role": "user", "content": "Hello!"}]
  }'
```

[`POST /v1/responses`](https://platform.openai.com/docs/api-reference/responses/create):

```bash
curl http://127.0.0.1:8082/v1/responses \
  -H "Content-Type: application/json" \
  -d '{
    "model": "nvidia_nim/nvidia/nemotron-3-ultra-550b-a55b",
    "input": "Hello!"
  }'
```

[`POST /v1/models/{model}:generateContent`](https://ai.google.dev/api/generate-content):

```bash
curl http://127.0.0.1:8082/v1/models/nvidia_nim/nvidia/nemotron-3-ultra-550b-a55b:generateContent \
  -H "Content-Type: application/json" \
  -d '{"contents": [{"role": "user", "parts": [{"text": "Hello!"}]}]}'
```

Responses replies always stream. To stream Messages or Chat Completions replies, add
`"stream": true`. For Gemini, call `:streamGenerateContent?alt=sse` instead of
`:generateContent`.

<a id="use-from-your-phone"></a>
<a id="optional-integrations"></a>

## Use From Your Phone

Run coding tasks from Discord or Telegram, and send voice notes. Configure them
in **Admin UI → Messaging**, then click **Apply**.

<details>
<summary><strong>Discord bot</strong></summary>

1. Create a bot in the [Discord Developer Portal](https://discord.com/developers/applications).
2. Enable **Message Content Intent** and invite it with read, send,
   message-history, and **Manage Messages** permissions so `/clear` can remove
   user prompts.
3. Set **Messaging Platform** to **discord**.
4. Enter **Discord Bot Token**, **Allowed Discord Channels**, and an absolute **Allowed Directory**.
5. Apply the settings and restart the server if requested.

</details>

<details>
<summary><strong>Telegram bot</strong></summary>

1. Create a bot with [@BotFather](https://t.me/BotFather).
2. Get your numeric user ID from [@userinfobot](https://t.me/userinfobot).
   In groups, grant the bot permission to delete messages.
3. Set **Messaging Platform** to **telegram**.
4. Enter **Telegram Bot Token**, **Allowed Telegram User ID**, and an absolute **Allowed Directory**.
5. Apply the settings and restart the server if requested.

</details>

### Messaging commands

| Usage | Behavior |
| --- | --- |
| `/stats` | Show session state. |
| Standalone `/stop` | Cancel all work. |
| Reply with `/stop` | Cancel only the selected request while other queued requests continue. |
| Standalone `/clear` | Clear all FCC conversations and their messages in this chat. |
| Reply with `/clear` | Delete the selected message and the chain of replies to it. Keep other conversations. |

<details>
<summary><strong>Voice notes</strong></summary>

For NVIDIA NIM transcription, configure its API key under **Providers**, then
in **Messaging → Voice** enable voice notes and choose `nvidia_nim`.

Local Whisper is optional and needs substantially more disk space and memory.
From this source checkout, install the `voice_local` extra:

```bash
uv sync --locked --extra voice_local
```

The configured PyTorch index targets CUDA 13.0; availability and CPU/GPU
compatibility depend on platform, driver and package wheels. Verify those
requirements before installing. Start with `uv run --extra voice_local fcc-server`,
enable voice in **Messaging → Voice**, select `cpu` or `cuda` and
choose a Whisper model. Gated models need `HUGGINGFACE_API_KEY`. Repeat the
extra on `uv sync` when updating; a plain sync can remove optional packages.
This extra has not been validated on every platform.

</details>

## Manage Your Installation

For the source checkout, run `uv run fcc-server --version` or `uv run fcc-doctor` from the `FCC` directory. To update the same checked-out branch, stop FCC and run:

```bash
git pull --ff-only
uv sync --locked
```

Review local changes and confirm the branch and remote before updating. If using
local voice, replace the sync command with `uv sync --locked --extra voice_local`.
In this revision, `fcc-update` is deliberately disabled: it exits without
changing files and prints the manual steps above. **Older installed copies** may
still run the upstream installer; do not use them for this fork. The inherited
`scripts/install.sh` and `scripts/install.ps1` also target upstream releases,
not this source checkout. Automatic fork release publication is disabled until
the package identity and release process have been reviewed.

To remove the source checkout, stop FCC and connected clients, back up any
needed data, then remove the checkout directory manually. Its `.venv` is inside
the checkout. FCC data such as `~/.fcc/` is separate and is not removed; delete
it only if you intentionally want to remove configuration, credentials and data.
The source install does not install or remove coding agents, RTK or Headroom.

## Project Links

For a problem with this checkout, run the command below from the repository and redact credentials and private data before sharing the output.

```bash
uv run fcc-doctor
```

- [This repository and its pull requests](https://github.com/MoLik-Labs/FCC) — issues are currently disabled on this fork.
- [Upstream issues](https://github.com/Alishahryar1/free-claude-code/issues) — for the upstream project and its releases; follow its reporting rules.
- [Contributing guide](CONTRIBUTING.md)

## Maintainers

See [MAINTAINERS.md](MAINTAINERS.md) for the project's maintainers and their areas of responsibility.

## License

GNU Affero General Public License version 3 only (`AGPL-3.0-only`). See [LICENSE](LICENSE) for details.
