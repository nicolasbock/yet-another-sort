package main

import (
	"bytes"
	"io"
	"os"
	"unsafe"

	"github.com/rs/zerolog/log"
)

// LoadFile loads the text file `filename` and returns an array of strings, the
// lines in that text file. The special filename `-` means standard input.
//
// The entire file is read into a single []byte buffer. Line strings are created
// as zero-copy views into that buffer via unsafe.String, so no per-line heap
// allocation is needed. Because those strings point at the original bytes, the
// buffer must not be modified for as long as any returned string is in use.
// The backing array remains reachable through the returned strings, so it stays
// alive for the lifetime of those strings.
func LoadFile(filename string) []string {
	var data []byte
	var err error

	if filename != "-" {
		log.Debug().Msgf("Loading contents of file %s", filename)
		data, err = os.ReadFile(filename)
		if err != nil {
			log.Fatal().Msgf("Error opening file %s: %s\n", filename, err.Error())
			os.Exit(1)
		}
	} else {
		log.Debug().Msgf("Reading from standard input")
		data, err = io.ReadAll(os.Stdin)
		if err != nil {
			log.Fatal().Msgf("Error reading standard input: %s", err.Error())
		}
	}

	// Strip a trailing newline so we don't produce a spurious empty line at
	// the end of the slice.
	if len(data) > 0 && data[len(data)-1] == '\n' {
		data = data[:len(data)-1]
	}

	if len(data) == 0 {
		return []string{}
	}

	// Count newlines to pre-allocate the lines slice with the exact capacity,
	// avoiding any reallocation during the split loop.
	numLines := bytes.Count(data, []byte{'\n'}) + 1
	lines := make([]string, 0, numLines)

	// Walk the buffer and produce zero-copy string views into it.
	// unsafe.String(ptr, len) constructs a string header pointing directly at
	// the existing bytes — no copy is made.
	rest := data
	for {
		idx := bytes.IndexByte(rest, '\n')
		if idx < 0 {
			lines = append(lines, unsafe.String(unsafe.SliceData(rest), len(rest)))
			break
		}
		lines = append(lines, unsafe.String(unsafe.SliceData(rest), idx))
		rest = rest[idx+1:]
	}

	return lines
}

// LoadInputFiles loads the input file(s) and returns a concatenated list of lines.
func LoadInputFiles(filenames []string, key KeyType) ContentType {
	var lines []string

	// Pre-allocate for multiple files
	if len(filenames) == 1 {
		lines = LoadFile(filenames[0])
	} else {
		// Estimate total capacity for multiple files
		lines = make([]string, 0, 16384*len(filenames))
		for _, file := range filenames {
			lines = append(lines, LoadFile(file)...)
		}
	}

	contents := ProcessInputFiles(lines, key)
	log.Debug().Msgf("Read %d files", len(filenames))
	return contents
}
