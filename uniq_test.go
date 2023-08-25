package main

import (
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
	if !got.isEqual(expected) {
		t.Errorf("got\n%s\nexpected\n%s", got, expected)
	}
	if input.isEqual(got) {
		t.Errorf("the input was modified")
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
	if !got.isEqual(expected) {
		t.Errorf("got\n%s\nexpected\n%s", got, expected)
	}
	if input.isEqual(got) {
		t.Errorf("the input was modified")
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
	if !got.isEqual(expected) {
		t.Errorf("got\n%s\nexpected\n%s", got, expected)
	}
	if input.isEqual(got) {
		t.Errorf("the input was modified")
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
	if !got.isEqual(expected) {
		t.Errorf("got\n%s\nexpected\n%s", got, expected)
	}
	if input.isEqual(got) {
		t.Errorf("the input was modified")
	}
}
