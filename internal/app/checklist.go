package app

import (
	"github.com/brohd11/bubblestack/components"
	"github.com/brohd11/bubblestack/core"

	"charm.land/bubbles/v2/list"
	tea "charm.land/bubbletea/v2"
)

// The Build and Refine screens are both a filterable list of [x]/[ ] rows where enter
// toggles the highlighted row and esc leaves. This file holds what they share.

// checklist is the screen plumbing Build and Refine embed; each supplies its own Update.
type checklist struct {
	list  list.Model
	crumb string
}

func (s *checklist) Init(*core.Shared) tea.Cmd        { return nil }
func (s *checklist) Filtering() bool                  { return s.list.FilterState() == list.Filtering }
func (s *checklist) View(*core.Shared) string         { return core.RenderList(s.list) }
func (s *checklist) HelpView(*core.Shared) string     { return core.ShortHelp(s.list, core.HelpMinimal) }
func (s *checklist) CrumbLabel(bool) string           { return s.crumb }
func (s *checklist) SetSize(_ *core.Shared, w, h int) { s.list.SetSize(w, h) }

// checkRow is one checklist row. idx maps it back to the source of truth. filter is separate
// from label because Refine matches the full path.
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

// checklistUpdate is both checklists' shared Update; onSelect is what enter does.
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

// setRowsKeepCursor swaps in rebuilt rows without resetting the highlight.
func setRowsKeepCursor(l *list.Model, rows []list.Item) {
	idx := l.Index()
	l.SetItems(rows)
	l.Select(idx)
}
