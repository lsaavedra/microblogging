package utils

import (
	"sort"
)

func MergeAndSort[T any](list1, list2 []T, applyFn func(a, b T) bool) []T {
	merged := append(list1, list2...) // Merge both lists

	// Sort using sort.Slice with the less function
	sort.Slice(merged, func(i, j int) bool {
		return applyFn(merged[i], merged[j])
	})

	return merged
}
