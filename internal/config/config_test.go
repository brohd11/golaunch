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

func TestExpandHome(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	if got := expandHome("~/scripts"); got != filepath.Join(home, "scripts") {
		t.Errorf("expandHome(~/scripts) = %q", got)
	}
	if got := expandHome("/abs/path"); got != "/abs/path" {
		t.Errorf("expandHome left absolute path unchanged failed: %q", got)
	}
}
