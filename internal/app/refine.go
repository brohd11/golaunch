package app

import (
	"path/filepath"

	"github.com/brohd11/bubblestack/core"

	"charm.land/bubbles/v2/key"
	"charm.land/bubbles/v2/list"
	tea "charm.land/bubbletea/v2"
)

// refineKey (shift+R) opens the Refine checklist from the Selection and Scripts screens.
var refineKey = key.NewBinding(key.WithKeys("R"), key.WithHelp("R", "refine"))

// pushRefine opens the Refine checklist over the current stack, or reports that there is nothing to
// refine when no selection has been built yet.
func pushRefine(sh *core.Shared) core.Action {
	if Of(sh).Sel.Empty() {
		return core.SetStatus("nothing to refine — build a selection first")
	}
	return core.Push(NewRefineScreen(sh))
}

// RefineScreen is a filterable checklist over the built selection. Toggles apply
// immediately; esc exits.
type RefineScreen struct {
	checklist
}

var (
	_ core.Filterer = (*RefineScreen)(nil)
	_ core.Crumber  = (*RefineScreen)(nil)
)

func NewRefineScreen(sh *core.Shared) *RefineScreen {
	return &RefineScreen{checklist{list: core.NewSelectList(refineItems(sh), "Refine selection"), crumb: "Refine"}}
}

// refineItems builds rows from Selection.Items (index-aligned): basename as label, full path
// as description and filter text.
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

func (s *RefineScreen) Update(sh *core.Shared, msg tea.Msg) (core.Screen, core.Action) {
	return s, checklistUpdate(&s.list, sh, msg, s.toggleSelected)
}

// toggleSelected flips the row's On state and rebuilds the rows.
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
