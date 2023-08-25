package main

import (
	/*
		#include <stdlib.h>
		#include <string.h>
	*/
	"C"
	"unsafe"
)

// bubbleSort sorts contents using bubble sort [1].
//
// Sorting is not done in place. A copy of the sorted ContentType is returned.
//
// [1] https://en.wikipedia.org/wiki/Bubble_sort
func bubbleSort(contents ContentType) ContentType {
	var sortedContents ContentType = append(ContentType{}, contents...)
	for i := range sortedContents {
		var a = C.CString(sortedContents[i].CompareLine)
		defer C.free(unsafe.Pointer(a))
		for j := range sortedContents {
			var b = C.CString(sortedContents[j].CompareLine)
			defer C.free(unsafe.Pointer(b))
			if int(C.strcoll(a, b)) < 0 {
				var temp ContentLineType = sortedContents[i]
				sortedContents[i] = sortedContents[j]
				sortedContents[j] = temp
			}
		}
	}
	return sortedContents
}
