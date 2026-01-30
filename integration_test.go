package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestIntegrationBasicSort(t *testing.T) {
	tmpDir := t.TempDir()
	inputFile := filepath.Join(tmpDir, "input.txt")
	outputFile := filepath.Join(tmpDir, "output.txt")

	// Create test input
	input := "line 3\nline 1\nline 2\n"
	err := os.WriteFile(inputFile, []byte(input), 0644)
	if err != nil {
		t.Fatalf("Failed to create input file: %v", err)
	}

	// Setup options
	options = NewConfigurationOptions()
	options.files = []string{inputFile}
	options.output = outputFile
	options.multiline = 1
	options.fieldSeparator = " "

	// Run sort
	contents := LoadInputFiles(options.files, options.key)
	sorted := SortContents(contents)
	uniqued := UniqContents(sorted)

	// Write output
	fd, err := os.Create(outputFile)
	if err != nil {
		t.Fatalf("Failed to create output file: %v", err)
	}
	for _, multiline := range uniqued {
		for _, line := range multiline.Lines {
			fd.WriteString(line + "\n")
		}
	}
	fd.Close()

	// Verify output
	output, err := os.ReadFile(outputFile)
	if err != nil {
		t.Fatalf("Failed to read output file: %v", err)
	}

	expected := "line 1\nline 2\nline 3\n"
	if string(output) != expected {
		t.Errorf("Expected:\n%s\nGot:\n%s", expected, string(output))
	}
}

func TestIntegrationMultilineSort(t *testing.T) {
	tmpDir := t.TempDir()
	inputFile := filepath.Join(tmpDir, "input.txt")

	// Create test input with timestamps like bash history
	input := "#1692110033\nls\n#1692110031\ncd /tmp\n#1692110035\npwd\n"
	err := os.WriteFile(inputFile, []byte(input), 0644)
	if err != nil {
		t.Fatalf("Failed to create input file: %v", err)
	}

	// Setup options
	options = NewConfigurationOptions()
	options.files = []string{inputFile}
	options.multiline = 2
	options.fieldSeparator = " "

	// Run sort
	contents := LoadInputFiles(options.files, options.key)
	sorted := SortContents(contents)

	// Verify sorting
	if len(sorted) != 3 {
		t.Errorf("Expected 3 multilines, got %d", len(sorted))
	}

	// Check order (should be sorted by timestamp)
	expected := []string{"#1692110031", "#1692110033", "#1692110035"}
	for i, exp := range expected {
		if i >= len(sorted) || len(sorted[i].Lines) < 1 {
			t.Errorf("Missing multiline %d", i)
			continue
		}
		if sorted[i].Lines[0] != exp {
			t.Errorf("Multiline %d: expected %q, got %q", i, exp, sorted[i].Lines[0])
		}
	}
}

func TestIntegrationKeySort(t *testing.T) {
	tmpDir := t.TempDir()
	inputFile := filepath.Join(tmpDir, "input.txt")

	// Create test input
	input := "user3 100 active\nuser1 50 inactive\nuser2 75 active\n"
	err := os.WriteFile(inputFile, []byte(input), 0644)
	if err != nil {
		t.Fatalf("Failed to create input file: %v", err)
	}

	// Setup options to sort by second field (numeric values)
	options = NewConfigurationOptions()
	options.files = []string{inputFile}
	options.multiline = 1
	options.fieldSeparator = " "
	key := KeyType{Type: SingleField, Keys: []int{2}}

	// Run sort
	contents := LoadInputFiles(options.files, key)
	sorted := SortContents(contents)

	// Verify sorting by field 2
	if len(sorted) != 3 {
		t.Errorf("Expected 3 lines, got %d", len(sorted))
	}

	// Check order (should be sorted by second field: 100, 50, 75)
	expectedCompare := []string{"100", "50", "75"}
	for i, exp := range expectedCompare {
		if i >= len(sorted) {
			t.Errorf("Missing line %d", i)
			continue
		}
		if sorted[i].CompareLine != exp {
			t.Errorf("Line %d: expected compare %q, got %q", i, exp, sorted[i].CompareLine)
		}
	}
}

