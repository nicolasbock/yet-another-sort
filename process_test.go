package main

import (
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	options = NewConfigurationOptions()
	exitCode := m.Run()
	os.Exit(exitCode)
}

func TestProcessInputFiles1(t *testing.T) {
	var input []string = []string{
		"Line one",
		"Line two",
		"Line three",
		"Line four",
	}
	options.multiline = 1
	var key KeyType = KeyType{
		Type: NoKey,
	}
	var expected ContentType = ContentType{
		{Lines: []string{"Line one"}, Fields: nil, CompareLine: "Line one"},
		{Lines: []string{"Line two"}, Fields: nil, CompareLine: "Line two"},
		{Lines: []string{"Line three"}, Fields: nil, CompareLine: "Line three"},
		{Lines: []string{"Line four"}, Fields: nil, CompareLine: "Line four"},
	}
	var got ContentType = ProcessInputFiles(input, key)
	if !got.isEqual(expected) {
		t.Errorf("got\n%s\nexpected\n%s", got, expected)
	}
}

func TestProcessInputFiles2(t *testing.T) {
	var input []string = []string{
		"Line one",
		"Line two",
		"Line three",
		"Line four",
	}
	options.multiline = 2
	var key KeyType = KeyType{
		Type: NoKey,
	}
	var expected ContentType = ContentType{
		{Lines: []string{"Line one", "Line two"}, Fields: nil, CompareLine: "Line one Line two"},
		{Lines: []string{"Line three", "Line four"}, Fields: nil, CompareLine: "Line three Line four"},
	}
	var got ContentType = ProcessInputFiles(input, key)
	if !got.isEqual(expected) {
		t.Errorf("got\n%s\nexpected\n%s", got, expected)
	}
}

func TestProcessInputFiles3(t *testing.T) {
	var input []string = []string{
		"Line one",
		" Line two",
		" Line three",
		"Line four",
	}
	options.multiline = 2
	var key KeyType = KeyType{
		Type: SingleField,
		Keys: []int{1},
	}
	var expected ContentType = ContentType{
		{Lines: []string{"Line one", " Line two"}, Fields: []string{"Line", "one", "Line", "two"}, CompareLine: "Line"},
		{Lines: []string{" Line three", "Line four"}, Fields: []string{"Line", "three", "Line", "four"}, CompareLine: "Line"},
	}
	var got ContentType = ProcessInputFiles(input, key)
	if !got.isEqual(expected) {
		t.Errorf("got\n%s\nexpected\n%s", got, expected)
	}
}

func TestProcessInputFiles4(t *testing.T) {
	var input []string = []string{
		"line 1",
		"Line 2",
		"line 3",
		"Line 4",
		"line 5",
	}
	var expected ContentType = ContentType{
		{Lines: []string{"line 1"}, Fields: nil, CompareLine: "line 1"},
		{Lines: []string{"Line 2"}, Fields: nil, CompareLine: "line 2"},
		{Lines: []string{"line 3"}, Fields: nil, CompareLine: "line 3"},
		{Lines: []string{"Line 4"}, Fields: nil, CompareLine: "line 4"},
		{Lines: []string{"line 5"}, Fields: nil, CompareLine: "line 5"},
	}
	var key = KeyType{}
	options.multiline = 1
	options.ignoreCase = true
	var got ContentType = ProcessInputFiles(input, key)
	if !got.isEqual(expected) {
		t.Errorf("got\n%s\nexpected\n%s", got, expected)
	}
}

func TestProcessInputFilesWithRemainder(t *testing.T) {
	var input []string = []string{
		"field1 field2 field3 field4",
		"a b c d e",
	}
	options.multiline = 1
	options.fieldSeparator = " "
	options.ignoreCase = false
	options.ignoreLeadingBlanks = false
	var key KeyType = KeyType{
		Type: Remainder,
		Keys: []int{2},
	}
	var expected ContentType = ContentType{
		{Lines: []string{"field1 field2 field3 field4"}, Fields: []string{"field1", "field2", "field3", "field4"}, CompareLine: "field2 field3 field4"},
		{Lines: []string{"a b c d e"}, Fields: []string{"a", "b", "c", "d", "e"}, CompareLine: "b c d e"},
	}
	var got ContentType = ProcessInputFiles(input, key)
	if !got.isEqual(expected) {
		t.Errorf("got\n%s\nexpected\n%s", got, expected)
	}
}

