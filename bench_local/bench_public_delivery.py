"""Measure packet-based public delivery on a deterministic Messages stream.

Run from the repository root after installing the Python project:
    .venv/bin/python bench_local/bench_public_delivery.py

This is a component benchmark, not an end-to-end provider benchmark. To compare
revisions, run the same script with PYTHONPATH pointing to each revision's src.
"""

import asyncio
import hashlib
import json
import statistics
import time

from free_claude_code.api.stream_delivery import PublicResponseStream
from free_claude_code.core.json_types import JsonObject
from free_claude_code.core.stream_delivery import StreamDeliveryState
from free_claude_code.core.stream_output import StreamOutput


def frame(payload: JsonObject) -> str:
    return f"event: {payload['type']}\ndata: {json.dumps(payload)}\n\n"


FRAMES = (
    frame({"type": "message_start", "message": {"id": "msg", "content": []}}),
    frame(
        {
            "type": "content_block_start",
            "index": 0,
            "content_block": {"type": "text", "text": ""},
        }
    ),
    *(
        frame(
            {
                "type": "content_block_delta",
                "index": 0,
                "delta": {"type": "text_delta", "text": f"chunk {i} unicode 中文"},
            }
        )
        for i in range(2000)
    ),
    frame({"type": "content_block_stop", "index": 0}),
    frame({"type": "message_stop"}),
)


async def run():
    async def source():
        for value in FRAMES:
            yield StreamOutput((value,))

    stream = PublicResponseStream(
        source(),
        StreamDeliveryState(),
        wire_api="messages",
        first_chunk=StreamOutput((FRAMES[0],)),
        terminal_frame=None,
        terminal_failure_observer=None,
    )
    try:
        return "".join([value async for value in stream])
    finally:
        await stream.aclose()


async def main():
    expected = await run()
    timings = []
    for _ in range(7):
        start = time.perf_counter()
        output = await run()
        timings.append(time.perf_counter() - start)
        if output != expected:
            raise AssertionError("non-deterministic public output")
    print(f"frames: {len(FRAMES)}; median: {statistics.median(timings) * 1000:.2f} ms")
    print(f"output sha256: {hashlib.sha256(expected.encode()).hexdigest()}")


if __name__ == "__main__":
    asyncio.run(main())
