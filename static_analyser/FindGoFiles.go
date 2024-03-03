package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// FindGoFilesWithFunctions searches through all Go files in a given directory tree
// starting from a root directory, to find occurrences of specified functions. It then
// maps each function to a list of file paths where the function is called.
//
// Parameters:
// - root: A string specifying the root directory from which the search begins. The function
//   will recursively search through all subdirectories from this point.
// - fn_list: A slice of strings containing the names of functions to search for within the Go files.
//
// Returns:
// - A map where each key is a function name (from the provided fn_list) and the value is a slice
//   of strings, each representing a path to a Go file where that function is called. If a function
//   is not found in any file, it will not appear in the map. If an error occurs during file
//   traversal, the function prints the error and returns nil.
//
// Note:
// This function performs a simple string-based search to identify function calls, which means it
// does not parse the Go files for syntax or semantics. As a result, it may also match comments or
// strings containing the function name pattern. This function assumes that a function call is
// made when it finds the function name followed by an opening parenthesis '(' in the file content.

func FindGoFilesWithFunctions(root string, fn_list []string) map[string][]string {
	var occurrences = make(map[string][]string)

	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		// Only consider .go files that are not directories
		if !info.IsDir() && strings.HasSuffix(path, ".go") {
			content, err := os.ReadFile(path)
			if err != nil {
				return err
			}
			contentStr := string(content)
			// Check each file for the presence of each function in the list
			for _, funcName := range fn_list {
				if strings.Contains(contentStr, funcName+"(") { // Simple check for function call
					occurrences[funcName] = append(occurrences[funcName], path)
				}
			}
		}
		return nil
	})

	if err != nil {
		fmt.Printf("Error walking the file tree: %v\n", err)
		return nil
	}

	return occurrences
}


func FindGoFiles(root string) ([]string) {
    var files []string

    err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
        if err != nil {
            return err
        }
        // Only consider .go files that are not directories
        if !info.IsDir() && strings.HasSuffix(path, ".go") {
            files = append(files, path)
        }
        return nil
    })

    if err != nil {
        fmt.Printf("Error walking the file tree: %v\n", err)
		return nil
    }

    return files
}