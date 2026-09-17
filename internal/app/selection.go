package app

import (
	"github.com/brohd11/bubblestack/components"
	"github.com/brohd11/bubblestack/core"

	"charm.land/bubbles/v2/list"
)

// NewSelectionScreen builds the selection tab using shared root-list behavior.
func NewSelectionScreen(sh *core.Shared) *components.RootListScreen {
	return components.NewRootList(selectionItems(), rootListOpts(sh, TitleSelection))
}

// selectionItems builds the two self-dispatching rows: Build opens the flag checklist, Refine opens
// the checklist over whatever the last build captured.
func selectionItems() []list.Item {
	return []list.Item{
		components.Item{
			Name: "Build selection",
			Desc: "toggle dirs / files / recursive / current — the paths resolve live",
			Pick: func(sh *core.Shared) core.Action { return pushBuild(sh) },
		},
		components.Item{
			Name: "Refine selection",
			Desc: "toggle each captured path on/off (or press R anywhere)",
			Pick: func(sh *core.Shared) core.Action { return pushRefine(sh) },
		},
	}
}
