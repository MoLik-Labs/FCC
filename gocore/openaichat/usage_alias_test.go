package openaichat

import "testing"

func TestUsageRequestAndRetryOwnership(t *testing.T) {
	options := map[string]any{"other": "kept", "include_usage": false}
	body := map[string]any{"stream_options": options, "model": "m"}
	RequestStreamUsage(body)
	if options["include_usage"] != true {
		t.Fatal("request did not update the original options map")
	}
	retry := CloneWithoutStreamUsage(body)
	if retry == nil {
		t.Fatal("expected retry body")
	}
	retryOptions, ok := retry["stream_options"].(map[string]any)
	if !ok {
		t.Fatal("expected retry options")
	}
	if _, exists := retryOptions["include_usage"]; exists {
		t.Fatal("include_usage retained in retry")
	}
	retryOptions["other"] = "changed"
	retry["model"] = "other"
	if options["other"] != "kept" || body["model"] != "m" || options["include_usage"] != true {
		t.Fatal("mutating retry changed original request")
	}
}