func TestIntegrationUniqFirst(t *testing.T) {
	tmpDir := t.TempDir()
	inputFile := filepath.Join(tmpDir, "input.txt")

	// Create test input with duplicates
	input := "apple\napple\nbanana\nbanana\nbanana\ncherry\n"
	err := os.WriteFile(inputFile, []byte(input), 0644)
	if err != nil {
		t.Fatalf("Failed to create input file: %v", err)
	}

	// Setup options
	options = NewConfigurationOptions()
	options.files = []string{inputFile}
	options.multiline = 1
	options.fieldSeparator = " "
	options.uniq = first

	// Run sort and uniq
	contents := LoadInputFiles(options.files, options.key)
	sorted := SortContents(contents)
	uniqued := UniqContents(sorted)

	// Verify output
	if len(uniqued) != 3 {
		t.Errorf("Expected 3 unique lines, got %d", len(uniqued))
	}

	expected := []string{"apple", "banana", "cherry"}
	for i, exp := range expected {
		if i >= len(uniqued) || len(uniqued[i].Lines) < 1 {
			t.Errorf("Missing line %d", i)
			continue
		}
		if uniqued[i].Lines[0] != exp {
			t.Errorf("Line %d: expected %q, got %q", i, exp, uniqued[i].Lines[0])
		}
	}
}

func TestIntegrationUniqLast(t *testing.T) {
	tmpDir := t.TempDir()
	inputFile := filepath.Join(tmpDir, "input.txt")

	// Create test input with duplicates
	input := "apple\napple\nbanana\n"
	err := os.WriteFile(inputFile, []byte(input), 0644)
	if err != nil {
		t.Fatalf("Failed to create input file: %v", err)
	}

	// Setup options
	options = NewConfigurationOptions()
	options.files = []string{inputFile}
	options.multiline = 1
	options.fieldSeparator = " "
	options.uniq = last

	// Run sort and uniq
	contents := LoadInputFiles(options.files, options.key)
	sorted := SortContents(contents)
	uniqued := UniqContents(sorted)

	// Verify output
	if len(uniqued) != 2 {
		t.Errorf("Expected 2 unique lines, got %d", len(uniqued))
	}
}

func TestIntegrationIgnoreCase(t *testing.T) {
	tmpDir := t.TempDir()
	inputFile := filepath.Join(tmpDir, "input.txt")

	// Create test input with mixed case
	input := "Zebra\napple\nBanana\n"
	err := os.WriteFile(inputFile, []byte(input), 0644)
	if err != nil {
		t.Fatalf("Failed to create input file: %v", err)
	}

	// Setup options with case-insensitive sorting
	options = NewConfigurationOptions()
	options.files = []string{inputFile}
	options.multiline = 1
	options.fieldSeparator = " "
	options.ignoreCase = true

	// Run sort
	contents := LoadInputFiles(options.files, options.key)
	sorted := SortContents(contents)

	// Verify sorting (case-insensitive)
	if len(sorted) != 3 {
		t.Errorf("Expected 3 lines, got %d", len(sorted))
	}

	// Check order (should be: apple, Banana, Zebra)
	expectedCompare := []string{"apple", "banana", "zebra"}
	for i, exp := range expectedCompare {
		if i >= len(sorted) {
			t.Errorf("Missing line %d", i)
			continue
		}
		if sorted[i].CompareLine != exp {
			t.Errorf("Line %d: expected compare %q, got %q", i, exp, sorted[i].CompareLine)
		}
	}
}

