package toolbuf

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"testing"

	"github.com/NggKhaiz/FCC/gocore/retention"
)

type currentQueueEntry struct {
	Value     any  `json:"value"`
	Ready     bool `json:"ready"`
	Discarded bool `json:"discarded"`
	Size      int  `json:"size"`
	Records   int  `json:"records"`
}

type currentBufferStep struct {
	Op           string              `json:"op"`
	Key          string              `json:"key"`
	Value        any                 `json:"value"`
	Size         int                 `json:"size"`
	Records      int                 `json:"records"`
	Output       any                 `json:"output"`
	HasValue     bool                `json:"has_value"`
	Error        *string             `json:"error"`
	Pending      bool                `json:"pending"`
	Keys         []string            `json:"keys"`
	Bytes        int                 `json:"bytes"`
	RecordsTotal int                 `json:"records_total"`
	Queue        []currentQueueEntry `json:"queue"`
}

func TestCurrentPythonOrderedBuffer(t *testing.T) {
	file, err := os.Open(filepath.Join("..", "testdata", "golden_ordered_buffer_current.json"))
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()
	var fixture struct {
		Schema    string `json:"schema"`
		Scenarios []struct {
			Name       string              `json:"name"`
			MaxBytes   int                 `json:"max_bytes"`
			MaxRecords int                 `json:"max_records"`
			Steps      []currentBufferStep `json:"steps"`
		} `json:"scenarios"`
	}
	decoder := json.NewDecoder(file)
	decoder.UseNumber()
	if err := decoder.Decode(&fixture); err != nil {
		t.Fatal(err)
	}
	if fixture.Schema != "ordered-buffer-current-v1" || len(fixture.Scenarios) == 0 {
		t.Fatal("invalid fixture")
	}
	for _, scenario := range fixture.Scenarios {
		t.Run(scenario.Name, func(t *testing.T) {
			budget := retention.NewPendingBudget(scenario.MaxBytes, scenario.MaxRecords)
			buffer := NewOrderedBuffer[string, any](budget)
			for index, step := range scenario.Steps {
				var value any
				var hasValue bool
				var err error
				switch step.Op {
				case "reserve":
					err = buffer.Reserve(step.Key, step.Value, step.Size)
				case "append":
					err = buffer.Append(step.Value, step.Size)
				case "grow":
					err = buffer.Grow(step.Key, step.Size, step.Records)
				case "value":
					value, err = buffer.Value(step.Key)
					hasValue = err == nil
				case "complete":
					err = buffer.Complete(step.Key)
				case "discard":
					err = buffer.Discard(step.Key)
				case "discard_pending":
					buffer.DiscardPending()
				case "clear":
					buffer.Clear()
				case "external_add":
					err = budget.Add(step.Size, step.Records)
				case "external_release":
					budget.Release(step.Size, step.Records)
				case "next":
					value, hasValue = buffer.Next()
				default:
					t.Fatalf("unknown op %q", step.Op)
				}
				var category string
				var limit *retention.PendingLimitError
				switch {
				case errors.As(err, &limit):
					category = "limit"
				case errors.Is(err, ErrMissingKey):
					category = "missing_key"
				case errors.Is(err, ErrDuplicateKey):
					category = "duplicate_key"
				case err != nil:
					t.Fatalf("step %d: unexpected error %v", index, err)
				}
				wantError := ""
				if step.Error != nil {
					wantError = *step.Error
				}
				if category != wantError || hasValue != step.HasValue || !reflect.DeepEqual(value, step.Output) {
					t.Fatalf("step %d (%s): result=%#v/%t/%s want=%#v/%t/%s", index, step.Op, value, hasValue, category, step.Output, step.HasValue, wantError)
				}
				if buffer.Pending() != step.Pending || budget.Bytes() != step.Bytes || budget.Records() != step.RecordsTotal {
					t.Fatalf("step %d: state pending=%t bytes=%d records=%d, want %t/%d/%d", index, buffer.Pending(), budget.Bytes(), budget.Records(), step.Pending, step.Bytes, step.RecordsTotal)
				}
				keys := make([]string, 0, len(buffer.pending))
				for key := range buffer.pending {
					keys = append(keys, key)
				}
				sort.Strings(keys)
				if !reflect.DeepEqual(keys, step.Keys) {
					t.Fatalf("step %d: pending keys=%v, want %v", index, keys, step.Keys)
				}
				if len(buffer.queue) != len(step.Queue) {
					t.Fatalf("step %d: queue size %d want %d", index, len(buffer.queue), len(step.Queue))
				}
				for i, item := range buffer.queue {
					got := currentQueueEntry{Value: item.value, Ready: item.ready, Discarded: item.discarded, Size: item.size, Records: item.records}
					if !reflect.DeepEqual(got, step.Queue[i]) {
						t.Fatalf("step %d queue %d: got %#v want %#v", index, i, got, step.Queue[i])
					}
				}
			}
		})
	}
}
