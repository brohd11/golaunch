package app

import (
	"fmt"

	"github.com/brohd11/bubblestack/components"
	"github.com/brohd11/bubblestack/core"
	"github.com/brohd11/golaunch/internal/selection"

	"github.com/charmbracelet/bubbles/key"
)

// buildForm is the "Build selection" form: four Yes/No toggles (Include dirs, Include files,
// Recursive, Include current) seeded from the current spec. Submitting resolves the candidate
// paths under the root and stores them as the new selection (all enabled), then pops back to the
// Selection root. Toggle fields cycle with ◄ ► (Left/Right); enter builds, esc cancels.
func buildForm(sh *core.Shared) core.Screen {
	spec := Of(sh).Sel.Spec
	// A never-built selection defaults to "files here" — the common case — so the form isn't all
	// No on first open.
	if !spec.Any() && !spec.Recursive {
		spec.Files = true
	}

	yesNo := func(k, label string, on bool) *components.ToggleField {
		f := components.NewToggleField(k, label, []string{"No", "Yes"})
		if on {
			f.SetIndex(1)
		}
		return f
	}
	dirs := yesNo("dirs", "Include dirs", spec.Dirs)
	files := yesNo("files", "Include files", spec.Files)
	recursive := yesNo("recursive", "Recursive", spec.Recursive)
	current := yesNo("current", "Include current", spec.Current)

	return components.NewForm(components.FormOpts{
		Title:  "Build selection",
		Crumb:  "Build",
		Fields: []components.FormField{dirs, files, recursive, current},
		Help: []key.Binding{
			core.Hint("toggle", core.Keys.Left, core.Keys.Right),
			core.Hint("build", core.Keys.Select),
			core.Hint("cancel", core.Keys.Back),
		},
		OnSubmit: func(sh *core.Shared, _ *components.FormScreen) core.Action {
			spec := selection.Spec{
				Dirs:      dirs.Value() == "Yes",
				Files:     files.Value() == "Yes",
				Recursive: recursive.Value() == "Yes",
				Current:   current.Value() == "Yes",
			}
			if !spec.Any() {
				return core.SetStatus("nothing selected — enable dirs, files, or current")
			}
			c := Of(sh)
			items, err := selection.Resolve(c.Root, spec)
			if err != nil {
				return core.StatusErr(err)
			}
			c.Sel = selection.Selection{Spec: spec, Items: items}
			return core.Seq(
				core.SetStatus(fmt.Sprintf("built: %d paths", len(items))),
				core.Pop(),
			)
		},
	})
}
