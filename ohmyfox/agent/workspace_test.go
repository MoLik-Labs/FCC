package agent

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestWorkspaceReadListAndTraversal(t *testing.T) {
	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, "note.txt"), []byte("text"), 0600); err != nil {
		t.Fatal(err)
	}
	if err := os.Mkdir(filepath.Join(root, "sub"), 0700); err != nil {
		t.Fatal(err)
	}
	if err := os.Mkdir(filepath.Join(root, ".git"), 0700); err != nil {
		t.Fatal(err)
	}
	w, err := newWorkspace(root)
	if err != nil {
		t.Fatal(err)
	}
	if got := w.call("read_file", `{"path":"note.txt"}`); got != "text" {
		t.Fatalf("read=%q", got)
	}
	listing := w.call("list_files", `{"path":"."}`)
	var decoded struct {
		Entries []struct {
			Name string `json:"name"`
			Type string `json:"type"`
		} `json:"entries"`
		Truncated bool `json:"truncated"`
	}
	if err := json.Unmarshal([]byte(listing), &decoded); err != nil || len(decoded.Entries) != 2 || decoded.Entries[0].Name != "note.txt" || decoded.Entries[1].Name != "sub" {
		t.Fatalf("listing=%s err=%v", listing, err)
	}
	for _, input := range []string{`{"path":"../secret"}`, `{"path":"/etc/passwd"}`, `{"path":".git"}`, `{"path":"note.txt","extra":true}`, `{"path":"note.txt"} {}`, `{}`, `null`} {
		if got := w.call("read_file", input); !strings.HasPrefix(got, "error:") {
			t.Fatalf("input %s returned %q", input, got)
		}
	}
	if got := w.call("execute", `{"path":"note.txt"}`); got != "error: unknown tool" {
		t.Fatalf("unknown tool=%q", got)
	}
}

func TestWorkspaceRejectsSymlinksBinaryAndLargeFiles(t *testing.T) {
	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, "binary"), []byte{0, 1, 2}, 0600); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "large"), []byte(strings.Repeat("a", maxFileBytes+1)), 0600); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(t.TempDir(), filepath.Join(root, "link")); err != nil {
		t.Skipf("symlinks unavailable: %v", err)
	}
	w, err := newWorkspace(root)
	if err != nil {
		t.Fatal(err)
	}
	for _, path := range []string{"binary", "large", "link", "link/file"} {
		args, _ := json.Marshal(map[string]string{"path": path})
		if got := w.call("read_file", string(args)); !strings.HasPrefix(got, "error:") {
			t.Fatalf("path %s returned %q", path, got)
		}
	}
}
