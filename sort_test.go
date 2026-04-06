package main

import (
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
	var got ContentType

	got = SortContents(input)
	if !expected.isEqual(got) {
		t.Errorf("got %s\nexpected %s", got, expected)
	}
	if input.isEqual(got) {
		t.Errorf("input list was modified during sort")
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
	var got ContentType

	got = SortContents(input)
	if !expected.isEqual(got) {
		t.Errorf("got %s\nexpected %s", got, expected)
	}
	if input.isEqual(got) {
		t.Errorf("input list was modified during sort")
	}
}

func TestSort3(t *testing.T) {
	var input ContentType = ContentType{
		{Lines: []string{"line 1"}, Fields: []string{"line", "1"}, CompareLine: "line 1"},
		{Lines: []string{"Line 2"}, Fields: []string{"Line", "2"}, CompareLine: "Line 2"},
		{Lines: []string{"line 3"}, Fields: []string{"line", "3"}, CompareLine: "line 3"},
		{Lines: []string{"Line 4"}, Fields: []string{"Line", "4"}, CompareLine: "Line 4"},
		{Lines: []string{"line 5"}, Fields: []string{"line", "5"}, CompareLine: "line 5"},
	}
	var expected ContentType = ContentType{
		{Lines: []string{"Line 2"}, Fields: []string{"Line", "2"}, CompareLine: "Line 2"},
		{Lines: []string{"Line 4"}, Fields: []string{"Line", "4"}, CompareLine: "Line 4"},
		{Lines: []string{"line 1"}, Fields: []string{"line", "1"}, CompareLine: "line 1"},
		{Lines: []string{"line 3"}, Fields: []string{"line", "3"}, CompareLine: "line 3"},
		{Lines: []string{"line 5"}, Fields: []string{"line", "5"}, CompareLine: "line 5"},
	}
	var got ContentType

	got = SortContents(input)
	if !expected.isEqual(got) {
		t.Errorf("got\n%s\nexpected\n%s", got, expected)
	}
	if input.isEqual(got) {
		t.Errorf("input list was modified during sort")
	}
}

func TestSortEmpty(t *testing.T) {
	var input ContentType = ContentType{}
	var expected ContentType = ContentType{}
	var got ContentType

	got = SortContents(input)
	if !expected.isEqual(got) {
		t.Errorf("got\n%s\nexpected\n%s", got, expected)
	}
}

func TestSortStable(t *testing.T) {
	// When --stable-sort is enabled, entries with equal CompareLine values must
	// preserve their original relative (input) order.
	savedStableSort := options.stableSort
	options.stableSort = true
	defer func() { options.stableSort = savedStableSort }()

	var input ContentType = ContentType{
		{Lines: []string{"#1773948684", "juju controllers"}, Fields: []string{"#1773948684", "juju", "controllers"}, CompareLine: "juju controllers"},
		{Lines: []string{"#1773949015", "ll"}, Fields: []string{"#1773949015", "ll"}, CompareLine: "ll"},
		{Lines: []string{"#1773949035", "juju controllers"}, Fields: []string{"#1773949035", "juju", "controllers"}, CompareLine: "juju controllers"},
	}
	// After a stable sort by CompareLine the two "juju controllers" entries must
	// remain in their original order: older timestamp first, newer timestamp second.
	var expected ContentType = ContentType{
		{Lines: []string{"#1773948684", "juju controllers"}, Fields: []string{"#1773948684", "juju", "controllers"}, CompareLine: "juju controllers"},
		{Lines: []string{"#1773949035", "juju controllers"}, Fields: []string{"#1773949035", "juju", "controllers"}, CompareLine: "juju controllers"},
		{Lines: []string{"#1773949015", "ll"}, Fields: []string{"#1773949015", "ll"}, CompareLine: "ll"},
	}

	got := SortContents(input)
	if !expected.isEqual(got) {
		t.Errorf("stable sort: got\n%s\nexpected\n%s", got, expected)
	}
}

func TestSortSingleElement(t *testing.T) {
	var input ContentType = ContentType{
		{Lines: []string{"only line"}, Fields: []string{"only", "line"}, CompareLine: "only line"},
	}
	var expected ContentType = ContentType{
		{Lines: []string{"only line"}, Fields: []string{"only", "line"}, CompareLine: "only line"},
	}
	var got ContentType

	got = SortContents(input)
	if !expected.isEqual(got) {
		t.Errorf("got\n%s\nexpected\n%s", got, expected)
	}
	// For single element, the function returns the same slice (no copy needed)
	// so we don't check for modification in this case
}
