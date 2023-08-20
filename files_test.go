package main

import (
	"reflect"
	"testing"
)

func compareContentTypes(a, b ContentType, t *testing.T) {
	if !reflect.DeepEqual(a, b) {
		t.Errorf("got\n%s\nexpected\n%s", a, b)
	}
}

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
	compareContentTypes(got, expected, t)
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
	compareContentTypes(got, expected, t)
}

// func TestProcessInputFiles3(t *testing.T) {
	// var input []string = []string{
		// "Line 1 a",
		// "Line 2 b",
		// "Line 3 c",
		// "Line 4 d",
	// }
	// multiline = 3
	// var key KeyType = KeyType{
		// Type: SingleField,
		// Keys: []int{2},
	// }
	// var expected ContentType = ContentType{
		// {Lines: []string{"Line one", "Line two"}, Fields: []string{"Line", "one", "Line", "two"}, CompareLine: "Line one Line two"},
		// {Lines: []string{"Line three", "Line four"}, Fields: []string{"Line", "three", "Line", "four"}, CompareLine: "Line three Line four"},
	// }
	// var got ContentType = ProcessInputFiles(input, key)
	// compareContentTypes(got, expected, t)
// }
