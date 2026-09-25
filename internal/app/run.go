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

// Options describes one golaunch session. Preselected stays separate from the selection so
// disabling every item in Refine doesn't bring back the Build UI.
type Options struct {
	Root        string
	Version     string
	Selection   selection.Selection
	Preselected bool
}

// Run ensures the config exists, then launches the Selection/Scripts TUI (or Scripts-only
// when preselected).
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
		TerminalAction:       func(dir string) core.Action { return sysopen.TerminalInlineFor("golaunch", dir) },
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

// refreshAction rescans the scripts and rebuilds the tab roots. Each scan problem is logged
// and counted in the status.
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
