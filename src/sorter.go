package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type SortOptions struct {
	Order         string
	FailOnMissing bool
	Debug         bool
}

func Sort(repoPath, branch string, inputHashes []string, opts SortOptions) ([]string, error) {
	// Step 1: Filter and expand input hashes to only valid ones
	validHashes := make(map[string]string) // fullHash -> originalHash
	var missingHashes []string

	for _, hash := range inputHashes {
		fullHash, err := ExpandCommitHash(repoPath, hash)
		if err != nil {
			missingHashes = append(missingHashes, hash)
		} else {
			validHashes[fullHash] = hash
		}
	}

	if len(missingHashes) > 0 && opts.FailOnMissing {
		return nil, fmt.Errorf("missing commits not found in branch: %v", missingHashes)
	}

	// Step 2: Stream commits and build ordered result
	orderedHashes, err := streamAndSort(repoPath, branch, validHashes, opts.Order, opts.Debug)
	if err != nil {
		return nil, err
	}

	// Step 3: Combine ordered and missing hashes
	result := make([]string, 0, len(inputHashes))
	result = append(result, orderedHashes...)
	result = append(result, missingHashes...)

	return result, nil
}

// streamAndSort streams commits from git and builds ordered result, stopping early when all hashes found
func streamAndSort(repoPath, branch string, validHashes map[string]string, order string, debug bool) ([]string, error) {
	cmd, err := StreamCommits(repoPath, branch, order)
	if err != nil {
		return nil, err
	}

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, fmt.Errorf("failed to create stdout pipe: %v", err)
	}

	if err := cmd.Start(); err != nil {
		return nil, fmt.Errorf("failed to start git rev-list: %v", err)
	}

	if debug {
		fmt.Fprintf(os.Stderr, "[DEBUG] Streaming commits in %s order from branch %s\n", order, branch)
		fmt.Fprintf(os.Stderr, "[DEBUG] Looking for %d commits: ", len(validHashes))
		for fullHash, origHash := range validHashes {
			fmt.Fprintf(os.Stderr, "%s->%s ", origHash, fullHash[:8])
		}
		fmt.Fprintf(os.Stderr, "\n")
	}

	var orderedResult []string
	remainingHashes := make(map[string]string)
	for k, v := range validHashes {
		remainingHashes[k] = v
	}

	commitCount := 0
	scanner := bufio.NewScanner(stdout)
	for scanner.Scan() {
		commit := strings.TrimSpace(scanner.Text())
		if commit == "" {
			continue
		}

		commitCount++
		if debug && commitCount <= 10 {
			fmt.Fprintf(os.Stderr, "[DEBUG] Commit %d: %s\n", commitCount, commit[:8])
		}

		// Check if this commit is one we're looking for
		if originalHash, exists := remainingHashes[commit]; exists {
			if debug {
				fmt.Fprintf(os.Stderr, "[DEBUG] FOUND: %s -> %s (position %d)\n", originalHash, commit[:8], commitCount)
			}
			orderedResult = append(orderedResult, originalHash)
			delete(remainingHashes, commit)

			// Early termination: stop when we've found all commits
			if len(remainingHashes) == 0 {
				if debug {
					fmt.Fprintf(os.Stderr, "[DEBUG] All commits found, stopping after %d commits\n", commitCount)
				}
				break
			}
		}
	}

	// Wait for git process to finish
	if err := cmd.Wait(); err != nil {
		// If we terminated early, the process might have been killed
		// Only return error if we didn't find all our commits
		if len(remainingHashes) > 0 {
			return nil, fmt.Errorf("git rev-list failed: %v", err)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading git output: %v", err)
	}

	return orderedResult, nil
}
