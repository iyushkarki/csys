package git

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func RunGit(args ...string) error {
	cmd := exec.Command("git", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func RunGitOutput(args ...string) (string, error) {
	cmd := exec.Command("git", args...)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		return "", fmt.Errorf("%s", strings.TrimSpace(stderr.String()))
	}
	return strings.TrimSpace(stdout.String()), nil
}

func GetCurrentBranch() (string, error) {
	return RunGitOutput("rev-parse", "--abbrev-ref", "HEAD")
}

func GetAllBranches() ([]string, error) {
	output, err := RunGitOutput("branch", "--format=%(refname:short)")
	if err != nil {
		return nil, err
	}
	if output == "" {
		return []string{}, nil
	}
	return strings.Split(output, "\n"), nil
}

func IsGitRepo() bool {
	cmd := exec.Command("git", "rev-parse", "--git-dir")
	return cmd.Run() == nil
}

func IsGitInstalled() bool {
	_, err := exec.LookPath("git")
	return err == nil
}
