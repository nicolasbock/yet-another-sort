package main

import (
	"fmt"
	"os"
	"runtime/pprof"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// Program version
var Version = "unknown"

// Configuration options.
var cpuprofile string
var debug bool
var fieldSeparator string = " "
var files []string = []string{}
var forceOutput bool
var ignoreCase bool
var ignoreLeadingBlanks bool
var key KeyType
var memprofile string
var multiline int = 1
var output string
var printVersion bool
var sortMode SortMode
var uniq UniqMode

// initializeLogging initializes the logger.
func initializeLogging() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
}

func main() {
	initializeLogging()
	parseCommandLine()
	if debug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}
	if cpuprofile != "" {
		f, err := os.Create(cpuprofile)
		if err != nil {
			log.Fatal().Msg(err.Error())
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

	var concatenatedContents ContentType = LoadInputFiles(files, key)
	var sortedContents ContentType = SortContents(concatenatedContents)
	var uniqContents ContentType = UniqContents(sortedContents)

	if memprofile != "" {
		f, err := os.Create(memprofile)
		if err != nil {
			log.Fatal().Msg(err.Error())
		}
		pprof.WriteHeapProfile(f)
		f.Close()
		return
	}

	var fd *os.File
	if output != "" {
		_, err := os.Stat(output)
		if err == nil && !forceOutput {
			log.Fatal().Msgf("output file %s already exists", output)
		}
		fd, err = os.Create(output)
		if err != nil {
			log.Fatal().Msgf("cannot open output file %s: %s", output, err.Error())
		}
		defer fd.Close()
	} else {
		fd = os.Stdout
	}
	for _, multiline := range uniqContents {
		for _, line := range multiline.Lines {
			fmt.Fprintln(fd, line)
		}
	}
}
