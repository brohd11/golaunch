package app

import (
	"github.com/brohd11/bubblestack/components"
	"github.com/brohd11/bubblestack/core"
	bsupdate "github.com/brohd11/bubblestack/selfupdate"

	tea "charm.land/bubbletea/v2"
)

// selfUpdateRepo is golaunch's own GitHub repo slug, passed to the shared self-update bridge.
const selfUpdateRepo = "brohd11/golaunch"

// selfUpdateHooks builds the shared self-update flow's (bubblestack/components) hook
// set for golaunch. The goutil↔components wiring — the Check/Apply closures and the
// conversion between goutil's selfupdate.Info and the flow's app-agnostic SelfUpdateInfo
// (field-identical by design) — lives in the bubblestack/selfupdate bridge, which every
// app in the monorepo used to hand-roll a copy of.
func selfUpdateHooks(version string) components.SelfUpdateHooks {
	return bsupdate.Hooks("golaunch", selfUpdateRepo, version)
}

// SelfUpdateCheckCmd is the app-level startup command (wired onto bubblestack Config.Init):
// it checks golaunch's own repo for a newer release off the UI thread and, only when an
// update is available, writes an "update available" line to the shared status line and log.
// Anything else (up to date, dev build, fetch error) is silent. The flow and timeout are
// the shared ones in bubblestack/components; only the hooks are golaunch's.
func SelfUpdateCheckCmd(sh *core.Shared) tea.Cmd {
	return components.SelfUpdateCheckCmd(selfUpdateHooks(Of(sh).Version))
}
