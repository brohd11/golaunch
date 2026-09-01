package app

import (
	"fmt"

	"github.com/brohd11/bubblestack/components"
	"github.com/brohd11/bubblestack/core"
	"github.com/brohd11/golaunch/internal/scripts"

	"charm.land/bubbles/v2/key"
	"charm.land/bubbles/v2/list"
	tea "charm.land/bubbletea/v2"
)

// groupGlyph marks a row that opens a submenu rather than launching a script.
const groupGlyph = " ▸"

// ScriptsScreen is golaunch's second tab: the scanned scripts, grouped into submenus by their
// path= metadata. A group row opens a PickerScreen for that subtree; a script row launches the
// script against the current selection (streaming into the log, or in an external terminal when the
// script opts in). The tree is rebuilt whenever the tab root is (startup and Refresh).
type ScriptsScreen struct {
	list list.Model
	root string
}

var (
	_ core.Filterer   = (*ScriptsScreen)(nil)
	_ core.Crumber    = (*ScriptsScreen)(nil)
	_ core.DirLocator = (*ScriptsScreen)(nil)
)

func NewScriptsScreen(sh *core.Shared) *ScriptsScreen {
	c := Of(sh)
	tree := scripts.BuildTree(c.Scripts)
	return &ScriptsScreen{
		list: core.NewSelectList(nodeItems(c.Root, tree), TitleScripts, refineKey, keys.Actions),
		root: c.Root,
	}
}

// nodeItems builds the rows for one tree node: a submenu-opening row per child group, then a
// launch row per script, with a placeholder when the node is empty. Recursion builds each subtree's
// picker lazily inside the group row's Pick.
func nodeItems(root string, node *scripts.Tree) []list.Item {
	var items []list.Item
	for _, child := range node.Children {
		child := child
		items = append(items, components.Item{
			Name: child.Name + groupGlyph,
			Desc: groupDesc(child),
			Pick: func(sh *core.Shared) core.Action {
				return core.Push(components.NewPicker(nodeItems(root, child), components.PickerOpts{
					Title: child.Name,
					Crumb: child.Name,
					Dir:   root,
					Help:  []key.Binding{refineKey},
					// R (shift+R) opens the refine checklist from inside a script submenu — a quick
					// tweak without backing out to the Selection tab.
					OnKey: func(sh *core.Shared, k string, _ list.Item) (core.Action, bool) {
						if core.MatchKey(k, refineKey) {
							return pushRefine(sh), true
						}
						return core.Action{}, false
					},
				}))
			},
		})
	}
	for _, s := range node.Scripts {
		s := s
		items = append(items, components.Item{
			Name: s.DisplayName(),
			Desc: scriptDesc(s),
			Pick: func(sh *core.Shared) core.Action {
				c := Of(sh)
				return scripts.Launch(sh, s, c.Root, c.Sel.Paths())
			},
		})
	}
	return components.EnsurePlaceholder(items, "no scripts",
		"add scripts to a directory in ~/.golaunch/config.yml, then Refresh")
}

// groupDesc summarizes a submenu's contents for its row.
func groupDesc(node *scripts.Tree) string {
	n := len(node.Scripts)
	if g := len(node.Children); g > 0 {
		return fmt.Sprintf("%d script(s), %d subgroup(s)", n, g)
	}
	return fmt.Sprintf("%d script(s)", n)
}

// scriptDesc is a script row's description: its metadata desc, with a marker when it launches in an
// external terminal instead of streaming into the TUI.
func scriptDesc(s scripts.Script) string {
	desc := s.Meta.Desc
	if desc == "" {
		desc = s.File
	}
	if s.Meta.Terminal {
		desc += "  [terminal]"
	}
	return desc
}

func (s *ScriptsScreen) Init(*core.Shared) tea.Cmd        { return nil }
func (s *ScriptsScreen) Filtering() bool                  { return s.list.FilterState() == list.Filtering }
func (s *ScriptsScreen) View(*core.Shared) string         { return core.RenderList(s.list) }
func (s *ScriptsScreen) HelpView(*core.Shared) string     { return core.ShortHelp(s.list, core.HelpTabbed) }
func (s *ScriptsScreen) SetSize(_ *core.Shared, w, h int) { s.list.SetSize(w, h) }
func (s *ScriptsScreen) CrumbLabel(bool) string           { return TitleScripts }
func (s *ScriptsScreen) LocateDir() (string, bool)        { return s.root, s.root != "" }

func (s *ScriptsScreen) Update(sh *core.Shared, msg tea.Msg) (core.Screen, core.Action) {
	return s, tabRootUpdate(sh, &s.list, msg)
}
