package scripts

import (
	"context"
	"strings"

	"github.com/brohd11/bubblestack/components"
	"github.com/brohd11/bubblestack/core"
	"github.com/brohd11/bubblestack/sysopen"
	"github.com/brohd11/goutil/stream"
)

// Launch runs a script against paths with root as the working directory: in an external
// terminal when terminal=true, otherwise streamed into the TUI. An empty selection is refused.
func Launch(sh *core.Shared, s Script, root string, paths []string) core.Action {
	if len(paths) == 0 {
		return core.SetStatus("no paths selected")
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
		done <- core.TaskEvent{Done: true, Err: stream.Cmd(ctx, root, nil, report, argv...)}
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
