package main

import "testing"

func TestKeyTypeString(t *testing.T) {
	var kt KeyType = KeyType{
		Type: NoKey,
	}
	var expected string = "whole line"
	var got string = kt.String()
	if got != expected {
		t.Errorf("Expected %s, got %s", expected, got)
	}
}
