package main

import (
	"strings"

	"github.com/rs/zerolog/log"
)

// UniqContents compares adjacent entries in contents and keeps only the first
// or the last of identical entries.
//
// The uniq operations is not done in place. A copy of the uniqified ContentType is returned.
func UniqContents(contents ContentType) ContentType {
	// If uniq is disabled, return the original slice without copying
	if options.uniq == no_uniq {
		return contents
	}

	var result ContentType = append(ContentType{}, contents...)
	if options.uniq != no_uniq {
		for i := 0; ; {
			if i == len(result)-1 {
				break
			}
			if strings.Compare(result[i].CompareLine, result[i+1].CompareLine) == 0 {
				log.Debug().Msgf("Found identical lines: '%s' <-> '%s'", result[i].CompareLine, result[i+1].CompareLine)
				switch options.uniq {
				case first:
					result = append(result[:i+1], result[i+2:]...)
				case last:
					result = append(result[:i], result[i+1:]...)
				}
			} else {
				i++
			}
		}
	}
	return result
}
