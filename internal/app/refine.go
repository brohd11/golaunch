package app

import (
	"path/filepath"

	"github.com/brohd11/bubblestack/components"
	"github.com/brohd11/bubblestack/core"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

// refineKey (shift+R) opens the Refine checklist from the menus. Matched per-screen (the Selection
// and Scripts roots, and the script submenu pickers) so a quick refinement is a keystroke away
// without leaving where you are. Distinct from Refresh ("r") and jump-to-bottom ("G").
var refineKey = key.NewBinding(key.WithKeys("R"), key.WithHelp("R", "refine"))

// pushRefine opens the Refine checklist over the current stack, or reports that there is nothing to
// refine when no selection has been built yet.
func pushRefine(sh *core.Shared) core.Action {
	if Of(sh).Sel.Empty() {
		return core.SetStatus("nothing to refine — build a selection first")
	}
	return core.Push(NewRefineScreen(sh))
}

// refineRow is one checklist row: a path with a [x]/[ ] marker (dirs suffixed "/"). idx maps the
// row back to its Selection.Items entry, which stays the source of truth for the On state.
type refineRow struct {
	idx        int
	path, name string
	isDir, on  bool
}

func (r refineRow) Title() string {
	mark := "[ ] "
	if r.on {
		mark = "[x] "
	}
	name := r.name
	if r.isDir {
		name += "/"
	}
	return mark + name
}

func (r refineRow) Description() string { return r.path }
func (r refineRow) FilterValue() string { return r.path }

// RefineScreen is a scrollable, filterable checklist over the built selection: enter toggles the
// highlighted path's inclusion, applied immediately to the shared Selection; esc exits, keeping the
// edits. There is no separate apply/cancel step — every toggle is already live.
type RefineScreen struct {
	list list.Model
}

var (
	_ core.Filterer = (*RefineScreen)(nil)
	_ core.Crumber  = (*RefineScreen)(nil)
)

func NewRefineScreen(sh *core.Shared) *RefineScreen {
	return &RefineScreen{list: core.NewSelectList(refineItems(sh), "Refine selection")}
}

// refineItems builds the checklist rows from the current Selection.Items (index-aligned).
func refineItems(sh *core.Shared) []list.Item {
	sel := Of(sh).Sel.Items
	rows := make([]list.Item, len(sel))
	for i, it := range sel {
		rows[i] = refineRow{idx: i, path: it.Path, name: filepath.Base(it.Path), isDir: it.IsDir, on: it.On}
	}
	return rows
}

func (s *RefineScreen) Init(*core.Shared) tea.Cmd    { return nil }
func (s *RefineScreen) Filtering() bool              { return s.list.FilterState() == list.Filtering }
func (s *RefineScreen) View(*core.Shared) string     { return s.list.View() }
func (s *RefineScreen) HelpView(*core.Shared) string { return core.ShortHelp(s.list, core.HelpMinimal) }
func (s *RefineScreen) CrumbLabel(bool) string       { return "Refine" }

func (s *RefineScreen) SetSize(_ *core.Shared, w, h int) { s.list.SetSize(w, h) }

func (s *RefineScreen) Update(sh *core.Shared, msg tea.Msg) (core.Screen, core.Action) {
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
			// enter toggles the highlighted path's inclusion, applied live; the screen stays open
			// so several rows can be flipped in a row.
			s.toggleSelected(sh)
			return s, core.Action{}
		case core.MatchKey(k, core.Keys.Back):
			// esc exits, keeping the live edits; confirm the result on the status line.
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

// toggleSelected flips the highlighted row's On state on the shared selection and rebuilds the rows
// so the [x]/[ ] marker updates, keeping the cursor where it was.
func (s *RefineScreen) toggleSelected(sh *core.Shared) {
	row, ok := s.list.SelectedItem().(refineRow)
	if !ok {
		return
	}
	c := Of(sh)
	if row.idx < 0 || row.idx >= len(c.Sel.Items) {
		return
	}
	c.Sel.Items[row.idx].On = !c.Sel.Items[row.idx].On
	idx := s.list.Index()
	s.list.SetItems(refineItems(sh))
	s.list.Select(idx)
}
