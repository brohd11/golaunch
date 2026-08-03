package app

import (
	"github.com/brohd11/bubblestack/components"
	"github.com/brohd11/bubblestack/core"
	"github.com/brohd11/golaunch/internal/selection"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

// SelectionScreen is golaunch's first tab: a short list of the ways to build a selection from the
// root directory (the dir itself, its child dirs, or its child files). Picking a row resolves the
// paths and stores them on the shared Ctx; the header reflects the new selection immediately, and
// the Scripts tab reads it when a script is launched.
type SelectionScreen struct {
	list list.Model
	root string // the directory this tab concerns; enables the global Terminal/OpenDir keys
}

var (
	_ core.Filterer   = (*SelectionScreen)(nil)
	_ core.Crumber    = (*SelectionScreen)(nil)
	_ core.DirLocator = (*SelectionScreen)(nil)
)

func NewSelectionScreen(sh *core.Shared) *SelectionScreen {
	return &SelectionScreen{
		list: core.NewSelectList(selectionItems(), TitleSelection),
		root: Of(sh).Root,
	}
}

// selectionItems builds the three self-dispatching mode rows. Each Pick resolves the mode against
// the root and records it as the current selection.
func selectionItems() []list.Item {
	mk := func(mode selection.Mode, desc string) list.Item {
		return components.Item{
			Name: mode.Label(),
			Desc: desc,
			Pick: func(sh *core.Shared) core.Action { return applySelection(sh, mode) },
		}
	}
	return []list.Item{
		mk(selection.CurrentDir, "the root directory itself"),
		mk(selection.ChildDirs, "the root's immediate subdirectories"),
		mk(selection.ChildFiles, "the root's immediate files"),
	}
}

// applySelection resolves mode against the root and stores it, reporting the outcome on the status
// line (the header shows the running summary).
func applySelection(sh *core.Shared, mode selection.Mode) core.Action {
	c := Of(sh)
	sel, err := selection.Resolve(c.Root, mode)
	if err != nil {
		return core.StatusErr(err)
	}
	c.Sel = sel
	return core.SetStatus("selection: " + sel.Summary())
}

func (s *SelectionScreen) Init(*core.Shared) tea.Cmd { return nil }
func (s *SelectionScreen) Filtering() bool           { return s.list.FilterState() == list.Filtering }
func (s *SelectionScreen) View(*core.Shared) string  { return s.list.View() }
func (s *SelectionScreen) HelpView(*core.Shared) string {
	return core.ShortHelp(s.list, core.HelpTabbed)
}
func (s *SelectionScreen) SetSize(_ *core.Shared, w, h int) { s.list.SetSize(w, h) }
func (s *SelectionScreen) CrumbLabel(bool) string           { return TitleSelection }
func (s *SelectionScreen) LocateDir() (string, bool)        { return s.root, s.root != "" }

func (s *SelectionScreen) Update(sh *core.Shared, msg tea.Msg) (core.Screen, core.Action) {
	return s, components.RootUpdate(sh, &s.list, msg)
}
