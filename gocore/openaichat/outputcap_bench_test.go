package openaichat

import "testing"

var benchmarkOutputCap *int

func BenchmarkParseOutputTokenCap(b *testing.B) {
	for _, tc := range []struct{ name, text string }{
		{"standard", "max_completion_tokens must be less than or equal to 40960"},
		{"unicode", "max_tokens\u00a0at\u00a0most\u00a0\uff11\uff12"},
		{"no_field", "request failed without an output token limit"},
	} {
		b.Run(tc.name, func(b *testing.B) {
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				benchmarkOutputCap = parseCap(tc.text)
			}
		})
	}
}
