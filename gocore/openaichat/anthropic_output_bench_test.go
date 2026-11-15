package openaichat

import "testing"

var anthropicFrameSink string

func BenchmarkAnthropicChatTextStream(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		output := NewAnthropicChatStreamOutput("msg", "model", 20)
		output.StartEvents()
		output.EnsureTextBlock()
		for j := 0; j < 100; j++ {
			anthropicFrameSink = output.EmitTextDelta("example output fragment")
		}
		output.FinishSuccess("end_turn", ChatStreamUsage{InputTokens: 20, OutputTokens: 100})
	}
}
