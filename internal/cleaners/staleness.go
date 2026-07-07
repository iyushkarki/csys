package cleaners

import (
	"context"
	"os"
	"os/exec"
	"strings"
	"sync"
	"time"
)

const mdlsTimeout = 500 * time.Millisecond

func lastUsed(paths ...string) time.Time {
	var newest time.Time
	for _, p := range paths {
		ts := spotlightLastUsed(p)
		if ts.IsZero() {
			if info, err := os.Stat(p); err == nil {
				ts = info.ModTime()
			}
		}
		if ts.After(newest) {
			newest = ts
		}
	}
	return newest
}

var mdlsAvailable = sync.OnceValue(func() bool {
	_, err := exec.LookPath("mdls")
	return err == nil
})

func spotlightLastUsed(path string) time.Time {
	if !mdlsAvailable() {
		return time.Time{}
	}
	ctx, cancel := context.WithTimeout(context.Background(), mdlsTimeout)
	defer cancel()

	out, err := exec.CommandContext(ctx, "mdls", "-name", "kMDItemLastUsedDate", "-raw", path).Output()
	if err != nil {
		return time.Time{}
	}
	raw := strings.TrimSpace(string(out))
	if raw == "" || raw == "(null)" {
		return time.Time{}
	}
	ts, err := time.Parse("2006-01-02 15:04:05 -0700", raw)
	if err != nil {
		return time.Time{}
	}
	return ts
}
