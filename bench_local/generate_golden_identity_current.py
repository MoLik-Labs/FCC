"""Observe current Python Chat identity resolution on the existing input corpus.

This new fixture is separate from the historical leaf fixture: the current
resolver returns both the slot and its learned ID, has different failure
messages, and no longer exposes explicit allocation or ordering operations.
"""

import json
from pathlib import Path

from free_claude_code.core.chat_completions.tool_calls import (
    AmbiguousToolCallIdentityError,
    ToolCallIdentityResolver,
)

DATA = Path("gocore/testdata")


def main():
    historical = json.loads((DATA / "golden_openaichat4c4a.json").read_text())
    inputs = [
        previous
        for previous in historical["identity"]
        if all(step["op"] == "resolve" for step in previous["steps"])
    ]
    inputs.extend(
        {
            "name": f"python_control_whitespace_{code:02x}",
            "steps": [
                {"index": None, "id": chr(code)},
                {"index": 1, "id": chr(code) + "call"},
            ],
        }
        for code in range(0x1C, 0x20)
    )
    scenarios = []
    for previous in inputs:
        resolver = ToolCallIdentityResolver()
        steps = []
        for previous_step in previous["steps"]:
            index, call_id = previous_step["index"], previous_step["id"]
            try:
                slot, learned_id = resolver.resolve(index, call_id)
                result = {"slot": slot, "tool_id": learned_id, "error": None}
            except AmbiguousToolCallIdentityError as exc:
                result = {"slot": None, "tool_id": None, "error": str(exc)}
            steps.append({"index": index, "id": call_id, **result})
        scenarios.append({"name": previous["name"], "steps": steps})
    fixture = {"schema": "chat-identity-current-v1", "scenarios": scenarios}
    path = DATA / "golden_chat_identity_current.json"
    path.write_text(json.dumps(fixture, ensure_ascii=False, indent=2) + "\n")
    print(f"wrote {path}: {len(scenarios)} scenarios")


if __name__ == "__main__":
    main()
