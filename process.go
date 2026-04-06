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

	// Pre-allocate contents slice with exact size to avoid reallocation
	numMultilines := (len(lines) + options.multiline - 1) / options.multiline
	contents := make(ContentType, 0, numMultilines)

	for i := 0; i < len(lines); i += options.multiline {
		var contentLine ContentLineType

		// Pre-allocate slices with known capacity
		contentLine.Lines = make([]string, 0, options.multiline)

		// For NoKey with the default space separator we never need the Fields
		// slice: CompareLine can be set directly to Lines[0] (multiline=1) or
		// joined lines (multiline>1), avoiding both the field-split allocations
		// and the strings.Builder reconstruction.
		// For other key types, or when a non-space separator is used (where
		// empty fields are collapsed and affect CompareLine), we must still
		// split into fields.
		needFields := key.Type != NoKey || options.fieldSeparator != " "

		var allFields []string

		// Process all lines in this multiline group
		for j := 0; j < options.multiline; j++ {
			linenumber := i + j
			if linenumber >= len(lines) {
				break
			}

			line := lines[linenumber]
			contentLine.Lines = append(contentLine.Lines, line)

			if needFields {
				// Split fields and filter empty ones
				// Using strings.FieldsFunc is more efficient for custom separators
				if options.fieldSeparator == " " {
					// Fast path for space separator - strings.Fields handles multiple spaces
					fields := strings.Fields(line)
					allFields = append(allFields, fields...)
				} else if len(options.fieldSeparator) == 1 {
					// Single character separator - optimize by avoiding repeated string operations
					rawFields := strings.Split(line, options.fieldSeparator)
					for _, field := range rawFields {
						if len(field) > 0 {
							allFields = append(allFields, field)
						}
					}
				} else {
					// Multi-character separator (rare case)
					rawFields := strings.Split(line, options.fieldSeparator)
					for _, field := range rawFields {
						if len(field) > 0 {
							allFields = append(allFields, field)
						}
					}
				}
			}
		}

		contentLine.Fields = allFields

		// Build compare line based on key type
		switch key.Type {
		case NoKey:
			if !needFields {
				// Fast path: space separator, NoKey.
				// CompareLine is the raw line (multiline=1) or lines joined by
				// a single space (multiline>1). No extra allocation needed for
				// the single-line case.
				if len(contentLine.Lines) == 0 {
					contentLine.CompareLine = ""
				} else if len(contentLine.Lines) == 1 {
					contentLine.CompareLine = contentLine.Lines[0]
				} else {
					var builder strings.Builder
					builder.Grow(len(contentLine.Lines) * (len(contentLine.Lines[0]) + 1))
					builder.WriteString(contentLine.Lines[0])
					for k := 1; k < len(contentLine.Lines); k++ {
						builder.WriteString(options.fieldSeparator)
						builder.WriteString(contentLine.Lines[k])
					}
					contentLine.CompareLine = builder.String()
				}
			} else {
				// General path: non-space separator (empty fields are collapsed,
				// so CompareLine differs from the raw line).
				if len(allFields) == 0 {
					contentLine.CompareLine = ""
				} else if len(allFields) == 1 {
					contentLine.CompareLine = allFields[0]
				} else {
					var builder strings.Builder
					builder.Grow(len(allFields) * 10)
					builder.WriteString(allFields[0])
					for i := 1; i < len(allFields); i++ {
						builder.WriteString(options.fieldSeparator)
						builder.WriteString(allFields[i])
					}
					contentLine.CompareLine = builder.String()
				}
			}

		case SingleField:
			if len(allFields) < key.Keys[0] {
				log.Fatal().Msgf("multiline %s does not have enough keys", contentLine.Lines)
			}
			contentLine.CompareLine = allFields[key.Keys[0]-1]

		case Remainder:
			if len(allFields) < key.Keys[0] {
				log.Fatal().Msgf("multiline %s does not have enough keys", contentLine.Lines)
			}
			remainderFields := allFields[key.Keys[0]-1:]
			if len(remainderFields) == 1 {
				contentLine.CompareLine = remainderFields[0]
			} else {
				var builder strings.Builder
				builder.Grow(len(remainderFields) * 10)
				builder.WriteString(remainderFields[0])
				for i := 1; i < len(remainderFields); i++ {
					builder.WriteString(options.fieldSeparator)
					builder.WriteString(remainderFields[i])
				}
				contentLine.CompareLine = builder.String()
			}

		case SubSet:
			if len(allFields) < key.Keys[1] {
				log.Fatal().Msgf("multiline %s does not have enough keys", contentLine.Lines)
			}
			subsetFields := allFields[key.Keys[0]-1 : key.Keys[1]]
			if len(subsetFields) == 1 {
				contentLine.CompareLine = subsetFields[0]
			} else {
				var builder strings.Builder
				builder.Grow(len(subsetFields) * 10)
				builder.WriteString(subsetFields[0])
				for i := 1; i < len(subsetFields); i++ {
					builder.WriteString(options.fieldSeparator)
					builder.WriteString(subsetFields[i])
				}
				contentLine.CompareLine = builder.String()
			}
		}

		// Apply transformations
		if options.ignoreLeadingBlanks {
			contentLine.CompareLine = strings.TrimLeft(contentLine.CompareLine, " \t")
		}

		if options.ignoreCase {
			contentLine.CompareLine = strings.ToLower(contentLine.CompareLine)
		}

		contents = append(contents, contentLine)
	}

	return contents
}
