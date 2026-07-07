package cleaners

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"strconv"
	"strings"

	"github.com/dustin/go-humanize"
)

func Registry() []*Target {
	targets := crossPlatformTargets()
	if runtime.GOOS == "darwin" {
		targets = append(targets, darwinTargets()...)
	}
	return targets
}

func crossPlatformTargets() []*Target {
	return []*Target{
		{ID: "bun", Name: "Bun install cache", Tier: Safe,
			globs: []string{"~/.bun/install/cache"},
			Explain: Explain{What: "Bun's downloaded package archive.",
				After: "Next `bun install` re-downloads what it needs. Nothing breaks."}},
		{ID: "npm", Name: "npm cache", Tier: Safe,
			globs: []string{"~/.npm/_cacache", "~/.npm/_npx"},
			Explain: Explain{What: "npm's downloaded package archive and npx cache.",
				After: "Next `npm install` re-downloads what it needs. Nothing breaks."}},
		{ID: "pnpm", Name: "pnpm cache", Tier: Safe,
			globs: []string{"~/Library/Caches/pnpm", "~/.cache/pnpm"},
			Explain: Explain{What: "pnpm's metadata and download cache (not the content store).",
				After: "pnpm refetches metadata on demand. Installed projects are untouched."}},
		{ID: "yarn", Name: "Yarn cache", Tier: Safe,
			globs: []string{"~/Library/Caches/Yarn", "~/.cache/yarn"},
			Explain: Explain{What: "Yarn's downloaded package archive.",
				After: "Next `yarn install` re-downloads what it needs."}},
		{ID: "pip", Name: "pip cache", Tier: Safe,
			globs: []string{"~/Library/Caches/pip", "~/.cache/pip"},
			Explain: Explain{What: "pip's downloaded wheels and HTTP cache.",
				After: "Next `pip install` re-downloads what it needs."}},
		{ID: "uv", Name: "uv cache", Tier: Safe,
			globs: []string{"~/.cache/uv", "~/Library/Caches/uv"},
			Explain: Explain{What: "uv's downloaded wheels and build cache.",
				After: "Next `uv` run re-downloads what it needs."}},
		{ID: "gobuild", Name: "Go build cache", Tier: Safe,
			globs: []string{"~/Library/Caches/go-build", "~/.cache/go-build"},
			Explain: Explain{What: "Compiled Go build artifacts.",
				After: "Next `go build` recompiles from scratch — slower once, then rebuilt."}},
		{ID: "gomod", Name: "Go module downloads", Tier: Safe,
			globs: []string{"~/go/pkg/mod/cache/download"},
			Explain: Explain{What: "Downloaded Go module archives.",
				After: "Next `go build` re-downloads missing modules."}},
		{ID: "cargo", Name: "Cargo registry cache", Tier: Safe,
			globs: []string{"~/.cargo/registry/cache", "~/.cargo/registry/src"},
			Explain: Explain{What: "Downloaded and extracted Rust crates.",
				After: "Next `cargo build` re-downloads what it needs."}},
		{ID: "brew", Name: "Homebrew cache", Tier: Safe,
			globs:     []string{"~/Library/Caches/Homebrew", "~/.cache/Homebrew"},
			postClean: []string{"brew", "cleanup", "--prune=all"},
			Explain: Explain{What: "Homebrew's downloads plus outdated formula versions.",
				After: "Brew re-downloads on demand. Installed packages are untouched."}},
		{ID: "gradle", Name: "Gradle caches", Tier: Safe,
			globs: []string{"~/.gradle/caches"},
			Explain: Explain{What: "Gradle's dependency and build caches.",
				After: "Next Gradle build re-downloads dependencies — slower once."}},
		{ID: "cocoapods", Name: "CocoaPods cache", Tier: Safe,
			globs: []string{"~/Library/Caches/CocoaPods"},
			Explain: Explain{What: "Downloaded pod specs and archives.",
				After: "Next `pod install` re-downloads what it needs."}},
		{ID: "composer", Name: "Composer cache", Tier: Safe,
			globs: []string{"~/.composer/cache", "~/.cache/composer"},
			Explain: Explain{What: "Composer's downloaded PHP package archive.",
				After: "Next `composer install` re-downloads what it needs."}},
		{ID: "gem", Name: "RubyGems cache", Tier: Safe,
			globs: []string{"~/.gem/ruby/*/cache"},
			Explain: Explain{What: "Downloaded gem archives.",
				After: "Next `gem install` re-downloads what it needs."}},
		{ID: "nodegyp", Name: "node-gyp cache", Tier: Safe,
			globs: []string{"~/Library/Caches/node-gyp", "~/.node-gyp"},
			Explain: Explain{What: "Node.js headers used to compile native addons.",
				After: "Re-downloaded automatically on the next native module build."}},
		{ID: "swiftpm", Name: "SwiftPM cache", Tier: Safe,
			globs: []string{"~/Library/Caches/org.swift.swiftpm", "~/.cache/org.swift.swiftpm"},
			Explain: Explain{What: "Swift Package Manager's checkout and manifest cache.",
				After: "Xcode/SwiftPM re-resolves packages on next build."}},
		{ID: "electron", Name: "Electron cache", Tier: Safe,
			globs: []string{"~/Library/Caches/electron", "~/.cache/electron", "~/.electron"},
			Explain: Explain{What: "Downloaded Electron binaries used by build tools.",
				After: "Re-downloaded on the next Electron app build."}},
		{ID: "jetbrains", Name: "JetBrains logs", Tier: Safe,
			globs: []string{"~/Library/Logs/JetBrains", "~/.cache/JetBrains"},
			Explain: Explain{What: "IDE log files from PhpStorm, PyCharm, IntelliJ, etc.",
				After: "IDEs start fresh log files. Settings and projects are untouched."}},
		{ID: "hprof", Name: "Java heap dumps", Tier: Safe,
			globs: []string{"~/*.hprof", "~/java_error_in_*.log"},
			Explain: Explain{What: "Crash heap dumps Java apps left in your home folder.",
				After: "Gone. Only needed if you were actively debugging that crash."}},
		{ID: "android", Name: "Android cache", Tier: Safe,
			globs: []string{"~/.android/cache"},
			Explain: Explain{What: "Android tooling's download cache.",
				After: "Re-downloaded on demand by the Android SDK."}},
		dockerPrune(),
		{ID: "playwright", Name: "Playwright browsers", Tier: Careful,
			Note:  "re-download with `npx playwright install`",
			globs: []string{"~/Library/Caches/ms-playwright*", "~/.cache/ms-playwright*"},
			Explain: Explain{What: "Browser builds Playwright downloaded for testing.",
				After: "Playwright tests fail until you run `npx playwright install`."}},
		{ID: "trash", Name: "Trash", Tier: Careful,
			Note:      "permanently deletes trashed files",
			globs:     []string{"~/.Trash", "~/.local/share/Trash/files"},
			keepDir:   true,
			permanent: true,
			Explain: Explain{What: "Everything currently sitting in your Trash.",
				After: "Gone for good — this is the one that can't be undone."}},
	}
}

