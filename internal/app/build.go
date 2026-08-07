package app

import (
	"github.com/brohd11/bubblestack/components"
	"github.com/brohd11/bubblestack/core"
	"github.com/brohd11/golaunch/internal/selection"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

// buildFlags is the Spec's four booleans as checklist rows. field hands back a pointer into a Spec
// so a row can read and flip its own flag — the alternative, a kind enum plus a switch in the
// screen, would put the same four cases in three places.
var buildFlags = []struct {
	label, desc string
	field       func(*selection.Spec) *bool
}{
	{"Include dirs", "gather directories under the root", func(s *selection.Spec) *bool { return &s.Dirs }},
	{"Include files", "gather files under the root", func(s *selection.Spec) *bool { return &s.Files }},
	{"Recursive", "descend into subdirectories", func(s *selection.Spec) *bool { return &s.Recursive }},
	{"Include current", "add the root directory itself", func(s *selection.Spec) *bool { return &s.Current }},
}

// buildRow is one build-flag row: a [x]/[ ] marker plus the flag's label, matching refineRow's
// look so the two checklists read as the same screen at different altitudes. idx maps the row back
// to its buildFlags entry; the Spec on the shared Selection stays the source of truth for on.
type buildRow struct {
	idx         int
	label, desc string
	on          bool
}

func (r buildRow) Title() string {
	mark := "[ ] "
	if r.on {
		mark = "[x] "
	}
	return mark + r.label
}

func (r buildRow) Description() string { return r.desc }
func (r buildRow) FilterValue() string { return r.label }

// BuildScreen is the checklist over the four build flags: enter toggles the highlighted flag and
// re-resolves the candidate paths immediately, so the header's selection summary moves with every
// keystroke and there is no separate build step. esc exits, keeping the result.
type BuildScreen struct {
	list list.Model
}

var (
	_ core.Filterer = (*BuildScreen)(nil)
	_ core.Crumber  = (*BuildScreen)(nil)
)

// pushBuild opens the Build checklist, resolving on the way in so the header count already matches
// the boxes on arrival. A resolve error is reported but doesn't block the push: an unreadable root
// is exactly when the flags need to be reachable to pick something that works.
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
	return &BuildScreen{list: core.NewSelectList(buildItems(spec), "Build selection")}
}

// buildItems builds the checklist rows from a spec (index-aligned with buildFlags).
func buildItems(spec selection.Spec) []list.Item {
	rows := make([]list.Item, len(buildFlags))
	for i, f := range buildFlags {
		rows[i] = buildRow{idx: i, label: f.label, desc: f.desc, on: *f.field(&spec)}
	}
	return rows
}

func (s *BuildScreen) Init(*core.Shared) tea.Cmd    { return nil }
func (s *BuildScreen) Filtering() bool              { return s.list.FilterState() == list.Filtering }
func (s *BuildScreen) View(*core.Shared) string     { return s.list.View() }
func (s *BuildScreen) HelpView(*core.Shared) string { return core.ShortHelp(s.list, core.HelpMinimal) }
func (s *BuildScreen) CrumbLabel(bool) string       { return "Build" }

func (s *BuildScreen) SetSize(_ *core.Shared, w, h int) { s.list.SetSize(w, h) }

func (s *BuildScreen) Update(sh *core.Shared, msg tea.Msg) (core.Screen, core.Action) {
	if m, ok := msg.(tea.MouseMsg); ok {
		if components.WheelNav(&s.list, m) {
			return s, core.Action{}
		}
	}
	// While actively typing a filter, every key belongs to the filter input.
	if s.Filtering() {
		var cmd tea.Cmd
		s.list, cmd = s.list.Update(msg)
		return s, core.Async(cmd)
	}
	if km, ok := msg.(tea.KeyMsg); ok {
		k := km.String()
		switch {
		case core.MatchKey(k, core.Keys.Select):
			// enter flips the highlighted flag and rebuilds the selection on the spot; the screen
			// stays open so the paths can be narrowed a flag at a time.
			return s, s.toggleSelected(sh)
		case core.MatchKey(k, core.Keys.Back):
			// esc exits, keeping whatever the last toggle resolved; confirm it on the status line.
			return s, core.Seq(core.SetStatus("selection: "+Of(sh).Sel.Summary()), core.Pop())
		default:
			if components.WrapNav(&s.list, k) {
				return s, core.Action{}
			}
		}
	}
	var cmd tea.Cmd
	s.list, cmd = s.list.Update(msg)
	return s, core.Async(cmd)
}

// toggleSelected flips the highlighted row's flag, re-resolves the paths under it, and rebuilds the
// rows so the [x]/[ ] marker updates, keeping the cursor where it was. A failed resolve leaves the
// spec untouched, so the rows still describe the selection that's actually loaded.
func (s *BuildScreen) toggleSelected(sh *core.Shared) core.Action {
	row, ok := s.list.SelectedItem().(buildRow)
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

	idx := s.list.Index()
	s.list.SetItems(buildItems(spec))
	s.list.Select(idx)

	// Recursive on its own gathers nothing — Any() ignores it. What was a submit-blocking error on
	// the old form is just a hint now: there's no submit left to block, and the header already
	// reads "none".
	if !spec.Any() {
		return core.SetStatus("nothing selected — enable dirs, files, or current")
	}
	return core.Action{}
}
