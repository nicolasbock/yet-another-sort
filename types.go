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
