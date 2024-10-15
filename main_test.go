package main

import (
	"flag"
	"os"
	"testing"
)

func equalSlices(a, b []int) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func TestParseCommandLine(t *testing.T) {
	tests := []struct {
		args     []string
		expected ConfigurationOptions
	}{
		{
			args: []string{"cmd", "--debug", "--field-separator", ",", "--force", "--ignore-case", "--ignore-leading-blanks", "--key", "1,2", "--multiline", "2", "--output", "output.txt", "--uniq", "first", "--version"},
			expected: ConfigurationOptions{
				debug:               true,
				fieldSeparator:      ",",
				forceOutput:         true,
				ignoreCase:          true,
				ignoreLeadingBlanks: true,
				key:                 KeyType{NoKey, []int{1, 2}},
				multiline:           2,
				output:              "output.txt",
				printVersion:        true,
				uniq:                first,
			},
		},
		{
			args: []string{"cmd", "--ignore-leading-whitespace"},
			expected: ConfigurationOptions{
				ignoreLeadingBlanks: true,
			},
		},
	}

	for _, test := range tests {
		// Reset the options and flag.CommandLine for each test
		options = NewConfigurationOptions()
		flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)

		// Set the command line arguments
		os.Args = test.args

		// Parse the command line
		parseCommandLine()

		// Check the parsed options
		if options.debug != test.expected.debug {
			t.Errorf("expected debug %v, got %v", test.expected.debug, options.debug)
		}
		if options.fieldSeparator != test.expected.fieldSeparator {
			t.Errorf("expected fieldSeparator %v, got %v", test.expected.fieldSeparator, options.fieldSeparator)
		}
		if options.forceOutput != test.expected.forceOutput {
			t.Errorf("expected forceOutput %v, got %v", test.expected.forceOutput, options.forceOutput)
		}
		if options.ignoreCase != test.expected.ignoreCase {
			t.Errorf("expected ignoreCase %v, got %v", test.expected.ignoreCase, options.ignoreCase)
		}
		if options.ignoreLeadingBlanks != test.expected.ignoreLeadingBlanks {
			t.Errorf("expected ignoreLeadingBlanks %v, got %v", test.expected.ignoreLeadingBlanks, options.ignoreLeadingBlanks)
		}
		if options.key.Type != test.expected.key.Type {
			t.Errorf("expected key %v, got %v", test.expected.key, options.key)
		}
		if equalSlices(options.key.Keys, test.expected.key.Keys) {
			t.Errorf("expected key %v, got %v", test.expected.key, options.key)
		}
		if options.multiline != test.expected.multiline {
			t.Errorf("expected multiline %v, got %v", test.expected.multiline, options.multiline)
		}
		if options.output != test.expected.output {
			t.Errorf("expected output %v, got %v", test.expected.output, options.output)
		}
		if options.printVersion != test.expected.printVersion {
			t.Errorf("expected printVersion %v, got %v", test.expected.printVersion, options.printVersion)
		}
		if options.uniq != test.expected.uniq {
			t.Errorf("expected uniq %v, got %v", test.expected.uniq, options.uniq)
		}
	}
}
