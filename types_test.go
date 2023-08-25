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
	var expected string = "NoKey"
	var got string = kt.String()
	if strings.Compare(got, expected) != 0 {
		t.Errorf("Expected\n%s\ngot\n%s", expected, got)
	}
	kt = KeyType{
		Type: SingleField,
		Keys: []int{1},
	}
	expected = "SingleField, [1]"
	got = kt.String()
	if strings.Compare(got, expected) != 0 {
		t.Errorf("Expected\n%s\ngot\n%s", expected, got)
	}
	kt = KeyType{
		Type: Remainder,
		Keys: []int{1},
	}
	expected = "Remainder, [1]"
	got = kt.String()
	if strings.Compare(got, expected) != 0 {
		t.Errorf("Expected\n%s\ngot\n%s", expected, got)
	}
	kt = KeyType{
		Type: SubSet,
		Keys: []int{1, 3},
	}
	expected = "SubSet, [1, 3]"
	got = kt.String()
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

func TestIsEqualContentTypeLine(t *testing.T) {
	var a ContentLineType = ContentLineType{
		Lines:       []string{"Line", "another line", "third line"},
		Fields:      []string{"Line", "another", "line", "third", "line"},
		CompareLine: "Line another line",
	}
	var b ContentLineType = ContentLineType{
		Lines:       append([]string{}, a.Lines...),
		Fields:      append([]string{}, a.Fields...),
		CompareLine: a.CompareLine,
	}
	if !a.isEqual(b) {
		t.Errorf("%s\nis not equal to\n%s", a, b)
	}
}

func TestIsEqualContentType(t *testing.T) {
	var a ContentType = []ContentLineType{
		{
			Lines:       []string{"Line", "another line", "third line"},
			Fields:      []string{"Line", "another", "line", "third", "line"},
			CompareLine: "Line another line",
		},
		{
			Lines:       []string{"Line", "another line", "third line"},
			Fields:      []string{"Line", "another", "line", "third", "line"},
			CompareLine: "Line another line",
		},
	}
	var b ContentType = []ContentLineType{
		{
			Lines:       append([]string{}, a[0].Lines...),
			Fields:      append([]string{}, a[0].Fields...),
			CompareLine: a[0].CompareLine,
		},
		{
			Lines:       append([]string{}, a[1].Lines...),
			Fields:      append([]string{}, a[1].Fields...),
			CompareLine: a[1].CompareLine,
		},
	}
	if !a.isEqual(b) {
		t.Errorf("%s\nis not equal to\n%s", a, b)
	}
}

func TestSortModeString(t *testing.T) {
	var sm SortMode
	var expected = "bubble sort"
	var got string = sm.String()
	if strings.Compare(got, expected) != 0 {
		t.Errorf("Expected %s got %s", expected, got)
	}
	sm = bubble
	got = sm.String()
	if strings.Compare(got, expected) != 0 {
		t.Errorf("Expected %s got %s", expected, got)
	}
	sm = merge
	expected = "merge sort"
	got = sm.String()
	if strings.Compare(got, expected) != 0 {
		t.Errorf("Expected %s got %s", expected, got)
	}
}
