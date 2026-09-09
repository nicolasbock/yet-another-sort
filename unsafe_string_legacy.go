//go:build !go1.20
// +build !go1.20

package main

import "unsafe"

// bytesToString returns a string that shares storage with b instead of copying
// it. The caller must not modify b (or the buffer it points into) for as long
// as the returned string is in use.
//
// This is the fallback for Go releases before 1.20, which do not provide
// unsafe.String/unsafe.SliceData. A slice header starts with the same
// {data pointer, length} fields as a string header, so reinterpreting the
// slice header as a string header yields the same zero-copy result. This
// avoids uintptr round-trips, so the backing array stays reachable by the GC.
func bytesToString(b []byte) string {
	if len(b) == 0 {
		return ""
	}
	return *(*string)(unsafe.Pointer(&b))
}
