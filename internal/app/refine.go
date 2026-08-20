package app

import (
	"path/filepath"

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

// refineItems builds the checklist rows from the current Selection.Items (index-aligned). The
// label is the basename — directories keep a trailing "/" so they read as directories — while
// the full path is both the description and what the filter matches, so typing a parent
// directory's name finds rows whose basename says nothing about where they live.
func refineItems(sh *core.Shared) []list.Item {
	sel := Of(sh).Sel.Items
	rows := make([]list.Item, len(sel))
	for i, it := range sel {
		name := filepath.Base(it.Path)
		if it.IsDir {
			name += "/"
		}
		rows[i] = checkRow{idx: i, label: name, desc: it.Path, filter: it.Path, on: it.On}
	}
	return rows
}

func (s *RefineScreen) Init(*core.Shared) tea.Cmd    { return nil }
func (s *RefineScreen) Filtering() bool              { return s.list.FilterState() == list.Filtering }
func (s *RefineScreen) View(*core.Shared) string     { return core.RenderList(s.list) }
func (s *RefineScreen) HelpView(*core.Shared) string { return core.ShortHelp(s.list, core.HelpMinimal) }
func (s *RefineScreen) CrumbLabel(bool) string       { return "Refine" }

func (s *RefineScreen) SetSize(_ *core.Shared, w, h int) { s.list.SetSize(w, h) }

func (s *RefineScreen) Update(sh *core.Shared, msg tea.Msg) (core.Screen, core.Action) {
	return s, checklistUpdate(&s.list, sh, msg, s.toggleSelected)
}

// toggleSelected flips the highlighted row's On state on the shared selection and rebuilds the rows
// so the [x]/[ ] marker updates, keeping the cursor where it was. Toggling a path is always live and
// can't fail, so there is nothing to report back: the action is empty either way.
func (s *RefineScreen) toggleSelected(sh *core.Shared) core.Action {
	row, ok := s.list.SelectedItem().(checkRow)
	if !ok {
		return core.Action{}
	}
	c := Of(sh)
	if row.idx < 0 || row.idx >= len(c.Sel.Items) {
		return core.Action{}
	}
	c.Sel.Items[row.idx].On = !c.Sel.Items[row.idx].On
	setRowsKeepCursor(&s.list, refineItems(sh))
	return core.Action{}
}
