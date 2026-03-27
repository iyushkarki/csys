package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/iyushkarki/csys/internal/display"
	"github.com/iyushkarki/csys/internal/git"
	"github.com/spf13/cobra"
)

func branchCompletions(_ *cobra.Command, _ []string, _ string) ([]string, cobra.ShellCompDirective) {
	branches, err := git.GetAllBranches()
	if err != nil {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}
	return branches, cobra.ShellCompDirectiveNoFileComp
}

var gsyncCmd = &cobra.Command{
	Use:               "gsync [branch]",
	Short:             display.GSyncShort,
	Long:              display.GSyncLong,
	Args:              cobra.MaximumNArgs(1),
	ValidArgsFunction: branchCompletions,
	Run: func(cmd *cobra.Command, args []string) {
		if !git.IsGitRepo() {
			fmt.Fprintln(os.Stderr, "Error: not a git repository")
			return
		}

		branch := "main"
		if len(args) > 0 {
			branch = args[0]
		}

		fmt.Printf("Fetching origin/%s...\n", branch)
		if err := git.RunGit("fetch", "origin", branch); err != nil {
			fmt.Fprintf(os.Stderr, "Error fetching: %v\n", err)
			return
		}

		fmt.Printf("Resetting to origin/%s...\n", branch)
		if err := git.RunGit("reset", "--hard", "origin/"+branch); err != nil {
			fmt.Fprintf(os.Stderr, "Error resetting: %v\n", err)
			return
		}

		fmt.Printf("✓ Reset to origin/%s\n", branch)
	},
}

var gcleanCmd = &cobra.Command{
	Use:   "gclean",
	Short: display.GCleanShort,
	Long:  display.GCleanLong,
	Run: func(cmd *cobra.Command, args []string) {
		if !git.IsGitRepo() {
			fmt.Fprintln(os.Stderr, "Error: not a git repository")
			return
		}

		force, _ := cmd.Flags().GetBool("force")

		currentBranch, err := git.GetCurrentBranch()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error getting current branch: %v\n", err)
			return
		}

		branches, err := git.GetAllBranches()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error getting branches: %v\n", err)
			return
		}

		var toDelete []string
		for _, b := range branches {
			if b != currentBranch && b != "" {
				toDelete = append(toDelete, b)
			}
		}

		if len(toDelete) == 0 {
			fmt.Println("No branches to delete")
			return
		}

		fmt.Printf("Current branch: %s\n", currentBranch)
		fmt.Printf("Branches to delete:\n")
		for _, b := range toDelete {
			fmt.Printf("  - %s\n", b)
		}

		if !force {
			fmt.Print("\nDelete these branches? [y/N]: ")
			reader := bufio.NewReader(os.Stdin)
			text, _ := reader.ReadString('\n')
			if strings.ToLower(strings.TrimSpace(text)) != "y" {
				fmt.Println("Cancelled")
				return
			}
		}

		deleted := 0
		for _, b := range toDelete {
			if err := git.RunGit("branch", "-D", b); err == nil {
				deleted++
			}
		}

		fmt.Printf("✓ Deleted %d branches\n", deleted)
	},
}

var gsoftCmd = &cobra.Command{
	Use:   "gsoft [n]",
	Short: display.GSoftShort,
	Long:  display.GSoftLong,
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if !git.IsGitRepo() {
			fmt.Fprintln(os.Stderr, "Error: not a git repository")
			return
		}

		n := 1
		if len(args) > 0 {
			var err error
			n, err = strconv.Atoi(args[0])
			if err != nil || n < 1 {
				fmt.Fprintln(os.Stderr, "Error: n must be a positive integer")
				return
			}
		}

		ref := fmt.Sprintf("HEAD~%d", n)
		if err := git.RunGit("reset", "--soft", ref); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			return
		}

		fmt.Printf("✓ Soft reset %s\n", ref)
	},
}

var gacCmd = &cobra.Command{
	Use:   "gac <message>",
	Short: display.GAcShort,
	Long:  display.GAcLong,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if !git.IsGitRepo() {
			fmt.Fprintln(os.Stderr, "Error: not a git repository")
			return
		}

		message := args[0]

		if err := git.RunGit("add", "."); err != nil {
			fmt.Fprintf(os.Stderr, "Error staging: %v\n", err)
			return
		}

		if err := git.RunGit("commit", "-m", message); err != nil {
			fmt.Fprintf(os.Stderr, "Error committing: %v\n", err)
			return
		}

		fmt.Printf("✓ Committed: %s\n", message)
	},
}

