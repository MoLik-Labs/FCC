"""Record current Python source-order buffer transitions without old fixtures."""

import json
import random
from contextlib import suppress
from pathlib import Path

from free_claude_code.core.json_types import JsonObject, JsonValue, json_objects
from free_claude_code.core.tool_call_buffer import (
    OrderedToolCallBuffer,
    RetentionBudget,
    ToolBufferLimitError,
)

OUTPUT = Path("gocore/testdata/golden_ordered_buffer_current.json")


def main() -> None:
    scenarios: list[JsonObject] = [
        {
            "name": "source_order_and_lazy_release",
            "max_bytes": 100,
            "max_records": 10,
            "operations": [
                {"op": "reserve", "key": "a", "value": "A", "size": 2},
                {"op": "append", "value": "text", "size": 3},
                {"op": "reserve", "key": "b", "value": "B", "size": 4},
                {"op": "complete", "key": "b"},
                {"op": "next"},
                {"op": "complete", "key": "a"},
                {"op": "next"},
                {"op": "next"},
                {"op": "next"},
                {"op": "next"},
            ],
        },
        {
            "name": "discard_only_pending",
            "max_bytes": 100,
            "max_records": 10,
            "operations": [
                {"op": "reserve", "key": "a", "value": "unfinished", "size": 2},
                {"op": "append", "value": "text", "size": 3},
                {"op": "reserve", "key": "b", "value": "done", "size": 4},
                {"op": "complete", "key": "b"},
                {"op": "discard_pending"},
                {"op": "next"},
                {"op": "next"},
                {"op": "next"},
            ],
        },
        {
            "name": "failed_charges_leave_queue_unchanged",
            "max_bytes": 4,
            "max_records": 2,
            "operations": [
                {"op": "reserve", "key": "a", "value": "A", "size": 2},
                {"op": "grow", "key": "a", "size": 2, "records": 1},
                {"op": "grow", "key": "a", "size": 1, "records": 0},
                {"op": "append", "value": "blocked", "size": 0},
                {"op": "complete", "key": "a"},
                {"op": "next"},
                {"op": "append", "value": "after", "size": 4},
                {"op": "next"},
            ],
        },
        {
            "name": "duplicate_and_missing_keys",
            "max_bytes": 100,
            "max_records": 10,
            "operations": [
                {"op": "reserve", "key": "a", "value": "A", "size": 2},
                {"op": "reserve", "key": "a", "value": "other", "size": 90},
                {"op": "value", "key": "missing"},
                {"op": "grow", "key": "missing", "size": 99, "records": 1},
                {"op": "complete", "key": "missing"},
                {"op": "discard", "key": "missing"},
                {"op": "value", "key": "a"},
                {"op": "clear"},
                {"op": "value", "key": "a"},
            ],
        },
        {
            "name": "reuse_key_before_drain_and_clear",
            "max_bytes": 100,
            "max_records": 10,
            "operations": [
                {"op": "reserve", "key": "x", "value": "old", "size": 1},
                {"op": "complete", "key": "x"},
                {"op": "reserve", "key": "x", "value": "new", "size": 2},
                {"op": "next"},
                {"op": "value", "key": "x"},
                {"op": "grow", "key": "x", "size": 3, "records": 2},
                {"op": "discard", "key": "x"},
                {"op": "next"},
                {"op": "append", "value": "kept", "size": 5},
                {"op": "reserve", "key": "z", "value": "pending", "size": 4},
                {"op": "clear"},
                {"op": "next"},
            ],
        },
        {
            "name": "shared_budget_is_not_reset_by_clear",
            "max_bytes": 10,
            "max_records": 4,
            "operations": [
                {"op": "external_add", "size": 3, "records": 1},
                {"op": "reserve", "key": "a", "value": "A", "size": 4},
                {"op": "append", "value": "B", "size": 3},
                {"op": "grow", "key": "a", "size": 1, "records": 0},
                {"op": "clear"},
                {"op": "external_release", "size": 3, "records": 1},
            ],
        },
        {
            "name": "null_and_structured_values",
            "max_bytes": 12,
            "max_records": 4,
            "operations": [
                {"op": "append", "value": None, "size": 1},
                {"op": "reserve", "key": "a", "value": {"x": [1, "é"]}, "size": 5},
                {"op": "append", "value": ["after", 2], "size": 3},
                {"op": "next"},
                {"op": "value", "key": "a"},
                {"op": "complete", "key": "a"},
                {"op": "next"},
                {"op": "next"},
                {"op": "next"},
            ],
        },
        {
            "name": "zero_limits_and_default_zero_size",
            "max_bytes": 0,
            "max_records": 1,
            "operations": [
                {"op": "append", "value": "zero", "size": 0},
                {"op": "reserve", "key": "x", "value": "blocked", "size": 0},
                {"op": "next"},
                {"op": "reserve", "key": "x", "value": "ok", "size": 0},
                {"op": "grow", "key": "x", "size": 0, "records": 0},
                {"op": "complete", "key": "x"},
                {"op": "next"},
            ],
        },
    ]
    randomizer = random.Random(20261027)
    for case in range(16):
        operations: list[JsonValue] = []
        for step in range(35):
            op = randomizer.choice(
                (
                    "reserve",
                    "append",
                    "grow",
                    "value",
                    "complete",
                    "discard",
                    "discard_pending",
                    "next",
                    "clear",
                )
            )
            operation: JsonObject = {"op": op}
            if op in ("reserve", "grow", "value", "complete", "discard"):
                operation["key"] = randomizer.choice(("a", "b", "c"))
            if op in ("reserve", "append"):
                operation["value"] = None if step % 7 == 0 else f"item-{case}-{step}"
            if op in ("reserve", "append", "grow"):
                operation["size"] = randomizer.randrange(5)
            if op == "grow":
                operation["records"] = randomizer.randrange(3)
            operations.append(operation)
        scenarios.append(
            {
                "name": f"seeded_{case:02d}",
                "max_bytes": 12,
                "max_records": 6,
                "operations": operations,
            }
        )
    observed: list[JsonValue] = []
    for scenario in scenarios:
        max_bytes, max_records = scenario["max_bytes"], scenario["max_records"]
        if not isinstance(max_bytes, int) or not isinstance(max_records, int):
            raise TypeError("Expected integer budget limits")
        budget = RetentionBudget(max_bytes=max_bytes, max_records=max_records)
        buffer: OrderedToolCallBuffer[JsonValue] = OrderedToolCallBuffer(budget)
        steps: list[JsonValue] = []
        for operation in json_objects(scenario["operations"]):
            op = operation["op"]
            key = operation.get("key", "")
            size = operation.get("size", 0)
            records = operation.get("records", 1)
            if (
                not isinstance(key, str)
                or not isinstance(size, int)
                or not isinstance(records, int)
            ):
                raise TypeError("Invalid scenario operation")
            output: JsonValue = None
            has_value = False
            error: str | None = None
            try:
                if op == "reserve":
                    buffer.reserve(key, operation["value"], size=size)
                elif op == "append":
                    buffer.append(operation["value"], size=size)
                elif op == "grow":
                    buffer.grow(key, size, records)
                elif op == "value":
                    output, has_value = buffer.value(key), True
                elif op == "complete":
                    buffer.complete(key)
                elif op == "discard":
                    buffer.discard(key)
                elif op == "discard_pending":
                    buffer.discard_pending()
                elif op == "clear":
                    buffer.clear()
                elif op == "external_add":
                    budget.add(size, records)
                elif op == "external_release":
                    budget.release(size, records)
                elif op == "next":
                    iterator = buffer.drain()
                    with suppress(StopIteration):
                        output, has_value = next(iterator), True
                else:
                    raise ValueError(f"Unknown operation: {op}")
            except ToolBufferLimitError:
                error = "limit"
            except KeyError:
                error = "missing_key"
            except ValueError:
                error = "duplicate_key"
            keys = sorted(key for key in buffer._pending if isinstance(key, str))
            if len(keys) != len(buffer._pending):
                raise TypeError("Expected string keys")
            steps.append(
                {
                    **operation,
                    "output": output,
                    "has_value": has_value,
                    "error": error,
                    "pending": buffer.pending,
                    "keys": keys,
                    "bytes": budget.bytes,
                    "records_total": budget.records,
                    "queue": [
                        {
                            "value": slot.value,
                            "ready": slot.ready,
                            "discarded": slot.discarded,
                            "size": slot.size,
                            "records": slot.records,
                        }
                        for slot in buffer._queue
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
    OUTPUT.write_text(
        json.dumps(
            {"schema": "ordered-buffer-current-v1", "scenarios": observed},
            ensure_ascii=False,
            indent=2,
        )
        + "\n"
    )
    print(f"wrote {OUTPUT}: {len(observed)} scenarios")


if __name__ == "__main__":
    main()
