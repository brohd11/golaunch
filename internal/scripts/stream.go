package scripts

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"strings"
)

// reporter is the line sink a streaming run writes progress to — the same shape components.RunFunc
// hands its report closure, so a run can pass it straight through.
type reporter func(format string, args ...any)

// streamCmd runs argv (argv[0] is the program) in dir, relaying its stdout+stderr to report one
// line at a time as it arrives, and returns a non-nil error when the program exits non-zero (with
// its last line of output folded in, the part usually worth reading). ctx cancellation kills the
// subprocess, which is how the TaskScreen's esc-abort works. It's a domain-neutral port of
// gitstack's GitStream — no git specifics, any program.
func streamCmd(ctx context.Context, dir string, report reporter, argv ...string) error {
	if len(argv) == 0 {
		return fmt.Errorf("no command to run")
	}
	w := &lineWriter{report: report}
	cmd := exec.CommandContext(ctx, argv[0], argv[1:]...)
	cmd.Dir = dir
	cmd.Stdout = w
	cmd.Stderr = w

	err := cmd.Run()
	w.flush() // a final line with no trailing newline (an error message often has none)
	if err != nil {
		if w.last != "" {
			return fmt.Errorf("%w: %s", err, w.last)
		}
		return err
	}
	return nil
}

// lineWriter turns a subprocess's byte stream into whole lines for the reporter, breaking on \r as
// well as \n (progress bars are carriage-return delimited) and dropping empty lines.
type lineWriter struct {
	report reporter
	buf    []byte
	last   string
}

func (w *lineWriter) Write(p []byte) (int, error) {
	w.buf = append(w.buf, p...)
	for {
		i := bytes.IndexAny(w.buf, "\r\n")
		if i < 0 {
			break
		}
		w.emit(string(w.buf[:i]))
		w.buf = w.buf[i+1:]
	}
	return len(p), nil
}

func (w *lineWriter) flush() {
	if len(w.buf) > 0 {
		w.emit(string(w.buf))
		w.buf = nil
	}
}

func (w *lineWriter) emit(line string) {
	line = strings.TrimRight(line, " \t")
	if line == "" {
		return
	}
	w.last = line
	// "%s" rather than the line as a format string: output can contain a literal %.
	w.report("%s", line)
}
