// Package selection models the paths a script runs against: a Spec resolves to Items, and
// each Item's On flag is toggled in Refine. Paths returns the enabled subset.
package selection

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
)

// Spec is the Build flags: Dirs and Files are independent; Current adds the root itself.
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

// FromPaths builds a selection from argv paths. Invalid paths are skipped and returned as
// errors.
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

// Resolve gathers the candidate paths under root, all enabled, sorted and absolute. Only a
// read error on the root is returned.
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

// Rebuild resolves spec and carries over the On flags of surviving paths. On error the
// receiver is unchanged.
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

// walkDescendants appends matching descendants of root, skipping unreadable subtrees.
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
