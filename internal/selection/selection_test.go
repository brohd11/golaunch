package selection

import (
	"os"
	"path/filepath"
	"testing"
)

func setupTree(t *testing.T) string {
	t.Helper()
	root := t.TempDir()
	for _, d := range []string{"sub1", "sub2"} {
		if err := os.Mkdir(filepath.Join(root, d), 0o755); err != nil {
			t.Fatal(err)
		}
	}
	for _, f := range []string{"a.txt", "b.txt", "c.txt"} {
		if err := os.WriteFile(filepath.Join(root, f), nil, 0o644); err != nil {
			t.Fatal(err)
		}
	}
	return root
}

func TestResolveCurrentDir(t *testing.T) {
	root := setupTree(t)
	sel, err := Resolve(root, CurrentDir)
	if err != nil {
		t.Fatal(err)
	}
	if len(sel.Paths) != 1 || sel.Paths[0] != root {
		t.Errorf("Paths = %v, want [%s]", sel.Paths, root)
	}
}

func TestResolveChildDirs(t *testing.T) {
	root := setupTree(t)
	sel, err := Resolve(root, ChildDirs)
	if err != nil {
		t.Fatal(err)
	}
	if len(sel.Paths) != 2 {
		t.Fatalf("got %d dirs, want 2: %v", len(sel.Paths), sel.Paths)
	}
	if sel.Paths[0] != filepath.Join(root, "sub1") {
		t.Errorf("Paths not sorted/absolute: %v", sel.Paths)
	}
}

func TestResolveChildFiles(t *testing.T) {
	root := setupTree(t)
	sel, err := Resolve(root, ChildFiles)
	if err != nil {
		t.Fatal(err)
	}
	if len(sel.Paths) != 3 {
		t.Errorf("got %d files, want 3: %v", len(sel.Paths), sel.Paths)
	}
	if sel.Summary() != "Child files (3)" {
		t.Errorf("Summary = %q, want Child files (3)", sel.Summary())
	}
}

func TestResolveMissingDir(t *testing.T) {
	if _, err := Resolve(filepath.Join(t.TempDir(), "nope"), ChildFiles); err == nil {
		t.Error("expected error for missing directory")
	}
}
