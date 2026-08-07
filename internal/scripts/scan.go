package scripts

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// Script is one launchable script: its absolute file path, the interpreter argv to run it with
// (empty when the file is executed directly), and its parsed header metadata.
type Script struct {
	File   string
	Interp []string // leading argv, e.g. ["python3"]; empty ⇒ run File directly
	Meta   Meta
}

// DisplayName is the menu label: the metadata name, or the filename without its extension.
func (s Script) DisplayName() string {
	if s.Meta.Name != "" {
		return s.Meta.Name
	}
	return strings.TrimSuffix(filepath.Base(s.File), filepath.Ext(s.File))
}

// interpByExt maps a script extension to the interpreter argv that runs it. Extend it (or add an
// interp= metadata key) to support more languages later.
var interpByExt = map[string][]string{
	".py": {"python3"},
	".sh": {"bash"},
}

// Scan reads each directory (non-recursively) for launchable scripts: files with a known
// extension, or any executable file (run directly). Results are deduplicated by absolute
// path and sorted by display name. What couldn't be read — an unreadable directory, a
// script whose header fails to parse — comes back as one error per problem rather than
// failing the scan or vanishing silently: a typo in one script's header should cost that
// script a menu row, not its invisibility without a trace.
func Scan(dirs []string) ([]Script, []error) {
	seen := map[string]bool{}
	var out []Script
	var problems []error
	for _, dir := range dirs {
		entries, err := os.ReadDir(dir)
		if err != nil {
			problems = append(problems, fmt.Errorf("reading %s: %w", dir, err))
			continue
		}
		for _, e := range entries {
			if e.IsDir() {
				continue
			}
			abs, err := filepath.Abs(filepath.Join(dir, e.Name()))
			if err != nil || seen[abs] {
				continue
			}
			interp, ok := interpFor(abs, e)
			if !ok {
				continue
			}
			meta, err := ParseHeader(abs)
			if err != nil {
				problems = append(problems, fmt.Errorf("%s: %w", abs, err))
				continue
			}
			seen[abs] = true
			out = append(out, Script{File: abs, Interp: interp, Meta: meta})
		}
	}
	sort.Slice(out, func(i, j int) bool { return out[i].DisplayName() < out[j].DisplayName() })
	return out, problems
}

// interpFor returns the interpreter argv for a file and whether it counts as a script: a known
// extension wins; otherwise any file with an execute bit is run directly (nil interp).
func interpFor(path string, e os.DirEntry) ([]string, bool) {
	if in, ok := interpByExt[strings.ToLower(filepath.Ext(path))]; ok {
		return in, true
	}
	if info, err := e.Info(); err == nil && info.Mode()&0o111 != 0 {
		return nil, true
	}
	return nil, false
}

// Tree is a node in the script menu: named subgroups (from the "/"-separated path= metadata) and
// the scripts that live directly at this node. The root node has an empty Name.
type Tree struct {
	Name     string
	Children []*Tree
	Scripts  []Script
}

// BuildTree groups scripts into a menu tree by their Path metadata: "Image/Filters" nests the
// script two levels deep, an empty path lands it at the root. Children and scripts are sorted for
// a stable menu.
func BuildTree(scripts []Script) *Tree {
	root := &Tree{}
	for _, s := range scripts {
		node := root
		if s.Meta.Path != "" {
			for _, seg := range strings.Split(s.Meta.Path, "/") {
				seg = strings.TrimSpace(seg)
				if seg == "" {
					continue
				}
				node = node.child(seg)
			}
		}
		node.Scripts = append(node.Scripts, s)
	}
	root.sort()
	return root
}

// child returns the named subgroup of a node, creating it on first use.
func (t *Tree) child(name string) *Tree {
	for _, c := range t.Children {
		if c.Name == name {
			return c
		}
	}
	c := &Tree{Name: name}
	t.Children = append(t.Children, c)
	return c
}

// sort orders a node's children by name and its scripts by display name, recursively.
func (t *Tree) sort() {
	sort.Slice(t.Children, func(i, j int) bool { return t.Children[i].Name < t.Children[j].Name })
	sort.Slice(t.Scripts, func(i, j int) bool { return t.Scripts[i].DisplayName() < t.Scripts[j].DisplayName() })
	for _, c := range t.Children {
		c.sort()
	}
}
