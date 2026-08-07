package selection

import (
	"os"
	"path/filepath"
	"testing"
)

// setupTree builds: root/{sub1/,sub2/,a.txt,b.txt} with sub1/nested.txt one level down, so tests
// can distinguish immediate vs recursive gathering.
func setupTree(t *testing.T) string {
	t.Helper()
	root := t.TempDir()
	for _, d := range []string{"sub1", "sub2"} {
		if err := os.Mkdir(filepath.Join(root, d), 0o755); err != nil {
			t.Fatal(err)
		}
	}
	for _, f := range []string{"a.txt", "b.txt", filepath.Join("sub1", "nested.txt")} {
		if err := os.WriteFile(filepath.Join(root, f), nil, 0o644); err != nil {
			t.Fatal(err)
		}
	}
	return root
}

func paths(items []Item) []string {
	out := make([]string, len(items))
	for i, it := range items {
		out[i] = it.Path
	}
	return out
}

func TestResolveImmediateFiles(t *testing.T) {
	root := setupTree(t)
	items, err := Resolve(root, Spec{Files: true})
	if err != nil {
		t.Fatal(err)
	}
	if len(items) != 2 { // a.txt, b.txt — not sub1/nested.txt
		t.Fatalf("got %d, want 2: %v", len(items), paths(items))
	}
	for _, it := range items {
		if it.IsDir || !it.On {
			t.Errorf("file item wrong flags: %+v", it)
		}
	}
}

func TestResolveImmediateDirs(t *testing.T) {
	root := setupTree(t)
	items, err := Resolve(root, Spec{Dirs: true})
	if err != nil {
		t.Fatal(err)
	}
	if len(items) != 2 { // sub1, sub2
		t.Fatalf("got %d, want 2: %v", len(items), paths(items))
	}
}

func TestResolveRecursiveFiles(t *testing.T) {
	root := setupTree(t)
	items, err := Resolve(root, Spec{Files: true, Recursive: true})
	if err != nil {
		t.Fatal(err)
	}
	if len(items) != 3 { // a.txt, b.txt, sub1/nested.txt
		t.Fatalf("got %d, want 3: %v", len(items), paths(items))
	}
}

func TestResolveDirsAndFiles(t *testing.T) {
	root := setupTree(t)
	items, err := Resolve(root, Spec{Dirs: true, Files: true})
	if err != nil {
		t.Fatal(err)
	}
	if len(items) != 4 { // sub1, sub2, a.txt, b.txt
		t.Fatalf("got %d, want 4: %v", len(items), paths(items))
	}
}

func TestResolveCurrentAddsRoot(t *testing.T) {
	root := setupTree(t)
	items, err := Resolve(root, Spec{Current: true})
	if err != nil {
		t.Fatal(err)
	}
	if len(items) != 1 || items[0].Path != root || !items[0].IsDir {
		t.Fatalf("got %v, want [%s] (dir)", paths(items), root)
	}
}

func TestResolveMissingDir(t *testing.T) {
	if _, err := Resolve(filepath.Join(t.TempDir(), "nope"), Spec{Files: true}); err == nil {
		t.Error("expected error for missing directory")
	}
}

// TestRebuildPreservesOffFlags is the contract the Build checklist leans on: it re-resolves on
// every toggle, so a refinement made in the Refine checklist has to survive a flag being flipped
// on and back off.
func TestRebuildPreservesOffFlags(t *testing.T) {
	root := setupTree(t)
	items, err := Resolve(root, Spec{Files: true})
	if err != nil {
		t.Fatal(err)
	}
	sel := Selection{Spec: Spec{Files: true}, Items: items}
	off := filepath.Join(root, "b.txt")
	for i := range sel.Items {
		if sel.Items[i].Path == off {
			sel.Items[i].On = false
		}
	}

	// Recursive on: b.txt survives and must still be off, while the newly reached nested.txt
	// arrives enabled.
	sel, err = sel.Rebuild(root, Spec{Files: true, Recursive: true})
	if err != nil {
		t.Fatal(err)
	}
	if len(sel.Items) != 3 {
		t.Fatalf("got %d items, want 3: %v", len(sel.Items), paths(sel.Items))
	}
	for _, it := range sel.Items {
		want := it.Path != off
		if it.On != want {
			t.Errorf("%s: On = %v, want %v", it.Path, it.On, want)
		}
	}
	if sel.Spec != (Spec{Files: true, Recursive: true}) {
		t.Errorf("Rebuild should store the new spec, got %+v", sel.Spec)
	}

	// ...and back off: b.txt is still there, still disabled.
	sel, err = sel.Rebuild(root, Spec{Files: true})
	if err != nil {
		t.Fatal(err)
	}
	if got := len(sel.Paths()); got != 1 {
		t.Errorf("after flipping Recursive on and off, Paths len = %d, want 1: %v", got, sel.Paths())
	}
}

func TestRebuildFromEmpty(t *testing.T) {
	root := setupTree(t)
	sel, err := Selection{}.Rebuild(root, Spec{Files: true})
	if err != nil {
		t.Fatal(err)
	}
	if got := len(sel.Paths()); got != 2 {
		t.Errorf("a fresh rebuild should enable everything it resolves, got %d of %d", got, len(sel.Items))
	}
}

// TestRebuildKeepsSelectionOnError pins that a failed resolve leaves the caller with what it had —
// the Build checklist keeps showing a loaded selection rather than blanking it.
func TestRebuildKeepsSelectionOnError(t *testing.T) {
	root := setupTree(t)
	sel, err := Selection{}.Rebuild(root, Spec{Files: true})
	if err != nil {
		t.Fatal(err)
	}
	got, err := sel.Rebuild(filepath.Join(root, "nope"), Spec{Dirs: true})
	if err == nil {
		t.Fatal("expected an error rebuilding under a missing directory")
	}
	if len(got.Items) != len(sel.Items) || got.Spec != sel.Spec {
		t.Errorf("a failed rebuild should return the receiver unchanged, got %+v", got)
	}
}

func TestPathsAndSummary(t *testing.T) {
	root := setupTree(t)
	items, _ := Resolve(root, Spec{Files: true, Dirs: true})
	sel := Selection{Items: items}
	if got := len(sel.Paths()); got != 4 {
		t.Fatalf("all enabled: Paths len = %d, want 4", got)
	}
	if sel.Summary() != "4 of 4 paths" {
		t.Errorf("Summary = %q, want 4 of 4 paths", sel.Summary())
	}
	sel.Items[0].On = false
	if got := len(sel.Paths()); got != 3 {
		t.Errorf("one disabled: Paths len = %d, want 3", got)
	}
	if sel.Summary() != "3 of 4 paths" {
		t.Errorf("Summary = %q, want 3 of 4 paths", sel.Summary())
	}
}
