"""Capture current Python Chat identity and pending-budget transitions."""

import json
from pathlib import Path

from free_claude_code.core.chat_completions.tool_calls import (
    AmbiguousToolCallIdentityError,
    ToolCallIdentityResolver,
)
from free_claude_code.core.json_types import JsonObject, JsonValue, json_objects
from free_claude_code.core.tool_call_buffer import RetentionBudget, ToolBufferLimitError

OUTPUT = Path("gocore/testdata/golden_chat_identity_budget.json")


def main() -> None:
    scenarios: list[JsonObject] = [
        {
            "name": "new_and_repeated_identity",
            "max_bytes": 100,
            "max_records": 2,
            "inputs": [
                {"index": 0, "id": "alpha"},
                {"index": 0, "id": "alpha"},
                {"index": 1, "id": "β"},
            ],
        },
        {
            "name": "exact_slot_limit",
            "max_bytes": 32,
            "max_records": 1,
            "inputs": [{"index": 2, "id": None}, {"index": 2, "id": None}],
        },
        {
            "name": "slot_byte_limit_then_retry",
            "max_bytes": 31,
            "max_records": 1,
            "inputs": [{"index": 0, "id": None}, {"index": 0, "id": None}],
        },
        {
            "name": "id_limit_after_new_slot",
            "max_bytes": 34,
            "max_records": 1,
            "inputs": [{"index": 0, "id": "abc"}, {"index": 0, "id": "ab"}],
        },
        {
            "name": "utf8_id_exact_and_over_limit",
            "max_bytes": 36,
            "max_records": 1,
            "inputs": [
                {"index": 0, "id": "🦊x"},
                {"index": 0, "id": "🦊"},
                {"index": 0, "id": "🦊"},
            ],
        },
        {
            "name": "record_limit_retains_first_slot",
            "max_bytes": 1000,
            "max_records": 1,
            "inputs": [{"index": 0, "id": None}, {"index": 1, "id": None}],
        },
        {
            "name": "id_added_to_existing_slot_without_record_charge",
            "max_bytes": 35,
            "max_records": 1,
            "inputs": [
                {"index": 0, "id": None},
                {"index": 0, "id": "abc"},
                {"index": 0, "id": "abc"},
            ],
        },
        {
            "name": "invalid_identity_does_not_charge",
            "max_bytes": 100,
            "max_records": 1,
            "inputs": [{"index": -1, "id": "  "}, {"index": 0, "id": None}],
        },
        {
            "name": "ambiguous_identity_does_not_charge",
            "max_bytes": 200,
            "max_records": 3,
            "inputs": [
                {"index": 0, "id": None},
                {"index": None, "id": "x"},
                {"index": 0, "id": "x"},
            ],
        },
        {
            "name": "zero_record_budget",
            "max_bytes": 100,
            "max_records": 0,
            "inputs": [{"index": 0, "id": None}],
        },
    ]
    observed: list[JsonValue] = []
    for scenario in scenarios:
        max_bytes, max_records = scenario["max_bytes"], scenario["max_records"]
        if not isinstance(max_bytes, int) or not isinstance(max_records, int):
            raise TypeError("Expected integer budget limits")
        budget = RetentionBudget(max_bytes=max_bytes, max_records=max_records)
        resolver = ToolCallIdentityResolver(budget)
        steps: list[JsonValue] = []
        for item in json_objects(scenario["inputs"]):
            try:
                slot, tool_id = resolver.resolve(item["index"], item["id"])
                result: JsonObject = {"slot": slot, "tool_id": tool_id, "error": None}
            except (AmbiguousToolCallIdentityError, ToolBufferLimitError) as exc:
                result = {"slot": None, "tool_id": None, "error": str(exc)}
            steps.append(
                {
                    **item,
                    **result,
                    "bytes": budget.bytes,
                    "records": budget.records,
                    "calls": [
                        {"index": call.index, "tool_id": call.tool_id}
                        for call in resolver._calls
                    ],
                }
            )
        observed.append(
            {
                "name": scenario["name"],
                "max_bytes": max_bytes,
                "max_records": max_records,
                "steps": steps,
            }
        )
    fixture = {"schema": "chat-identity-budget-v1", "scenarios": observed}
    OUTPUT.write_text(json.dumps(fixture, ensure_ascii=False, indent=2) + "\n")
    print(f"wrote {OUTPUT}: {len(observed)} scenarios")


if __name__ == "__main__":
    main()
