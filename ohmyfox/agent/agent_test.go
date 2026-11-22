package agent

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestToolLoopAgainstChatEndpoint(t *testing.T) {
	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, "hello.txt"), []byte("hello workspace"), 0600); err != nil {
		t.Fatal(err)
	}
	var requests int
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requests++
		if r.Method != "POST" || r.URL.Path != "/v1/chat/completions" || r.Header.Get("Authorization") != "Bearer test-secret" {
			t.Errorf("unexpected request: %s %s", r.Method, r.URL.Path)
		}
		var body request
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Error(err)
			return
		}
		if body.Model != "provider/model" || body.Stream || len(body.Tools) != 2 {
			t.Errorf("unexpected request configuration")
		}
		w.Header().Set("Content-Type", "application/json")
		if requests == 1 {
			fmt.Fprint(w, `{"choices":[{"message":{"role":"assistant","content":null,"tool_calls":[{"id":"call-1","type":"function","function":{"name":"read_file","arguments":"{\"path\":\"hello.txt\"}"}}]}}]}`)
		} else {
			if len(body.Messages) != 4 || body.Messages[3].Role != "tool" || body.Messages[3].ToolCallID != "call-1" || body.Messages[3].Content == nil || *body.Messages[3].Content != "hello workspace" {
				t.Errorf("unexpected tool replay: %#v", body.Messages)
			}
			fmt.Fprint(w, `{"choices":[{"message":{"role":"assistant","content":"Read hello.txt."}}]}`)
		}
	}))
	defer srv.Close()
	fox, err := New(Config{Endpoint: srv.URL + "/v1/chat/completions", Model: "provider/model", APIKey: "test-secret", Workspace: root, MaxTurns: 3})
	if err != nil {
		t.Fatal(err)
	}
	result, err := fox.Run(context.Background(), "Inspect the project")
	if err != nil || result != "Read hello.txt." || requests != 2 {
		t.Fatalf("result=%q err=%v requests=%d", result, err, requests)
	}
}

func TestConfigAcceptsRemoteHTTPAndRejectsCredentialsInURL(t *testing.T) {
	root := t.TempDir()
	for _, endpoint := range []string{"https://fcc.example.trycloudflare.com/v1/chat/completions", "https://fcc.example.ngrok.app/v1/chat/completions", "http://fcc.example.ts.net/v1/chat/completions"} {
		if _, err := New(Config{Endpoint: endpoint, Model: "m", Workspace: root, MaxTurns: 2}); err != nil {
			t.Errorf("%s: %v", endpoint, err)
		}
	}
	for _, endpoint := range []string{"", "ftp://example.com", "https://user:secret@example.com/v1", "https://example.com/path?key=secret", "https://example.com/#fragment"} {
		if _, err := New(Config{Endpoint: endpoint, Model: "m", Workspace: root, MaxTurns: 2}); err == nil {
			t.Errorf("accepted %q", endpoint)
		}
	}
}

func TestNoRedirectOrCredentialLeak(t *testing.T) {
	seen := false
	target := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { seen = true; w.WriteHeader(http.StatusOK) }))
	defer target.Close()
	origin := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, target.URL, http.StatusTemporaryRedirect)
	}))
	defer origin.Close()
	fox, err := New(Config{Endpoint: origin.URL, Model: "m", APIKey: "secret", Workspace: t.TempDir(), MaxTurns: 2})
	if err != nil {
		t.Fatal(err)
	}
	_, err = fox.Run(context.Background(), "hello")
	if err == nil || !strings.Contains(err.Error(), "HTTP 307") || seen {
		t.Fatalf("err=%v redirected=%t", err, seen)
	}
}

func TestRejectInvalidOrOversizeResponsesAndLimitTurns(t *testing.T) {
	cases := []struct{ name, body, want string }{
		{"invalid", `{"choices":[]}`, "expected one Chat choice"},
		{"empty", `{"choices":[{"message":{"role":"assistant","content":" "}}]}`, "empty assistant message"},
		{"truncated", `{"choices":[{"finish_reason":"length","message":{"role":"assistant","content":"partial"}}]}`, "truncated"},
		{"filtered", `{"choices":[{"finish_reason":"content_filter","message":{"role":"assistant","content":"partial"}}]}`, "filtered"},
		{"unsupported", `{"choices":[{"finish_reason":"function_call","message":{"role":"assistant","content":"partial"}}]}`, "unsupported"},
		{"oversize", strings.Repeat("x", maxResponse+1), "size limit"},
		{"turns", `{"choices":[{"message":{"role":"assistant","content":null,"tool_calls":[{"id":"a","type":"function","function":{"name":"list_files","arguments":"{\"path\":\".\"}"}}]}}]}`, "maximum model turns"},
		{"duplicate", `{"choices":[{"message":{"role":"assistant","tool_calls":[{"id":"a","type":"function"},{"id":"a","type":"function"}]}}]}`, "invalid tool call"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { fmt.Fprint(w, tc.body) }))
			defer srv.Close()
			fox, err := New(Config{Endpoint: srv.URL, Model: "m", Workspace: t.TempDir(), MaxTurns: 1})
			if err != nil {
				t.Fatal(err)
			}
			_, err = fox.Run(context.Background(), "hi")
			if err == nil || !strings.Contains(err.Error(), tc.want) {
				t.Fatalf("err=%v want %q", err, tc.want)
			}
		})
	}
}

func TestCancelledContextAndPromptLimits(t *testing.T) {
	fox, err := New(Config{Endpoint: "https://example.com/v1/chat/completions", Model: "m", Workspace: t.TempDir(), MaxTurns: 2})
	if err != nil {
		t.Fatal(err)
	}
	for _, p := range []string{"", strings.Repeat("x", maxPrompt+1)} {
		if _, err := fox.Run(context.Background(), p); err == nil {
			t.Fatal("expected prompt error")
		}
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if _, err := fox.Run(ctx, "hello"); err == nil {
		t.Fatal("expected cancellation")
	}
}
