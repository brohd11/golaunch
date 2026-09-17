package app

import (
	"testing"

	"github.com/brohd11/bubblestack/components"
	"github.com/brohd11/bubblestack/core"
	"github.com/brohd11/golaunch/internal/scripts"

	tea "charm.land/bubbletea/v2"
)

func TestRootAndScriptGroupDensity(t *testing.T) {
	c := &Ctx{Root: t.TempDir(), Scripts: []scripts.Script{{File: "example.sh", Meta: scripts.Meta{Path: "Group"}}}}
	sh := core.NewShared(c)
	r := core.NewRouter(sh, tabs(false))
	step := func(msg tea.Msg) {
		m, _ := r.Update(msg)
		r = m.(core.Router)
	}
	step(tea.WindowSizeMsg{Width: 80, Height: 24})
	selectionRoot := r.Top().(*components.RootListScreen)
	step(keyMsg("D"))
	step(keyMsg("]"))
	scriptsRoot := r.Top().(*components.RootListScreen)
	if !scriptsRoot.Compact() || scriptsRoot.CrumbLabel(false) != TitleScripts {
		t.Fatal("Scripts must inherit density and retain its breadcrumb")
	}
	if dir, ok := scriptsRoot.LocateDir(); !ok || dir != c.Root {
		t.Fatal("Scripts must retain its directory capability")
	}
	step(keyMsg("enter")) // open a group, without launching its script
	group, ok := r.Top().(*components.PickerScreen)
	if !ok || !group.Compact() || group.CrumbLabel(false) != "Group" {
		t.Fatal("script group picker must inherit density and navigation")
	}
	step(keyMsg("D"))
	if c.ListCompact || scriptsRoot.Compact() || selectionRoot.Compact() {
		t.Fatal("a group picker must update both cached roots")
	}
	step(keyMsg("esc"))
	if r.Top() != scriptsRoot {
		t.Fatal("back must return to the same Scripts root")
	}
}
