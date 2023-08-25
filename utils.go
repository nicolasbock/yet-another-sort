package main

import "reflect"

// isEqual return true if a == b.
func (a ContentType) isEqual(b ContentType) bool {
	return reflect.DeepEqual(a, b)
}
