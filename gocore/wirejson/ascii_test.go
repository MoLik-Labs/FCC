package wirejson

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestASCIIJSONGolden(t *testing.T) {
	raw, err := os.ReadFile(filepath.Join("..", "testdata", "golden_ascii_json.json"))
	if err != nil {
		t.Fatal(err)
	}
	var fixture struct {
		Schema string                           `json:"schema"`
		Cases  []struct{ Input, Output string } `json:"cases"`
	}
	if err := json.Unmarshal(raw, &fixture); err != nil {
		t.Fatal(err)
	}
	if fixture.Schema != "ascii-json-v1" || len(fixture.Cases) != 1005 {
		t.Fatal("invalid fixture")
	}
	for i, tc := range fixture.Cases {
		value, err := DecodeOrdered(tc.Input)
		if err != nil {
			t.Fatalf("case %d: %v", i, err)
		}
		got, err := EncodeASCII(value)
		if err != nil {
			t.Fatalf("case %d: %v", i, err)
		}
		if got != tc.Output {
			t.Errorf("case %d: input %s got %s want %s", i, tc.Input, got, tc.Output)
		}
	}
}

func TestOrderedJSONRejectsInvalid(t *testing.T) {
	for _, raw := range []string{"", "{}{}", "[1,]", "{\"key\":}", "null true"} {
		if _, err := DecodeOrdered(raw); err == nil {
			t.Errorf("accepted %q", raw)
		}
	}
	for _, value := range []any{json.Number("true"), json.Number("01"), make(chan int)} {
		if _, err := EncodeASCII(value); err == nil {
			t.Errorf("encoded invalid value %v", value)
		}
	}
}
