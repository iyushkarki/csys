package cleaners

import (
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"syscall"
	"time"
)

const (
	preselectAge  = 7 * 24 * time.Hour
	detectWorkers = 8
)

func DetectStream(targets []*Target, out chan<- *Target) {
	sem := make(chan struct{}, detectWorkers)
	var wg sync.WaitGroup
	for _, t := range targets {
		wg.Add(1)
		go func(t *Target) {
			defer wg.Done()
			sem <- struct{}{}
			defer func() { <-sem }()
			t.detect()
			if t.Size > 0 {
				out <- t
			}
		}(t)
	}
	wg.Wait()
	close(out)
}

func DetectAll(targets []*Target) []*Target {
	ch := make(chan *Target)
	go DetectStream(targets, ch)

	var found []*Target
	for t := range ch {
		found = append(found, t)
	}
	Sort(found)
	return found
}

func Less(a, b *Target) bool {
	if a.Tier != b.Tier {
		return a.Tier < b.Tier
	}
	return a.Size > b.Size
}

func Sort(targets []*Target) {
	sort.SliceStable(targets, func(i, j int) bool {
		return Less(targets[i], targets[j])
	})
}

func (t *Target) detect() {
	switch {
	case t.detectFn != nil:
		t.detectFn(t)
	case len(t.globs) == 0:
		for _, p := range t.paths {
			t.Size += DirSize(p)
		}
	default:
		for _, g := range t.globs {
			matches, _ := filepath.Glob(expand(g))
			for _, m := range matches {
				if _, err := os.Lstat(m); err != nil {
					continue
				}
				t.paths = append(t.paths, m)
				t.Size += DirSize(m)
			}
		}
	}

	if t.Size == 0 {
		return
	}
	if t.LastUsed.IsZero() && len(t.paths) > 0 {
		t.LastUsed = lastUsed(t.paths...)
	}
	if t.Tier == Safe && time.Since(t.LastUsed) > preselectAge {
		t.Preselect = true
	}
}

func expand(p string) string {
	if strings.HasPrefix(p, "~/") {
		if home, err := os.UserHomeDir(); err == nil {
			return filepath.Join(home, p[2:])
		}
	}
	return p
}

func shortPath(p string) string {
	if home, err := os.UserHomeDir(); err == nil && strings.HasPrefix(p, home) {
		return "~" + p[len(home):]
	}
	return p
}

func commonDir(paths []string) string {
	if len(paths) == 0 {
		return ""
	}
	dir := filepath.Dir(paths[0])
	for _, p := range paths[1:] {
		for !strings.HasPrefix(p, dir+string(filepath.Separator)) && dir != "/" {
			dir = filepath.Dir(dir)
		}
	}
	return dir
}

type fileID struct {
	dev uint64
	ino uint64
}

func DirSize(root string) uint64 {
	var total uint64
	seen := make(map[fileID]struct{})
	_ = filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil || d == nil || !d.Type().IsRegular() {
			return nil
		}
		info, err := d.Info()
		if err != nil {
			return nil
		}
		if st, ok := info.Sys().(*syscall.Stat_t); ok {
			if st.Nlink > 1 {
				id := fileID{uint64(st.Dev), uint64(st.Ino)}
				if _, dup := seen[id]; dup {
					return nil
				}
				seen[id] = struct{}{}
			}
			total += uint64(st.Blocks) * 512
		} else {
			total += uint64(info.Size())
		}
		return nil
	})
	return total
}
