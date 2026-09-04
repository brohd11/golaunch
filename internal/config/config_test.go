package config

import (
	"os"
	"path/filepath"
	"testing"
)

// TestEnsureFirstRun verifies the first-run flow in an isolated HOME: config.yml is written with
// the default scripts dir, and the bundled example scripts are materialized into it. A second call
// is a no-op (returns created=false) and leaves user files alone.
func TestEnsureFirstRun(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv("USERPROFILE", os.Getenv("HOME"))

	created, err := Ensure()
	if err != nil {
		t.Fatal(err)
	}
	if !created {
		t.Fatal("Ensure reported not-created on first run")
	}

	cfgPath := filepath.Join(home, ".golaunch", "config.yml")
	if _, err := os.Stat(cfgPath); err != nil {
		t.Fatalf("config.yml not written: %v", err)
	}

	cfg, err := Load()
	if err != nil {
		t.Fatal(err)
	}
	if len(cfg.ScriptDirs) != 1 {
		t.Fatalf("ScriptDirs = %v, want one default", cfg.ScriptDirs)
	}
	entries, err := os.ReadDir(cfg.ScriptDirs[0])
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) == 0 {
		t.Error("no example scripts materialized")
	}

	// Second run must be a no-op.
	created, err = Ensure()
	if err != nil {
		t.Fatal(err)
	}
	if created {
		t.Error("Ensure reported created on a second run")
	}
}

// TestLoadExpandsScriptDirs pins what Load owes its callers about hand-edited script_dirs:
// a tilde entry resolves against HOME, an absolute one is passed through, and an entry the
// shared expander refuses ("~other-user") survives as typed rather than being dropped from
// the scan. The tilde mechanics themselves belong to goutil/strutil and are tested there.
func TestLoadExpandsScriptDirs(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv("USERPROFILE", os.Getenv("HOME"))

	dir := filepath.Join(home, ".golaunch")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	body := "script_dirs:\n  - ~/scripts\n  - /abs/path\n  - ~other/scripts\n"
	if err := os.WriteFile(filepath.Join(dir, "config.yml"), []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}

	cfg, err := Load()
	if err != nil {
		t.Fatal(err)
	}
	want := []string{filepath.Join(home, "scripts"), "/abs/path", "~other/scripts"}
	if len(cfg.ScriptDirs) != len(want) {
		t.Fatalf("ScriptDirs = %v, want %v", cfg.ScriptDirs, want)
	}
	for i, w := range want {
		if cfg.ScriptDirs[i] != w {
			t.Errorf("ScriptDirs[%d] = %q, want %q", i, cfg.ScriptDirs[i], w)
		}
	}
}

// TestLoadMissingFile covers first-run-before-Ensure: no config.yml is not an error, and the
// zero Config (no script dirs) is what the Scripts tab scans with.
func TestLoadMissingFile(t *testing.T) {
	t.Setenv("HOME", t.TempDir())
	t.Setenv("USERPROFILE", os.Getenv("HOME"))
	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load on a missing file: %v", err)
	}
	if len(cfg.ScriptDirs) != 0 {
		t.Errorf("ScriptDirs = %v, want empty", cfg.ScriptDirs)
	}
}
