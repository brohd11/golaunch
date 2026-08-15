// Package config is golaunch's own user config, kept in ~/.golaunch/config.yml. Its one setting
// for now is script_dirs — the directories the Scripts tab scans for launchable scripts. A
// missing file is not an error: Ensure creates it on first run, seeds it with a default scripts
// directory (~/.golaunch/scripts), and materializes the bundled example scripts there so the tool
// works out of the box.
//
// It mirrors bubblestack/config's shape (a directory, read per call, no process-wide cache) but
// owns golaunch's own dotdir rather than the shared ~/.bubblestack one, which stays reserved for
// framework-level settings (the theme).
package config

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/brohd11/golaunch/internal/examples"

	"gopkg.in/yaml.v3"
)

// Config is the parsed ~/.golaunch/config.yml. Every field is optional; a missing file yields the
// zero value.
type Config struct {
	ScriptDirs []string `yaml:"script_dirs,omitempty"` // directories the Scripts tab scans
}

// Dir is ~/.golaunch, the home for config.yml and the default scripts directory.
func Dir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".golaunch"), nil
}

// Path is ~/.golaunch/config.yml.
func Path() (string, error) {
	dir, err := Dir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "config.yml"), nil
}

// ScriptsDir is ~/.golaunch/scripts, the default scan location seeded on first run.
func ScriptsDir() (string, error) {
	dir, err := Dir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "scripts"), nil
}

// Load reads ~/.golaunch/config.yml. A missing file is not an error — it returns the zero Config.
// A malformed file returns the parse error. Tilde-prefixed script dirs are expanded to the user's
// home so a hand-edited "~/scripts" resolves.
func Load() (*Config, error) {
	path, err := Path()
	if err != nil {
		return nil, err
	}
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return &Config{}, nil
		}
		return nil, err
	}
	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}
	for i, d := range cfg.ScriptDirs {
		cfg.ScriptDirs[i] = expandHome(d)
	}
	return &cfg, nil
}

// Ensure creates ~/.golaunch and its config.yml on first run, seeding config.yml with the default
// scripts directory and materializing the bundled example scripts into it. An already-present
// config.yml is left untouched (a returning user's edits survive). It returns whether it created
// the config (first run) so the caller can, e.g., surface a welcome.
func Ensure() (created bool, err error) {
	dir, err := Dir()
	if err != nil {
		return false, err
	}
	path := filepath.Join(dir, "config.yml")
	if _, err := os.Stat(path); err == nil {
		return false, nil // already configured
	} else if !os.IsNotExist(err) {
		return false, err
	}

	scriptsDir := filepath.Join(dir, "scripts")
	if err := os.MkdirAll(scriptsDir, 0o755); err != nil {
		return false, err
	}
	if err := examples.Materialize(scriptsDir); err != nil {
		return false, err
	}
	cfg := Config{ScriptDirs: []string{scriptsDir}}
	out, err := yaml.Marshal(&cfg)
	if err != nil {
		return false, err
	}
	if err := os.WriteFile(path, out, 0o644); err != nil {
		return false, err
	}
	return true, nil
}

// expandHome resolves a leading ~ (or ~/) to the user's home directory, leaving other paths
// unchanged. A home-lookup failure degrades to the original string rather than failing a scan.
func expandHome(p string) string {
	if p != "~" && !strings.HasPrefix(p, "~/") {
		return p
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return p
	}
	if p == "~" {
		return home
	}
	return filepath.Join(home, p[2:])
}
