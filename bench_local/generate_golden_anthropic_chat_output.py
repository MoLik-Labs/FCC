"""Generate exact Anthropic Chat writer SSE frames from the Python reference."""

import json
import random
from pathlib import Path

from free_claude_code.providers.openai_chat.stream_output import (
    AnthropicChatStreamOutput,
    ChatStreamUsage,
)

OUTPUT = Path("gocore/testdata/golden_anthropic_chat_output.json")


def apply(output, action):
    op = action["op"]
    index = action.get("index", 0)
    value = action.get("value", "")
    if op == "start":
        return output.start_events()
    if op == "reasoning":
        return output.ensure_reasoning_block(source=index)
    if op == "begin":
        return output.begin_reasoning_record(source=index)
    if op == "reasoning_delta":
        return [
            output.emit_reasoning_delta(
                value, summary=action.get("summary", False), source=index
            )
        ]
    if op == "complete":
        return output.complete_reasoning_record(
            list(action.get("opaque", [])), source=index
        )
    if op == "opaque":
        return output.emit_opaque_reasoning(value)
    if op == "defer":
        output.defer_opaque_reasoning(value)
        return []
    if op == "text":
        return output.ensure_text_block()
    if op == "text_delta":
        return [output.emit_text_delta(value)]
    if op == "tool_start":
        return [
            output.start_tool_block(
                index, action["id"], value, extra_content=action.get("extra")
            )
        ]
    if op == "tool_delta":
        return [output.emit_tool_delta(index, value)]
    if op == "tool_stop":
        return output.stop_tool_block(index)
    if op == "close_content":
        return output.close_content_blocks()
    if op == "close_all":
        return output.close_all_blocks()
    if op == "close_unclosed":
        return output.close_unclosed_blocks()
    if op == "success":
        return output.finish_success(
            stop_reason=value,
            usage=ChatStreamUsage(
                action.get("input_tokens", 7),
                action.get("output_tokens", 13),
                anthropic_fields=action.get("usage_fields", {}),
            ),
        )
    if op == "failure":
        return output.finish_failure(None)
    raise ValueError(op)


def main():
    scenarios = [
        (
            "text_success",
            [
                {"op": "start"},
                {"op": "text"},
                {"op": "text_delta", "value": 'Hi "<world>"!\n工具 😀 \u007f'},
                {"op": "text_delta", "value": ""},
                {"op": "success", "value": "end_turn"},
                {"op": "success", "value": "stop"},
            ],
        ),
        (
            "reasoning_interleaved",
            [
                {"op": "start"},
                {"op": "reasoning", "index": 4},
                {"op": "reasoning_delta", "index": 4, "value": "思考", "summary": True},
                {"op": "complete", "index": 4, "opaque": ["sig😀", "opaque"]},
                {"op": "text"},
                {"op": "text_delta", "value": "answer"},
                {"op": "reasoning", "index": 3},
                {"op": "reasoning_delta", "index": 3, "value": "r"},
                {"op": "text"},
                {"op": "success", "value": "tool_use"},
            ],
        ),
        (
            "tools_success",
            [
                {"op": "start"},
                {
                    "op": "tool_start",
                    "index": 4,
                    "id": "id-4",
                    "value": "工具",
                    "extra": {"z": {"b": 1, "a": [1.0, 1e-5, "é"]}, "a": "😀"},
                },
                {"op": "tool_delta", "index": 4, "value": '{"q": \u000301}'},
                {"op": "tool_start", "index": 1, "id": "id-1", "value": "read"},
                {"op": "tool_stop", "index": 1},
                {"op": "defer", "value": "secret"},
                {
                    "op": "success",
                    "value": "end_turn",
                    "usage_fields": {
                        "output_tokens": 99,
                        "z": 0,
                        "a": True,
                        "invalid": "skip",
                    },
                },
            ],
        ),
        (
            "failure_lifo",
            [
                {"op": "start"},
                {"op": "begin", "index": 3},
                {"op": "begin", "index": 1},
                {"op": "tool_start", "index": 7, "id": "call", "value": "tool"},
                {"op": "defer", "value": "dropped"},
                {"op": "failure"},
                {"op": "failure"},
            ],
        ),
        (
            "close_and_restart",
            [
                {"op": "start"},
                {"op": "text"},
                {"op": "close_unclosed"},
                {"op": "text"},
                {"op": "text_delta", "value": "new"},
                {"op": "close_all"},
                {"op": "opaque", "value": "sealed"},
                {"op": "success", "value": "stop"},
            ],
        ),
    ]
    randomizer = random.Random(44009)
    alphabet = [
        "a",
        "Z",
        " ",
        '"',
        "\\",
        "<",
        "&",
        "\n",
        "\t",
        "\x00",
        "\x7f",
        "é",
        "漢",
        "😀",
        "\u2028",
    ]
    for index in range(60):
        content = "".join(
            randomizer.choice(alphabet) for _ in range(randomizer.randrange(60))
        )
        scenarios.append(
            (
                f"text_unicode_{index}",
                [
                    {"op": "start"},
                    {"op": "text"},
                    {"op": "text_delta", "value": content},
                    {"op": "success", "value": "end_turn"},
                ],
            )
        )
    results = []
    for name, actions in scenarios:
        output = AnthropicChatStreamOutput(
            message_id="msg-test", model="model-é", input_tokens=7
        )
        steps = [
            {"input": action, "events": apply(output, action)} for action in actions
        ]
        results.append({"name": name, "steps": steps})
    OUTPUT.write_text(
        json.dumps(
            {"schema": "anthropic-chat-output-v1", "scenarios": results},
            ensure_ascii=False,
            indent=2,
        )
        + "\n",
        encoding="utf-8",
    )
    print(f"wrote {OUTPUT}: {len(results)} scenarios")


if __name__ == "__main__":
    main()
