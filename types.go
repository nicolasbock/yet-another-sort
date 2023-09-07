package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/rs/zerolog/log"
)

// KeyT is the type of key.
//
// The general structure of a key specification is `F[,[F]]` where F is the
// field number starting from 1 separated by the field-separator.
type KeyT int

const (
	// NoKey (the default): use all fields, i.e. the whole line
	NoKey KeyT = iota
	// SingleField: use only one field
	SingleField
	// Remainder: open set [F,), i.e. starting with field F the remainder of the
	// line
	Remainder
	// SubSet: all fields in closed set (F1,F2)
	SubSet
)

// KeyType stores the type of key and the keys.
func (k KeyT) String() string {
	switch k {
	case NoKey:
		return "NoKey"
	case SingleField:
		return "SingleField"
	case Remainder:
		return "Remainder"
	case SubSet:
		return "SubSet"
	}
	return ""
}

type KeyType struct {
	Type KeyT
	Keys []int
}

func (k KeyType) String() string {
	var result string = k.Type.String()
	if len(k.Keys) > 0 {
		var keyStrings []string = []string{}
		for _, key := range k.Keys {
			keyStrings = append(keyStrings, fmt.Sprint(key))
		}
		result += ", [" + strings.Join(keyStrings, ", ") + "]"
	}
	return result
}

func (k KeyType) Representation() string {
	var result string = fmt.Sprintf("KeyType{Type: %s", k.Type)
	if len(k.Keys) > 0 {
		result += fmt.Sprintf(", Keys: []int{%s}",
			strings.Join(strings.Fields(strings.Trim(fmt.Sprint(k.Keys), "[]")), ", "))
	}
	result += "}"
	return result
}

func (k *KeyType) Set(s string) error {
	var fields []string = strings.Split(s, ",")
	if len(fields) == 1 {
		k.Type = SingleField
		keyValue, err := strconv.ParseInt(fields[0], 10, 32)
		if err != nil {
			log.Fatal().Msgf("could not parse key %s (%s): %s", fields[0], s, err.Error())
		}
		k.Keys = append(k.Keys, int(keyValue))
	} else if len(fields) == 2 {
		if fields[1] == "" {
			k.Type = Remainder
			keyValue, err := strconv.ParseInt(fields[0], 10, 32)
			if err != nil {
				log.Fatal().Msgf("could not parse key %s (%s): %s", fields[0], s, err.Error())
			}
			k.Keys = append(k.Keys, int(keyValue))
		} else {
			k.Type = SubSet
			for i := 0; i < 2; i++ {
				keyValue, err := strconv.ParseInt(fields[i], 10, 32)
				if err != nil {
					log.Fatal().Msgf("could not parse key %s (%s): %s", fields[i], s, err.Error())
				}
				k.Keys = append(k.Keys, int(keyValue))
			}
			if k.Keys[0] > k.Keys[1] {
				log.Fatal().Msgf("key specification inverted: %d > %d (%s)", k.Keys[0], k.Keys[1], s)
			}
		}
	} else {
		log.Fatal().Msgf("illegal key specification %s", s)
	}
	return nil
}

type ContentLineType struct {
	// Lines are the lines in this multiline
	Lines []string
	// Fields are the fields of the lines split by field-separator
	Fields []string
	// CompareLine is the part of the multiline that is used for comparison.
	CompareLine string
}

func (l ContentLineType) String() string {
	var result string
	result = "multiline\n"
	for i := range l.Lines {
		result += fmt.Sprintf("  line: \"%s\"\n", l.Lines[i])
	}
	var fields []string = []string{}
	for _, field := range l.Fields {
		fields = append(fields, fmt.Sprintf("\"%s\"", field))
	}
	result += fmt.Sprintf("  fields: [ %s ]\n", strings.Join(fields, ", "))
	result += fmt.Sprintf("  compare: \"%s\"", l.CompareLine)
	return result
}

func (a ContentLineType) isEqual(b ContentLineType) bool {
	if len(a.Lines) != len(b.Lines) {
		return false
	}
	if len(a.Fields) != len(b.Fields) {
		return false
	}
	for i := range a.Lines {
		if a.Lines[i] != b.Lines[i] {
			return false
		}
	}
	for i := range a.Fields {
		if a.Fields[i] != b.Fields[i] {
			return false
		}
	}
	if a.CompareLine != b.CompareLine {
		return false
	}
	return true
}

type ContentType []ContentLineType

func (c ContentType) String() string {
	var result string
	result = fmt.Sprintf("%d multilines\n", len(c))
	for _, line := range c {
		result += fmt.Sprintf("%s\n", line)
	}
	return strings.Trim(result, "\n")
}

func (a ContentType) isEqual(b ContentType) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if !a[i].isEqual(b[i]) {
			return false
		}
	}
	return true
}

type UniqMode int

const (
	no_uniq UniqMode = iota
	first
	last
)

func (um UniqMode) String() string {
	switch um {
	case no_uniq:
		return "no uniq"
	case first:
		return "first"
	case last:
		return "last"
	}
	return "FIXME"
}

func (um *UniqMode) Set(s string) error {
	switch s {
	case "none":
		*um = no_uniq
	case "no uniq":
		*um = no_uniq
	case "first":
		*um = first
	case "last":
		*um = last
	default:
		return fmt.Errorf("unknown value %s for UniqMode", s)
	}
	return nil
}

type SortMode int

const (
	bubble = iota
	merge
)

func (sm SortMode) String() string {
	switch sm {
	case bubble:
		return "bubble sort"
	case merge:
		return "merge sort"
	}
	return "FIXME"
}

func (sm *SortMode) Set(s string) error {
	switch s {
	case "bubble":
		*sm = bubble
	case "merge":
		*sm = merge
	default:
		return fmt.Errorf("unknown value %s for SortMode", s)
	}
	return nil
}
