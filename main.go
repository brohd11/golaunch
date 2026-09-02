// Command golaunch is a small TUI for launching scripts against a file selection. With no
// positional paths, the Selection tab builds a set under the current root; with paths, they seed
// the selection and the TUI opens directly on Scripts while retaining its Refine shortcut. Script
// output streams into the TUI unless the script opts into an external terminal. It's the
// manifest-free sibling of repoview/gdaddon, built on the same bubblestack framework.
package main

import "github.com/brohd11/golaunch/cmd"

func main() {
	cmd.Execute()
}
