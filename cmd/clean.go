package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/iyushkarki/csys/internal/cleaners"
	"github.com/iyushkarki/csys/internal/display"
	"github.com/iyushkarki/csys/internal/system"
	"github.com/iyushkarki/csys/internal/tui"
	"github.com/mattn/go-isatty"
	"github.com/spf13/cobra"
)

var (
	cleanDryRun bool
	cleanSafe   bool
	cleanJSON   bool
	cleanNuke   bool
	cleanTrash  bool
	modulePaths []string
	moduleOlder string
)

var cleanCmd = &cobra.Command{
	Use:   "clean",
	Short: display.CleanShort,
	Long:  display.CleanLong,
	Run: func(cmd *cobra.Command, args []string) {
		runClean(func() []*cleaners.Target {
			return append(cleaners.Registry(), cleaners.LeftoverTargets()...)
		})
	},
}

var cleanModulesCmd = &cobra.Command{
	Use:   "modules",
	Short: display.CleanModulesShort,
	Long:  display.CleanModulesLong,
	Run: func(cmd *cobra.Command, args []string) {
		older, err := parseOlder(moduleOlder)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			return
		}
		roots := modulePaths
		if len(roots) == 0 {
			roots = cleaners.ModuleRoots()
		}
		runClean(func() []*cleaners.Target {
			var targets []*cleaners.Target
			for _, t := range cleaners.ModuleTargets(roots) {
				if older == 0 || time.Since(t.LastUsed) > older {
					targets = append(targets, t)
				}
			}
			return targets
		})
	},
}

func init() {
	rootCmd.AddCommand(cleanCmd)
	cleanCmd.AddCommand(cleanModulesCmd)

	for _, c := range []*cobra.Command{cleanCmd, cleanModulesCmd} {
		c.Flags().BoolVarP(&cleanDryRun, "dry-run", "n", false, "Show what is reclaimable without deleting anything")
		c.Flags().BoolVarP(&cleanSafe, "safe", "s", false, "Clean everything in the safe tier without prompting")
		c.Flags().BoolVar(&cleanJSON, "json", false, "Print detected targets as JSON")
		c.Flags().BoolVar(&cleanNuke, "nuke", false, "Delete permanently instead of moving careful items to Trash")
		c.Flags().BoolVar(&cleanTrash, "trash", false, "Move everything to Trash, even safe caches")
	}
	cleanModulesCmd.Flags().StringArrayVarP(&modulePaths, "path", "p", nil, "Root to scan (repeatable; default: ~/Documents, ~/dev, ~/code, ~/projects)")
	cleanModulesCmd.Flags().StringVar(&moduleOlder, "older", "", "Only projects untouched for this long (e.g. 30d, 8w, 6m)")
}

func runClean(produce func() []*cleaners.Target) {
	opts := cleaners.Options{Nuke: cleanNuke, Trash: cleanTrash}

	interactive := !cleanDryRun && !cleanSafe && !cleanJSON &&
		isatty.IsTerminal(os.Stdout.Fd()) && isatty.IsTerminal(os.Stdin.Fd())
	if interactive {
		if err := tui.Run(produce, opts); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		}
		return
	}

	targets := cleaners.DetectAll(produce())

	switch {
	case cleanJSON:
		printJSON(targets)
	case cleanSafe:
		cleanBatch(targets, opts)
	default:
		if len(targets) == 0 {
			fmt.Println(display.RenderNothingToClean())
			return
		}
		fmt.Println(display.RenderCleanList(targets))
	}
}

func cleanBatch(targets []*cleaners.Target, opts cleaners.Options) {
	var selected []*cleaners.Target
	for _, t := range targets {
		if t.Tier == cleaners.Safe {
			selected = append(selected, t)
		}
	}
	if len(selected) == 0 {
		fmt.Println(display.RenderNothingToClean())
		return
	}

	before, _ := system.GetDiskInfo()
	var freed uint64
	var entries []cleaners.HistoryEntry
	for _, t := range selected {
		err := t.Clean(opts)
		fmt.Println(display.FormatCleanProgress(t.Name, t.Size, err))
		if err == nil {
			freed += t.Size
			entries = append(entries, cleaners.NewHistoryEntry(t, opts))
		}
	}
	cleaners.AppendHistory(entries)
	after, _ := system.GetDiskInfo()
	fmt.Println(display.RenderCleanSummary(freed, before, after))
}

func printJSON(targets []*cleaners.Target) {
	type row struct {
		ID       string    `json:"id"`
		Name     string    `json:"name"`
		Tier     string    `json:"tier"`
		Size     uint64    `json:"size"`
		LastUsed time.Time `json:"lastUsed,omitzero"`
		Note     string    `json:"note,omitempty"`
		Paths    []string  `json:"paths,omitempty"`
	}
	rows := make([]row, 0, len(targets))
	for _, t := range targets {
		tier := "safe"
		if t.Tier == cleaners.Careful {
			tier = "careful"
		}
		rows = append(rows, row{
			ID: t.ID, Name: t.Name, Tier: tier, Size: t.Size,
			LastUsed: t.LastUsed, Note: t.Note, Paths: t.Paths(),
		})
	}
	out, _ := json.MarshalIndent(rows, "", "  ")
	fmt.Println(string(out))
}

func parseOlder(s string) (time.Duration, error) {
	if s == "" {
		return 0, nil
	}
	errInvalid := fmt.Errorf("invalid --older value %q (use e.g. 30d, 8w, 6m, 1y)", s)
	if len(s) < 2 {
		return 0, errInvalid
	}
	n, err := strconv.Atoi(s[:len(s)-1])
	if err != nil || n <= 0 {
		return 0, errInvalid
	}
	day := 24 * time.Hour
	switch s[len(s)-1] {
	case 'd':
		return time.Duration(n) * day, nil
	case 'w':
		return time.Duration(n) * 7 * day, nil
	case 'm':
		return time.Duration(n) * 30 * day, nil
	case 'y':
		return time.Duration(n) * 365 * day, nil
	}
	return 0, errInvalid
}
