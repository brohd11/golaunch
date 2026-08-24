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
	valWidth := core.HeaderValueWidth(sh.Width(), "Root:   ")
	body := core.Label("Root:   ") + core.Value(core.TruncLeft(c.Root, valWidth)) + "\n" +
		core.Label("Select: ") + core.Value(c.Sel.Summary()) + "\n" +
		core.Label("Scripts:") + core.Value(fmt.Sprintf(" %d found", len(c.Scripts)))
	return core.HeaderBox(sh.Width(), body)
}
