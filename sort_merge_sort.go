package main

import (
	// "fmt"
	"strings"
)

// mergeSort sorts contents using merge sort [1].
//
// Sorting is not done in place. A copy of the sorted ContentType is returned.
//
// [1] https://en.wikipedia.org/wiki/Merge_sort
func mergeSort(contents ContentType) ContentType {
	var sortedContents ContentType = append(ContentType{}, contents...)
	var work ContentType = append(ContentType{}, contents...)

	splitContents(sortedContents, work, 0, len(sortedContents))

	return sortedContents
}

// splitContents splits a into two pieces, then sorts both pieces into b, and
// finally merges the sorted pieces of b into a. The set [iBegin, iEnd) marks
// the subset of a and b to consider.
func splitContents(a, b ContentType, iBegin, iEnd int) {
	// fmt.Printf("[splitContents] splitting [%d, %d], ", iBegin, iEnd)
	if iEnd-iBegin <= 1 {
		// We are done.
		// fmt.Println("Sublist has only one element, done")
		return
	}
	var iMiddle = iBegin + (iEnd-iBegin)/2
	// fmt.Printf("iMiddle: %d\n", iMiddle)
	splitContents(a, b, iBegin, iMiddle)       // Sort the left side
	splitContents(a, b, iMiddle, iEnd)         // Sort the right side
	mergeContents(a, b, iBegin, iMiddle, iEnd) // Merge the results
	copySublist(a, b, iBegin, iEnd)            // Copy updated sublist from a to b
}

// mergeContents merges the two halves given by [iBegin, iMiddle) and [iMiddle,
// iEnd) of b into a.
func mergeContents(a, b ContentType, iBegin, iMiddle, iEnd int) {
	var i = iBegin
	var j = iMiddle
	// fmt.Printf("[mergeContents] merging [%d:%d:%d]\n", iBegin, iMiddle, iEnd)
	for k := iBegin; k < iEnd; k++ {
		// fmt.Printf("[mergeContents] [%d:%d:%d], comparing i = %d with j = %d, storing in k = %d, ", iBegin, iMiddle, iEnd, i, j, k)
		if i < iMiddle && j < iEnd {
			var comparison = strings.Compare(b[i].CompareLine, b[j].CompareLine)
			// fmt.Printf("comparison = %d, ", comparison)
			if comparison > 0 {
				// fmt.Printf("storing a[%d] <- b[%d]\n", k, j)
				a[k] = b[j]
				j++
			} else {
				// fmt.Printf("storing a[%d] <- b[%d]\n", k, i)
				a[k] = b[i]
				i++
			}
		} else {
			if i < iMiddle {
				// fmt.Printf("storing left side a[%d] <- b[%d]\n", k, i)
				a[k] = b[i]
				i++
			} else {
				// fmt.Printf("storing right side a[%d] <- b[%d]\n", k, j)
				a[k] = b[j]
				j++
			}
		}
	}
}

// copySublist copies [iBeing:iEnd) from a into b.
func copySublist(a, b ContentType, iBegin, iEnd int) {
	// fmt.Printf("[copySublist] updated b[%d:%d]\n", iBegin, iEnd)
	for i := range a {
		b[i] = a[i]
	}
}
