package main

import (
	"os"
	"testing"
)

// parseArgs simulates parsing command-line arguments using the cobra root command.
// It resets options before each parse and only parses flags without executing run().
func parseArgs(args []string) {
	options = NewConfigurationOptions()
	cmd := buildRootCmd()
	// ParseFlags stops before RunE, so no file I/O happens.
	cmd.ParseFlags(args) //nolint:errcheck
	// Populate options.files from remaining args (mirroring run() logic)
	remaining := cmd.Flags().Args()
	if len(remaining) == 0 {
		options.files = append(options.files, "-")
	} else {
		options.files = append(options.files, remaining...)
	}
}

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
			options = NewConfigurationOptions()
			cmd := buildRootCmd()

			output := captureOutput(func() {
				cmd.ParseFlags(tt.args) //nolint:errcheck
				remaining := cmd.Flags().Args()
				if len(remaining) == 0 {
					options.files = append(options.files, "-")
				} else {
					options.files = append(options.files, remaining...)
				}
			})

			if output != tt.expectedOutput {
				t.Errorf("expected output %q, got %q", tt.expectedOutput, output)
			}

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
	parseArgs([]string{"--field-separator", ",", "test.txt"})
	if options.fieldSeparator != "," {
		t.Errorf("Expected field separator ',', got %q", options.fieldSeparator)
	}

	// Test multiline
	parseArgs([]string{"--multiline", "3", "test.txt"})
	if options.multiline != 3 {
		t.Errorf("Expected multiline 3, got %d", options.multiline)
	}

	// Test ignore-case
	parseArgs([]string{"--ignore-case", "test.txt"})
	if !options.ignoreCase {
		t.Errorf("Expected ignoreCase to be true")
	}

	// Test ignore-leading-blanks
	parseArgs([]string{"--ignore-leading-blanks", "test.txt"})
	if !options.ignoreLeadingBlanks {
		t.Errorf("Expected ignoreLeadingBlanks to be true")
	}

	// Test ignore-leading-whitespace (alias)
	parseArgs([]string{"--ignore-leading-whitespace", "test.txt"})
	if !options.ignoreLeadingBlanks {
		t.Errorf("Expected ignoreLeadingBlanks to be true via alias")
	}

	// Test force
	parseArgs([]string{"--force", "test.txt"})
	if !options.forceOutput {
		t.Errorf("Expected forceOutput to be true")
	}

	// Test output
	parseArgs([]string{"--output", "output.txt", "test.txt"})
	if options.output != "output.txt" {
		t.Errorf("Expected output 'output.txt', got %q", options.output)
	}

	// Test key with single field
	parseArgs([]string{"--key", "2", "test.txt"})
	if options.key.KeyKind != SingleField || len(options.key.Keys) != 1 || options.key.Keys[0] != 2 {
		t.Errorf("Expected key SingleField [2], got %s", options.key)
	}

	// Test uniq with first
	parseArgs([]string{"--uniq", "first", "test.txt"})
	if options.uniq != first {
		t.Errorf("Expected uniq mode 'first', got %s", options.uniq)
	}

	// Test cpuprofile
	parseArgs([]string{"--cpuprofile", "cpu.prof", "test.txt"})
	if options.cpuprofile != "cpu.prof" {
		t.Errorf("Expected cpuprofile 'cpu.prof', got %q", options.cpuprofile)
	}

	// Test memprofile
	parseArgs([]string{"--memprofile", "mem.prof", "test.txt"})
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
