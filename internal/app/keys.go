package app

import (
	"github.com/brohd11/bubblestack/components"
	"github.com/brohd11/bubblestack/core"

	"charm.land/bubbles/v2/key"
	"charm.land/bubbles/v2/list"
	tea "charm.land/bubbletea/v2"
)

// keys are golaunch's screen-level bindings that aren't part of bubblestack's framework
// keymap (core.Keys). Only the Actions menu for now; refineKey lives beside the refine
// screen it belongs to.
var keys = struct {
	Actions key.Binding // open the Actions menu (theme, update, refresh)
}{
	Actions: key.NewBinding(key.WithKeys("a"), key.WithHelp("a", "actions")),
}

// tabRootUpdate is the Update body both tab roots (Selection and Scripts) share: golaunch's
// two tab-level keys, then the framework's root list handling. Adding a third tab-level key
// should reach both tabs, which is why this is one function and not a copy in each.
//
// The keys are gated behind the filter guard so they don't hijack filter typing — an "a"
// typed into a filter is a letter, not the Actions menu.
func tabRootUpdate(sh *core.Shared, l *list.Model, msg tea.Msg) core.Action {
	if k, ok := msg.(tea.KeyPressMsg); ok && l.FilterState() != list.Filtering {
		switch {
		case core.MatchKey(k.String(), refineKey):
			return pushRefine(sh)
		case core.MatchKey(k.String(), keys.Actions):
			return core.Push(actionsMenu(sh))
		}
	}
	return components.RootUpdate(sh, l, msg)
}