func darwinTargets() []*Target {
	return []*Target{
		{ID: "updaters", Name: "App updater leftovers", Tier: Safe,
			globs: []string{"~/Library/Caches/*ShipIt*", "~/Library/Caches/*-updater"},
			Explain: Explain{What: "Downloaded update payloads apps forgot to remove.",
				After: "Nothing — the updates are already installed."}},
		{ID: "derived", Name: "Xcode DerivedData", Tier: Safe,
			globs: []string{"~/Library/Developer/Xcode/DerivedData"},
			Explain: Explain{What: "Xcode's build products, indexes, and logs.",
				After: "Next build recompiles and re-indexes — slower once."}},
		{ID: "previews", Name: "Xcode Previews cache", Tier: Safe,
			globs: []string{"~/Library/Developer/Xcode/UserData/Previews"},
			Explain: Explain{What: "Simulators and builds backing SwiftUI previews.",
				After: "Previews rebuild automatically the next time you use them."}},
		{ID: "ibsupport", Name: "Xcode IB Support", Tier: Safe,
			globs: []string{"~/Library/Developer/Xcode/UserData/IB Support", "~/Library/Developer/Xcode/UserData/IB%20Support"},
			Explain: Explain{What: "Interface Builder's simulator caches.",
				After: "Rebuilt automatically when you next edit a storyboard/xib."}},
		{ID: "xcodedocs", Name: "Xcode documentation cache", Tier: Safe,
			globs: []string{"~/Library/Developer/Xcode/DocumentationCache"},
			Explain: Explain{What: "Rendered documentation Xcode caches locally.",
				After: "Re-rendered on demand when you browse docs."}},
		{ID: "simcaches", Name: "Simulator dyld caches", Tier: Safe,
			globs: []string{"~/Library/Developer/CoreSimulator/Caches"},
			Explain: Explain{What: "Simulator runtime link caches.",
				After: "Rebuilt on next simulator boot — first boot is slower."}},
		simUnavailable(),
		deviceSupport(),
		oldRuntimes(),
		{ID: "archives", Name: "Xcode Archives", Tier: Careful,
			Note:  "needed to symbolicate crash logs of shipped builds",
			globs: []string{"~/Library/Developer/Xcode/Archives"},
			Explain: Explain{What: "Archived app builds with their debug symbols.",
				After: "You can't re-export or symbolicate those shipped builds anymore."}},
	}
}

func simUnavailable() *Target {
	t := &Target{ID: "sims", Name: "Orphaned iOS simulators", Tier: Safe,
		Note: "runs `xcrun simctl delete unavailable`",
		Explain: Explain{What: "Simulator devices whose runtime is no longer installed.",
			After: "Nothing — they couldn't boot anyway."}}
	t.detectFn = func(t *Target) {
		out, err := xcrunOutput("simctl", "list", "devices", "-j")
		if err != nil {
			return
		}
		var parsed struct {
			Devices map[string][]struct {
				IsAvailable bool   `json:"isAvailable"`
				DataPath    string `json:"dataPath"`
			} `json:"devices"`
		}
		if json.Unmarshal(out, &parsed) != nil {
			return
		}
		for _, devs := range parsed.Devices {
			for _, d := range devs {
				if !d.IsAvailable && d.DataPath != "" {
					t.Size += DirSize(filepath.Dir(d.DataPath))
				}
			}
		}
	}
	t.cleanFn = func(*Target, Options) error {
		return exec.Command("xcrun", "simctl", "delete", "unavailable").Run()
	}
	return t
}

