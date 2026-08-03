package cmd

import (
	"os"
	"path/filepath"

	"github.com/brohd11/golaunch/internal/app"

	"github.com/spf13/cobra"
)

// version is the binary version; defaults to "dev" for a plain `go build`. A makefile can stamp
// it via -X ldflags later, matching the sibling tools.
var version = "dev"

var rootCmd = &cobra.Command{
	Use:   "golaunch [dir]",
	Short: "Launch scripts against a file selection (TUI)",
	Long: `golaunch opens a TUI rooted at a directory. The Selection tab builds a set of paths
(the current dir, its child dirs, or its child files); the Scripts tab lists scripts scanned
from configured directories and runs the selected one against that selection.

  golaunch          # current directory
  golaunch /path    # an explicit root`,
	Version:       version,
	Args:          cobra.MaximumNArgs(1),
	SilenceUsage:  true,
	SilenceErrors: false,
	RunE:          runRoot,
}

func init() {
	rootCmd.SetVersionTemplate("golaunch {{.Version}}\n")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

// runRoot resolves the optional root directory (default: cwd) to an absolute path and launches
// the TUI.
func runRoot(cmd *cobra.Command, args []string) error {
	root := "."
	if len(args) > 0 {
		root = args[0]
	}
	abs, err := filepath.Abs(root)
	if err != nil {
		return err
	}
	return app.Run(abs, version)
}
