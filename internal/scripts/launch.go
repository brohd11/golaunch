package scripts

import (
	"context"
	"strings"

	"github.com/brohd11/bubblestack/components"
	"github.com/brohd11/bubblestack/core"
	"github.com/brohd11/bubblestack/sysopen"
)

// Launch runs a script against the selected paths, rooted at root (the script's working directory).
// The argv is the interpreter (if any), the script file, then the selected paths as positional
// arguments. A script with terminal=true metadata is handed to an external terminal; otherwise it
// streams into the TUI as a stay-task the user dismisses with esc. An empty selection is refused
// with a status hint rather than running a script against nothing.
func Launch(sh *core.Shared, s Script, root string, paths []string) core.Action {
	if len(paths) == 0 {
		return core.SetStatus("no selection — pick paths in the Selection tab first")
	}
	argv := make([]string, 0, len(s.Interp)+1+len(paths))
	argv = append(argv, s.Interp...)
	argv = append(argv, s.File)
	argv = append(argv, paths...)

	if s.Meta.Terminal {
		// terminal=true: launch in its own window (interactive/GUI/long-running tools). Output is
		// not captured in the TUI; the terminal is rooted at the selection's root directory.
		return sysopen.Terminal(root, argv...)
	}

	label := "run " + s.DisplayName()
	run := func(ctx context.Context, sh *core.Shared, report func(string, ...any), done chan<- core.TaskEvent) {
		report("$ %s", strings.Join(argv, " "))
		done <- core.TaskEvent{Done: true, Err: streamCmd(ctx, root, report, argv...)}
	}
	onDone := func(sh *core.Shared, ev core.TaskEvent) core.Action {
		if ev.Err != nil {
			return core.SetStatusAndLog(s.DisplayName() + " failed: " + ev.Err.Error())
		}
		return core.SetStatus(s.DisplayName() + " — done")
	}
	onDismiss := func(*core.Shared) core.Action { return core.Pop() }
	return core.Push(components.NewStayTask(label, "done — esc to go back", run, onDone, onDismiss))
}