func TestIntegrationCustomSeparator(t *testing.T) {
	tmpDir := t.TempDir()
	inputFile := filepath.Join(tmpDir, "input.csv")

	// Create CSV input
	input := "name,age,city\nAlice,30,NYC\nBob,25,LA\nCharlie,35,SF\n"
	err := os.WriteFile(inputFile, []byte(input), 0644)
	if err != nil {
		t.Fatalf("Failed to create input file: %v", err)
	}

	// Setup options with comma separator, sort by field 2 (age)
	options = NewConfigurationOptions()
	options.files = []string{inputFile}
	options.multiline = 1
	options.fieldSeparator = ","
	key := KeyType{Type: SingleField, Keys: []int{2}}

	// Run sort (skip header)
	contents := LoadInputFiles(options.files, key)
	sorted := SortContents(contents)

	// Verify sorting
	if len(sorted) != 4 {
		t.Errorf("Expected 4 lines, got %d", len(sorted))
	}

	// Sorted lexicographically by age field: "25", "30", "35", "age"
	expected := []string{"25", "30", "35", "age"}
	for i, exp := range expected {
		if i >= len(sorted) {
			t.Errorf("Missing line %d", i)
			continue
		}
		if sorted[i].CompareLine != exp {
			t.Errorf("Line %d: expected compare %q, got %q", i, exp, sorted[i].CompareLine)
		}
	}
}

func TestIntegrationRemainderKey(t *testing.T) {
	tmpDir := t.TempDir()
	inputFile := filepath.Join(tmpDir, "input.txt")

	// Create test input
	input := "id name age city\n3 Charlie 35 SF\n1 Alice 30 NYC\n2 Bob 25 LA\n"
	err := os.WriteFile(inputFile, []byte(input), 0644)
	if err != nil {
		t.Fatalf("Failed to create input file: %v", err)
	}

	// Setup options with remainder key (sort from field 2 onwards)
	options = NewConfigurationOptions()
	options.files = []string{inputFile}
	options.multiline = 1
	options.fieldSeparator = " "
	key := KeyType{Type: Remainder, Keys: []int{2}}

	// Run sort
	contents := LoadInputFiles(options.files, key)
	sorted := SortContents(contents)

	// Verify sorting
	if len(sorted) != 4 {
		t.Errorf("Expected 4 lines, got %d", len(sorted))
	}

	// Should be sorted lexicographically by "name age city"
	// Order: "Alice 30 NYC", "Bob 25 LA", "Charlie 35 SF", "name age city"
	expectedCompares := []string{"Alice 30 NYC", "Bob 25 LA", "Charlie 35 SF", "name age city"}
	for i, exp := range expectedCompares {
		if i >= len(sorted) {
			t.Errorf("Missing line %d", i)
			continue
		}
		if sorted[i].CompareLine != exp {
			t.Errorf("Line %d: expected compare %q, got %q", i, exp, sorted[i].CompareLine)
		}
	}
}

func TestIntegrationSubsetKey(t *testing.T) {
	tmpDir := t.TempDir()
	inputFile := filepath.Join(tmpDir, "input.txt")

	// Create test input
	input := "id name age city\n3 Charlie 35 SF\n1 Alice 30 NYC\n2 Bob 25 LA\n"
	err := os.WriteFile(inputFile, []byte(input), 0644)
	if err != nil {
		t.Fatalf("Failed to create input file: %v", err)
	}

	// Setup options with subset key (sort by fields 2-3: name and age)
	options = NewConfigurationOptions()
	options.files = []string{inputFile}
	options.multiline = 1
	options.fieldSeparator = " "
	key := KeyType{Type: SubSet, Keys: []int{2, 3}}

	// Run sort
	contents := LoadInputFiles(options.files, key)
	sorted := SortContents(contents)

	// Verify sorting
	if len(sorted) != 4 {
		t.Errorf("Expected 4 lines, got %d", len(sorted))
	}

	// Should be sorted lexicographically by "name age"
	// Order: "Alice 30", "Bob 25", "Charlie 35", "name age"
	expectedCompares := []string{"Alice 30", "Bob 25", "Charlie 35", "name age"}
	for i, exp := range expectedCompares {
		if i >= len(sorted) {
			t.Errorf("Missing line %d", i)
			continue
		}
		if sorted[i].CompareLine != exp {
			t.Errorf("Line %d: expected compare %q, got %q", i, exp, sorted[i].CompareLine)
		}
	}
}

