package openaichat

import (
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

func TestOutputCapUnicodeAndQuoteDifferential(t *testing.T) {
	data, err := os.ReadFile(filepath.Join("..", "testdata", "golden_outputcap_edges.json"))
	if err != nil {
		t.Fatal(err)
	}
	var fixture struct {
		Schema string `json:"schema"`
		Cases  []struct {
			Text   string `json:"text"`
			Output *int   `json:"output"`
		} `json:"cases"`
	}
	if err := json.Unmarshal(data, &fixture); err != nil {
		t.Fatal(err)
	}
	if fixture.Schema != "outputcap-edges-v1" || len(fixture.Cases) != 1570 {
		t.Fatal("invalid corpus")
	}
	for i, tc := range fixture.Cases {
		got := parseCap(strings.ToLower(tc.Text))
		if !reflect.DeepEqual(got, tc.Output) {
			t.Errorf("case %d %q: got %v, want %v", i, tc.Text, intPointerValue(got), intPointerValue(tc.Output))
		}
	}
}
