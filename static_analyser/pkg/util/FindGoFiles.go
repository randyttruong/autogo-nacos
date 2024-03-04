package util

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

func FindGoFilesWithFunctions(root string, fn_list []string) (map[string][]string, error) {
	// FindGoFilesWithFunctions searches for all .go files starting from the root directory and checks if they contain certain function calls.
	// It returns a map where the keys are the function names and the values are lists of paths to the .go files that contain those function calls.
	//
	// root: The starting directory for the search.
	// fn_list: A list of function names to search for in the .go files.
	//
	// Returns:
	// A map where the keys are the function names and the values are lists of paths to the .go files that contain those function calls.
	// An error if there was a problem searching for .go files or reading a file.
	var occurrences = make(map[string][]string)

	// Find all the Go files in this directory
	files, err := FindGoFiles(root)
	if err != nil {
		return nil, fmt.Errorf("error finding go files in %s: %v", root, err)
	}

	// Check each file for the function calls
	for _, file := range files {
		// Read the file
		content, err := os.ReadFile(file)
		if err != nil {
			fmt.Printf("Error reading file %s: %v", file, err)
			continue
		}

		// Convert file to a string
		contentStr := string(content)

		// Check for each function in the list if it exists in the file
		for _, funcName := range fn_list {
			if !strings.Contains(contentStr, funcName+"(") { // Simple check for function call
				occurrences[funcName] = append(occurrences[funcName], file)
			}
		}
	}

	return occurrences, nil
}
