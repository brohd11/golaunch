package app

import (
	"fmt"

	"github.com/brohd11/bubblestack/core"
)

// Header renders golaunch's persistent context box: the root directory, the current selection
// summary, and how many scripts were scanned. Wired onto core.Chrome.Header, so the router draws
// it above every screen — and because it reads Ctx live each frame, the selection line updates the
// moment a mode is picked in the Selection tab.
func Header(sh *core.Shared) string {
	c := Of(sh)
	inner := core.HeaderInnerWidth(sh.Width())
	// Value budget: inner width minus the box's horizontal padding (2) and the label (8).
	valWidth := inner - 10
	body := core.Label("Root:   ") + core.Value(core.TruncLeft(c.Root, valWidth)) + "\n" +
		core.Label("Select: ") + core.Value(c.Sel.Summary()) + "\n" +
		core.Label("Scripts:") + core.Value(fmt.Sprintf(" %d found", len(c.Scripts)))
	return core.HeaderBox(sh.Width(), body)
}
