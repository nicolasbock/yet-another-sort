package main

import (
	"strings"

	"github.com/rs/zerolog/log"
)

// ProcessInputFiles analyzes the input and splits the lines into multilines. It
// uses the key to construct the part of the multiline that will be used for
// comparisons.
func ProcessInputFiles(lines []string, key KeyType) ContentType {
	log.Debug().Msgf("Processing %d lines, key: %s", len(lines), key)
	var contents ContentType = ContentType{}
	var lastContentLine *ContentLineType
	for i := 0; i < len(lines); i += options.multiline {
		contents = append(contents, ContentLineType{})
		lastContentLine = &contents[len(contents)-1]
		for j := 0; j < options.multiline; j++ {
			var linenumber int = i + j
			if linenumber >= len(lines) {
				break
			}

			var rawFields []string = strings.Split(lines[linenumber], options.fieldSeparator)
			var fields []string = []string{}
			for _, field := range rawFields {
				if len(field) > 0 {
					fields = append(fields, field)
				}
			}
			lastContentLine.Lines = append(lastContentLine.Lines, lines[linenumber])
			lastContentLine.Fields = append(lastContentLine.Fields, fields...)
		}

		switch key.Type {
		case NoKey:
			lastContentLine.CompareLine = strings.Join(lastContentLine.Fields, options.fieldSeparator)
		case SingleField:
			if len(lastContentLine.Fields) < key.Keys[0] {
				log.Fatal().Msgf("multiline %s does not have enough keys", lastContentLine.Lines)
			}
			lastContentLine.CompareLine = lastContentLine.Fields[key.Keys[0]-1]
		case Remainder:
			if len(lastContentLine.Fields) < key.Keys[0] {
				log.Fatal().Msgf("multiline %s does not have enough keys", lastContentLine.Lines)
			}
			lastContentLine.CompareLine = strings.Join(
				lastContentLine.Fields[key.Keys[0]-1:], options.fieldSeparator)
		case SubSet:
			if len(lastContentLine.Fields) < key.Keys[1] {
				log.Fatal().Msgf("multiline %s does not have enough keys", lastContentLine.Lines)
			}
			lastContentLine.CompareLine = strings.Join(
				lastContentLine.Fields[key.Keys[0]-1:key.Keys[1]], options.fieldSeparator)
		}

		if options.ignoreLeadingBlanks {
			lastContentLine.CompareLine = strings.TrimLeft(lastContentLine.CompareLine, " \t")
		}

		if options.ignoreCase {
			lastContentLine.CompareLine = strings.ToLower(lastContentLine.CompareLine)
		}
	}
	return contents
}