func TestIntegrationMultipleFiles(t *testing.T) {
	tmpDir := t.TempDir()
	file1 := filepath.Join(tmpDir, "file1.txt")
	file2 := filepath.Join(tmpDir, "file2.txt")

	// Create test inputs
	input1 := "zebra\napple\n"
	input2 := "mango\nbanana\n"

	err := os.WriteFile(file1, []byte(input1), 0644)
	if err != nil {
		t.Fatalf("Failed to create file1: %v", err)
	}

	err = os.WriteFile(file2, []byte(input2), 0644)
	if err != nil {
		t.Fatalf("Failed to create file2: %v", err)
	}

	// Setup options
	options = NewConfigurationOptions()
	options.files = []string{file1, file2}
	options.multiline = 1
	options.fieldSeparator = " "

	// Run sort
	contents := LoadInputFiles(options.files, options.key)
	sorted := SortContents(contents)

	// Verify sorting
	if len(sorted) != 4 {
		t.Errorf("Expected 4 lines from 2 files, got %d", len(sorted))
	}

	// Check order (alphabetically sorted)
	expected := []string{"apple", "banana", "mango", "zebra"}
	for i, exp := range expected {
		if i >= len(sorted) || len(sorted[i].Lines) < 1 {
			t.Errorf("Missing line %d", i)
			continue
		}
		if sorted[i].Lines[0] != exp {
			t.Errorf("Line %d: expected %q, got %q", i, exp, sorted[i].Lines[0])
		}
	}
}

func TestIntegrationEmptyFile(t *testing.T) {
	tmpDir := t.TempDir()
	inputFile := filepath.Join(tmpDir, "empty.txt")

	// Create empty file
	err := os.WriteFile(inputFile, []byte(""), 0644)
	if err != nil {
		t.Fatalf("Failed to create empty file: %v", err)
	}

	// Setup options
	options = NewConfigurationOptions()
	options.files = []string{inputFile}
	options.multiline = 1
	options.fieldSeparator = " "

	// Run sort
	contents := LoadInputFiles(options.files, options.key)
	sorted := SortContents(contents)
	uniqued := UniqContents(sorted)

	// Verify empty output
	if len(uniqued) != 0 {
		t.Errorf("Expected 0 lines for empty file, got %d", len(uniqued))
	}
}

func TestIntegrationIgnoreLeadingBlanks(t *testing.T) {
	tmpDir := t.TempDir()
	inputFile := filepath.Join(tmpDir, "input.txt")

	// Create test input with leading spaces
	input := "  zebra\n apple\n  banana\n"
	err := os.WriteFile(inputFile, []byte(input), 0644)
	if err != nil {
		t.Fatalf("Failed to create input file: %v", err)
	}

	// Setup options with ignore leading blanks
	options = NewConfigurationOptions()
	options.files = []string{inputFile}
	options.multiline = 1
	options.fieldSeparator = " "
	options.ignoreLeadingBlanks = true

	// Run sort
	contents := LoadInputFiles(options.files, options.key)
	sorted := SortContents(contents)

	// Verify sorting (leading blanks ignored)
	if len(sorted) != 3 {
		t.Errorf("Expected 3 lines, got %d", len(sorted))
	}

	// Check order (should be: apple, banana, zebra - ignoring leading spaces)
	expectedCompare := []string{"apple", "banana", "zebra"}
	for i, exp := range expectedCompare {
		if i >= len(sorted) {
			t.Errorf("Missing line %d", i)
			continue
		}
		if sorted[i].CompareLine != exp {
			t.Errorf("Line %d: expected compare %q, got %q", i, exp, sorted[i].CompareLine)
		}
	}
}
