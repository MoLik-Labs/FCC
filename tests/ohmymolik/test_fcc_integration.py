"""Exercise a built Molik binary through FCC with a mocked upstream provider."""

import json
import os
import socket
import subprocess
import threading
from pathlib import Path
from typing import cast, override

import httpx2
import pytest
import uvicorn

from free_claude_code.config.settings import Settings
from tests.api.support import create_test_app
from tests.providers.test_custom_provider import definition, provider
from tests.wire_payload import decode_object, object_array, string_payload


def test_molik_read_tool_roundtrip_via_fcc(tmp_path: Path) -> None:
    binary = os.environ.get("MOLIK_BINARY")
    if not binary:
        pytest.skip("Set MOLIK_BINARY to a built Oh My Molik executable")
    binary_path = Path(binary).resolve()
    assert binary_path.is_file()
    (tmp_path / "demo.txt").write_text("FCC fixture content", encoding="utf-8")
    seen: list[str] = []

    def reply(request: httpx2.Request) -> httpx2.Response:
        payload = decode_object(request.content)
        tool_result = next(
            (m for m in object_array(payload["messages"]) if m["role"] == "tool"),
            None,
        )
        if tool_result is None:
            delta = {
                "role": "assistant",
                "tool_calls": [
                    {
                        "index": 0,
                        "id": "call_a",
                        "type": "function",
                        "function": {
                            "name": "read_file",
                            "arguments": '{"path":"demo.txt"}',
                        },
                    }
                ],
            }
            reason = "tool_calls"
        else:
            seen.append(string_payload(tool_result["content"]))
            delta = {"role": "assistant", "content": "FCC mock tool loop succeeded"}
            reason = "stop"
        event = {
            "id": "cmpl-test",
            "object": "chat.completion.chunk",
            "created": 1,
            "model": "backend",
            "choices": [{"index": 0, "delta": delta, "finish_reason": reason}],
        }
        return httpx2.Response(
            200,
            text=f"data: {json.dumps(event)}\n\ndata: [DONE]\n\n",
            headers={"content-type": "text/event-stream"},
        )

    configured = definition(api_format="openai_chat", api_keys=["upstream-key"])
    instance = provider(configured, openai_handler=reply)
    app = create_test_app(
        Settings(
            MODEL=f"{configured.provider_id}/backend",
            custom_providers=(configured,),
            PROXY_AUTH_ENABLED=True,
            proxy_auth_token="ingress-key",
        ),
        providers={configured.provider_id: instance},
    )
    ready = threading.Event()

    class ReadyServer(uvicorn.Server):
        @override
        async def startup(self, sockets: list[socket.socket] | None = None) -> None:
            await super().startup(sockets=sockets)
            ready.set()

    with socket.socket() as sock:
        sock.bind(("127.0.0.1", 0))
        sock.listen()
        port = cast(tuple[str, int], sock.getsockname())[1]
        server = ReadyServer(uvicorn.Config(app, log_level="error", loop="asyncio"))
        thread = threading.Thread(target=server.run, kwargs={"sockets": [sock]})
        thread.start()
        try:
            assert ready.wait(20), "FCC test server did not start"
            result = subprocess.run(
                [
                    str(binary_path),
                    "run",
                    "--endpoint",
                    f"http://127.0.0.1:{port}/v1/chat/completions",
                    "--model",
                    "public",
                    "--workspace",
                    str(tmp_path),
                    "--prompt-stdin",
                ],
                input="Read the demo file",
                text=True,
                capture_output=True,
                timeout=30,
                env={**os.environ, "OMM_API_KEY": "ingress-key"},
                check=False,
            )
            assert result.returncode == 0, result.stderr
            assert result.stdout.strip() == "FCC mock tool loop succeeded"
            assert seen == ["FCC fixture content"]
            assert "ingress-key" not in result.stderr
        finally:
            server.should_exit = True
            thread.join(timeout=20)
            assert not thread.is_alive()
