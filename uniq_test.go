package main

import (
	"reflect"
	"testing"
)

func TestUniq1(t *testing.T) {
	var input ContentType = ContentType{
		ContentLineType{Lines: []string{"First line"}, Fields: []string{"First", "line"}, CompareLine: "First"},
		ContentLineType{Lines: []string{"First line"}, Fields: []string{"First", "line"}, CompareLine: "First"},
		ContentLineType{Lines: []string{"Second line"}, Fields: []string{"Second", "line"}, CompareLine: "Second"},
	}
	var expected ContentType = ContentType{
		ContentLineType{Lines: []string{"First line"}, Fields: []string{"First", "line"}, CompareLine: "First"},
		ContentLineType{Lines: []string{"Second line"}, Fields: []string{"Second", "line"}, CompareLine: "Second"},
	}
	uniq = first
	var got = UniqContents(input)
	if !reflect.DeepEqual(got, expected) {
		t.Errorf("got\n%s\nexpected\n%s", got, expected)
	}
}

func TestUniq2(t *testing.T) {
	var input ContentType = ContentType{
		ContentLineType{Lines: []string{"First line"}, Fields: []string{"First", "line"}, CompareLine: "line"},
		ContentLineType{Lines: []string{"First line"}, Fields: []string{"First", "line"}, CompareLine: "line"},
		ContentLineType{Lines: []string{"Second line"}, Fields: []string{"Second", "line"}, CompareLine: "line"},
	}
	var expected ContentType = ContentType{
		ContentLineType{Lines: []string{"First line"}, Fields: []string{"First", "line"}, CompareLine: "line"},
	}
	uniq = first
	var got = UniqContents(input)
	if !reflect.DeepEqual(got, expected) {
		t.Errorf("got\n%s\nexpected\n%s", got, expected)
	}
}

func TestUniq3(t *testing.T) {
	var input ContentType = ContentType{
		ContentLineType{Lines: []string{"First line"}, Fields: []string{"First", "line"}, CompareLine: "line"},
		ContentLineType{Lines: []string{"First line "}, Fields: []string{"First", "line"}, CompareLine: "line"},
		ContentLineType{Lines: []string{"Second line"}, Fields: []string{"Second", "line"}, CompareLine: "line "},
	}
	var expected ContentType = ContentType{
		ContentLineType{Lines: []string{"First line"}, Fields: []string{"First", "line"}, CompareLine: "line"},
		ContentLineType{Lines: []string{"Second line"}, Fields: []string{"Second", "line"}, CompareLine: "line "},
	}
	uniq = first
	var got = UniqContents(input)
	if !reflect.DeepEqual(got, expected) {
		t.Errorf("got\n%s\nexpected\n%s", got, expected)
	}
}

func TestUniq4(t *testing.T) {
	var input ContentType = ContentType{
		ContentLineType{Lines: []string{"First line"}, Fields: []string{"First", "line"}, CompareLine: "line"},
		ContentLineType{Lines: []string{"First line "}, Fields: []string{"First", "line"}, CompareLine: "line"},
		ContentLineType{Lines: []string{"Second line"}, Fields: []string{"Second", "line"}, CompareLine: "line"},
	}
	var expected ContentType = ContentType{
		ContentLineType{Lines: []string{"Second line"}, Fields: []string{"Second", "line"}, CompareLine: "line"},
	}
	uniq = last
	var got = UniqContents(input)
	if !reflect.DeepEqual(got, expected) {
		t.Errorf("got\n%s\nexpected\n%s", got, expected)
	}
}
