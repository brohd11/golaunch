package app

import (
	"fmt"

	"github.com/brohd11/bubblestack/core"
	"github.com/brohd11/golaunch/internal/config"
	"github.com/brohd11/golaunch/internal/scripts"
	"github.com/brohd11/golaunch/internal/selection"
)

// Ctx is golaunch's app context, stored on core.Shared.App and recovered with Of. It holds the
// root directory, the current file selection (built in the Selection tab or supplied on argv), and
// the scripts found by scanning the configured directories. There is no manifest — like repoview,
// the state is whatever a fresh scan turns up.
type Ctx struct {
	Root        string
	Version     string
	Sel         selection.Selection
	Preselected bool
	Scripts     []scripts.Script
}

// New builds the context and performs the initial script scan, so the Scripts tab has rows on
// first render. Startup scan problems go unreported here (no status line exists yet); the same
// problems surface on the first manual refresh, which reports what Rescan returns.
func New(opts Options) *Ctx {
	c := &Ctx{
		Root:        opts.Root,
		Version:     opts.Version,
		Sel:         opts.Selection,
		Preselected: opts.Preselected,
	}
	c.Rescan()
	return c
}

// Of recovers the golaunch context from a Shared. Screens call c := app.Of(sh).
func Of(sh *core.Shared) *Ctx { return core.App[Ctx](sh) }

// Rescan re-reads the configured script directories and returns what went wrong: a config
// read error (the previous list stays intact rather than blanking the Scripts tab) or one
// error per unreadable directory / unparsable script header. Nil means a clean scan.
func (c *Ctx) Rescan() []error {
	cfg, err := config.Load()
	if err != nil {
		return []error{fmt.Errorf("loading config: %w", err)}
	}
	found, problems := scripts.Scan(cfg.ScriptDirs)
	c.Scripts = found
	return problems
}

// Receive handles app-level broadcasts. On a theme change it returns RefreshRoots so the tab roots
// rebuild and re-bake their list styles from the new palette (the router-drawn chrome repaints on
// its own). Everything else is handled by the screens.
func (c *Ctx) Receive(sh *core.Shared, payload any) core.Action {
	return core.OnThemeChange(payload)
}
