package main

import (
	"reflect"
	"testing"
)

func TestSort1(t *testing.T) {
	var expected ContentType = ContentType{
		ContentLineType{Lines: []string{"1. First line"}, Fields: []string{"1.", "First", "line"}},
		ContentLineType{Lines: []string{"2. Second line"}, Fields: []string{"2.", "Second", "line"}},
		ContentLineType{Lines: []string{"3. Third line"}, Fields: []string{"3.", "Third", "line"}},
		ContentLineType{Lines: []string{"4. Fourth line"}, Fields: []string{"4.", "Fourth", "line"}},
		ContentLineType{Lines: []string{"5. Fifth line"}, Fields: []string{"5.", "Fifth", "line"}},
		}
	var got ContentType = SortContents(input, 1)
	if !reflect.DeepEqual(got, expected) {
		t.Errorf("got %s, expected %s", got, expected)
	}
}

func TestSort2(t *testing.T) {
	var expected ContentType = ContentType{
		ContentLineType{Lines: []string{"5. Fifth line"}, Fields: []string{"5.", "Fifth", "line"}},
		ContentLineType{Lines: []string{"1. First line"}, Fields: []string{"1.", "First", "line"}},
		ContentLineType{Lines: []string{"4. Fourth line"}, Fields: []string{"4.", "Fourth", "line"}},
		ContentLineType{Lines: []string{"2. Second line"}, Fields: []string{"2.", "Second", "line"}},
		ContentLineType{Lines: []string{"3. Third line"}, Fields: []string{"3.", "Third", "line"}},
		}
	var got ContentType = SortContents(input, 2)
	if !reflect.DeepEqual(got, expected) {
		t.Errorf("got %s, expected %s", got, expected)
	}
}
