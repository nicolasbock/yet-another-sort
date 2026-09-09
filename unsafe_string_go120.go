//go:build go1.20
// +build go1.20

package main

import "unsafe"

// bytesToString returns a string that shares storage with b instead of copying
// it. The caller must not modify b (or the buffer it points into) for as long
// as the returned string is in use.
//
// This is the Go 1.20+ implementation, which uses the supported
// unsafe.String/unsafe.SliceData helpers.
func bytesToString(b []byte) string {
	return unsafe.String(unsafe.SliceData(b), len(b))
}
