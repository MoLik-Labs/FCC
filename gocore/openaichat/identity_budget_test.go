package openaichat

import (
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"testing"

	"github.com/NggKhaiz/FCC/gocore/retention"
)

func TestCurrentPythonIdentityBudget(t *testing.T) {
	f, err := os.Open(filepath.Join("..", "testdata", "golden_chat_identity_budget.json"))
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	var fixture struct {
		Schema    string `json:"schema"`
		Scenarios []struct {
			Name       string `json:"name"`
			MaxBytes   int    `json:"max_bytes"`
			MaxRecords int    `json:"max_records"`
			Steps      []struct {
				Index   any     `json:"index"`
				ID      any     `json:"id"`
				Slot    *int    `json:"slot"`
				ToolID  *string `json:"tool_id"`
				Error   *string `json:"error"`
				Bytes   int     `json:"bytes"`
				Records int     `json:"records"`
				Calls   []struct {
					Index  *int    `json:"index"`
					ToolID *string `json:"tool_id"`
				} `json:"calls"`
			} `json:"steps"`
		} `json:"scenarios"`
	}
	decoder := json.NewDecoder(f)
	decoder.UseNumber()
	if err := decoder.Decode(&fixture); err != nil {
		t.Fatal(err)
	}
	if fixture.Schema != "chat-identity-budget-v1" || len(fixture.Scenarios) == 0 {
		t.Fatal("invalid budget fixture")
	}
	for _, scenario := range fixture.Scenarios {
		t.Run(scenario.Name, func(t *testing.T) {
			budget := retention.NewPendingBudget(scenario.MaxBytes, scenario.MaxRecords)
			resolver := NewToolCallIdentityResolverWithBudget(budget)
			for i, step := range scenario.Steps {
				label := "step " + strconv.Itoa(i)
				slot, id, err := resolver.ResolveCurrent(step.Index, step.ID)
				if step.Error != nil {
					if err == nil || err.Error() != *step.Error {
						t.Fatalf("%s: error=%v want=%q", label, err, *step.Error)
					}
				} else if err != nil || step.Slot == nil || slot != *step.Slot || !reflect.DeepEqual(id, step.ToolID) {
					t.Fatalf("%s: slot=%d id=%v err=%v want slot=%v id=%v", label, slot, id, err, step.Slot, step.ToolID)
				}
				if budget.Bytes() != step.Bytes || budget.Records() != step.Records {
					t.Fatalf("%s: bytes=%d records=%d want=%d/%d", label, budget.Bytes(), budget.Records(), step.Bytes, step.Records)
				}
				snapshot := resolver.Snapshot()
				if len(snapshot) != len(step.Calls) {
					t.Fatalf("%s: %d calls want %d", label, len(snapshot), len(step.Calls))
				}
				for j, want := range step.Calls {
					if !reflect.DeepEqual(snapshot[j].Index, want.Index) || !reflect.DeepEqual(snapshot[j].ToolID, want.ToolID) {
						t.Fatalf("%s: call %d = %#v, want %#v", label, j, snapshot[j], want)
					}
				}
			}
		})
	}
}
