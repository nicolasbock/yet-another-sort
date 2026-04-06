package main

import (
	"sort"
	"strconv"
	"strings"
)

// isBashTimestamp reports whether s matches the bash history timestamp format:
// a '#' character followed by one or more ASCII decimal digits that represent a
// value in the valid int64 range, and nothing else. The empty string is rejected.
func isBashTimestamp(s string) bool {
	if len(s) < 2 || s[0] != '#' {
		return false
	}
	for i := 1; i < len(s); i++ {
		if s[i] < '0' || s[i] > '9' {
			return false
		}
	}
	_, err := strconv.ParseInt(s[1:], 10, 64)
	return err == nil
}

// ParseBashHistory parses lines from a bash history file into ContentType.
// Each ContentLineType has:
//
//	Lines[0]    = the "#timestamp" line
//	Lines[1:]   = the command lines (one or more)
//	CompareLine = command lines joined by "\n" (used for deduplication)
//
// Records with no command lines (a timestamp immediately followed by another
// timestamp or EOF) are silently dropped.
// Leading non-timestamp lines (before the first timestamp is seen) are
// silently skipped.
func ParseBashHistory(lines []string) ContentType {
	contents := make(ContentType, 0)

	// Index of the current timestamp line; -1 means we haven't seen one yet.
	tsIdx := -1

	for i := 0; i <= len(lines); i++ {
		// Treat the position just past the last line as a virtual "new timestamp"
		// so we can flush the final record without duplicating logic.
		isTS := false
		if i < len(lines) {
			isTS = isBashTimestamp(lines[i])
		}

		if i == len(lines) || isTS {
			// Flush the previous record if we have one.
			if tsIdx >= 0 {
				cmdLines := lines[tsIdx+1 : i]
				if len(cmdLines) > 0 {
					record := ContentLineType{
						Lines:       make([]string, 0, 1+len(cmdLines)),
						CompareLine: strings.Join(cmdLines, "\n"),
					}
					record.Lines = append(record.Lines, lines[tsIdx])
					record.Lines = append(record.Lines, cmdLines...)
					contents = append(contents, record)
				}
				// else: no command lines — silently drop
			}
			// Start a new record at the current timestamp line.
			if i < len(lines) {
				tsIdx = i
			}
		}
		// Non-timestamp lines before the first timestamp are skipped because
		// tsIdx is still -1 and we only flush when tsIdx >= 0.
	}

	return contents
}

// DeduplicateBashHistory keeps, for each unique command block (CompareLine),
// only the record with the numerically largest timestamp.
// Input need not be sorted. The returned ContentType preserves the relative
// order of the surviving (first-seen) entries, though callers typically
// sort afterwards.
func DeduplicateBashHistory(contents ContentType) ContentType {
	// Map from command block to the winning record.
	best := make(map[string]ContentLineType, len(contents))

	for _, record := range contents {
		key := record.CompareLine
		existing, seen := best[key]
		if !seen {
			best[key] = record
			continue
		}
		// Compare timestamps numerically. Lines[0] is "#NNNNNNNNNN".
		curTS := parseTimestamp(record.Lines[0])
		existTS := parseTimestamp(existing.Lines[0])
		if curTS > existTS {
			best[key] = record
		}
	}

	// Rebuild in a stable order: iterate the original slice and emit each
	// command block the first time its winner is encountered.
	seen := make(map[string]bool, len(best))
	result := make(ContentType, 0, len(best))
	for _, record := range contents {
		key := record.CompareLine
		if seen[key] {
			continue
		}
		seen[key] = true
		result = append(result, best[key])
	}
	return result
}

// parseTimestamp extracts the integer value from a bash history timestamp
// string of the form "#NNNNNNNNNN". Returns 0 if the string is malformed.
func parseTimestamp(s string) int64 {
	if len(s) < 2 || s[0] != '#' {
		return 0
	}
	v, err := strconv.ParseInt(s[1:], 10, 64)
	if err != nil {
		return 0
	}
	return v
}

// bashHistoryByTimestamp is a thin sort.Interface wrapper for ContentType that
// orders records by their timestamp (Lines[0]) numerically ascending.
// This follows the same pattern as the sort.Interface implementation on
// ContentType in sort_merge_sort.go, but uses numeric rather than lexicographic
// comparison so that timestamps of varying digit lengths sort correctly.
type bashHistoryByTimestamp struct {
	records ContentType
}

// Len returns the number of records.
func (b bashHistoryByTimestamp) Len() int { return len(b.records) }

// Less reports whether record i should sort before record j by comparing their
// timestamps as int64 values.
func (b bashHistoryByTimestamp) Less(i, j int) bool {
	return parseTimestamp(b.records[i].Lines[0]) < parseTimestamp(b.records[j].Lines[0])
}

// Swap exchanges records at positions i and j.
func (b bashHistoryByTimestamp) Swap(i, j int) {
	b.records[i], b.records[j] = b.records[j], b.records[i]
}

// SortBashHistory sorts records by timestamp ascending.
// Timestamps are "#NNNNNNNNNN"; the leading '#' is stripped and the remainder
// is compared as an int64, ensuring correct ordering even if timestamp strings
// differ in length.
func SortBashHistory(contents ContentType) ContentType {
	if len(contents) <= 1 {
		return contents
	}
	sorted := make(ContentType, len(contents))
	copy(sorted, contents)
	records := bashHistoryByTimestamp{records: sorted}
	if options.stableSort {
		sort.Stable(records)
	} else {
		sort.Sort(records)
	}
	return sorted
}

// RunBashHistoryMode is called from main() when --bash-history is set.
// It loads all input files, parses them as bash history, deduplicates command
// blocks (keeping the record with the latest timestamp for each), sorts the
// surviving records by timestamp ascending, and returns the final ContentType.
func RunBashHistoryMode(filenames []string) ContentType {
	var lines []string
	if len(filenames) == 1 {
		lines = LoadFile(filenames[0])
	} else {
		lines = make([]string, 0, 16384*len(filenames))
		for _, f := range filenames {
			lines = append(lines, LoadFile(f)...)
		}
	}

	parsed := ParseBashHistory(lines)
	deduped := DeduplicateBashHistory(parsed)
	sorted := SortBashHistory(deduped)
	return sorted
}
