package app

import (
	"fmt"

	"github.com/brohd11/bubblestack"
	"github.com/brohd11/bubblestack/components"
	"github.com/brohd11/bubblestack/core"
	"github.com/brohd11/bubblestack/sysopen"
	"github.com/brohd11/golaunch/internal/config"
	"github.com/brohd11/golaunch/internal/selection"
)

// Tab titles, shared by the tab wiring and the (single) place each screen names itself.
const (
	TitleSelection = "Selection"
	TitleScripts   = "Scripts"
)

// Options describes one golaunch session. Preselected is deliberately separate from the current
// selection: every supplied item may later be disabled in Refine without making the Build selection
// UI appear.
type Options struct {
	Root        string
	Version     string
	Selection   selection.Selection
	Preselected bool
}

// Run ensures golaunch's config exists (first run materializes ~/.golaunch and the example
// scripts), then launches either the normal Selection/Scripts TUI or the Scripts-only preselected
// mode. It wires the persistent header, a log/output pane (streamed script output lands there),
// and a status line. The global Refresh key rescans the script directories; the global
// Terminal/OpenDir keys act on the root directory (each tab root is a DirLocator). "a" opens the
// Actions menu (theme, self-update, rescan). The Init startup command runs a background self-update
// check that notes "update available" on the status line (silent otherwise).
func Run(opts Options) error {
	if _, err := config.Ensure(); err != nil {
		return err
	}
	c := New(opts)
	return bubblestack.Run(bubblestack.Config{
		App:                  c,
		Header:               Header,
		Output:               components.NewLogPane(),
		Status:               components.NewStatusLine(),
		Tabs:                 tabs(c.Preselected),
		Init:                 SelfUpdateCheckCmd,
		RefreshAction:        func(sh *core.Shared) core.Action { return refreshAction(sh) },
		TerminalAction:       func(dir string) core.Action { return sysopen.TerminalInline(dir) },
		TerminalWindowAction: func(dir string) core.Action { return sysopen.Terminal(dir) },
		OpenDirAction:        func(dir string) core.Action { return sysopen.Path(dir, false) },
	})
}

// tabs omits selection building for an argv-supplied selection. Refine remains reachable from the
// Scripts root and its submenus through refineKey.
func tabs(preselected bool) []bubblestack.TabEntry {
	scriptsTab := bubblestack.TabEntry{
		Title: TitleScripts,
		New:   func(sh *core.Shared) core.Screen { return NewScriptsScreen(sh) },
	}
	if preselected {
		return []bubblestack.TabEntry{scriptsTab}
	}
	return []bubblestack.TabEntry{
		{Title: TitleSelection, New: func(sh *core.Shared) core.Screen { return NewSelectionScreen(sh) }},
		scriptsTab,
	}
}

// refreshAction rescans the script directories and rebuilds the tab roots so the Scripts tab picks
// up added/removed scripts and edited metadata. Scan problems (an unreadable directory, a script
// whose header fails to parse, a broken config) are logged individually and counted in the status,
// so a script that vanishes from the menu leaves a trace.
func refreshAction(sh *core.Shared) core.Action {
	problems := Of(sh).Rescan()
	for _, p := range problems {
		sh.Log("rescan: " + p.Error())
	}
	status := "rescanned scripts"
	if len(problems) > 0 {
		status = fmt.Sprintf("rescanned scripts — %d problem(s), see log", len(problems))
	}
	return core.Seq(
		core.SetStatus(status),
		core.RefreshRoots(),
	)
}
