// Package examples ships golaunch's starter scripts inside the binary and writes them out on
// first run, so a fresh install has something to launch. Each script carries a metadata header
// (see the scripts package) demonstrating the recognized keys — name/desc/path/terminal.
package examples

import (
	"embed"
	"os"
	"path/filepath"
)

//go:embed scripts/*.py scripts/*.sh
var files embed.FS

// Materialize writes every bundled script into dir (executable, 0o755), skipping any that already
// exist so a user's edits to a previously materialized script are never clobbered. dir is expected
// to exist (config.Ensure creates it).
func Materialize(dir string) error {
	entries, err := files.ReadDir("scripts")
	if err != nil {
		return err
	}
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		dst := filepath.Join(dir, e.Name())
		if _, err := os.Stat(dst); err == nil {
			continue // leave an existing script alone
		}
		data, err := files.ReadFile("scripts/" + e.Name())
		if err != nil {
			return err
		}
		if err := os.WriteFile(dst, data, 0o755); err != nil {
			return err
		}
	}
	return nil
}
