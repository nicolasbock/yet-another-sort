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
