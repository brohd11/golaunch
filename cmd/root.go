package cmd

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/brohd11/golaunch/internal/app"
	"github.com/brohd11/golaunch/internal/selection"

	"github.com/spf13/cobra"
)

// version is the binary version; defaults to "dev" for a plain `go build`. A makefile can stamp
// it via -X ldflags later, matching the sibling tools.
var version = "dev"

var rootPath string

var rootCmd = &cobra.Command{
	Use:   "golaunch [flags] [paths...]",
	Short: "Launch scripts against a file selection (TUI)",
	Long: `golaunch launches configured scripts against a file selection. With no paths it opens
the Selection tab to build a selection under the root. With paths it uses those paths as the
selection and opens the Scripts view directly; press R there to refine the supplied selection.

  golaunch                         # build a selection under the current directory
  golaunch file.txt photos/        # launch against those two paths
  golaunch --root /work file.txt   # use /work as the scripts' working directory`,
	Version:       version,
	Args:          cobra.ArbitraryArgs,
	SilenceUsage:  true,
	SilenceErrors: false,
	RunE:          runRoot,
}

func init() {
	rootCmd.SetVersionTemplate("golaunch {{.Version}}\n")
	rootCmd.Flags().StringVar(&rootPath, "root", ".", "script working directory and selection-building root")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

// runRoot resolves the working root and any argv selection to absolute paths, then launches the
// normal two-tab TUI or the preselected Scripts-only mode.
func runRoot(cmd *cobra.Command, args []string) error {
	abs, err := filepath.Abs(rootPath)
	if err != nil {
		return err
	}

	sel, preselected, err := resolveSelection(args, cmd.ErrOrStderr())
	if err != nil {
		return err
	}

	return app.Run(app.Options{
		Root:        abs,
		Version:     version,
		Selection:   sel,
		Preselected: preselected,
	})
}

// resolveSelection prepares argv paths without making one stale or unmounted Finder item prevent
// the valid remainder from opening. Supplying only invalid paths is still an error: preselected
// mode has no Build screen and would otherwise have no usable selection.
func resolveSelection(args []string, stderr io.Writer) (selection.Selection, bool, error) {
	sel, problems := selection.FromPaths(args)
	for _, problem := range problems {
		fmt.Fprintln(stderr, "skipping:", problem)
	}
	preselected := len(args) > 0
	if preselected && sel.Empty() {
		return sel, true, errors.New("no valid paths selected")
	}
	return sel, preselected, nil
}
