package openaichat

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strconv"
	"testing"
)

func TestCurrentPythonChatIdentity(t *testing.T) {
	f, err := os.Open(filepath.Join("..", "testdata", "golden_chat_identity_current.json"))
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	var fixture struct {
		Schema    string `json:"schema"`
		Scenarios []struct {
			Name  string `json:"name"`
			Steps []struct {
				Index  any     `json:"index"`
				ID     any     `json:"id"`
				Slot   *int    `json:"slot"`
				ToolID *string `json:"tool_id"`
				Error  *string `json:"error"`
			} `json:"steps"`
		} `json:"scenarios"`
	}
	decoder := json.NewDecoder(f)
	decoder.UseNumber()
	if err := decoder.Decode(&fixture); err != nil {
		t.Fatal(err)
	}
	if fixture.Schema != "chat-identity-current-v1" || len(fixture.Scenarios) == 0 {
		t.Fatalf("unexpected fixture: %s, %d scenarios", fixture.Schema, len(fixture.Scenarios))
	}
	for _, scenario := range fixture.Scenarios {
		t.Run(scenario.Name, func(t *testing.T) {
			resolver := NewToolCallIdentityResolver()
			for i, step := range scenario.Steps {
				label := "step " + strconv.Itoa(i)
				slot, toolID, err := resolver.ResolveCurrent(step.Index, step.ID)
				if step.Error != nil {
					if err == nil || err.Error() != *step.Error {
						t.Fatalf("%s: error = %v, want %q", label, err, *step.Error)
					}
					continue
				}
				if err != nil || step.Slot == nil || slot != *step.Slot {
					t.Fatalf("%s: slot=%d error=%v, want slot=%v", label, slot, err, step.Slot)
				}
				jsonEqual(t, label+"/tool_id", step.ToolID, toolID)
			}
		})
	}
}
