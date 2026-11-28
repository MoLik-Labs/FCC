package openaichat

import (
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/NggKhaiz/FCC/gocore/wirejson"
)

type anthropicAction struct {
	Op           string          `json:"op"`
	Index        int             `json:"index"`
	Value        string          `json:"value"`
	ID           string          `json:"id"`
	Extra        json.RawMessage `json:"extra"`
	Opaque       []string        `json:"opaque"`
	Summary      bool            `json:"summary"`
	InputTokens  *int            `json:"input_tokens"`
	OutputTokens *int            `json:"output_tokens"`
	UsageFields  json.RawMessage `json:"usage_fields"`
}

func orderedFixtureValue(t *testing.T, raw json.RawMessage) any {
	t.Helper()
	if len(raw) == 0 {
		return nil
	}
	value, err := wirejson.DecodeOrdered(string(raw))
	if err != nil {
		t.Fatal(err)
	}
	return value
}

func TestAnthropicChatOutputGolden(t *testing.T) {
	data, err := os.ReadFile(filepath.Join("..", "testdata", "golden_anthropic_chat_output.json"))
	if err != nil {
		t.Fatal(err)
	}
	var fixture struct {
		Schema    string `json:"schema"`
		Scenarios []struct {
			Name  string `json:"name"`
			Steps []struct {
				Input  anthropicAction `json:"input"`
				Events []string        `json:"events"`
			} `json:"steps"`
		} `json:"scenarios"`
	}
	if err := json.Unmarshal(data, &fixture); err != nil {
		t.Fatal(err)
	}
	if fixture.Schema != "anthropic-chat-output-v1" || len(fixture.Scenarios) == 0 {
		t.Fatal("invalid fixture")
	}
	for _, scenario := range fixture.Scenarios {
		t.Run(scenario.Name, func(t *testing.T) {
			output := NewAnthropicChatStreamOutput("msg-test", "model-é", 7)
			for i, step := range scenario.Steps {
				got := runAnthropicAction(t, output, step.Input)
				if len(got) == 0 {
					got = []string{}
				}
				if !reflect.DeepEqual(got, step.Events) {
					t.Fatalf("step %d %s:\ngot  %#v\nwant %#v", i, step.Input.Op, got, step.Events)
				}
			}
		})
	}
}

func runAnthropicAction(t *testing.T, output *AnthropicChatStreamOutput, a anthropicAction) []string {
	t.Helper()
	switch a.Op {
	case "start":
		return output.StartEvents()
	case "reasoning":
		return output.EnsureReasoningBlock(a.Index)
	case "begin":
		return output.BeginReasoningRecord(a.Index)
	case "reasoning_delta":
		return []string{output.EmitReasoningDelta(a.Value, a.Summary, a.Index)}
	case "complete":
		return output.CompleteReasoningRecord(&a.Opaque, a.Index)
	case "opaque":
		return output.EmitOpaqueReasoning(a.Value)
	case "defer":
		output.DeferOpaqueReasoning(a.Value)
	case "text":
		return output.EnsureTextBlock()
	case "text_delta":
		return []string{output.EmitTextDelta(a.Value)}
	case "tool_start":
		return []string{output.StartToolBlock(a.Index, a.ID, a.Value, orderedFixtureValue(t, a.Extra))}
	case "tool_delta":
		return []string{output.EmitToolDelta(a.Index, a.Value)}
	case "tool_stop":
		return output.StopToolBlock(a.Index)
	case "close_content":
		return output.CloseContentBlocks()
	case "close_all":
		return output.CloseAllBlocks()
	case "close_unclosed":
		return output.CloseUnclosedBlocks()
	case "success":
		in, out := 7, 13
		if a.InputTokens != nil {
			in = *a.InputTokens
		}
		if a.OutputTokens != nil {
			out = *a.OutputTokens
		}
		return output.FinishSuccess(a.Value, ChatStreamUsage{InputTokens: in, OutputTokens: out, AnthropicFields: orderedFixtureValue(t, a.UsageFields)})
	case "failure":
		return output.FinishFailure(nil)
	default:
		t.Fatalf("unknown op %q", a.Op)
	}
	return nil
}

func TestAnthropicFailureClosesStateWithoutClosingTools(t *testing.T) {
	output := NewAnthropicChatStreamOutput("msg", "model", 1)
	output.BeginReasoningRecord(7)
	output.StartToolBlock(2, "tool-id", "tool", nil)
	events := output.FinishFailure(nil)
	if len(events) != 1 || events[0] != "event: content_block_stop\ndata: {\"type\": \"content_block_stop\", \"index\": 0}\n\n" {
		t.Fatalf("unexpected failure cleanup: %#v", events)
	}
	state := output.ToolState(2)
	if state == nil || state.Open || !state.Started || output.FinalStopReason("stop") != "tool_use" {
		t.Fatal("failure cleanup lost the tool state or left it open")
	}
	sink := output.ChatStreamOutput.sink.(*anthropicChatSink)
	if len(sink.openStack) != 0 || len(sink.blockTypes) != 0 || len(sink.reasoningIndices) != 0 {
		t.Fatal("failure cleanup left open ledger state")
	}
	if len(output.FinishFailure(nil)) != 0 || len(output.FinishSuccess("stop", nil)) != 0 {
		t.Fatal("terminal transition was not idempotent")
	}
}
