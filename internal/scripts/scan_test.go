package scripts

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"
)

func TestBuildTreeNesting(t *testing.T) {
	scripts := []Script{
		{File: "/a/resize.py", Meta: Meta{Name: "Resize", Path: "Image/Filters"}},
		{File: "/a/crop.py", Meta: Meta{Name: "Crop", Path: "Image"}},
		{File: "/a/top.sh", Meta: Meta{Name: "Top"}}, // no path ⇒ root
	}
	root := BuildTree(scripts)

	if len(root.Scripts) != 1 || root.Scripts[0].Meta.Name != "Top" {
		t.Fatalf("root scripts = %+v, want [Top]", root.Scripts)
	}
	if len(root.Children) != 1 || root.Children[0].Name != "Image" {
		t.Fatalf("root children = %+v, want [Image]", root.Children)
	}
	image := root.Children[0]
	if len(image.Scripts) != 1 || image.Scripts[0].Meta.Name != "Crop" {
		t.Errorf("Image scripts = %+v, want [Crop]", image.Scripts)
	}
	if len(image.Children) != 1 || image.Children[0].Name != "Filters" {
		t.Fatalf("Image children = %+v, want [Filters]", image.Children)
	}
	filters := image.Children[0]
	if len(filters.Scripts) != 1 || filters.Scripts[0].Meta.Name != "Resize" {
		t.Errorf("Filters scripts = %+v, want [Resize]", filters.Scripts)
	}
}

func TestDisplayNameFallsBackToFilename(t *testing.T) {
	s := Script{File: "/a/my_tool.py"}
	if got := s.DisplayName(); got != "my_tool" {
		t.Errorf("DisplayName = %q, want my_tool", got)
	}
	s.Meta.Name = "Pretty"
	if got := s.DisplayName(); got != "Pretty" {
		t.Errorf("DisplayName = %q, want Pretty", got)
	}
}

func TestScanReportsProblems(t *testing.T) {
	dir := t.TempDir()
	good := filepath.Join(dir, "good.sh")
	if err := os.WriteFile(good, []byte("# name=Good\n"), 0o755); err != nil {
		t.Fatal(err)
	}
	bad := filepath.Join(dir, "bad.sh")
	// A token larger than bufio.Scanner's limit makes header parsing fail on
	// every platform. chmod(000) is not a reliable unreadable-file fixture on
	// Windows, where the owner can still read it.
	if err := os.WriteFile(bad, bytes.Repeat([]byte("x"), 128*1024), 0o644); err != nil {
		t.Fatal(err)
	}
	missing := filepath.Join(dir, "no-such-dir")

	got, problems := Scan([]string{dir, missing})
	if len(got) != 1 || got[0].Meta.Name != "Good" {
		t.Errorf("scripts = %+v, want just Good", got)
	}
	if len(problems) != 2 {
		t.Fatalf("problems = %v, want 2 (invalid header + missing dir)", problems)
	}
}
