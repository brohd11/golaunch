package app

import (
	"fmt"

	"github.com/brohd11/bubblestack/core"
)

// Header renders the root, the selection summary and the script count, read live from Ctx.
func Header(sh *core.Shared) string {
	c := Of(sh)
	valWidth := core.HeaderValueWidth(sh.Width(), "Root:   ")
	body := core.Label("Root:   ") + core.Value(core.TruncLeft(c.Root, valWidth)) + "\n" +
		core.Label("Select: ") + core.Value(c.Sel.Summary()) + "\n" +
		core.Label("Scripts:") + core.Value(fmt.Sprintf(" %d found", len(c.Scripts)))
	return core.HeaderBox(sh.Width(), body)
}
