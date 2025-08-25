package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

type Config struct {
	RepoPath      string
	CommitFile    string
	Branch        string
	Order         string
	FailOnMissing bool
	Debug         bool
}

func NewRootCmd() *cobra.Command {
	config := &Config{}

	cmd := &cobra.Command{
		Use:   "gil <repo-path> <commit-file>",
		Short: "Sort commit hashes by their order in a git branch (git-in-line)",
		Long: `gil (git-in-line) takes a git repository path and a file containing newline-delimited 
commit hashes, then outputs the commit hashes sorted by their order in the specified branch.`,
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			config.RepoPath = args[0]
			config.CommitFile = args[1]
			return Run(config)
		},
	}

	cmd.Flags().StringVarP(&config.Branch, "branch", "b", "main", "Git branch to use for ordering")
	cmd.Flags().StringVarP(&config.Order, "order", "o", "desc", "Sort order: 'asc' or 'desc'")
	cmd.Flags().BoolVarP(&config.FailOnMissing, "fail-on-missing", "f", false, "Fail if commits are not found in branch")
	cmd.Flags().BoolVarP(&config.Debug, "debug", "d", false, "Show debug output for troubleshooting")

	return cmd
}

func Run(config *Config) error {
	if config.Order != "desc" && config.Order != "asc" {
		return fmt.Errorf("order must be 'desc' or 'asc', got: %s", config.Order)
	}

	commitHashes, err := ReadCommitHashes(config.CommitFile)
	if err != nil {
		return fmt.Errorf("error reading commit file: %v", err)
	}

	// Expand short hashes to full hashes
	expandedHashes := make([]string, 0, len(commitHashes))
	originalHashes := make(map[string]string) // full -> original mapping for output

	for _, hash := range commitHashes {
		fullHash, err := ExpandCommitHash(config.RepoPath, hash)
		if err != nil {
			// If expansion fails, keep the original hash (might be invalid or not in repo)
			expandedHashes = append(expandedHashes, hash)
			originalHashes[hash] = hash
		} else {
			expandedHashes = append(expandedHashes, fullHash)
			originalHashes[fullHash] = hash
		}
	}

	sortedHashes, err := Sort(config.RepoPath, config.Branch, expandedHashes, SortOptions{
		Order:         config.Order,
		FailOnMissing: config.FailOnMissing,
		Debug:         config.Debug,
	})
	if err != nil {
		return fmt.Errorf("error sorting commits: %v", err)
	}

	// Output using original hash format from input
	for _, fullHash := range sortedHashes {
		if originalHash, exists := originalHashes[fullHash]; exists {
			fmt.Println(originalHash)
		} else {
			fmt.Println(fullHash)
		}
	}

	return nil
}

func ReadCommitHashes(filename string) ([]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var hashes []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" {
			hashes = append(hashes, line)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return hashes, nil
}
