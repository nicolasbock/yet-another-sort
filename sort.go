package main

import (
	/*
		#include <stdlib.h>
		#include <string.h>
	*/
	"C"
	"unsafe"
)

// SortContents sorts the content lines and returns a sorted list.
func SortContents(contents ContentType) (sortedContents ContentType) {
	sortedContents = append(sortedContents, contents...)

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
