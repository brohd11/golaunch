package app

import (
	"testing"

	"github.com/brohd11/bubblestack/core"
	"github.com/brohd11/golaunch/internal/selection"
)

func TestTabsForSelectionMode(t *testing.T) {
	normal := tabs(false)
	if len(normal) != 2 || normal[0].Title != TitleSelection || normal[1].Title != TitleScripts {
		t.Fatalf("normal tabs = %+v, want Selection then Scripts", normal)
	}

	preselected := tabs(true)
	if len(preselected) != 1 || preselected[0].Title != TitleScripts {
		t.Fatalf("preselected tabs = %+v, want Scripts only", preselected)
	}
}

func TestPreselectedScriptsCanOpenRefine(t *testing.T) {
	sh := core.NewShared(&Ctx{
		Root:        ".",
		Preselected: true,
		Sel: selection.Selection{Items: []selection.Item{
			{Path: "/tmp/report.txt", On: true},
		}},
	})
	screen := NewScriptsScreen(sh)
	_, act := screen.Update(sh, keyMsg("R"))
	if got := msgType(act); got != "core.pushMsg" {
		t.Fatalf("R should push Refine in preselected mode, got %s", got)
	}
}
