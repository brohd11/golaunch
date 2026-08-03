package scripts

import "testing"

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
