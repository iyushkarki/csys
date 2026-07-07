package cleaners

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/shirou/gopsutil/v3/process"
)

const leftoverMinSize = 5 << 20

var leftoverScanDirs = []string{
	"~/Library/Application Support",
	"~/Library/Containers",
	"~/Library/Group Containers",
	"~/Library/HTTPStorages",
	"~/Library/Saved Application State",
	"~/Library/Logs",
}

var (
	bundleIDPattern = regexp.MustCompile(`^[a-zA-Z0-9-]+(\.[a-zA-Z0-9_-]+){2,}$`)
	teamIDPrefix    = regexp.MustCompile(`^[A-Z0-9]{10}\.`)
)

var protectedPrefixes = []string{"com.apple.", "group.com.apple."}

func LeftoverTargets() []*Target {
	if runtime.GOOS != "darwin" {
		return nil
	}
	installed := installedBundleIDs()
	if len(installed) == 0 {
		return nil
	}

	byID := map[string]*Target{}
	for _, dir := range leftoverScanDirs {
		base := expand(dir)
		entries, err := os.ReadDir(base)
		if err != nil {
			continue
		}
		for _, e := range entries {
			if !e.IsDir() {
				continue
			}
			id, ok := leftoverCandidate(e.Name(), installed)
			if !ok {
				continue
			}
			t := byID[id]
			if t == nil {
				t = &Target{
					ID:       "leftover:" + id,
					Name:     id,
					Tier:     Careful,
					detectFn: detectLeftover,
					Explain: Explain{
						After: "If you reinstall the app later it starts fresh. Recoverable from Trash.",
					},
				}
				byID[id] = t
			}
			t.paths = append(t.paths, filepath.Join(base, e.Name()))
		}
	}

	var targets []*Target
	for id, t := range byID {
		t.Explain.What = fmt.Sprintf("Data left behind by '%s' — no installed app matches it (%d locations).",
			id, len(t.paths))
		targets = append(targets, t)
	}
	return targets
}

func detectLeftover(t *Target) {
	id := strings.TrimPrefix(t.ID, "leftover:")
	if appInstalledAnywhere(id) {
		return
	}
	for _, p := range t.paths {
		t.Size += DirSize(p)
	}
	if t.Size < leftoverMinSize {
		t.Size = 0
		return
	}
	if token := lastToken(id); runningNames()[token] {
		t.Note = "⚠ a process named '" + token + "' is running"
	} else {
		t.Note = "app not installed"
	}
}

func appInstalledAnywhere(id string) bool {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	query := fmt.Sprintf("kMDItemCFBundleIdentifier ==[c] '%s'", id)
	out, err := exec.CommandContext(ctx, "mdfind", query).Output()
	if err != nil {
		return true
	}
	return strings.TrimSpace(string(out)) != ""
}

func leftoverCandidate(name string, installed map[string]bool) (string, bool) {
	id := strings.TrimSuffix(name, ".savedState")
	id = strings.TrimPrefix(id, "group.")
	id = teamIDPrefix.ReplaceAllString(id, "")
	if !bundleIDPattern.MatchString(id) {
		return "", false
	}
	lower := strings.ToLower(id)
	for _, prefix := range protectedPrefixes {
		if strings.HasPrefix(lower, prefix) {
			return "", false
		}
	}
	if installed[lower] {
		return "", false
	}
	for bid := range installed {
		if strings.HasPrefix(lower, bid+".") || strings.HasPrefix(bid, lower+".") {
			return "", false
		}
	}
	return id, true
}

func installedBundleIDs() map[string]bool {
	var apps []string
	for _, dir := range []string{"/Applications", "/Applications/Utilities", "/System/Applications", "/System/Applications/Utilities", expand("~/Applications")} {
		matches, _ := filepath.Glob(filepath.Join(dir, "*.app"))
		apps = append(apps, matches...)
	}
	if len(apps) == 0 {
		return nil
	}

	ids := make(map[string]bool)
	args := append([]string{"-name", "kMDItemCFBundleIdentifier"}, apps...)
	out, err := exec.Command("mdls", args...).Output()
	if err != nil {
		return nil
	}
	re := regexp.MustCompile(`kMDItemCFBundleIdentifier = "([^"]+)"`)
	for _, m := range re.FindAllStringSubmatch(string(out), -1) {
		ids[strings.ToLower(m[1])] = true
	}
	return ids
}

var runningNames = sync.OnceValue(func() map[string]bool {
	names := make(map[string]bool)
	procs, err := process.Processes()
	if err != nil {
		return names
	}
	for _, p := range procs {
		if name, err := p.Name(); err == nil {
			names[strings.ToLower(name)] = true
		}
	}
	return names
})

func lastToken(id string) string {
	parts := strings.Split(id, ".")
	return strings.ToLower(parts[len(parts)-1])
}
