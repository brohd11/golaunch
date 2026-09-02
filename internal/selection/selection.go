// Package selection models the set of paths a script runs against. A Selection is built in two
// stages: a Spec (which kinds of path under the root to gather) resolves to a candidate list of
// Items, and each Item then carries an On flag the Refine checklist toggles. The enabled subset
// (Paths) is what actually reaches a script.
package selection

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
)

// Spec is the set of build flags chosen in the Build form: which kinds of path to gather from the
// root, and whether to descend recursively. Dirs and Files are independent (either or both);
// Current adds the root directory itself.
type Spec struct {
	Dirs      bool
	Files     bool
	Recursive bool
	Current   bool
}

// Any reports whether the spec would gather anything at all.
func (s Spec) Any() bool { return s.Dirs || s.Files || s.Current }

// Item is one candidate path plus whether it is currently enabled (toggled in the Refine
// checklist). IsDir drives the checklist's trailing-"/" marker.
type Item struct {
	Path  string
	IsDir bool
	On    bool
}

// Selection is the resolved candidate set (Items) plus the Spec that produced it. The zero value
// is empty. Paths returns the enabled subset — the list handed to a script.
type Selection struct {
	Spec  Spec
	Items []Item
}

// FromPaths builds a selection from explicit command-line paths. Valid paths are made absolute,
// classified as files or directories, and enabled in argument order. Invalid paths are omitted and
// returned as individual problems so a file-manager launch can continue with the rest.
func FromPaths(paths []string) (Selection, []error) {
	items := make([]Item, 0, len(paths))
	var problems []error
	for _, path := range paths {
		abs, err := filepath.Abs(path)
		if err != nil {
			problems = append(problems, fmt.Errorf("%s: %w", path, err))
			continue
		}
		info, err := os.Stat(abs)
		if err != nil {
			problems = append(problems, fmt.Errorf("%s: %w", path, err))
			continue
		}
		items = append(items, Item{Path: abs, IsDir: info.IsDir(), On: true})
	}
	return Selection{Items: items}, problems
}

// Resolve gathers the candidate paths under root for spec, each enabled by default. Current adds
// the root itself; otherwise the root's immediate children are read, or every descendant when
// Recursive. Dirs keeps directories, Files keeps files. Results are sorted and absolute. A read
// error on the root is returned; errors walking individual descendants are skipped so one
// unreadable subdir doesn't abort the whole build.
func Resolve(root string, spec Spec) ([]Item, error) {
	var items []Item
	if spec.Current {
		items = append(items, Item{Path: root, IsDir: true, On: true})
	}

	if spec.Dirs || spec.Files {
		var err error
		if spec.Recursive {
			err = walkDescendants(root, spec, &items)
		} else {
			err = readImmediate(root, spec, &items)
		}
		if err != nil {
			return nil, err
		}
	}

	sort.Slice(items, func(i, j int) bool { return items[i].Path < items[j].Path })
	return items, nil
}

// Rebuild resolves spec under root and returns the replacement Selection, carrying the receiver's
// per-path On flags onto any path that survives — so flipping Recursive on and back off doesn't
// silently undo a refinement. Paths new to the selection keep Resolve's default of enabled. On
// error the receiver is left untouched (the caller keeps what it had).
func (s Selection) Rebuild(root string, spec Spec) (Selection, error) {
	items, err := Resolve(root, spec)
	if err != nil {
		return s, err
	}
	return Selection{Spec: spec, Items: carryFlags(s.Items, items)}, nil
}

// carryFlags re-applies the disabled paths of prev to next. Only the *off* set is carried, which is
// what makes a path new to the selection default to on for free — Resolve already enabled it.
func carryFlags(prev, next []Item) []Item {
	off := make(map[string]bool)
	for _, it := range prev {
		if !it.On {
			off[it.Path] = true
		}
	}
	if len(off) == 0 {
		return next
	}
	for i := range next {
		if off[next[i].Path] {
			next[i].On = false
		}
	}
	return next
}

// readImmediate appends the root's immediate children matching spec.
func readImmediate(root string, spec Spec, items *[]Item) error {
	entries, err := os.ReadDir(root)
	if err != nil {
		return err
	}
	for _, e := range entries {
		if keep(e.IsDir(), spec) {
			*items = append(*items, Item{Path: filepath.Join(root, e.Name()), IsDir: e.IsDir(), On: true})
		}
	}
	return nil
}

// walkDescendants appends every descendant of root matching spec (the root itself is excluded —
// Current handles it). A per-entry error is swallowed so an unreadable subtree is skipped rather
// than failing the whole build.
func walkDescendants(root string, spec Spec, items *[]Item) error {
	return filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return nil // skip unreadable entries
		}
		if path == root {
			return nil
		}
		if keep(d.IsDir(), spec) {
			*items = append(*items, Item{Path: path, IsDir: d.IsDir(), On: true})
		}
		return nil
	})
}

// keep reports whether an entry of the given kind is wanted by spec.
func keep(isDir bool, spec Spec) bool {
	if isDir {
		return spec.Dirs
	}
	return spec.Files
}

// Paths returns the enabled paths — the final list a script receives.
func (s Selection) Paths() []string {
	var out []string
	for _, it := range s.Items {
		if it.On {
			out = append(out, it.Path)
		}
	}
	return out
}

// Summary is a one-line description for the header: enabled count over candidate count, or "none".
func (s Selection) Summary() string {
	if len(s.Items) == 0 {
		return "none"
	}
	return fmt.Sprintf("%d of %d paths", len(s.Paths()), len(s.Items))
}

// Empty reports whether there are no candidate paths at all.
func (s Selection) Empty() bool { return len(s.Items) == 0 }
