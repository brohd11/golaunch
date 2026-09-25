package app

import (
	"fmt"

	"github.com/brohd11/bubblestack/core"
	"github.com/brohd11/golaunch/internal/config"
	"github.com/brohd11/golaunch/internal/scripts"
	"github.com/brohd11/golaunch/internal/selection"
)

// Ctx is golaunch's app context: the root directory, the current file selection and the
// scanned scripts.
type Ctx struct {
	// ListCompact is the session density shared by standard roots and pickers.
	ListCompact bool

	Root        string
	Version     string
	Sel         selection.Selection
	Preselected bool
	Scripts     []scripts.Script
}

// New builds the context and scans scripts so the first render has rows. Scan problems are
// reported on the first manual refresh.
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

// Rescan re-reads the script directories and returns every problem (nil when clean). A
// config error keeps the previous list.
func (c *Ctx) Rescan() []error {
	cfg, err := config.Load()
	if err != nil {
		return []error{fmt.Errorf("loading config: %w", err)}
	}
	found, problems := scripts.Scan(cfg.ScriptDirs)
	c.Scripts = found
	return problems
}

// Receive rebuilds the tab roots on a theme change.
func (c *Ctx) Receive(sh *core.Shared, payload any) core.Action {
	return core.OnThemeChange(payload)
}

// ListDensity opts standard lists into the app-wide session preference.
func (c *Ctx) ListDensity() *bool { return &c.ListCompact }
