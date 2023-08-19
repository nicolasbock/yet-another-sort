package main

import (
	"reflect"
	"testing"
)

func TestSort1(t *testing.T) {
	var input ContentType = ContentType{
		ContentLineType{Lines: []string{"2. Second line"}, Fields: []string{"2.", "Second", "line"}, CompareLine: "2. Second line"},
		ContentLineType{Lines: []string{"1. First line"}, Fields: []string{"1.", "First", "line"}, CompareLine: "1. First line"},
		ContentLineType{Lines: []string{"4. Fourth line"}, Fields: []string{"4.", "Fourth", "line"}, CompareLine: "4. Fourth line"},
		ContentLineType{Lines: []string{"5. Fifth line"}, Fields: []string{"5.", "Fifth", "line"}, CompareLine: "5. Fifth line"},
		ContentLineType{Lines: []string{"3. Third line"}, Fields: []string{"3.", "Third", "line"}, CompareLine: "3. Third line"},
	}
	var expected ContentType = ContentType{
		ContentLineType{Lines: []string{"1. First line"}, Fields: []string{"1.", "First", "line"}, CompareLine: "1. First line"},
		ContentLineType{Lines: []string{"2. Second line"}, Fields: []string{"2.", "Second", "line"}, CompareLine: "2. Second line"},
		ContentLineType{Lines: []string{"3. Third line"}, Fields: []string{"3.", "Third", "line"}, CompareLine: "3. Third line"},
		ContentLineType{Lines: []string{"4. Fourth line"}, Fields: []string{"4.", "Fourth", "line"}, CompareLine: "4. Fourth line"},
		ContentLineType{Lines: []string{"5. Fifth line"}, Fields: []string{"5.", "Fifth", "line"}, CompareLine: "5. Fifth line"},
	}
	var got ContentType = SortContents(input)
	if !reflect.DeepEqual(got, expected) {
		t.Errorf("got %s\nexpected %s", got, expected)
	}
}

func TestSort2(t *testing.T) {
	var input ContentType = ContentType{
		ContentLineType{Lines: []string{"2. Second line"}, Fields: []string{"2.", "Second", "line"}, CompareLine: "2."},
		ContentLineType{Lines: []string{"1. First line"}, Fields: []string{"1.", "First", "line"}, CompareLine: "1."},
		ContentLineType{Lines: []string{"4. Fourth line"}, Fields: []string{"4.", "Fourth", "line"}, CompareLine: "4."},
		ContentLineType{Lines: []string{"5. Fifth line"}, Fields: []string{"5.", "Fifth", "line"}, CompareLine: "5."},
		ContentLineType{Lines: []string{"3. Third line"}, Fields: []string{"3.", "Third", "line"}, CompareLine: "3."},
	}
	var expected ContentType = ContentType{
		ContentLineType{Lines: []string{"1. First line"}, Fields: []string{"1.", "First", "line"}, CompareLine: "1."},
		ContentLineType{Lines: []string{"2. Second line"}, Fields: []string{"2.", "Second", "line"}, CompareLine: "2."},
		ContentLineType{Lines: []string{"3. Third line"}, Fields: []string{"3.", "Third", "line"}, CompareLine: "3."},
		ContentLineType{Lines: []string{"4. Fourth line"}, Fields: []string{"4.", "Fourth", "line"}, CompareLine: "4."},
		ContentLineType{Lines: []string{"5. Fifth line"}, Fields: []string{"5.", "Fifth", "line"}, CompareLine: "5."},
	}
	var got ContentType = SortContents(input)
	if !reflect.DeepEqual(got, expected) {
		t.Errorf("got %s\nexpected %s", got, expected)
	}
}
