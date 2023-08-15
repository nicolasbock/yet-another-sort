package main

import (
	"fmt"
	"github.com/akamensky/argparse"
	"os"
)

func main() {
	var parser *argparse.Parser = argparse.NewParser("yet-another-sort",
		`Yet-another-sort is yet another command line utility that sorts lines in a file.`)
	var err error = parser.Parse(os.Args)
	if err != nil {
		fmt.Println(parser.Usage(err))
	}
}
