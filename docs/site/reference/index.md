---
title: API overview
---

# FCC API overview

A concise reference to the Python gateway's compatible API surfaces. FCC routes a request to the configured provider; supported model features, availability, pricing, and provider terms vary. The Go port and Oh My Molik remain experimental.

<a href="../#quickstart">← Installation and operational guide</a> · <a href="../#security">Security and remote access</a>

## Before you call

1. Start FCC and configure a provider and model in Admin. The default local base URL is `http://127.0.0.1:8082`; use the actual address in your server log if it differs.
2. List configured models with `GET /v1/models?view=messages`. Use an available model ID in requests.
3. If Proxy Authentication is enabled, send `Authorization: Bearer <token>` or configure the client to use that token. Keep credentials out of public logs and examples.

::: warning Remote access
Proxy Authentication protects the proxy API; it is not separate authentication for Admin or Code. Restrict remote access with independent network or tunnel controls. Never expose Admin or Code to untrusted networks.
:::

## API endpoints

| Method and path | Purpose | Streaming |
| --- | --- | --- |
| `GET /health` | Basic server health | No |
| `GET /v1/models` | List advertised models | No |
| `POST /v1/messages` | Anthropic-compatible messages | Optional `stream: true` |
| `POST /v1/messages/count_tokens` | Count tokens for a messages request | No |
| `POST /v1/chat/completions` | OpenAI-compatible chat | Optional `stream: true` |
| `POST /v1/responses` | OpenAI-compatible Responses | Streaming only |
| `POST /v1/models/{model}:generateContent` | Gemini-compatible generation | No |
| `POST /v1/models/{model}:streamGenerateContent?alt=sse` | Gemini-compatible streaming | SSE |

FCC also accepts the Gemini paths under `/v1beta/models/`. These compatibility surfaces do not imply every upstream field, tool, provider, or streaming behavior is supported. Consult the [repository API examples](https://github.com/MoLik-Labs/FCC/blob/main/README.md#use-the-api-directly) and source for current details.

## First request

Replace the example model with an ID available from your configured provider:

```bash
curl http://127.0.0.1:8082/v1/messages \
  -H "Content-Type: application/json" \
  -d '{
    "model": "nvidia_nim/nvidia/nemotron-3-ultra-550b-a55b",
    "max_tokens": 1024,
    "messages": [{"role": "user", "content": "Hello!"}]
  }'
```

The example assumes proxy authentication is disabled. If enabled, add the bearer header described above. A working provider, valid credentials, an available model, and any required account permissions are necessary.

## Errors and operations

Check the server log, model ID, provider credentials and account limits first; retry only where appropriate. A local health response does not verify provider connectivity. Review the <a href="../#troubleshooting">troubleshooting guide</a> and <a href="../#maintenance">maintenance notes</a>. This reference is static and does not send requests to your server.
