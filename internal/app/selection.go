package app

import (
	"github.com/brohd11/bubblestack/components"
	"github.com/brohd11/bubblestack/core"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

// SelectionScreen is golaunch's first tab: two rows — Build selection (a form of dir/file/recursive/
// current toggles that resolves the candidate paths) and Refine selection (a checklist that toggles
// each captured path on/off). The enabled subset is what the Scripts tab launches against; the
// header reflects the running selection. R (shift+R) opens the refine checklist from here too.
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
		list: core.NewSelectList(selectionItems(), TitleSelection, refineKey),
		root: Of(sh).Root,
	}
}

// selectionItems builds the two self-dispatching rows: Build opens the toggle form, Refine opens
// the checklist over whatever the last build captured.
func selectionItems() []list.Item {
	return []list.Item{
		components.Item{
			Name: "Build selection",
			Desc: "choose dirs / files / recursive / current, then resolve the paths",
			Pick: func(sh *core.Shared) core.Action { return core.Push(buildForm(sh)) },
		},
		components.Item{
			Name: "Refine selection",
			Desc: "toggle each captured path on/off (or press R anywhere)",
			Pick: func(sh *core.Shared) core.Action { return pushRefine(sh) },
		},
	}
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
	// R (shift+R) opens the refine checklist, gated behind the filter guard so it doesn't hijack
	// filter typing.
	if k, ok := msg.(tea.KeyMsg); ok && !s.Filtering() && core.MatchKey(k.String(), refineKey) {
		return s, pushRefine(sh)
	}
	return s, components.RootUpdate(sh, &s.list, msg)
}
