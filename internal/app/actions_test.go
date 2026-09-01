package app

import (
	"reflect"
	"testing"

	"github.com/brohd11/bubblestack/core"
)

// The Actions-key tests drive each tab root's Update directly with the "a" keypress —
// no router, no network (the menu's self-update hooks are inert closures until fired).

func actionsTestShared() *core.Shared {
	return core.NewShared(&Ctx{Root: ".", Version: "dev"})
}

// msgType names the action's navigation message type, "" when the action carries none.
func msgType(act core.Action) string {
	if act.Msg == nil {
		return ""
	}
	return reflect.TypeOf(act.Msg).String()
}

func TestActionsKeyPushesMenu(t *testing.T) {
	for _, newScreen := range map[string]func(*core.Shared) core.Screen{
		"Selection": func(sh *core.Shared) core.Screen { return NewSelectionScreen(sh) },
		"Scripts":   func(sh *core.Shared) core.Screen { return NewScriptsScreen(sh) },
	} {
		sh := actionsTestShared()
		s := newScreen(sh)
		_, act := s.Update(sh, keyMsg("a"))
		if got := msgType(act); got != "core.pushMsg" {
			t.Errorf("%s tab: \"a\" should push the Actions menu, got %s", reflect.TypeOf(s), got)
		}
	}
}

func TestActionsKeyYieldsToFilterTyping(t *testing.T) {
	sh := actionsTestShared()
	s := NewScriptsScreen(sh)
	// "/" opens the list's filter input (reached via RootUpdate's fall-through).
	s.Update(sh, keyMsg("/"))
	if !s.Filtering() {
		t.Fatal("\"/\" should put the list into its Filtering state")
	}
	_, act := s.Update(sh, keyMsg("a"))
	if got := msgType(act); got == "core.pushMsg" {
		t.Error("while filtering, \"a\" should be typed into the filter, not open the Actions menu")
	}
}
