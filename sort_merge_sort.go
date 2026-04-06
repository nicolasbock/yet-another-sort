package main

import (
	"sort"
)

// mergeSort sorts contents using Go's built-in optimized sort algorithm.
//
// This implementation uses sort.Interface which delegates to Go's highly
// optimized sorting implementation (pdqsort - pattern-defeating quicksort).
// This is significantly faster than a custom merge sort implementation.
//
// Go's sort.Sort uses an introsort variant that automatically switches between:
// - Quicksort for general cases
// - Heapsort when recursion depth is too high (to avoid worst-case O(n²))
// - Insertion sort for small slices (< 12 elements)
//
// For stable sorting (preserving order of equal elements), we use sort.Stable
// which implements a stable version of the sorting algorithm.
// Use --stable-sort to enable this; it is required when --uniq last/first must
// keep the duplicate that appeared last/first in the original input.
func mergeSort(contents ContentType) ContentType {
	// If the slice is empty or has one element, it's already sorted
	if len(contents) <= 1 {
		return contents
	}

	// Create a copy to avoid modifying the original
	sortedContents := make(ContentType, len(contents))
	copy(sortedContents, contents)

	if options.stableSort {
		// Stable sort preserves the original input order among equal keys.
		// Required for --uniq last/first to correctly identify which duplicate
		// appeared last/first in the original input.
		sort.Stable(sortedContents)
	} else {
		// Unstable sort (default): faster, but does not preserve input order
		// among entries with equal keys.
		sort.Sort(sortedContents)
	}

	return sortedContents
}

// Implement sort.Interface for ContentType
// These three methods allow ContentType to be used with Go's sort package

// Len returns the number of elements in the collection
func (c ContentType) Len() int {
	return len(c)
}

// Less reports whether the element with index i should sort before the element with index j
// This uses direct string comparison which is highly optimized in Go
func (c ContentType) Less(i, j int) bool {
	// Direct string comparison is faster than extracting to variables
	// Go's string comparison is optimized at the compiler level
	return c[i].CompareLine < c[j].CompareLine
}

// Swap swaps the elements with indexes i and j
// Go optimizes this into efficient memory operations
func (c ContentType) Swap(i, j int) {
	c[i], c[j] = c[j], c[i]
}
