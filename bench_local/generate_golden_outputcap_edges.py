"""Differential corpus for Python's output-token cap regex semantics."""

import json
import random
from pathlib import Path

from free_claude_code.providers.openai_chat.output_cap import _parse_cap

OUTPUT = Path("gocore/testdata/golden_outputcap_edges.json")


def main() -> None:
    randomizer = random.Random(44004)
    fields = ["max_tokens", "max_completion_tokens"]
    surrounds = ["", "`", "'", '"', "``", "'`"]
    prefixes = ["", "x", "é", "漢", "_", "-", " "]
    suffixes = ["", "x", "é", "_", "-", " "]
    spaces = [" ", "\t", "\n", "\u00a0", "\u2003"]
    numbers = ["0", "1", "256", "\uff11\uff12", "٠٩", "1024", "99999"]
    patterns = [
        "{p}range{w}of{w}{field}{w}should{w}be{w}[{lower},{w}{num}]{s}",
        "{p}{field}{w}must{w}not{w}exceed{w}{num}{s}",
        "{p}maximum{w}allowed{w}value{w}for{w}{field}{w}is{w}{num}{s}",
        "{p}maximum{w}of{w}{num}{w}for{w}{field}{s}",
        "{p}{field}:{w}<={w}{num}{s}",
    ]
    cases = []
    for _ in range(1500):
        field = (
            randomizer.choice(surrounds)
            + randomizer.choice(fields)
            + randomizer.choice(surrounds)
        )
        number = (
            randomizer.choice(surrounds)
            + randomizer.choice(numbers)
            + randomizer.choice(surrounds)
        )
        text = randomizer.choice(patterns).format(
            p=randomizer.choice(prefixes),
            s=randomizer.choice(suffixes),
            w=randomizer.choice(spaces),
            field=field,
            num=number,
            lower=randomizer.choice(["1", "1", "\uff11", "\u0661"]),
        )
        cases.append({"text": text, "output": _parse_cap(text.lower())})

    # CPython 3.14's Unicode 16 adds word and decimal characters not present
    # in the Go 1.22 Unicode 15 tables. Exercise both endpoints of every
    # additional word range and all of the new digit blocks.
    word_ranges = [
        (0x1C89, 0x1C8A),
        (0xA7CB, 0xA7CD),
        (0xA7DA, 0xA7DC),
        (0x105C0, 0x105F3),
        (0x10D40, 0x10D65),
        (0x10D6F, 0x10D85),
        (0x10EC2, 0x10EC4),
        (0x11380, 0x11389),
        (0x1138B, 0x1138B),
        (0x1138E, 0x1138E),
        (0x11390, 0x113B5),
        (0x113B7, 0x113B7),
        (0x113D1, 0x113D1),
        (0x113D3, 0x113D3),
        (0x116D0, 0x116E3),
        (0x11BC0, 0x11BE0),
        (0x11BF0, 0x11BF9),
        (0x13460, 0x143FA),
        (0x16100, 0x1611D),
        (0x16130, 0x16139),
        (0x16D40, 0x16D6C),
        (0x16D70, 0x16D79),
        (0x18CFF, 0x18CFF),
        (0x1CCF0, 0x1CCF9),
        (0x1E5D0, 0x1E5ED),
        (0x1E5F0, 0x1E5FA),
        (0x2EBF0, 0x2EE5D),
    ]
    digits = [0x10D40, 0x116D0, 0x116DA, 0x11BF0, 0x16130, 0x16D70, 0x1CCF0, 0x1E5F1]
    additional = [
        f"{chr(code)}max_tokens at most 5"
        for low, high in word_ranges
        for code in (low, high)
    ]
    additional += [
        f"max_tokens at most {chr(start + digit)}"
        for start in digits
        for digit in (1, 9)
    ]
    cases.extend(
        {"text": text, "output": _parse_cap(text.lower())} for text in additional
    )
    OUTPUT.write_text(
        json.dumps(
            {"schema": "outputcap-edges-v1", "cases": cases},
            ensure_ascii=False,
            separators=(",", ":"),
        )
        + "\n",
        encoding="utf-8",
    )
    print(f"wrote {OUTPUT}: {len(cases)} cases")


if __name__ == "__main__":
    main()
