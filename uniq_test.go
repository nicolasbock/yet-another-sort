package main

import (
	"reflect"
	"testing"
)

var inputUniq ContentType = ContentType{
	ContentLineType{Lines: []string{"First line"}, Fields: []string{"First", "line"}},
	ContentLineType{Lines: []string{"First line"}, Fields: []string{"First", "line"}},
	ContentLineType{Lines: []string{"Second line"}, Fields: []string{"Second", "line"}},
}

func TestUniq1(t *testing.T) {
	var expected ContentType = ContentType{
		ContentLineType{Lines: []string{"First line"}, Fields: []string{"First", "line"}},
		ContentLineType{Lines: []string{"Second line"}, Fields: []string{"Second", "line"}},
	}
	var got = UniqContents(inputUniq, 1)
	if !reflect.DeepEqual(got, expected) {
		t.Errorf("got %s\nexpected %s", got, expected)
	}
}

func TestUniq2(t *testing.T) {
	var expected ContentType = ContentType{
		ContentLineType{Lines: []string{"First line"}, Fields: []string{"First", "line"}},
	}
	var got = UniqContents(inputUniq, 2)
	if !reflect.DeepEqual(got, expected) {
		t.Errorf("got %s\nexpected %s", got, expected)
	}
}
