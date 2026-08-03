// Command golaunch is a small TUI for launching scripts against a file selection. It opens
// rooted at the current directory: the Selection tab builds a set of paths (the current dir,
// its child dirs, or its child files), and the Scripts tab lists scripts scanned from configured
// directories (nested into submenus by each script's metadata) and runs the selected one against
// the current selection — streaming its output into the TUI, or launching it in an external
// terminal when the script opts in. It's the manifest-free sibling of repoview/gdaddon, built on
// the same bubblestack framework.
package main

import "github.com/brohd11/golaunch/cmd"

func main() {
	cmd.Execute()
}
