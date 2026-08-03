package scripts

import (
	"bufio"
	"os"
	"strings"
)

// Meta is the launch metadata a script declares in its comment header. Every field is optional:
// an empty Name falls back to the filename, an empty Path puts the script at the top level, and
// Terminal defaults to false (stream in the TUI).
type Meta struct {
	Name     string // display name in the menu
	Desc     string // one-line description under the name
	Path     string // "/"-separated submenu route, e.g. "Image/Filters"
	Terminal bool   // true ⇒ launch in an external terminal instead of streaming in the TUI
}

// maxHeaderLines caps how far ParseHeader reads before giving up — the metadata is meant to sit at
// the very top, so a script with no header costs only a few lines of scanning.
const maxHeaderLines = 20

// ParseHeader reads the leading comment header of a script for its metadata. It skips a shebang,
// treats each subsequent "# key=value" comment line as a metadata pair (recognized keys only),
// and stops at the first line that is neither blank nor a comment (the first line of real code).
// A missing/unreadable file returns the zero Meta and the error; unknown keys are ignored.
func ParseHeader(path string) (Meta, error) {
	f, err := os.Open(path)
	if err != nil {
		return Meta{}, err
	}
	defer f.Close()

	var m Meta
	sc := bufio.NewScanner(f)
	for i := 0; i < maxHeaderLines && sc.Scan(); i++ {
		line := strings.TrimSpace(sc.Text())
		if i == 0 && strings.HasPrefix(line, "#!") {
			continue // shebang
		}
		if line == "" {
			continue // blank lines inside the header don't end it
		}
		if !strings.HasPrefix(line, "#") {
			break // first real code line — the header is over
		}
		key, val, ok := strings.Cut(strings.TrimSpace(strings.TrimPrefix(line, "#")), "=")
		if !ok {
			continue // a plain comment, not a key=value pair
		}
		applyMeta(&m, strings.TrimSpace(key), strings.TrimSpace(val))
	}
	if err := sc.Err(); err != nil {
		return Meta{}, err
	}
	return m, nil
}

// applyMeta sets one recognized key on m; unknown keys are ignored so a header can carry
// script-specific comments alongside golaunch's own.
func applyMeta(m *Meta, key, val string) {
	switch strings.ToLower(key) {
	case "name":
		m.Name = val
	case "desc", "description":
		m.Desc = val
	case "path":
		m.Path = strings.Trim(val, "/")
	case "terminal":
		m.Terminal = isTrue(val)
	}
}

// isTrue reads a boolean-ish metadata value; anything other than a recognized truthy token is
// false, so a malformed value fails safe to in-TUI streaming.
func isTrue(v string) bool {
	switch strings.ToLower(v) {
	case "true", "1", "yes", "y", "on":
		return true
	}
	return false
}