var gcoCmd = &cobra.Command{
	Use:               "gco <branch>",
	Short:             display.GCoShort,
	Long:              display.GCoLong,
	Args:              cobra.ExactArgs(1),
	ValidArgsFunction: branchCompletions,
	Run: func(cmd *cobra.Command, args []string) {
		if !git.IsGitRepo() {
			fmt.Fprintln(os.Stderr, "Error: not a git repository")
			return
		}

		if err := git.RunGit("checkout", args[0]); err != nil {
			fmt.Fprintf(os.Stderr, "Error switching branch: %v\n", err)
			return
		}

		fmt.Printf("✓ Switched to %s\n", args[0])
	},
}

var gpushCmd = &cobra.Command{
	Use:   "gpush",
	Short: display.GPushShort,
	Long:  display.GPushLong,
	Run: func(cmd *cobra.Command, args []string) {
		if !git.IsGitRepo() {
			fmt.Fprintln(os.Stderr, "Error: not a git repository")
			return
		}

		if err := git.RunGit("push"); err != nil {
			fmt.Fprintf(os.Stderr, "Error pushing: %v\n", err)
			return
		}

		fmt.Println("✓ Pushed to remote")
	},
}

var gpullCmd = &cobra.Command{
	Use:   "gpull",
	Short: display.GPullShort,
	Long:  display.GPullLong,
	Run: func(cmd *cobra.Command, args []string) {
		if !git.IsGitRepo() {
			fmt.Fprintln(os.Stderr, "Error: not a git repository")
			return
		}

		if err := git.RunGit("pull"); err != nil {
			fmt.Fprintf(os.Stderr, "Error pulling: %v\n", err)
			return
		}

		fmt.Println("✓ Pulled latest changes")
	},
}

var gfpCmd = &cobra.Command{
	Use:   "gfp",
	Short: display.GFpShort,
	Long:  display.GFpLong,
	Run: func(cmd *cobra.Command, args []string) {
		if !git.IsGitRepo() {
			fmt.Fprintln(os.Stderr, "Error: not a git repository")
			return
		}

		if err := git.RunGit("push", "--force-with-lease"); err != nil {
			fmt.Fprintf(os.Stderr, "Error force pushing: %v\n", err)
			return
		}

		fmt.Println("✓ Force pushed (with lease)")
	},
}

var gundoCmd = &cobra.Command{
	Use:   "gundo",
	Short: display.GUndoShort,
	Long:  display.GUndoLong,
	Run: func(cmd *cobra.Command, args []string) {
		if !git.IsGitRepo() {
			fmt.Fprintln(os.Stderr, "Error: not a git repository")
			return
		}

		if err := git.RunGit("reset", "--soft", "HEAD~1"); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			return
		}

		fmt.Println("✓ Undid last commit (changes staged)")
	},
}

var gwipCmd = &cobra.Command{
	Use:   "gwip",
	Short: display.GWipShort,
	Long:  display.GWipLong,
	Run: func(cmd *cobra.Command, args []string) {
		if !git.IsGitRepo() {
			fmt.Fprintln(os.Stderr, "Error: not a git repository")
			return
		}

		if err := git.RunGit("add", "."); err != nil {
			fmt.Fprintf(os.Stderr, "Error staging: %v\n", err)
			return
		}

		if err := git.RunGit("commit", "-m", "WIP"); err != nil {
			fmt.Fprintf(os.Stderr, "Error committing: %v\n", err)
			return
		}

		fmt.Println("✓ WIP commit created")
	},
}

var gamendCmd = &cobra.Command{
	Use:   "gamend [message]",
	Short: display.GAmendShort,
	Long:  display.GAmendLong,
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if !git.IsGitRepo() {
			fmt.Fprintln(os.Stderr, "Error: not a git repository")
			return
		}

		if err := git.RunGit("add", "."); err != nil {
			fmt.Fprintf(os.Stderr, "Error staging: %v\n", err)
			return
		}

		if len(args) > 0 {
			if err := git.RunGit("commit", "--amend", "-m", args[0]); err != nil {
				fmt.Fprintf(os.Stderr, "Error amending: %v\n", err)
				return
			}
			fmt.Printf("✓ Amended last commit: %s\n", args[0])
		} else {
			if err := git.RunGit("commit", "--amend", "--no-edit"); err != nil {
				fmt.Fprintf(os.Stderr, "Error amending: %v\n", err)
				return
			}
			fmt.Println("✓ Amended last commit")
		}
	},
}