func deviceSupport() *Target {
	t := &Target{ID: "devicesupport", Name: "Old device symbols", Tier: Safe,
		Note: "keeps the newest version per platform",
		Explain: Explain{What: "Debug symbols Xcode copied from devices on old OS versions.",
			After: "Re-copied from the device if you ever debug on that OS again."}}
	t.detectFn = func(t *Target) {
		base := expand("~/Library/Developer/Xcode")
		platforms, _ := filepath.Glob(filepath.Join(base, "* DeviceSupport"))
		for _, platform := range platforms {
			entries, err := os.ReadDir(platform)
			if err != nil {
				continue
			}
			var names []string
			for _, e := range entries {
				if e.IsDir() {
					names = append(names, e.Name())
				}
			}
			if len(names) < 2 {
				continue
			}
			newest := names[0]
			for _, name := range names[1:] {
				if versionLess(newest, name) {
					newest = name
				}
			}
			for _, name := range names {
				if name == newest {
					continue
				}
				p := filepath.Join(platform, name)
				t.paths = append(t.paths, p)
				t.Size += DirSize(p)
			}
		}
	}
	return t
}

func oldRuntimes() *Target {
	t := &Target{ID: "runtimes", Name: "Old simulator runtimes", Tier: Careful,
		Note: "keeps the newest runtime per platform; re-download is multi-GB",
		Explain: Explain{What: "Downloadable simulator OS images superseded by newer ones.",
			After: "Simulators on those OS versions need a multi-GB re-download."}}
	t.detectFn = func(t *Target) {
		out, err := xcrunOutput("simctl", "runtime", "list", "-j")
		if err != nil {
			return
		}
		var parsed map[string]struct {
			Platform  string `json:"platformIdentifier"`
			Version   string `json:"version"`
			SizeBytes uint64 `json:"sizeBytes"`
			State     string `json:"state"`
		}
		if json.Unmarshal(out, &parsed) != nil {
			return
		}
		newest := map[string]string{}
		for _, r := range parsed {
			if newest[r.Platform] == "" || versionLess(newest[r.Platform], r.Version) {
				newest[r.Platform] = r.Version
			}
		}
		for id, r := range parsed {
			if r.State == "Ready" && versionLess(r.Version, newest[r.Platform]) {
				t.paths = append(t.paths, id)
				t.Size += r.SizeBytes
			}
		}
	}
	t.cleanFn = func(t *Target, _ Options) error {
		var firstErr error
		for _, id := range t.paths {
			if err := exec.Command("xcrun", "simctl", "runtime", "delete", id).Run(); err != nil && firstErr == nil {
				firstErr = err
			}
		}
		return firstErr
	}
	return t
}

func dockerPrune() *Target {
	t := &Target{ID: "docker", Name: "Docker reclaimable", Tier: Careful,
		Note: "runs `docker system prune -f`: stopped containers, dangling images, unused networks",
		Explain: Explain{What: "Stopped containers, dangling images, and build cache Docker reports as reclaimable.",
			After: "Removed images re-pull on next use. Running containers are untouched."}}
	t.detectFn = func(t *Target) {
		docker, err := exec.LookPath("docker")
		if err != nil {
			return
		}
		out, err := exec.Command(docker, "system", "df", "--format", "{{.Reclaimable}}").Output()
		if err != nil {
			return
		}
		for _, line := range strings.Split(strings.TrimSpace(string(out)), "\n") {
			size := strings.Fields(line)
			if len(size) == 0 {
				continue
			}
			if b, err := humanize.ParseBytes(size[0]); err == nil {
				t.Size += b
			}
		}
	}
	t.cleanFn = func(*Target, Options) error {
		return exec.Command("docker", "system", "prune", "-f").Run()
	}
	return t
}

var versionPattern = regexp.MustCompile(`\d+(?:\.\d+)*`)

func versionLess(a, b string) bool {
	as := strings.Split(versionPattern.FindString(a), ".")
	bs := strings.Split(versionPattern.FindString(b), ".")
	for i := 0; i < len(as) || i < len(bs); i++ {
		var av, bv int
		if i < len(as) {
			av, _ = strconv.Atoi(as[i])
		}
		if i < len(bs) {
			bv, _ = strconv.Atoi(bs[i])
		}
		if av != bv {
			return av < bv
		}
	}
	return false
}

func xcrunOutput(args ...string) ([]byte, error) {
	if _, err := exec.LookPath("xcrun"); err != nil {
		return nil, err
	}
	return exec.Command("xcrun", args...).Output()
}
