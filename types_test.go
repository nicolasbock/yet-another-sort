package main

import (
	"strings"
	"testing"
)

func TestKeyTString(t *testing.T) {
	var expected string = "NoKey"
	var got string = NoKey.String()
	if strings.Compare(expected, got) != 0 {
		t.Errorf("Expected\n%s\ngot\n%s", expected, got)
	}
	expected = "SingleField"
	got = SingleField.String()
	if strings.Compare(expected, got) != 0 {
		t.Errorf("Expected\n%s\ngot\n%s", expected, got)
	}
	expected = "Remainder"
	got = Remainder.String()
	if strings.Compare(expected, got) != 0 {
		t.Errorf("Expected\n%s\ngot\n%s", expected, got)
	}
	expected = "SubSet"
	got = SubSet.String()
	if strings.Compare(expected, got) != 0 {
		t.Errorf("Expected\n%s\ngot\n%s", expected, got)
	}
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
	var got string = kt.Representation()
	if strings.Compare(got, expected) != 0 {
		t.Errorf("Expected\n%s\ngot\n%s", expected, got)
	}
	kt = KeyType{
		Type: SingleField,
		Keys: []int{1},
	}
	expected = "KeyType{Type: SingleField, Keys: []int{1}}"
	got = kt.Representation()
	if strings.Compare(got, expected) != 0 {
		t.Errorf("Expected\n%s\ngot\n%s", expected, got)
	}
	kt = KeyType{
		Type: Remainder,
		Keys: []int{1},
	}
	expected = "KeyType{Type: Remainder, Keys: []int{1}}"
	got = kt.Representation()
	if strings.Compare(got, expected) != 0 {
		t.Errorf("Expected\n%s\ngot\n%s", expected, got)
	}
	kt = KeyType{
		Type: SubSet,
		Keys: []int{1, 4},
	}
	expected = "KeyType{Type: SubSet, Keys: []int{1, 4}}"
	got = kt.Representation()
	if strings.Compare(got, expected) != 0 {
		t.Errorf("Expected\n%s\ngot\n%s", expected, got)
	}
}

func TestContentLineTypeString(t *testing.T) {
	var c ContentLineType = ContentLineType{Lines: []string{"first line"}, Fields: []string{"first", "line"}, CompareLine: "first"}
	var expected string = "multiline\n  line: \"first line\"\n  fields: [ \"first\", \"line\" ]\n  compare: \"first\""
	var got string = c.String()
	if strings.Compare(got, expected) != 0 {
		t.Errorf("Expected\n%s\ngot\n%s", expected, got)
	}
}

func TestUniqModeString(t *testing.T) {
	var um UniqMode
	var expected = "none"
	var got string = um.String()
	if strings.Compare(got, expected) != 0 {
		t.Errorf("Expected\n%s\ngot\n%s", expected, got)
	}
	um = first
	expected = "first"
	got = um.String()
	if strings.Compare(got, expected) != 0 {
		t.Errorf("Expected\n%s\ngot\n%s", expected, got)
	}
	um = last
	expected = "last"
	got = um.String()
	if strings.Compare(got, expected) != 0 {
		t.Errorf("Expected\n%s\ngot\n%s", expected, got)
	}
}
