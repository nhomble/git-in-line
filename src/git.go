package main

import (
	"fmt"
	"os/exec"
	"strings"
)

func ExpandCommitHash(repoPath, shortHash string) (string, error) {
	cmd := exec.Command("git", "-C", repoPath, "rev-parse", shortHash)
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to expand hash %s: %v", shortHash, err)
	}

	fullHash := strings.TrimSpace(string(output))
	return fullHash, nil
}

func StreamCommits(repoPath, branch, order string) (*exec.Cmd, error) {
	var cmd *exec.Cmd
	if order == "desc" {
		cmd = exec.Command("git", "-C", repoPath, "rev-list", branch) // newest first
	} else {
		cmd = exec.Command("git", "-C", repoPath, "rev-list", "--reverse", branch) // oldest first
	}
	return cmd, nil
}
