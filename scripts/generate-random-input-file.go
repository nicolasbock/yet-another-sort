package main

import (
	"bufio"
	"flag"
	"math/rand"
	"os"
	"strings"
	"time"
)

const (
	// Buffer size for output writer (256KB)
	bufferSize = 256 * 1024
	// ASCII letters for random generation
	asciiLetters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
)

type options struct {
	lines                int
	fields               int
	fieldLength          int
	addLeadingWhitespace bool
}

func parseOptions() options {
	var opts options

	flag.IntVar(&opts.lines, "lines", 1, "The number of lines to generate")
	flag.IntVar(&opts.fields, "fields", 1, "The number of fields per line")
	flag.IntVar(&opts.fieldLength, "field-length", 1, "The number of characters of a field")
	flag.BoolVar(&opts.addLeadingWhitespace, "add-leading-whitespace", false, "Add random numbers of leading whitespace")

	flag.Parse()
	return opts
}

// Pre-allocated slice for leading whitespace choices
var leadingChoices = []int{0, 0, 0, 0, 1, 2, 3}

// generateField generates a random field of the specified length
func generateField(fieldLength int) string {
	// Pre-allocate byte slice for better performance
	field := make([]byte, fieldLength)
	for i := 0; i < fieldLength; i++ {
		field[i] = asciiLetters[rand.Intn(len(asciiLetters))]
	}
	return string(field)
}

// generateLine generates a single line with the specified number of fields
func generateLine(opts options, builder *strings.Builder) {
	builder.Reset()

	// Add leading whitespace if requested
	if opts.addLeadingWhitespace {
		spaces := leadingChoices[rand.Intn(len(leadingChoices))]
		for i := 0; i < spaces; i++ {
			builder.WriteByte(' ')
		}
	}

	// Generate fields
	for i := 0; i < opts.fields; i++ {
		if i > 0 {
			builder.WriteByte(' ')
		}
		builder.WriteString(generateField(opts.fieldLength))
	}
}

func main() {
	// Seed the random number generator
	rand.Seed(time.Now().UnixNano())

	opts := parseOptions()

	// Use buffered writer for better I/O performance
	writer := bufio.NewWriterSize(os.Stdout, bufferSize)
	defer writer.Flush()

	// Pre-allocate string builder with estimated capacity
	// Estimated line length: (fieldLength * fields) + (fields - 1 spaces) + potential leading spaces
	estimatedLineLength := (opts.fieldLength * opts.fields) + opts.fields + 5
	builder := strings.Builder{}
	builder.Grow(estimatedLineLength)

	// Generate lines
	for i := 0; i < opts.lines; i++ {
		generateLine(opts, &builder)
		writer.WriteString(builder.String())
		writer.WriteByte('\n')
	}
}
