"""Opt-in interoperability test with the official Python MCP SDK."""

import os
import subprocess
from pathlib import Path

import pytest

SDK_SCRIPT = """
import asyncio
import sys
from mcp import ClientSession, StdioServerParameters
from mcp.client.stdio import stdio_client

async def main():
    params = StdioServerParameters(command=sys.argv[1], args=["mcp", "serve", "--workspace", sys.argv[2]])
    async with stdio_client(params) as (read, write):
        async with ClientSession(read, write) as session:
            init = await session.initialize()
            assert init.protocolVersion == "2025-06-18"
            tools = await session.list_tools()
            assert [tool.name for tool in tools.tools] == ["list_files", "read_file"]
            read_result = await session.call_tool("read_file", {"path": "hello.txt"})
            assert not read_result.isError and read_result.content[0].text == "SDK fixture content"
            escape = await session.call_tool("read_file", {"path": "../outside"})
            assert escape.isError
            listing = await session.call_tool("list_files", {"path": "."})
            assert not listing.isError and "hello.txt" in listing.content[0].text
            print("official MCP SDK roundtrip passed")

asyncio.run(main())
"""


def test_official_mcp_sdk_roundtrip(tmp_path: Path) -> None:
    binary = os.environ.get("MOLIK_BINARY")
    sdk_python = os.environ.get("MCP_SDK_PYTHON")
    if not binary or not sdk_python:
        pytest.skip(
            "Set MOLIK_BINARY and MCP_SDK_PYTHON to run the official SDK smoke test"
        )
    binary_path = Path(binary).resolve()
    assert binary_path.is_file()
    (tmp_path / "hello.txt").write_text("SDK fixture content", encoding="utf-8")
    result = subprocess.run(
        [sdk_python, "-c", SDK_SCRIPT, str(binary_path), str(tmp_path)],
        capture_output=True,
        text=True,
        timeout=30,
        check=False,
    )
    assert result.returncode == 0, result.stderr
    assert "official MCP SDK roundtrip passed" in result.stdout
