package util

import (
	"slices"
)

// Returns an alphabetically sorted list of the keys of a map[string]any
func SortedKeys[T any](m map[string]T) []string {
	keys := make([]string, 0, len(m))

	for k := range m {
		keys = append(keys, k)
	}
	slices.Sort(keys)
	return keys
}
