// Package config is golaunch's user config, ~/.golaunch/config.yml. Ensure seeds it on
// first run with ~/.golaunch/scripts and the bundled example scripts.
package config

import (
	"os"
	"path/filepath"

	"github.com/brohd11/golaunch/internal/examples"

	"github.com/brohd11/goutil/configdir"
	"github.com/brohd11/goutil/strutil"
)

// Config is the parsed ~/.golaunch/config.yml.
type Config struct {
	ScriptDirs []string `yaml:"script_dirs,omitempty"` // directories the Scripts tab scans
}

const configName = "config.yml"

// Dir is ~/.golaunch.
func Dir() (string, error) { return configdir.Dir("golaunch") }

// Path is ~/.golaunch/config.yml.
func Path() (string, error) {
	dir, err := Dir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, configName), nil
}

// ScriptsDir is ~/.golaunch/scripts, the default scan location seeded on first run.
func ScriptsDir() (string, error) {
	dir, err := Dir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "scripts"), nil
}

// Load reads the config; a missing file yields the zero Config. "~" in script_dirs is
// expanded, and an entry that cannot be expanded is kept as typed for the scan to report.
func Load() (*Config, error) {
	path, err := Path()
	if err != nil {
		return nil, err
	}
	var cfg Config
	if err := configdir.Load(path, &cfg); err != nil {
		return nil, err
	}
	for i, d := range cfg.ScriptDirs {
		if exp, err := strutil.ExpandHome(d); err == nil {
			cfg.ScriptDirs[i] = exp
		}
	}
	return &cfg, nil
}

// Ensure seeds config.yml, the scripts directory and the example scripts on first run,
// and reports whether it did. An existing config is left untouched.
func Ensure() (created bool, err error) {
	path, err := Path()
	if err != nil {
		return false, err
	}
	if _, err := os.Stat(path); err == nil {
		return false, nil
	} else if !os.IsNotExist(err) {
		return false, err
	}
	scriptsDir, err := ScriptsDir()
	if err != nil {
		return false, err
	}
	if err := os.MkdirAll(scriptsDir, 0o755); err != nil {
		return false, err
	}
	if err := examples.Materialize(scriptsDir); err != nil {
		return false, err
	}
	// Written last, so a failed seed is retried on the next run.
	return true, configdir.SaveAtomic(filepath.Dir(path), configName, &Config{ScriptDirs: []string{scriptsDir}})
}
