package util

func Contains(slice []string, item string) bool {
	// Contains checks if a slice contains a specific item.
	//
	// slice: The slice to search.
	// item: The item to search for in the slice.
	//
	// Returns:
	// True if the item is found in the slice, false otherwise.

	for _, a := range slice {
		if a == item {
			return true
		}
	}
	return false
}
