package main

import (
	"flag"
	"os"
	"testing"
)

func TestParseCommandLine(t *testing.T) {
	tests := []struct {
		name           string
		args           []string
		expectedOutput string
		expectedFiles  []string
	}{
		{
			name:           "No arguments",
			args:           []string{},
			expectedOutput: "",
			expectedFiles:  []string{"-"},
		},
		{
			name:           "Single file",
			args:           []string{"file1.txt"},
			expectedOutput: "",
			expectedFiles:  []string{"file1.txt"},
		},
		{
			name:           "Multiple files",
			args:           []string{"file1.txt", "file2.txt"},
			expectedOutput: "",
			expectedFiles:  []string{"file1.txt", "file2.txt"},
		},
		{
			name:           "Debug flag",
			args:           []string{"--debug"},
			expectedOutput: "",
			expectedFiles:  []string{"-"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset the flag.CommandLine to avoid conflicts with other tests
			flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
			options = NewConfigurationOptions()

			// Set the command line arguments
			os.Args = append([]string{os.Args[0]}, tt.args...)

			// Capture the output
			output := captureOutput(func() {
				parseCommandLine()
			})

			// Check the output
			if output != tt.expectedOutput {
				t.Errorf("expected output %q, got %q", tt.expectedOutput, output)
			}

			// Check the files
			if len(options.files) != len(tt.expectedFiles) {
				t.Errorf("expected files %v, got %v", tt.expectedFiles, options.files)
			}
			for i, file := range tt.expectedFiles {
				if options.files[i] != file {
					t.Errorf("expected file %q, got %q", file, options.files[i])
				}
			}
		})
	}
}

func TestParseCommandLineWithFlags(t *testing.T) {
	// Test field separator
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	options = NewConfigurationOptions()
	os.Args = []string{os.Args[0], "--field-separator", ",", "test.txt"}
	parseCommandLine()
	if options.fieldSeparator != "," {
		t.Errorf("Expected field separator ',', got %q", options.fieldSeparator)
	}

	// Test multiline
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	options = NewConfigurationOptions()
	os.Args = []string{os.Args[0], "--multiline", "3", "test.txt"}
	parseCommandLine()
	if options.multiline != 3 {
		t.Errorf("Expected multiline 3, got %d", options.multiline)
	}

	// Test ignore-case
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	options = NewConfigurationOptions()
	os.Args = []string{os.Args[0], "--ignore-case", "test.txt"}
	parseCommandLine()
	if !options.ignoreCase {
		t.Errorf("Expected ignoreCase to be true")
	}

	// Test ignore-leading-blanks
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	options = NewConfigurationOptions()
	os.Args = []string{os.Args[0], "--ignore-leading-blanks", "test.txt"}
	parseCommandLine()
	if !options.ignoreLeadingBlanks {
		t.Errorf("Expected ignoreLeadingBlanks to be true")
	}

	// Test ignore-leading-whitespace (alias)
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	options = NewConfigurationOptions()
	os.Args = []string{os.Args[0], "--ignore-leading-whitespace", "test.txt"}
	parseCommandLine()
	if !options.ignoreLeadingBlanks {
		t.Errorf("Expected ignoreLeadingBlanks to be true via alias")
	}

	// Test force
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	options = NewConfigurationOptions()
	os.Args = []string{os.Args[0], "--force", "test.txt"}
	parseCommandLine()
	if !options.forceOutput {
		t.Errorf("Expected forceOutput to be true")
	}

	// Test output
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	options = NewConfigurationOptions()
	os.Args = []string{os.Args[0], "--output", "output.txt", "test.txt"}
	parseCommandLine()
	if options.output != "output.txt" {
		t.Errorf("Expected output 'output.txt', got %q", options.output)
	}

	// Test key with single field
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	options = NewConfigurationOptions()
	os.Args = []string{os.Args[0], "--key", "2", "test.txt"}
	parseCommandLine()
	if options.key.Type != SingleField || len(options.key.Keys) != 1 || options.key.Keys[0] != 2 {
		t.Errorf("Expected key SingleField [2], got %s", options.key)
	}

	// Test uniq with first
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	options = NewConfigurationOptions()
	os.Args = []string{os.Args[0], "--uniq", "first", "test.txt"}
	parseCommandLine()
	if options.uniq != first {
		t.Errorf("Expected uniq mode 'first', got %s", options.uniq)
	}

	// Test cpuprofile
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	options = NewConfigurationOptions()
	os.Args = []string{os.Args[0], "--cpuprofile", "cpu.prof", "test.txt"}
	parseCommandLine()
	if options.cpuprofile != "cpu.prof" {
		t.Errorf("Expected cpuprofile 'cpu.prof', got %q", options.cpuprofile)
	}

	// Test memprofile
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	options = NewConfigurationOptions()
	os.Args = []string{os.Args[0], "--memprofile", "mem.prof", "test.txt"}
	parseCommandLine()
	if options.memprofile != "mem.prof" {
		t.Errorf("Expected memprofile 'mem.prof', got %q", options.memprofile)
	}
}

// captureOutput captures the output of a function
func captureOutput(f func()) string {
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	f()

	w.Close()
	os.Stdout = old

	var buf [1024]byte
	n, _ := r.Read(buf[:])
	return string(buf[:n])
}
