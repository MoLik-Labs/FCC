package openaichat

import (
	"encoding/json"
	"sort"
	"strconv"
	"strings"

	"github.com/NggKhaiz/FCC/gocore/wirejson"
	"github.com/NggKhaiz/FCC/gocore/wireuuid"
)

// ChatStreamUsage is the final Chat usage passed to the Anthropic writer.
// AnthropicFields may be a map[string]any or *wirejson.Object. Use ordered
// objects when the input key order must be retained in the emitted bytes.
type ChatStreamUsage struct {
	InputTokens     int
	OutputTokens    int
	AnthropicFields any
}

// AnthropicChatStreamOutput emits Anthropic Messages SSE from the shared
// Chat lifecycle. The embedded output is confined to one stream attempt.
type AnthropicChatStreamOutput struct {
	*ChatStreamOutput
}

func NewAnthropicChatStreamOutput(messageID, model string, inputTokens int) *AnthropicChatStreamOutput {
	if messageID == "" {
		messageID = "msg_" + wireuuid.MustV4().String()
	}
	sink := &anthropicChatSink{messageID: messageID, model: model, inputTokens: inputTokens,
		textIndex: -1, reasoningIndices: make(map[int]int), toolIndices: make(map[int]int), blockTypes: make(map[int]string)}
	output := NewChatStreamOutput(inputTokens, sink)
	sink.output = output
	return &AnthropicChatStreamOutput{ChatStreamOutput: output}
}

type anthropicChatSink struct {
	output                            *ChatStreamOutput
	messageID, model                  string
	inputTokens, nextIndex, textIndex int
	reasoningIndices                  map[int]int
	toolIndices                       map[int]int
	blockTypes                        map[int]string
	openStack                         []int
}

func ordered(fields ...any) *wirejson.Object {
	obj := wirejson.NewObject(len(fields) / 2)
	for i := 0; i < len(fields); i += 2 {
		obj.Set(fields[i].(string), fields[i+1])
	}
	return obj
}

func anthropicEvent(event string, data *wirejson.Object) string {
	encoded, err := wirejson.EncodeASCII(data)
	if err != nil {
		// Like Python json.dumps, reject non-JSON metadata.
		panic(err)
	}
	return "event: " + event + "\ndata: " + encoded + "\n\n"
}

func (s *anthropicChatSink) allocate() int { index := s.nextIndex; s.nextIndex++; return index }
func (s *anthropicChatSink) startBlock(kind string, block *wirejson.Object) string {
	index := s.allocate()
	s.blockTypes[index] = kind
	s.openStack = append(s.openStack, index)
	return anthropicEvent("content_block_start", ordered("type", "content_block_start", "index", index, "content_block", block))
}
func (s *anthropicChatSink) delta(index int, kind, key, content string) string {
	// Content deltas are the per-token hot path. The envelope keys and
	// delta kinds are internal constants; only content needs JSON escaping.
	var frame strings.Builder
	frame.Grow(160 + len(content))
	frame.WriteString("event: content_block_delta\ndata: {\"type\": \"content_block_delta\", \"index\": ")
	frame.WriteString(strconv.Itoa(index))
	frame.WriteString(`, "delta": {"type": "`)
	frame.WriteString(kind)
	frame.WriteString(`", "`)
	frame.WriteString(key)
	frame.WriteString(`": `)
	wirejson.WriteASCIIString(&frame, content)
	frame.WriteString("}}\n\n")
	return frame.String()
}
func (s *anthropicChatSink) stop(index int) string {
	for i, open := range s.openStack {
		if open == index {
			s.openStack = append(s.openStack[:i], s.openStack[i+1:]...)
			break
		}
	}
	delete(s.blockTypes, index)
	return anthropicEvent("content_block_stop", ordered("type", "content_block_stop", "index", index))
}

