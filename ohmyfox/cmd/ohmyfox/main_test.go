package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"strings"
	"testing"
)

func TestCLIWithPromptOnStdin(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != "Bearer test-token" {
			t.Error("missing bearer token")
		}
		fmt.Fprint(w, `{"choices":[{"message":{"role":"assistant","content":"Hello from the test endpoint."}}]}`)
	}))
	defer server.Close()
	cmd := exec.Command("go", "run", ".", "--endpoint", server.URL, "--model", "test", "--workspace", t.TempDir(), "--prompt-stdin")
	cmd.Env = append(os.Environ(), "OMF_API_KEY=test-token")
	cmd.Stdin = strings.NewReader("Explain the project")
	output, err := cmd.CombinedOutput()
	if err != nil || strings.TrimSpace(string(output)) != "Hello from the test endpoint." {
		t.Fatalf("CLI: output=%q err=%v", output, err)
	}
}
