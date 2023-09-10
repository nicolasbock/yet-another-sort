package main

// bubbleSort sorts contents using bubble sort [1].
//
// Sorting is not done in place. A copy of the sorted ContentType is returned.
//
// [1] https://en.wikipedia.org/wiki/Bubble_sort
func bubbleSort(contents ContentType) ContentType {
	var sortedContents ContentType = append(ContentType{}, contents...)
	for i := range sortedContents {
		var a = sortedContents[i].CompareLine
		for j := range sortedContents {
			var b = sortedContents[j].CompareLine
			if a < b {
				var temp ContentLineType = sortedContents[i]
				sortedContents[i] = sortedContents[j]
				sortedContents[j] = temp
			}
		}
	}
	return sortedContents
}
