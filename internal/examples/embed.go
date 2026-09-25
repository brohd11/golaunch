// Package examples ships golaunch's starter scripts and writes them out on first run.
package examples

import (
	"embed"
	"os"
	"path/filepath"
)

//go:embed scripts/*.py scripts/*.sh
var files embed.FS

// Materialize writes each bundled script into dir (0o755), never overwriting an existing one.
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
