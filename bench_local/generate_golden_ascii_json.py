"""Generate ordered ASCII JSON fixtures from Python's json module."""

import json
import math
import random
import struct
from pathlib import Path
from typing import cast

from free_claude_code.core.json_types import load_json_value

OUTPUT = Path("gocore/testdata/golden_ascii_json.json")


def main() -> None:
    inputs = [
        "{}",
        '{"z":1,"a":{"y":true,"x":["é😀\\u007f",null]}}',
        '{"x":1,"y":2,"x":3}',
        '["\\n\\r\\t\\b\\f\\u0000\\u007f <>& /", 1, 1.0, -0.0, 1e6, 1e16, 1e-4, 1e-5]',
        '{"big":123456789012345678901234567890,"float":9.999999999999999e15}',
        "-0",
        "1e999",
        "-1e999",
    ]
    randomizer = random.Random(61824)
    while len(inputs) < 1005:
        value = cast(
            float, struct.unpack(">d", randomizer.getrandbits(64).to_bytes(8, "big"))[0]
        )
        if math.isfinite(value):
            inputs.append(json.dumps(value))
    fixture = {
        "schema": "ascii-json-v1",
        "cases": [
            {"input": text, "output": json.dumps(load_json_value(text))}
            for text in inputs
        ],
    }
    OUTPUT.write_text(
        json.dumps(fixture, ensure_ascii=False, separators=(",", ":")) + "\n",
        encoding="utf-8",
    )
    print(f"wrote {OUTPUT}: {len(inputs)} cases")


if __name__ == "__main__":
    main()
