package cmd

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestResolveSelectionNoArgsUsesBuildMode(t *testing.T) {
	var stderr bytes.Buffer
	sel, preselected, err := resolveSelection(nil, &stderr)
	if err != nil {
		t.Fatal(err)
	}
	if preselected || !sel.Empty() {
		t.Fatalf("no args should use an empty build-mode selection, got preselected=%v sel=%+v", preselected, sel)
	}
	if stderr.Len() != 0 {
		t.Fatalf("unexpected stderr: %q", stderr.String())
	}
}

func TestResolveSelectionSkipsBadPaths(t *testing.T) {
	root := t.TempDir()
	valid := filepath.Join(root, "valid.txt")
	if err := os.WriteFile(valid, nil, 0o644); err != nil {
		t.Fatal(err)
	}
	missing := filepath.Join(root, "missing.txt")

	var stderr bytes.Buffer
	sel, preselected, err := resolveSelection([]string{missing, valid}, &stderr)
	if err != nil {
		t.Fatal(err)
	}
	if !preselected || len(sel.Items) != 1 || sel.Items[0].Path != valid {
		t.Fatalf("got preselected=%v sel=%+v", preselected, sel)
	}
	if got := stderr.String(); !strings.Contains(got, "skipping: "+missing) {
		t.Fatalf("stderr = %q, want skipped-path warning", got)
	}
}

func TestResolveSelectionRejectsAllBadPaths(t *testing.T) {
	missing := filepath.Join(t.TempDir(), "missing.txt")
	var stderr bytes.Buffer
	sel, preselected, err := resolveSelection([]string{missing}, &stderr)
	if err == nil || err.Error() != "no valid paths selected" {
		t.Fatalf("err = %v, want no valid paths selected", err)
	}
	if !preselected || !sel.Empty() {
		t.Fatalf("got preselected=%v sel=%+v", preselected, sel)
	}
	if !strings.Contains(stderr.String(), "skipping: "+missing) {
		t.Fatalf("stderr = %q, want skipped-path warning", stderr.String())
	}
}
