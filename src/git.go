package main

import (
	"fmt"
	"os/exec"
	"strings"
)

func GetBranchCommits(repoPath, branch string) ([]string, error) {
	cmd := exec.Command("git", "-C", repoPath, "rev-list", "--reverse", branch)
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to get commits from branch %s: %v", branch, err)
	}

	commits := strings.Split(strings.TrimSpace(string(output)), "\n")
	var result []string
	for _, commit := range commits {
		commit = strings.TrimSpace(commit)
		if commit != "" {
			result = append(result, commit)
		}
	}

	return result, nil
}

func ExpandCommitHash(repoPath, shortHash string) (string, error) {
	cmd := exec.Command("git", "-C", repoPath, "rev-parse", shortHash)
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to expand hash %s: %v", shortHash, err)
	}
	
	fullHash := strings.TrimSpace(string(output))
	return fullHash, nil
}

// CompareCommits returns -1 if commit1 comes before commit2 in the branch history,
// 1 if commit1 comes after commit2, and 0 if they're the same commit.
// Returns error if either commit is not found in the branch.
func CompareCommits(repoPath, branch, commit1, commit2 string) (int, error) {
	if commit1 == commit2 {
		return 0, nil
	}

	// Check if commit1 is an ancestor of commit2
	cmd := exec.Command("git", "-C", repoPath, "merge-base", "--is-ancestor", commit1, commit2)
	err := cmd.Run()
	if err == nil {
		// commit1 is ancestor of commit2, so commit1 comes first
		return -1, nil
	}

	// Check if commit2 is an ancestor of commit1
	cmd = exec.Command("git", "-C", repoPath, "merge-base", "--is-ancestor", commit2, commit1)
	err = cmd.Run()
	if err == nil {
		// commit2 is ancestor of commit1, so commit2 comes first
		return 1, nil
	}

	// Neither is ancestor of the other - they might be on different branches
	// or one/both commits don't exist in the branch
	return 0, fmt.Errorf("commits %s and %s are not in ancestor relationship in branch %s", commit1, commit2, branch)
}