package main

import (
	"strings"

	"github.com/rs/zerolog/log"
)

func UniqContents(contents ContentType) ContentType {
	var result ContentType = append(ContentType{}, contents...)
	if uniq {
		for i := 0; ; {
			if i == len(result)-1 {
				break
			}
			if strings.Compare(result[i].CompareLine, result[i+1].CompareLine) == 0 {
				log.Debug().Msgf("Found identical lines: '%s' <-> '%s'", result[i].CompareLine, result[i+1].CompareLine)
				result = append(result[:i+1], result[i+2:]...)
			} else {
				i++
			}
		}
	}
	return result
}
