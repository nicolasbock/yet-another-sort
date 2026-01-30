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
	var expected = "no uniq"
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

func TestKeyTypeSetSingleField(t *testing.T) {
	var kt KeyType
	err := kt.Set("3")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if kt.Type != SingleField {
		t.Errorf("Expected type SingleField, got %s", kt.Type)
	}
	if len(kt.Keys) != 1 || kt.Keys[0] != 3 {
		t.Errorf("Expected Keys [3], got %v", kt.Keys)
	}
}

func TestKeyTypeSetRemainder(t *testing.T) {
	var kt KeyType
	err := kt.Set("2,")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if kt.Type != Remainder {
		t.Errorf("Expected type Remainder, got %s", kt.Type)
	}
	if len(kt.Keys) != 1 || kt.Keys[0] != 2 {
		t.Errorf("Expected Keys [2], got %v", kt.Keys)
	}
}

func TestKeyTypeSetSubSet(t *testing.T) {
	var kt KeyType
	err := kt.Set("1,4")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if kt.Type != SubSet {
		t.Errorf("Expected type SubSet, got %s", kt.Type)
	}
	if len(kt.Keys) != 2 || kt.Keys[0] != 1 || kt.Keys[1] != 4 {
		t.Errorf("Expected Keys [1, 4], got %v", kt.Keys)
	}
}

func TestUniqModeSetNone(t *testing.T) {
	var um UniqMode
	err := um.Set("none")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if um != no_uniq {
		t.Errorf("Expected no_uniq, got %s", um)
	}
}

func TestUniqModeSetNoUniq(t *testing.T) {
	var um UniqMode
	err := um.Set("no uniq")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if um != no_uniq {
		t.Errorf("Expected no_uniq, got %s", um)
	}
}

func TestUniqModeSetFirst(t *testing.T) {
	var um UniqMode
	err := um.Set("first")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if um != first {
		t.Errorf("Expected first, got %s", um)
	}
}

func TestUniqModeSetLast(t *testing.T) {
	var um UniqMode
	err := um.Set("last")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if um != last {
		t.Errorf("Expected last, got %s", um)
	}
}

func TestUniqModeSetInvalid(t *testing.T) {
	var um UniqMode
	err := um.Set("invalid")
	if err == nil {
		t.Errorf("Expected error for invalid value, got nil")
	}
}

func TestContentTypeString(t *testing.T) {
	var ct ContentType = ContentType{
		{Lines: []string{"line1"}, Fields: []string{"line1"}, CompareLine: "line1"},
		{Lines: []string{"line2"}, Fields: []string{"line2"}, CompareLine: "line2"},
	}
	result := ct.String()
	if !strings.Contains(result, "2 multilines") {
		t.Errorf("Expected string to contain '2 multilines', got: %s", result)
	}
	if !strings.Contains(result, "line1") {
		t.Errorf("Expected string to contain 'line1', got: %s", result)
	}
	if !strings.Contains(result, "line2") {
		t.Errorf("Expected string to contain 'line2', got: %s", result)
	}
}

func TestContentLineTypeIsEqualDifferentLength(t *testing.T) {
	a := ContentLineType{
		Lines:       []string{"line1", "line2"},
		Fields:      []string{"field1"},
		CompareLine: "compare",
	}
	b := ContentLineType{
		Lines:       []string{"line1"},
		Fields:      []string{"field1"},
		CompareLine: "compare",
	}
	if a.isEqual(b) {
		t.Errorf("Expected contentlines with different line counts to not be equal")
	}
}

func TestContentLineTypeIsEqualDifferentFields(t *testing.T) {
	a := ContentLineType{
		Lines:       []string{"line1"},
		Fields:      []string{"field1", "field2"},
		CompareLine: "compare",
	}
	b := ContentLineType{
		Lines:       []string{"line1"},
		Fields:      []string{"field1"},
		CompareLine: "compare",
	}
	if a.isEqual(b) {
		t.Errorf("Expected contentlines with different field counts to not be equal")
	}
}

func TestContentLineTypeIsEqualDifferentLineContent(t *testing.T) {
	a := ContentLineType{
		Lines:       []string{"line1"},
		Fields:      []string{"field1"},
		CompareLine: "compare",
	}
	b := ContentLineType{
		Lines:       []string{"line2"},
		Fields:      []string{"field1"},
		CompareLine: "compare",
	}
	if a.isEqual(b) {
		t.Errorf("Expected contentlines with different line content to not be equal")
	}
}

func TestContentLineTypeIsEqualDifferentFieldContent(t *testing.T) {
	a := ContentLineType{
		Lines:       []string{"line1"},
		Fields:      []string{"field1"},
		CompareLine: "compare",
	}
	b := ContentLineType{
		Lines:       []string{"line1"},
		Fields:      []string{"field2"},
		CompareLine: "compare",
	}
	if a.isEqual(b) {
		t.Errorf("Expected contentlines with different field content to not be equal")
	}
}

func TestContentLineTypeIsEqualDifferentCompareLine(t *testing.T) {
	a := ContentLineType{
		Lines:       []string{"line1"},
		Fields:      []string{"field1"},
		CompareLine: "compare1",
	}
	b := ContentLineType{
		Lines:       []string{"line1"},
		Fields:      []string{"field1"},
		CompareLine: "compare2",
	}
	if a.isEqual(b) {
		t.Errorf("Expected contentlines with different compare lines to not be equal")
	}
}

func TestContentTypeIsEqualDifferentLength(t *testing.T) {
	a := ContentType{
		{Lines: []string{"line1"}, Fields: []string{"field1"}, CompareLine: "compare1"},
		{Lines: []string{"line2"}, Fields: []string{"field2"}, CompareLine: "compare2"},
	}
	b := ContentType{
		{Lines: []string{"line1"}, Fields: []string{"field1"}, CompareLine: "compare1"},
	}
	if a.isEqual(b) {
		t.Errorf("Expected content types with different lengths to not be equal")
	}
}

func TestContentTypeIsEqualDifferentContent(t *testing.T) {
	a := ContentType{
		{Lines: []string{"line1"}, Fields: []string{"field1"}, CompareLine: "compare1"},
	}
	b := ContentType{
		{Lines: []string{"line2"}, Fields: []string{"field1"}, CompareLine: "compare1"},
	}
	if a.isEqual(b) {
		t.Errorf("Expected content types with different content to not be equal")
	}
}

func TestKeyTStringInvalid(t *testing.T) {
	// Test an invalid KeyT value to cover the default case
	var invalidKey KeyT = KeyT(999)
	result := invalidKey.String()
	if result != "" {
		t.Errorf("Expected empty string for invalid KeyT, got %q", result)
	}
}

func TestUniqModeStringInvalid(t *testing.T) {
	// Test an invalid UniqMode value to cover the default case
	var invalidMode UniqMode = UniqMode(999)
	result := invalidMode.String()
	if result != "FIXME" {
		t.Errorf("Expected 'FIXME' for invalid UniqMode, got %q", result)
	}
}
