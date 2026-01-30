package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestLoadFile(t *testing.T) {
	// Create a temporary test file
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.txt")

	content := "line 1\nline 2\nline 3\n"
	err := os.WriteFile(testFile, []byte(content), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Test loading a regular file
	lines := LoadFile(testFile)

	expected := []string{"line 1", "line 2", "line 3"}
	if len(lines) != len(expected) {
		t.Errorf("Expected %d lines, got %d", len(expected), len(lines))
	}

	for i, line := range lines {
		if line != expected[i] {
			t.Errorf("Line %d: expected %q, got %q", i, expected[i], line)
		}
	}
}

func TestLoadFileWithEmptyLines(t *testing.T) {
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test_empty.txt")

	content := "line 1\n\nline 3\n"
	err := os.WriteFile(testFile, []byte(content), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	lines := LoadFile(testFile)

	expected := []string{"line 1", "", "line 3"}
	if len(lines) != len(expected) {
		t.Errorf("Expected %d lines, got %d", len(expected), len(lines))
	}

	for i, line := range lines {
		if line != expected[i] {
			t.Errorf("Line %d: expected %q, got %q", i, expected[i], line)
		}
	}
}

func TestLoadFileStdin(t *testing.T) {
	// Save original stdin
	oldStdin := os.Stdin
	defer func() { os.Stdin = oldStdin }()

	// Create a pipe to simulate stdin
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("Failed to create pipe: %v", err)
	}
	os.Stdin = r

	// Write test data to the pipe
	testData := "stdin line 1\nstdin line 2\n"
	go func() {
		w.WriteString(testData)
		w.Close()
	}()

	// Test loading from stdin using "-"
	lines := LoadFile("-")

	expected := []string{"stdin line 1", "stdin line 2"}
	if len(lines) != len(expected) {
		t.Errorf("Expected %d lines, got %d", len(expected), len(lines))
	}

	for i, line := range lines {
		if line != expected[i] {
			t.Errorf("Line %d: expected %q, got %q", i, expected[i], line)
		}
	}
}

func TestLoadFileLongLines(t *testing.T) {
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test_long.txt")

	// Create a long line that fits within buffer
	longLine := strings.Repeat("a", 200000) // 200KB line, within the 256KB buffer
	content := "short line\n" + longLine + "\nanother short line\n"
	err := os.WriteFile(testFile, []byte(content), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	lines := LoadFile(testFile)

	if len(lines) != 3 {
		t.Errorf("Expected 3 lines, got %d", len(lines))
	}

	if lines[0] != "short line" {
		t.Errorf("Line 0: expected %q, got %q", "short line", lines[0])
	}

	if lines[1] != longLine {
		t.Errorf("Line 1: long line doesn't match (length: %d vs %d)", len(longLine), len(lines[1]))
	}

	if lines[2] != "another short line" {
		t.Errorf("Line 2: expected %q, got %q", "another short line", lines[2])
	}
}

func TestLoadInputFilesSingle(t *testing.T) {
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "single.txt")

	content := "alpha\nbeta\ngamma\n"
	err := os.WriteFile(testFile, []byte(content), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	options = NewConfigurationOptions()
	options.multiline = 1
	options.fieldSeparator = " "

	key := KeyType{Type: NoKey}
	contents := LoadInputFiles([]string{testFile}, key)

	if len(contents) != 3 {
		t.Errorf("Expected 3 content lines, got %d", len(contents))
	}

	expectedLines := []string{"alpha", "beta", "gamma"}
	for i, expected := range expectedLines {
		if len(contents[i].Lines) != 1 || contents[i].Lines[0] != expected {
			t.Errorf("Content %d: expected line %q, got %v", i, expected, contents[i].Lines)
		}
	}
}

func TestLoadInputFilesMultiple(t *testing.T) {
	tmpDir := t.TempDir()
	testFile1 := filepath.Join(tmpDir, "file1.txt")
	testFile2 := filepath.Join(tmpDir, "file2.txt")

	content1 := "file1 line1\nfile1 line2\n"
	content2 := "file2 line1\nfile2 line2\n"

	err := os.WriteFile(testFile1, []byte(content1), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file 1: %v", err)
	}

	err = os.WriteFile(testFile2, []byte(content2), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file 2: %v", err)
	}

	options = NewConfigurationOptions()
	options.multiline = 1
	options.fieldSeparator = " "

	key := KeyType{Type: NoKey}
	contents := LoadInputFiles([]string{testFile1, testFile2}, key)

	if len(contents) != 4 {
		t.Errorf("Expected 4 content lines, got %d", len(contents))
	}

	expectedLines := []string{"file1 line1", "file1 line2", "file2 line1", "file2 line2"}
	for i, expected := range expectedLines {
		if i >= len(contents) {
			t.Errorf("Missing content line %d", i)
			continue
		}
		if len(contents[i].Lines) != 1 || contents[i].Lines[0] != expected {
			t.Errorf("Content %d: expected line %q, got %v", i, expected, contents[i].Lines)
		}
	}
}

func TestLoadInputFilesWithMultiline(t *testing.T) {
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "multiline.txt")

	content := "line1\nline2\nline3\nline4\n"
	err := os.WriteFile(testFile, []byte(content), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	options = NewConfigurationOptions()
	options.multiline = 2
	options.fieldSeparator = " "

	key := KeyType{Type: NoKey}
	contents := LoadInputFiles([]string{testFile}, key)

	if len(contents) != 2 {
		t.Errorf("Expected 2 multiline groups, got %d", len(contents))
	}

	if len(contents) >= 1 && len(contents[0].Lines) != 2 {
		t.Errorf("First group: expected 2 lines, got %d", len(contents[0].Lines))
	}

	if len(contents) >= 2 && len(contents[1].Lines) != 2 {
		t.Errorf("Second group: expected 2 lines, got %d", len(contents[1].Lines))
	}
}
