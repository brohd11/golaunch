package app

import (
	"github.com/brohd11/bubblestack/components"
	"github.com/brohd11/bubblestack/core"
)

// actionsMenu is the small Actions picker opened with "a": the shared bubblestack menu
// (theme, self-update, refresh), with golaunch's Refresh being the script rescan the
// global Refresh key fires. PopStop (set by the shared menu) makes it the hub its
// sub-flows return to.
func actionsMenu(sh *core.Shared) *components.PickerScreen {
	return components.NewActionsMenu(selfUpdateHooks(Of(sh).Version),
		"rescan the script directories", refreshAction, nil) // no docs compiled in
}
