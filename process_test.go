package main

import (
	"testing"
)

func TestProcessInputFiles1(t *testing.T) {
	var input []string = []string{
		"Line one",
		"Line two",
		"Line three",
		"Line four",
	}
	multiline = 1
	var key KeyType = KeyType{
		Type: NoKey,
	}
	var expected ContentType = ContentType{
		{Lines: []string{"Line one"}, Fields: []string{"Line", "one"}, CompareLine: "Line one"},
		{Lines: []string{"Line two"}, Fields: []string{"Line", "two"}, CompareLine: "Line two"},
		{Lines: []string{"Line three"}, Fields: []string{"Line", "three"}, CompareLine: "Line three"},
		{Lines: []string{"Line four"}, Fields: []string{"Line", "four"}, CompareLine: "Line four"},
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
	multiline = 2
	var key KeyType = KeyType{
		Type: NoKey,
	}
	var expected ContentType = ContentType{
		{Lines: []string{"Line one", "Line two"}, Fields: []string{"Line", "one", "Line", "two"}, CompareLine: "Line one Line two"},
		{Lines: []string{"Line three", "Line four"}, Fields: []string{"Line", "three", "Line", "four"}, CompareLine: "Line three Line four"},
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
	multiline = 2
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
		{Lines: []string{"line 1"}, Fields: []string{"line", "1"}, CompareLine: "line 1"},
		{Lines: []string{"Line 2"}, Fields: []string{"Line", "2"}, CompareLine: "line 2"},
		{Lines: []string{"line 3"}, Fields: []string{"line", "3"}, CompareLine: "line 3"},
		{Lines: []string{"Line 4"}, Fields: []string{"Line", "4"}, CompareLine: "line 4"},
		{Lines: []string{"line 5"}, Fields: []string{"line", "5"}, CompareLine: "line 5"},
	}
	var key = KeyType{}
	multiline = 1
	ignoreCase = true
	var got ContentType = ProcessInputFiles(input, key)
	if !got.isEqual(expected) {
		t.Errorf("got\n%s\nexpected\n%s", got, expected)
	}
}

func TestProcessInputFiles5(t *testing.T) {
	var input []string = []string{
		"#1727714684",
		"gh pr create --assignee nicolasbock --fill",
		"#1714578836",
		"gh pr create --assignee nicolasbock --fill",
	}
	var expected ContentType = ContentType{}
	var key = KeyType{}
	multiline = 2
	var got ContentType = ProcessInputFiles(input, key)
	if !got.isEqual(expected) {
		t.Errorf("got\n%s\nexpected\n%s", got, expected)
	}
}
