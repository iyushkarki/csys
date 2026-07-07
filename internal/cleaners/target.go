package cleaners

import (
	"fmt"
	"strings"
	"time"
)

type Tier int

const (
	Safe Tier = iota
	Careful
)

type Explain struct {
	What  string
	After string
}

type Target struct {
	ID        string
	Name      string
	Note      string
	Tier      Tier
	Size      uint64
	LastUsed  time.Time
	Preselect bool
	Explain   Explain

	globs     []string
	paths     []string
	keepDir   bool
	permanent bool
	postClean []string
	detectFn  func(*Target)
	cleanFn   func(*Target, Options) error
}

type Options struct {
	Nuke  bool
	Trash bool
}

func (t *Target) Paths() []string { return t.paths }

func (t *Target) UsesTrash(opts Options) bool {
	if t.cleanFn != nil || len(t.postClean) > 0 || t.permanent || opts.Nuke {
		return false
	}
	return t.Tier == Careful || opts.Trash
}

func (t *Target) Action(opts Options) string {
	if t.cleanFn != nil {
		return t.Note
	}
	if len(t.postClean) > 0 {
		return fmt.Sprintf("delete cache, then run `%s`", strings.Join(t.postClean, " "))
	}
	verb := "delete permanently"
	if t.UsesTrash(opts) {
		verb = "move to Trash (recoverable)"
	}
	if len(t.paths) == 1 {
		return fmt.Sprintf("%s: %s", verb, shortPath(t.paths[0]))
	}
	return fmt.Sprintf("%s: %d locations under %s", verb, len(t.paths), shortPath(commonDir(t.paths)))
}

func Ago(ts time.Time) string {
	if ts.IsZero() {
		return "—"
	}
	d := time.Since(ts)
	switch {
	case d < time.Hour:
		return "just now"
	case d < 24*time.Hour:
		return fmt.Sprintf("%dh ago", int(d.Hours()))
	case d < 7*24*time.Hour:
		return fmt.Sprintf("%dd ago", int(d.Hours()/24))
	case d < 30*24*time.Hour:
		return fmt.Sprintf("%dw ago", int(d.Hours()/(24*7)))
	case d < 365*24*time.Hour:
		return fmt.Sprintf("%dmo ago", int(d.Hours()/(24*30)))
	default:
		return fmt.Sprintf("%dy ago", int(d.Hours()/(24*365)))
	}
}