func (s *anthropicChatSink) StartEvents() []string {
	usage := ordered("input_tokens", s.inputTokens, "output_tokens", 1)
	message := ordered("id", s.messageID, "type", "message", "role", "assistant", "content", []any{}, "model", s.model, "stop_reason", nil, "stop_sequence", nil, "usage", usage)
	return []string{anthropicEvent("message_start", ordered("type", "message_start", "message", message))}
}
func (s *anthropicChatSink) StartReasoning(source int) []string {
	s.reasoningIndices[source] = s.nextIndex
	return []string{s.startBlock("thinking", ordered("type", "thinking", "thinking", ""))}
}
func (s *anthropicChatSink) ReasoningDelta(content string, summary bool, source int) string {
	return s.delta(s.reasoningIndex(source), "thinking_delta", "thinking", content)
}
func (s *anthropicChatSink) reasoningIndex(source int) int {
	index, ok := s.reasoningIndices[source]
	if !ok {
		panic("reasoning source is not open")
	}
	return index
}
func (s *anthropicChatSink) OpaqueReasoning(data string) []string {
	index := s.nextIndex
	start := s.startBlock("redacted_thinking", ordered("type", "redacted_thinking", "data", data))
	return []string{start, s.stop(index)}
}
func (s *anthropicChatSink) AttachReasoningReplay(data string, source int) []string {
	return []string{s.delta(s.reasoningIndex(source), "signature_delta", "signature", data)}
}
func (s *anthropicChatSink) StopReasoning(source int) []string {
	index := s.reasoningIndex(source)
	delete(s.reasoningIndices, source)
	return []string{s.stop(index)}
}
func (s *anthropicChatSink) StartText() []string {
	s.textIndex = s.nextIndex
	return []string{s.startBlock("text", ordered("type", "text", "text", ""))}
}
func (s *anthropicChatSink) TextDelta(content string) string {
	return s.delta(s.textIndex, "text_delta", "text", content)
}
func (s *anthropicChatSink) StopText() []string { return []string{s.stop(s.textIndex)} }
func (s *anthropicChatSink) StartTool(index int, state *ChatToolState) string {
	s.toolIndices[index] = s.nextIndex
	block := ordered("type", "tool_use", "id", state.ToolID, "name", state.Name, "input", ordered())
	if nonEmptyExtraContent(state.ExtraContent) {
		block.Set("extra_content", state.ExtraContent)
	}
	return s.startBlock("tool_use", block)
}
func (s *anthropicChatSink) ToolDelta(index int, state *ChatToolState, partialJSON string) string {
	return s.delta(s.toolIndex(index), "input_json_delta", "partial_json", partialJSON)
}
func (s *anthropicChatSink) toolIndex(index int) int {
	block, ok := s.toolIndices[index]
	if !ok {
		panic("tool is not started")
	}
	return block
}
func (s *anthropicChatSink) StopTool(index int, state *ChatToolState) []string {
	return []string{s.stop(s.toolIndex(index))}
}
func (s *anthropicChatSink) FinishSuccess(reason string, value any) []string {
	usage, ok := value.(ChatStreamUsage)
	if !ok {
		panic("Anthropic output requires ChatStreamUsage")
	}
	fields := ordered("input_tokens", usage.InputTokens, "output_tokens", usage.OutputTokens)
	appendUsageFields(fields, usage.AnthropicFields)
	delta := ordered("stop_reason", s.output.FinalStopReason(reason), "stop_sequence", nil)
	return []string{anthropicEvent("message_delta", ordered("type", "message_delta", "delta", delta, "usage", fields)), anthropicEvent("message_stop", ordered("type", "message_stop"))}
}
func appendUsageFields(target *wirejson.Object, value any) {
	var keys []string
	var fields map[string]any
	switch v := value.(type) {
	case *wirejson.Object:
		if v == nil {
			return
		}
		keys, fields = v.Keys, v.Fields
	case wirejson.Object:
		keys, fields = v.Keys, v.Fields
	case map[string]any:
		fields = v
		for key := range v {
			keys = append(keys, key)
		}
		sort.Strings(keys)
	default:
		return
	}
	for _, key := range keys {
		v := fields[key]
		switch number := v.(type) {
		case int, int64, bool:
			target.Set(key, v)
		case json.Number:
			if !strings.ContainsAny(string(number), ".eE") {
				target.Set(key, v)
			}
		}
	}
}
func (s *anthropicChatSink) FinishFailure(failure any) []string { return nil }

// closeUnclosedBlocks follows AnthropicStreamLedger's LIFO cleanup. Open tool
// blocks are deliberately left without a stop frame on failure.
func (s *anthropicChatSink) closeUnclosedBlocks(output *ChatStreamOutput) []string {
	events := output.FlushReasoningReplay(false)
	for len(s.openStack) > 0 {
		index := s.openStack[len(s.openStack)-1]
		s.openStack = s.openStack[:len(s.openStack)-1]
		kind := s.blockTypes[index]
		delete(s.blockTypes, index)
		if kind != "tool_use" {
			events = append(events, anthropicEvent("content_block_stop", ordered("type", "content_block_stop", "index", index)))
		}
	}
	output.textStarted = false
	clear(output.openReasoning)
	clear(s.reasoningIndices)
	for _, state := range output.toolStates {
		state.Open = false
	}
	return events
}
