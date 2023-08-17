package main

import (
	"fmt"
	"strings"
)

type ContentLineType struct {
	Lines  []string
	Fields []string
}

func (l ContentLineType) String() string {
	var result string
	result = "multiline\n"
	for i := range l.Lines {
		result += fmt.Sprintf("  line: '%s'\n", l.Lines[i])
	}
	result += fmt.Sprintf("  fields: %s\n", strings.Join(l.Fields, ", "))
	return result
}

type ContentType []ContentLineType

func (c ContentType) String() string {
	var result string
	result = fmt.Sprintf("%d multilines\n", len(c))
	for _, line := range c {
		result += fmt.Sprint(line)
	}
	return result
}

// SortContents sorts the content lines and returns a sorted list.
func SortContents(contents ContentType, key int) (sortedContents ContentType) {
	sortedContents = append(sortedContents, contents...)

	// Combine multilines with field-separator.
	// Sort multilined content
	// Split multilined content into original multilines

	for i := range sortedContents {
		for j := range sortedContents {
			if strings.Compare(sortedContents[i].Fields[key-1], sortedContents[j].Fields[key-1]) < 0 {
				var temp ContentLineType = sortedContents[i]
				sortedContents[i] = sortedContents[j]
				sortedContents[j] = temp
			}
		}
	}
	return sortedContents
}
