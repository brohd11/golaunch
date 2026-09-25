package app

import (
	"github.com/brohd11/bubblestack/components"
	"github.com/brohd11/bubblestack/core"
)

// actionsMenu is the shared Actions picker opened with "a".
func actionsMenu(sh *core.Shared) *components.PickerScreen {
	return components.NewActionsMenu(selfUpdateHooks(Of(sh).Version),
		"rescan the script directories", refreshAction, nil) // no docs compiled in
}
