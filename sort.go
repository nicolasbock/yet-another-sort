package main

// SortContents sorts the content lines and returns a sorted list.
func SortContents(contents ContentType) ContentType {
	return mergeSort(contents)
}
