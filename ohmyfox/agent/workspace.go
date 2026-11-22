package agent

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"unicode/utf8"
)

const maxFileBytes = 64 << 10
const maxEntries = 200

type workspace struct{ root string }

func newWorkspace(path string) (*workspace, error) {
	if path == "" {
		return nil, errors.New("workspace is required")
	}
	absolute, err := filepath.Abs(path)
	if err != nil {
		return nil, errors.New("invalid workspace path")
	}
	root, err := filepath.EvalSymlinks(absolute)
	if err != nil {
		return nil, errors.New("cannot resolve workspace")
	}
	info, err := os.Stat(root)
	if err != nil || !info.IsDir() {
		return nil, errors.New("workspace must be an existing directory")
	}
	return &workspace{root: root}, nil
}

// resolve rejects traversal and symlinks. This is a best-effort guard for a
// trusted workspace, not an OS sandbox: a concurrent actor can race path checks.
func (w *workspace) resolve(path string) (string, error) {
	if len(path) == 0 || len(path) > 4096 || strings.ContainsRune(path, 0) || !filepath.IsLocal(path) {
		return "", errors.New("path must be relative to the workspace")
	}
	clean := filepath.Clean(path)
	current := w.root
	if clean == "." {
		return current, nil
	}
	for _, part := range strings.Split(clean, string(filepath.Separator)) {
		if strings.EqualFold(part, ".git") {
			return "", errors.New("git metadata is not available")
		}
		current = filepath.Join(current, part)
		info, err := os.Lstat(current)
		if err != nil {
			return "", errors.New("path is unavailable")
		}
		if info.Mode()&os.ModeSymlink != 0 {
			return "", errors.New("symlinks are not supported")
		}
	}
	return current, nil
}

func (w *workspace) call(name, raw string) string {
	var args struct {
		Path string `json:"path"`
	}
	decoder := json.NewDecoder(strings.NewReader(raw))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&args); err != nil {
		return "error: invalid tool arguments"
	}
	if err := decoder.Decode(new(any)); err != io.EOF {
		return "error: invalid tool arguments"
	}
	path, err := w.resolve(args.Path)
	if err != nil {
		return "error: " + err.Error()
	}
	switch name {
	case "read_file":
		result, err := readFile(path)
		if err != nil {
			return "error: " + err.Error()
		}
		return result
	case "list_files":
		result, err := listFiles(path)
		if err != nil {
			return "error: " + err.Error()
		}
		return result
	default:
		return "error: unknown tool"
	}
}

func readFile(path string) (string, error) {
	info, err := os.Lstat(path)
	if err != nil || !info.Mode().IsRegular() {
		return "", errors.New("path is not a regular file")
	}
	f, err := os.Open(path)
	if err != nil {
		return "", errors.New("cannot open file")
	}
	defer f.Close()
	data, err := io.ReadAll(io.LimitReader(f, maxFileBytes+1))
	if err != nil {
		return "", errors.New("cannot read file")
	}
	if len(data) > maxFileBytes {
		return "", errors.New("file exceeds 65536 bytes")
	}
	if !utf8.Valid(data) || bytes.IndexByte(data, 0) >= 0 {
		return "", errors.New("file is not UTF-8 text")
	}
	return string(data), nil
}

func listFiles(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", errors.New("cannot open directory")
	}
	defer f.Close()
	info, err := f.Stat()
	if err != nil || !info.IsDir() {
		return "", errors.New("path is not a directory")
	}
	entries, err := f.ReadDir(maxEntries + 1)
	if err != nil && err != io.EOF {
		return "", errors.New("cannot read directory")
	}
	truncated := len(entries) > maxEntries
	if truncated {
		entries = entries[:maxEntries]
	}
	type entry struct {
		Name string `json:"name"`
		Type string `json:"type"`
	}
	items := make([]entry, 0, len(entries))
	for _, e := range entries {
		if strings.EqualFold(e.Name(), ".git") {
			continue
		}
		kind := "file"
		if e.Type()&os.ModeSymlink != 0 {
			kind = "symlink"
		} else if e.IsDir() {
			kind = "directory"
		} else if !e.Type().IsRegular() {
			kind = "other"
		}
		items = append(items, entry{Name: e.Name(), Type: kind})
	}
	sort.Slice(items, func(i, j int) bool { return items[i].Name < items[j].Name })
	encoded, err := json.Marshal(struct {
		Entries   []entry `json:"entries"`
		Truncated bool    `json:"truncated"`
	}{items, truncated})
	if err != nil {
		return "", fmt.Errorf("cannot encode directory listing")
	}
	return string(encoded), nil
}
