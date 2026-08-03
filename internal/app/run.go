package app

import (
	"github.com/brohd11/bubblestack"
	"github.com/brohd11/bubblestack/components"
	"github.com/brohd11/bubblestack/core"
	"github.com/brohd11/bubblestack/sysopen"
	"github.com/brohd11/golaunch/internal/config"
)

// Tab titles, shared by the tab wiring and the (single) place each screen names itself.
const (
	TitleSelection = "Selection"
	TitleScripts   = "Scripts"
)

// Run ensures golaunch's config exists (first run materializes ~/.golaunch and the example
// scripts), then launches the two-tab TUI: Selection (build a set of paths) and Scripts (launch
// one against that selection). It wires the persistent header, a log/output pane (streamed script
// output lands there), and a status line. The global Refresh key rescans the script directories;
// the global Terminal/OpenDir keys act on the root directory (each tab root is a DirLocator).
func Run(root, version string) error {
	if _, err := config.Ensure(); err != nil {
		return err
	}
	return bubblestack.Run(bubblestack.Config{
		App:    New(root, version),
		Header: Header,
		Output: components.NewLogPane(),
		Status: components.NewStatusLine(),
		Tabs: []bubblestack.TabEntry{
			{Title: TitleSelection, New: func(sh *core.Shared) core.Screen { return NewSelectionScreen(sh) }},
			{Title: TitleScripts, New: func(sh *core.Shared) core.Screen { return NewScriptsScreen(sh) }},
		},
		RefreshAction:  func(sh *core.Shared) core.Action { return refreshAction(sh) },
		TerminalAction: func(dir string) core.Action { return sysopen.Terminal(dir) },
		OpenDirAction:  func(dir string) core.Action { return sysopen.Path(dir, false) },
	})
}

// refreshAction rescans the script directories and rebuilds the tab roots so the Scripts tab picks
// up added/removed scripts and edited metadata.
func refreshAction(sh *core.Shared) core.Action {
	Of(sh).Rescan()
	return core.Seq(
		core.SetStatus("rescanned scripts"),
		core.RefreshRoots(),
	)
}
