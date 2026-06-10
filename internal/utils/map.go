package utils

import "slices"

// GetMapSortedKeys returns a sorted slice of keys for any map with string keys.
func GetMapSortedKeys[V any](m map[int]V) []int {
	keys := make([]int, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	slices.Sort(keys)
	return keys
}
