package scripts

import (
	"os"
	"path/filepath"
	"testing"
)

func writeScript(t *testing.T, name, content string) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), name)
	if err := os.WriteFile(path, []byte(content), 0o755); err != nil {
		t.Fatal(err)
	}
	return path
}

func TestParseHeader(t *testing.T) {
	path := writeScript(t, "s.py", `#!/usr/bin/env python3
# name=Resize
# desc=Resize images
# path=Image/Filters
# terminal=true
# a plain comment, not a pair
import sys
# name=ShouldBeIgnored
`)
	m, err := ParseHeader(path)
	if err != nil {
		t.Fatal(err)
	}
	if m.Name != "Resize" {
		t.Errorf("Name = %q, want Resize", m.Name)
	}
	if m.Desc != "Resize images" {
		t.Errorf("Desc = %q", m.Desc)
	}
	if m.Path != "Image/Filters" {
		t.Errorf("Path = %q", m.Path)
	}
	if !m.Terminal {
		t.Error("Terminal = false, want true")
	}
}

func TestParseHeaderStopsAtCode(t *testing.T) {
	// The metadata after the first code line must not be read.
	path := writeScript(t, "s.sh", `#!/usr/bin/env bash
# name=First
echo hi
# name=Second
`)
	m, err := ParseHeader(path)
	if err != nil {
		t.Fatal(err)
	}
	if m.Name != "First" {
		t.Errorf("Name = %q, want First (metadata after code must be ignored)", m.Name)
	}
}

func TestParseHeaderNoMetadata(t *testing.T) {
	path := writeScript(t, "s.py", "print('hi')\n")
	m, err := ParseHeader(path)
	if err != nil {
		t.Fatal(err)
	}
	if m != (Meta{}) {
		t.Errorf("Meta = %+v, want zero", m)
	}
}

func TestParseHeaderTerminalDefaultsFalse(t *testing.T) {
	path := writeScript(t, "s.py", "# name=X\n# terminal=maybe\n")
	m, err := ParseHeader(path)
	if err != nil {
		t.Fatal(err)
	}
	if m.Terminal {
		t.Error("Terminal = true for a non-truthy value, want false")
	}
}