func TestProcessInputFilesWithSubSet(t *testing.T) {
	var input []string = []string{
		"field1 field2 field3 field4",
		"a b c d e",
	}
	options.multiline = 1
	options.fieldSeparator = " "
	options.ignoreCase = false
	options.ignoreLeadingBlanks = false
	var key KeyType = KeyType{
		Type: SubSet,
		Keys: []int{2, 4},
	}
	var expected ContentType = ContentType{
		{Lines: []string{"field1 field2 field3 field4"}, Fields: []string{"field1", "field2", "field3", "field4"}, CompareLine: "field2 field3 field4"},
		{Lines: []string{"a b c d e"}, Fields: []string{"a", "b", "c", "d", "e"}, CompareLine: "b c d"},
	}
	var got ContentType = ProcessInputFiles(input, key)
	if !got.isEqual(expected) {
		t.Errorf("got\n%s\nexpected\n%s", got, expected)
	}
}

func TestProcessInputFilesWithIgnoreLeadingBlanks(t *testing.T) {
	var input []string = []string{
		"  leading spaces",
		"\tleading tab",
		" \t mixed whitespace",
	}
	options.multiline = 1
	options.fieldSeparator = " "
	options.ignoreCase = false
	options.ignoreLeadingBlanks = true
	var key KeyType = KeyType{
		Type: NoKey,
	}
	var expected ContentType = ContentType{
		{Lines: []string{"  leading spaces"}, Fields: nil, CompareLine: "leading spaces"},
		{Lines: []string{"\tleading tab"}, Fields: nil, CompareLine: "leading tab"},
		{Lines: []string{" \t mixed whitespace"}, Fields: nil, CompareLine: "mixed whitespace"},
	}
	var got ContentType = ProcessInputFiles(input, key)
	if !got.isEqual(expected) {
		t.Errorf("got\n%s\nexpected\n%s", got, expected)
	}
	options.ignoreLeadingBlanks = false
}

func TestProcessInputFilesWithCustomSeparator(t *testing.T) {
	var input []string = []string{
		"field1,field2,field3",
		"a,b,c",
	}
	options.multiline = 1
	options.fieldSeparator = ","
	options.ignoreCase = false
	options.ignoreLeadingBlanks = false
	var key KeyType = KeyType{
		Type: SingleField,
		Keys: []int{2},
	}
	var expected ContentType = ContentType{
		{Lines: []string{"field1,field2,field3"}, Fields: []string{"field1", "field2", "field3"}, CompareLine: "field2"},
		{Lines: []string{"a,b,c"}, Fields: []string{"a", "b", "c"}, CompareLine: "b"},
	}
	var got ContentType = ProcessInputFiles(input, key)
	if !got.isEqual(expected) {
		t.Errorf("got\n%s\nexpected\n%s", got, expected)
	}
	options.fieldSeparator = " "
}

func TestProcessInputFilesWithMultiCharSeparator(t *testing.T) {
	var input []string = []string{
		"field1::field2::field3",
		"a::b::c",
	}
	options.multiline = 1
	options.fieldSeparator = "::"
	options.ignoreCase = false
	options.ignoreLeadingBlanks = false
	var key KeyType = KeyType{
		Type: NoKey,
	}
	var expected ContentType = ContentType{
		{Lines: []string{"field1::field2::field3"}, Fields: []string{"field1", "field2", "field3"}, CompareLine: "field1::field2::field3"},
		{Lines: []string{"a::b::c"}, Fields: []string{"a", "b", "c"}, CompareLine: "a::b::c"},
	}
	// Note: non-space separator still populates Fields (empty fields are collapsed in CompareLine)
	var got ContentType = ProcessInputFiles(input, key)
	if !got.isEqual(expected) {
		t.Errorf("got\n%s\nexpected\n%s", got, expected)
	}
	options.fieldSeparator = " "
}

