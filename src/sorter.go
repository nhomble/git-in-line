package main

import (
	"fmt"
	"sort"
)

type SortOptions struct {
	Order         string
	FailOnMissing bool
}

func Sort(repoPath, branch string, inputHashes []string, opts SortOptions) ([]string, error) {
	// Separate valid commits from invalid ones
	var validHashes []string
	var invalidHashes []string

	for _, hash := range inputHashes {
		// Try to expand the hash to verify it exists in the repo
		_, err := ExpandCommitHash(repoPath, hash)
		if err != nil {
			invalidHashes = append(invalidHashes, hash)
		} else {
			validHashes = append(validHashes, hash)
		}
	}

	if len(invalidHashes) > 0 && opts.FailOnMissing {
		return nil, fmt.Errorf("missing commits not found in branch: %v", invalidHashes)
	}

	// Sort valid hashes using git comparison
	sort.Slice(validHashes, func(i, j int) bool {
		cmp, err := CompareCommits(repoPath, branch, validHashes[i], validHashes[j])
		if err != nil {
			// If comparison fails, maintain original order
			return false
		}
		
		if opts.Order == "desc" {
			return cmp > 0 // Later commits first
		}
		return cmp < 0 // Earlier commits first
	})

	// Combine results
	result := make([]string, 0, len(inputHashes))
	result = append(result, validHashes...)
	result = append(result, invalidHashes...)

	return result, nil
}
