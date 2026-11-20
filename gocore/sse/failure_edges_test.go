package sse

import "testing"

func TestFailureMalformedNamedEvent(t *testing.T) {
	failures := map[string]struct{}{"error": {}, "response.failed": {}}
	for _, tc := range []struct{ frame, want string }{
		{"event: error\ndata: {bad json}\n\n", "error"},
		{"event: response.failed\ndata: {\"type\":\"ok\"} trailing\n\n", "response.failed"},
		{"data: {\"error\":\"broken\"} trailing\n\n", ""},
		{"data: {\"type\":42,\"error\":{\"type\":\"limit\"}}\n\n", ""},
		{"data: {\"type\":\"\",\"error\":{\"code\":42,\"type\":\"limit\"}}\n\n", "error"},
		{"data: {\"error\":{\"code\":\"\",\"type\":\"limit\"}}\n\n", "limit"},
		{"event: error\ndata: {\"error\":{\"type\":\"\"}}\n\n", ""},
	} {
		if got := FailureReason(tc.frame, failures); got != tc.want {
			t.Errorf("FailureReason(%q) = %q, want %q", tc.frame, got, tc.want)
		}
	}
}