func TestProcessInputFilesEmptyFields(t *testing.T) {
	var input []string = []string{
		"field1,,field3",
		"a,,c",
	}
	options.multiline = 1
	options.fieldSeparator = ","
	options.ignoreCase = false
	options.ignoreLeadingBlanks = false
	var key KeyType = KeyType{
		Type: NoKey,
	}
	var expected ContentType = ContentType{
		{Lines: []string{"field1,,field3"}, Fields: []string{"field1", "field3"}, CompareLine: "field1,field3"},
		{Lines: []string{"a,,c"}, Fields: []string{"a", "c"}, CompareLine: "a,c"},
	}
	// Note: non-space separator still populates Fields (empty fields are collapsed in CompareLine)
	var got ContentType = ProcessInputFiles(input, key)
	if !got.isEqual(expected) {
		t.Errorf("got\n%s\nexpected\n%s", got, expected)
	}
	options.fieldSeparator = " "
}

func TestProcessInputFilesSingleFieldOnly(t *testing.T) {
	var input []string = []string{
		"onlyfield",
	}
	options.multiline = 1
	options.fieldSeparator = " "
	options.ignoreCase = false
	options.ignoreLeadingBlanks = false
	var key KeyType = KeyType{
		Type: NoKey,
	}
	var expected ContentType = ContentType{
		{Lines: []string{"onlyfield"}, Fields: nil, CompareLine: "onlyfield"},
	}
	var got ContentType = ProcessInputFiles(input, key)
	if !got.isEqual(expected) {
		t.Errorf("got\n%s\nexpected\n%s", got, expected)
	}
}

func TestProcessInputFilesNoFields(t *testing.T) {
	var input []string = []string{
		"",
	}
	options.multiline = 1
	options.fieldSeparator = " "
	options.ignoreCase = false
	options.ignoreLeadingBlanks = false
	var key KeyType = KeyType{
		Type: NoKey,
	}
	var expected ContentType = ContentType{
		{Lines: []string{""}, Fields: []string{}, CompareLine: ""},
	}
	var got ContentType = ProcessInputFiles(input, key)
	if !got.isEqual(expected) {
		t.Errorf("got\n%s\nexpected\n%s", got, expected)
	}
}

func TestProcessInputFilesMultilineWithOddNumber(t *testing.T) {
	var input []string = []string{
		"line1",
		"line2",
		"line3",
	}
	options.multiline = 2
	options.fieldSeparator = " "
	options.ignoreCase = false
	options.ignoreLeadingBlanks = false
	var key KeyType = KeyType{
		Type: NoKey,
	}
	var expected ContentType = ContentType{
		{Lines: []string{"line1", "line2"}, Fields: nil, CompareLine: "line1 line2"},
		{Lines: []string{"line3"}, Fields: nil, CompareLine: "line3"},
	}
	var got ContentType = ProcessInputFiles(input, key)
	if !got.isEqual(expected) {
		t.Errorf("got\n%s\nexpected\n%s", got, expected)
	}
}

func TestProcessInputFilesBothCaseAndLeadingBlanks(t *testing.T) {
	var input []string = []string{
		"  UPPER case",
		" \t Mixed CASE",
	}
	options.multiline = 1
	options.fieldSeparator = " "
	options.ignoreCase = true
	options.ignoreLeadingBlanks = true
	var key KeyType = KeyType{
		Type: NoKey,
	}
	var expected ContentType = ContentType{
		{Lines: []string{"  UPPER case"}, Fields: nil, CompareLine: "upper case"},
		{Lines: []string{" \t Mixed CASE"}, Fields: nil, CompareLine: "mixed case"},
	}
	var got ContentType = ProcessInputFiles(input, key)
	if !got.isEqual(expected) {
		t.Errorf("got\n%s\nexpected\n%s", got, expected)
	}
	options.ignoreCase = false
	options.ignoreLeadingBlanks = false
}