var grbCmd = &cobra.Command{
	Use:               "grb [base] [branch]",
	Short:             display.GRbShort,
	Long:              display.GRbLong,
	Args:              cobra.MaximumNArgs(2),
	ValidArgsFunction: branchCompletions,
	Run: func(cmd *cobra.Command, args []string) {
		if !git.IsGitInstalled() {
			fmt.Fprintln(os.Stderr, "Error: git is not installed")
			return
		}

		if !git.IsGitRepo() {
			fmt.Fprintln(os.Stderr, "Error: not a git repository")
			return
		}

		base := "main"
		var targetBranch string

		if len(args) == 1 {
			base = args[0]
		} else if len(args) == 2 {
			base = args[0]
			targetBranch = args[1]
		}

		currentBranch, err := git.GetCurrentBranch()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error getting current branch: %v\n", err)
			return
		}

		if targetBranch == "" && currentBranch == base {
			fmt.Fprintf(os.Stderr, "Error: cannot rebase %s onto itself\n", base)
			return
		}

		fmt.Printf("Fetching origin/%s...\n", base)
		if err := git.RunGit("fetch", "origin", base); err != nil {
			fmt.Fprintf(os.Stderr, "Error fetching: %v\n", err)
			return
		}

		if targetBranch != "" {
			fmt.Printf("Checking out %s...\n", targetBranch)
			if err := git.RunGit("checkout", targetBranch); err != nil {
				fmt.Fprintf(os.Stderr, "Error checking out: %v\n", err)
				return
			}
		}

		rebaseBranch := targetBranch
		if rebaseBranch == "" {
			rebaseBranch = currentBranch
		}

		fmt.Printf("Rebasing %s onto origin/%s...\n", rebaseBranch, base)
		if err := git.RunGit("rebase", "origin/"+base); err != nil {
			fmt.Fprintf(os.Stderr, "Error rebasing: %v\n", err)
			fmt.Fprintln(os.Stderr, "Resolve conflicts and run: git rebase --continue")
			return
		}

		fmt.Printf("✓ Rebased %s onto origin/%s\n", rebaseBranch, base)
	},
}

var glogCmd = &cobra.Command{
	Use:   "glog [n]",
	Short: display.GLogShort,
	Long:  display.GLogLong,
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if !git.IsGitRepo() {
			fmt.Fprintln(os.Stderr, "Error: not a git repository")
			return
		}

		n := 10
		if len(args) > 0 {
			parsed, err := strconv.Atoi(args[0])
			if err != nil || parsed < 1 {
				fmt.Fprintln(os.Stderr, "Error: n must be a positive integer")
				return
			}
			n = parsed
		}

		git.RunGit("log", "--oneline", "--graph", "--decorate", "-n", strconv.Itoa(n))
	},
}

var gstCmd = &cobra.Command{
	Use:   "gst",
	Short: display.GStShort,
	Long:  display.GStLong,
	Run: func(cmd *cobra.Command, args []string) {
		if !git.IsGitRepo() {
			fmt.Fprintln(os.Stderr, "Error: not a git repository")
			return
		}

		git.RunGit("status", "-sb")
	},
}

var gcbCmd = &cobra.Command{
	Use:   "gcb <branch>",
	Short: display.GCbShort,
	Long:  display.GCbLong,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if !git.IsGitRepo() {
			fmt.Fprintln(os.Stderr, "Error: not a git repository")
			return
		}

		branch := args[0]
		if err := git.RunGit("checkout", "-b", branch); err != nil {
			fmt.Fprintf(os.Stderr, "Error creating branch: %v\n", err)
			return
		}

		fmt.Printf("✓ Created and switched to %s\n", branch)
	},
}

var gbrnCmd = &cobra.Command{
	Use:   "gbrn <name>",
	Short: display.GBrnShort,
	Long:  display.GBrnLong,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if !git.IsGitRepo() {
			fmt.Fprintln(os.Stderr, "Error: not a git repository")
			return
		}

		name := args[0]
		oldName, err := git.GetCurrentBranch()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error getting current branch: %v\n", err)
			return
		}

		if err := git.RunGit("branch", "-M", name); err != nil {
			fmt.Fprintf(os.Stderr, "Error renaming branch: %v\n", err)
			return
		}

		fmt.Printf("✓ Renamed %s → %s\n", oldName, name)
	},
}

func init() {
	rootCmd.AddCommand(gsyncCmd)
	rootCmd.AddCommand(gcleanCmd)
	rootCmd.AddCommand(gsoftCmd)
	rootCmd.AddCommand(gacCmd)
	rootCmd.AddCommand(gcoCmd)
	rootCmd.AddCommand(gpushCmd)
	rootCmd.AddCommand(gpullCmd)
	rootCmd.AddCommand(gfpCmd)
	rootCmd.AddCommand(gundoCmd)
	rootCmd.AddCommand(gwipCmd)
	rootCmd.AddCommand(gamendCmd)
	rootCmd.AddCommand(grbCmd)
	rootCmd.AddCommand(glogCmd)
	rootCmd.AddCommand(gstCmd)
	rootCmd.AddCommand(gcbCmd)
	rootCmd.AddCommand(gbrnCmd)

	gcleanCmd.Flags().BoolP("force", "f", false, "Skip confirmation")
}
