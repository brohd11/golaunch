package app

import (
	"github.com/brohd11/bubblestack/core"
	"github.com/brohd11/golaunch/internal/selection"

	"charm.land/bubbles/v2/list"
	tea "charm.land/bubbletea/v2"
)

// buildFlags are the Spec's booleans as checklist rows; field points into a Spec.
var buildFlags = []struct {
	label, desc string
	field       func(*selection.Spec) *bool
}{
	{"Include dirs", "gather directories under the root", func(s *selection.Spec) *bool { return &s.Dirs }},
	{"Include files", "gather files under the root", func(s *selection.Spec) *bool { return &s.Files }},
	{"Recursive", "descend into subdirectories", func(s *selection.Spec) *bool { return &s.Recursive }},
	{"Include current", "add the root directory itself", func(s *selection.Spec) *bool { return &s.Current }},
}

// BuildScreen is the checklist over the build flags. Each toggle re-resolves the paths
// immediately; esc exits, keeping the result.
type BuildScreen struct {
	checklist
}

var (
	_ core.Filterer = (*BuildScreen)(nil)
	_ core.Crumber  = (*BuildScreen)(nil)
)

// pushBuild resolves on the way in so the header count matches. A resolve error doesn't
// block the push, so the flags stay reachable.
func pushBuild(sh *core.Shared) core.Action {
	c := Of(sh)
	spec := c.Sel.Spec
	// A never-built selection defaults to "files here" — the common case — so the checklist isn't
	// all-empty on first open.
	if !spec.Any() && !spec.Recursive {
		spec.Files = true
	}
	sel, err := c.Sel.Rebuild(c.Root, spec)
	if err != nil {
		// The rows still describe the spec that was asked for, so the failing flag is visible and
		// can be flipped back off; c.Sel keeps whatever it already had.
		return core.Seq(core.StatusErr(err), core.Push(newBuildScreen(spec)))
	}
	c.Sel = sel
	return core.Push(newBuildScreen(c.Sel.Spec))
}

func newBuildScreen(spec selection.Spec) *BuildScreen {
	return &BuildScreen{checklist{list: core.NewSelectList(buildItems(spec), "Build selection"), crumb: "Build"}}
}

// buildItems builds the checklist rows from a spec (index-aligned with buildFlags). There is
// nothing but the flag's label to filter on here, so that is what the row matches.
func buildItems(spec selection.Spec) []list.Item {
	rows := make([]list.Item, len(buildFlags))
	for i, f := range buildFlags {
		rows[i] = checkRow{idx: i, label: f.label, desc: f.desc, filter: f.label, on: *f.field(&spec)}
	}
	return rows
}

func (s *BuildScreen) Update(sh *core.Shared, msg tea.Msg) (core.Screen, core.Action) {
	return s, checklistUpdate(&s.list, sh, msg, s.toggleSelected)
}

// toggleSelected flips the row's flag and re-resolves. A failed resolve leaves the spec
// untouched.
func (s *BuildScreen) toggleSelected(sh *core.Shared) core.Action {
	row, ok := s.list.SelectedItem().(checkRow)
	if !ok {
		return core.Action{}
	}
	c := Of(sh)
	spec := c.Sel.Spec
	f := buildFlags[row.idx].field(&spec)
	*f = !*f

	sel, err := c.Sel.Rebuild(c.Root, spec)
	if err != nil {
		return core.StatusErr(err)
	}
	c.Sel = sel

	setRowsKeepCursor(&s.list, buildItems(spec))

	// Recursive alone gathers nothing; the header already reads "none".
	if !spec.Any() {
		return core.SetStatus("nothing selected — enable dirs, files, or current")
	}
	return core.Action{}
}
