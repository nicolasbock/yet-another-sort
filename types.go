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

// KeyType stores the type of key and the keys.
type KeyType struct {
	Type KeyT
	Keys []int
}

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

func (k KeyType) String() string {
	var result string
	switch k.Type {
	case NoKey:
		result = "whole line"
	}
	if len(k.Keys) > 0 {
		var keyStrings []string = []string{}
		for _, key := range k.Keys {
			keyStrings = append(keyStrings, fmt.Sprint(key))
		}
		result += ", [" + strings.Join(keyStrings, ", ") + "]"
	}
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
				log.Fatal().Msgf("key sepcification inverted: %d > %d (%s)", k.Keys[0], k.Keys[1], s)
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
		result += fmt.Sprintf("  line: '%s'\n", l.Lines[i])
	}
	result += fmt.Sprintf("  fields: %s\n", strings.Join(l.Fields, ", "))
	result += fmt.Sprintf("  compare: '%s'", l.CompareLine)
	return result
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
