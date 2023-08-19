package main

import (
	"strings"
)

// SortContents sorts the content lines and returns a sorted list.
func SortContents(contents ContentType, key int) (sortedContents ContentType) {
	sortedContents = append(sortedContents, contents...)

	// Combine multilines with field-separator.
	// Sort multilined content
	// Split multilined content into original multilines

	for i := range sortedContents {
		for j := range sortedContents {
			if strings.Compare(sortedContents[i].Fields[key-1], sortedContents[j].Fields[key-1]) < 0 {
				var temp ContentLineType = sortedContents[i]
				sortedContents[i] = sortedContents[j]
				sortedContents[j] = temp
			}
		}
	}
	return sortedContents
}
