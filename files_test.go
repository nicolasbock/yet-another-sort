package main

import (
	"reflect"
	"testing"
)

func TestProcessInputFiles(t *testing.T) {
	var input []string = []string{
		"Line one",
		"Line two",
		"Line three",
	}
	var key KeyType = KeyType{
		Type: NoKey,
	}
	var expected ContentType = ContentType{
		{Lines: []string{"Line one"}, Fields: []string{"Line", "one"}, CompareLine: "Line one"},
		{Lines: []string{"Line two"}, Fields: []string{"Line", "two"}, CompareLine: "Line two"},
		{Lines: []string{"Line three"}, Fields: []string{"Line", "three"}, CompareLine: "Line three"},
	}
	var got ContentType = ProcessInputFiles(input, key)
	if !reflect.DeepEqual(got, expected) {
		t.Errorf("got\n%s\nexpected\n%s", got, expected)
	}
}
