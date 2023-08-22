package main

import (
	"strings"
	"testing"
)

func compareStrings(expected, got string, t *testing.T) {
	if strings.Compare(expected, got) != 0 {
		t.Errorf("Expected\n%s\ngot\n%s", expected, got)
	}
}

func TestKeyTString(t *testing.T) {
	compareStrings("NoKey", NoKey.String(), t)
	compareStrings("SingleField", SingleField.String(), t)
	compareStrings("Remainder", Remainder.String(), t)
	compareStrings("SubSet", SubSet.String(), t)
}

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

func TestKeyTypeRepresentation(t *testing.T) {
	var kt KeyType = KeyType{
		Type: NoKey,
	}
	var expected string = "KeyType{Type: NoKey}"
	compareStrings(expected, kt.Representation(), t)
	kt = KeyType{
		Type: SingleField,
		Keys: []int{1},
	}
	expected = "KeyType{Type: SingleField, Keys: []int{1}}"
	compareStrings(expected, kt.Representation(), t)
	kt = KeyType{
		Type: Remainder,
		Keys: []int{1},
	}
	expected = "KeyType{Type: Remainder, Keys: []int{1}}"
	compareStrings(expected, kt.Representation(), t)
	kt = KeyType{
		Type: SubSet,
		Keys: []int{1, 4},
	}
	expected = "KeyType{Type: SubSet, Keys: []int{1, 4}}"
	compareStrings(expected, kt.Representation(), t)
}

func TestContentLineType(t *testing.T) {
	var c ContentLineType = ContentLineType{Lines: []string{"first line"}, Fields: []string{"first", "line"}, CompareLine: "first"}
	var expected string = "multiline\n  line: \"first line\"\n  fields: [ \"first\", \"line\" ]\n  compare: \"first\""
	var got string = c.String()
	if strings.Compare(got, expected) != 0 {
		t.Errorf("Expected\n%s\ngot\n%s", expected, got)
	}
}
