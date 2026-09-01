package app

import (
	"github.com/brohd11/bubblestack/components"
	"github.com/brohd11/bubblestack/core"

	"charm.land/bubbles/v2/list"
	tea "charm.land/bubbletea/v2"
)

// The Build and Refine screens are the same screen at two altitudes: a filterable list of
// [x]/[ ] rows where enter toggles the highlighted one in place and esc leaves with the
// result. This file holds the parts that were literally the same in both, so a change to how
// a checklist behaves lands on both instead of only whichever one was being edited.

// checkRow is one checklist row: a [x]/[ ] marker plus a label. idx maps the row back to
// whatever the screen is a view of (a buildFlags entry, a Selection item) — that, not the
// row, stays the source of truth for on, since every toggle rebuilds all the rows anyway.
//
// filter is separate from label because the two screens filter on different text: Build has
// nothing but the label to match, while Refine matches the full path it shows as the
// description, so typing a directory name finds rows whose basename doesn't contain it.
type checkRow struct {
	idx         int
	label, desc string
	filter      string
	on          bool
}

func (r checkRow) Title() string {
	mark := "[ ] "
	if r.on {
		mark = "[x] "
	}
	return mark + r.label
}

func (r checkRow) Description() string { return r.desc }
func (r checkRow) FilterValue() string { return r.filter }

// checklistUpdate is the shared body of both checklists' Update. onSelect is the only thing
// that differs between them — what enter does to the highlighted row — so it is the only
// thing passed in. The caller returns itself as the screen; this returns only the action.
func checklistUpdate(l *list.Model, sh *core.Shared, msg tea.Msg, onSelect func(*core.Shared) core.Action) core.Action {
	// v2 gives the wheel its own message type, so the kind is in the match rather
	// than in a field check inside WheelNav.
	if m, ok := msg.(tea.MouseWheelMsg); ok {
		if components.WheelNav(l, m.Mouse()) {
			return core.Action{}
		}
	}
	// While actively typing a filter, every key belongs to the filter input.
	if l.FilterState() == list.Filtering {
		var cmd tea.Cmd
		*l, cmd = l.Update(msg)
		return core.Async(cmd)
	}
	if km, ok := msg.(tea.KeyPressMsg); ok {
		k := km.String()
		switch {
		case core.MatchKey(k, core.Keys.Select):
			// enter toggles the highlighted row and applies it on the spot; the screen stays
			// open so the selection can be narrowed a row at a time.
			return onSelect(sh)
		case core.MatchKey(k, core.Keys.Back):
			// esc exits, keeping whatever the last toggle resolved; confirm it on the status line.
			return core.Seq(core.SetStatus("selection: "+Of(sh).Sel.Summary()), core.Pop())
		default:
			if components.WrapNav(l, k) {
				return core.Action{}
			}
		}
	}
	var cmd tea.Cmd
	*l, cmd = l.Update(msg)
	return core.Async(cmd)
}

// setRowsKeepCursor swaps in freshly built rows without moving the highlight. Both screens
// rebuild every row on each toggle — that is how the [x] marker updates — and SetItems on its
// own would send the cursor back to the top after every keystroke.
func setRowsKeepCursor(l *list.Model, rows []list.Item) {
	idx := l.Index()
	l.SetItems(rows)
	l.Select(idx)
}
