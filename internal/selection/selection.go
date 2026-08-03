// Package selection models the set of paths a script runs against. A Selection is one Mode (the
// whole root directory, its immediate child directories, or its immediate child files) resolved to
// a concrete []string. This is the first-steps shape: choosing a mode replaces the selection
// wholesale — per-item checklists and saved selections come later.
package selection

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
)

// Mode is which set of paths under the root a Selection covers.
type Mode int

const (
	None       Mode = iota // nothing selected yet
	CurrentDir             // the root directory itself (one path)
	ChildDirs              // the root's immediate subdirectories
	ChildFiles             // the root's immediate files
)

// Label is the human name for a mode, shown in the Selection list and header.
func (m Mode) Label() string {
	switch m {
	case CurrentDir:
		return "Current dir"
	case ChildDirs:
		return "Child dirs"
	case ChildFiles:
		return "Child files"
	default:
		return "none"
	}
}

// Selection is a resolved set of paths plus the mode that produced it. The zero value is an empty
// selection (Mode None), which the Scripts tab treats as "nothing to run against".
type Selection struct {
	Mode  Mode
	Paths []string
}

// Resolve reads root and returns the paths for mode: the root itself for CurrentDir, or its
// immediate child dirs/files (sorted, absolute) for the child modes. A read error is returned so
// the caller can surface it; hidden entries (dot-prefixed) are kept — the user chose this dir.
func Resolve(root string, mode Mode) (Selection, error) {
	if mode == CurrentDir {
		return Selection{Mode: mode, Paths: []string{root}}, nil
	}
	if mode == None {
		return Selection{}, nil
	}
	entries, err := os.ReadDir(root)
	if err != nil {
		return Selection{}, err
	}
	wantDir := mode == ChildDirs
	var paths []string
	for _, e := range entries {
		if e.IsDir() == wantDir {
			paths = append(paths, filepath.Join(root, e.Name()))
		}
	}
	sort.Strings(paths)
	return Selection{Mode: mode, Paths: paths}, nil
}

// Summary is a one-line description of the selection for the header, e.g. "Child files (5)". An
// empty selection reads as "none".
func (s Selection) Summary() string {
	if s.Mode == None {
		return "none"
	}
	return fmt.Sprintf("%s (%d)", s.Mode.Label(), len(s.Paths))
}

// Empty reports whether the selection holds no paths.
func (s Selection) Empty() bool { return len(s.Paths) == 0 }
