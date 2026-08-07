package app

import "github.com/charmbracelet/bubbles/key"

// keys are golaunch's screen-level bindings that aren't part of bubblestack's framework
// keymap (core.Keys). Only the Actions menu for now; refineKey lives beside the refine
// screen it belongs to.
var keys = struct {
	Actions key.Binding // open the Actions menu (theme, update, refresh)
}{
	Actions: key.NewBinding(key.WithKeys("a"), key.WithHelp("a", "actions")),
}
