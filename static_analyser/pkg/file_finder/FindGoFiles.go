package file_finder

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func FindGoFiles(root string) ([]string, error) {
	// FindGoFiles searches for all .go files starting from the root directory.
	// It returns a list of paths to the .go files.
	//
	// root: The starting directory for the search.
	//
	// Returns:
	// A list of paths to .go files, and an error if there was a problem searching.
	var files []string

	// Walk through each directory and file in the root
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("error accessing a path %q: %w", path, err)
		}
		// Only consider .go files that are not directories
		if !info.IsDir() && strings.HasSuffix(path, ".go") {
			files = append(files, path)
		}
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("error walking the file tree: %w", err)
	}

	return files, nil
}
