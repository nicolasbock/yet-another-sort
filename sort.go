package main

// SortContents sorts the content lines and returns a sorted list.
func SortContents(contents ContentType) ContentType {
	switch sortMode {
	case bubble:
		return bubbleSort(contents)
	case merge:
		return mergeSort(contents)
	}
	return ContentType{}
}
