package app

import (
	"github.com/brohd11/bubblestack/components"
	"github.com/brohd11/bubblestack/core"

	"charm.land/bubbles/v2/key"
	"charm.land/bubbles/v2/list"
)

// keys are golaunch's screen-level bindings. Actions is core.Keys.Actions rather than a
// copy of it — the picker key is shared by every app on the framework, and aliasing it
// here is what carries its ctrl+alt+a form in. refineKey lives beside the refine screen
// it belongs to.
var keys = struct {
	Actions key.Binding // open the Actions menu (theme, update, refresh)
}{
	Actions: core.Keys.Actions,
}

// rootListOpts shares the roots' directory capabilities and app commands.
// RootListScreen supplies filtering, density, tab help, and row-key fallback.
func rootListOpts(sh *core.Shared, title string) components.RootListOpts {
	return components.RootListOpts{PickerOpts: components.PickerOpts{
		Title: title,
		Crumb: title,
		Dir:   Of(sh).Root,
		Help:  []key.Binding{refineKey, keys.Actions},
		OnKey: func(sh *core.Shared, k string, _ list.Item) (core.Action, bool) {
			switch {
			case core.MatchKey(k, refineKey):
				return pushRefine(sh), true
			case core.MatchKey(k, keys.Actions):
				return core.Push(actionsMenu(sh)), true
			}
			return core.Action{}, false
		},
	}}
}
