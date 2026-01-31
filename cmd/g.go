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

var gCmd = &cobra.Command{
	Use:   "g",
	Short: display.GShort,
	Long:  display.GLong,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

var syncCmd = &cobra.Command{
	Use:   "sync [branch]",
	Short: display.GSyncShort,
	Long:  display.GSyncLong,
	Args:  cobra.MaximumNArgs(1),
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
		err := git.RunGit("fetch", "origin", branch)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error fetching: %v\n", err)
			return
		}

		fmt.Printf("Resetting to origin/%s...\n", branch)
		err = git.RunGit("reset", "--hard", "origin/"+branch)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error resetting: %v\n", err)
			return
		}

		fmt.Printf("✓ Reset to origin/%s\n", branch)
	},
}

var cleanCmd = &cobra.Command{
	Use:   "clean",
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
			err := git.RunGit("branch", "-D", b)
			if err == nil {
				deleted++
			}
		}

		fmt.Printf("✓ Deleted %d branches\n", deleted)
	},
}

var softCmd = &cobra.Command{
	Use:   "soft [n]",
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
		err := git.RunGit("reset", "--soft", ref)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			return
		}

		fmt.Printf("✓ Soft reset %s\n", ref)
	},
}

var acCmd = &cobra.Command{
	Use:   "ac <message>",
	Short: display.GAcShort,
	Long:  display.GAcLong,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if !git.IsGitRepo() {
			fmt.Fprintln(os.Stderr, "Error: not a git repository")
			return
		}

		message := args[0]

		err := git.RunGit("add", ".")
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error staging: %v\n", err)
			return
		}

		err = git.RunGit("commit", "-m", message)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error committing: %v\n", err)
			return
		}

		fmt.Printf("✓ Committed: %s\n", message)
	},
}

var acpCmd = &cobra.Command{
	Use:   "acp <message>",
	Short: display.GAcpShort,
	Long:  display.GAcpLong,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if !git.IsGitRepo() {
			fmt.Fprintln(os.Stderr, "Error: not a git repository")
			return
		}

		message := args[0]

		err := git.RunGit("add", ".")
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error staging: %v\n", err)
			return
		}

		err = git.RunGit("commit", "-m", message)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error committing: %v\n", err)
			return
		}

		err = git.RunGit("push")
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error pushing: %v\n", err)
			return
		}

		fmt.Printf("✓ Committed and pushed: %s\n", message)
	},
}

var undoCmd = &cobra.Command{
	Use:   "undo",
	Short: display.GUndoShort,
	Long:  display.GUndoLong,
	Run: func(cmd *cobra.Command, args []string) {
		if !git.IsGitRepo() {
			fmt.Fprintln(os.Stderr, "Error: not a git repository")
			return
		}

		err := git.RunGit("reset", "--soft", "HEAD~1")
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			return
		}

		fmt.Println("✓ Undid last commit (changes staged)")
	},
}

var wipCmd = &cobra.Command{
	Use:   "wip",
	Short: display.GWipShort,
	Long:  display.GWipLong,
	Run: func(cmd *cobra.Command, args []string) {
		if !git.IsGitRepo() {
			fmt.Fprintln(os.Stderr, "Error: not a git repository")
			return
		}

		err := git.RunGit("add", ".")
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error staging: %v\n", err)
			return
		}

		err = git.RunGit("commit", "-m", "WIP")
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error committing: %v\n", err)
			return
		}

		fmt.Println("✓ WIP commit created")
	},
}

var amendCmd = &cobra.Command{
	Use:   "amend",
	Short: display.GAmendShort,
	Long:  display.GAmendLong,
	Run: func(cmd *cobra.Command, args []string) {
		if !git.IsGitRepo() {
			fmt.Fprintln(os.Stderr, "Error: not a git repository")
			return
		}

		err := git.RunGit("add", ".")
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error staging: %v\n", err)
			return
		}

		err = git.RunGit("commit", "--amend", "--no-edit")
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error amending: %v\n", err)
			return
		}

		fmt.Println("✓ Amended last commit")
	},
}

var rbCmd = &cobra.Command{
	Use:   "rb [base] [branch]",
	Short: display.GRbShort,
	Long:  display.GRbLong,
	Args:  cobra.MaximumNArgs(2),
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
		err = git.RunGit("fetch", "origin", base)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error fetching: %v\n", err)
			return
		}

		if targetBranch != "" {
			fmt.Printf("Checking out %s...\n", targetBranch)
			err = git.RunGit("checkout", targetBranch)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error checking out: %v\n", err)
				return
			}
		}

		rebaseBranch := targetBranch
		if rebaseBranch == "" {
			rebaseBranch = currentBranch
		}

		fmt.Printf("Rebasing %s onto origin/%s...\n", rebaseBranch, base)
		err = git.RunGit("rebase", "origin/"+base)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error rebasing: %v\n", err)
			fmt.Fprintln(os.Stderr, "Resolve conflicts and run: git rebase --continue")
			return
		}

		fmt.Printf("✓ Rebased %s onto origin/%s\n", rebaseBranch, base)
	},
}

var logCmd = &cobra.Command{
	Use:   "log [n]",
	Short: display.GLogShort,
	Long:  display.GLogLong,
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if !git.IsGitRepo() {
			fmt.Fprintln(os.Stderr, "Error: not a git repository")
			return
		}

		n := "10"
		if len(args) > 0 {
			n = args[0]
		}

		git.RunGit("log", "--oneline", "--graph", "--decorate", "-n", n)
	},
}

var stCmd = &cobra.Command{
	Use:   "st",
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

func init() {
	rootCmd.AddCommand(gCmd)
	gCmd.AddCommand(syncCmd)
	gCmd.AddCommand(cleanCmd)
	gCmd.AddCommand(softCmd)
	gCmd.AddCommand(acCmd)
	gCmd.AddCommand(acpCmd)
	gCmd.AddCommand(undoCmd)
	gCmd.AddCommand(wipCmd)
	gCmd.AddCommand(amendCmd)
	gCmd.AddCommand(rbCmd)
	gCmd.AddCommand(logCmd)
	gCmd.AddCommand(stCmd)

	cleanCmd.Flags().BoolP("force", "f", false, "Skip confirmation")
}
