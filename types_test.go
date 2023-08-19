package main

import (
	"strings"
	"testing"
)

func TestKeyTypeString(t *testing.T) {
	var kt KeyType = KeyType{
		Type: NoKey,
	}
	var expected string = "whole line"
	var got string = kt.String()
	if strings.Compare(got, expected) != 0 {
		t.Errorf("Expected\n%s\ngot\n%s", expected, got)
	}
}

func TestContentLineType(t *testing.T) {
	var c ContentLineType = ContentLineType{Lines: []string{"first line"}, Fields: []string{"first", "line"}, CompareLine: "first"}
	var expected string = "multiline\n  line: 'first line'\n  fields: first, line\n  compare: 'first'\n"
	var got string = c.String()
	if strings.Compare(got, expected) != 0 {
		t.Errorf("Expected\n%s\ngot\n%s", expected, got)
	}
}
