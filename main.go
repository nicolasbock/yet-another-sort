package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"path"
)

var files []string = []string{}
var debug bool

var fileContents [][]string = [][]string{}

func main() {
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Usage: %s [OPTION]... [FILE]...\n", path.Clean(os.Args[0]))
		fmt.Fprintln(flag.CommandLine.Output(), `
Write sorted concatenation of all FILE(s) to standard output.

With no FILE, or when FILE is -, read standard input.

Options:`)
		flag.PrintDefaults()
	}
	flag.BoolVar(&debug, "debug", false, "Print debugging information.")
	flag.Parse()
	if flag.NArg() == 0 {
		files = append(files, "-")
	} else {
		files = append(files, flag.Args()...)
	}
	fmt.Printf("Sorting %s\n", files)

	for _, file := range files {
		fmt.Printf("Loading contents of file %s\n", file)
		fd, err := os.Open(file)
		if err != nil {
			fmt.Printf("Error openeing file %s: %s\n", file, err.Error())
			os.Exit(1)
		}
		defer fd.Close()
		fs := bufio.NewScanner(fd)
		fs.Split(bufio.ScanLines)
		fileContents = append(fileContents, []string{})
		for fs.Scan() {
			fileContents[len(fileContents)-1] = append(fileContents[len(fileContents)-1], fs.Text())
		}
	}

	fmt.Printf("Read %d files\n", len(files))

	for i := range fileContents {
		for _, line := range fileContents[i] {
			fmt.Printf("%s\n", line)
		}
	}
}
