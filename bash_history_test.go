package main

import (
	"os"
	"path/filepath"
	"testing"
)

// ---------------------------------------------------------------------------
// isBashTimestamp
// ---------------------------------------------------------------------------

func TestIsBashTimestamp_Valid(t *testing.T) {
	cases := []string{
		"#0",
		"#1",
		"#1773948684",
		"#9999999999",
		"#0000000000",
	}
	for _, s := range cases {
		if !isBashTimestamp(s) {
			t.Errorf("isBashTimestamp(%q) = false, want true", s)
		}
	}
}

func TestIsBashTimestamp_Invalid(t *testing.T) {
	cases := []string{
		"",           // empty string
		"#",          // hash with no digits
		"1773948684", // no leading hash
		"#abc",       // non-digit after hash
		"#123abc",    // digits then non-digit
		"# 123",      // space after hash
		"##123",      // double hash
		"#-1",        // negative sign
		"#1.0",       // decimal point
	}
	for _, s := range cases {
		if isBashTimestamp(s) {
			t.Errorf("isBashTimestamp(%q) = true, want false", s)
		}
	}
}

// ---------------------------------------------------------------------------
// ParseBashHistory
// ---------------------------------------------------------------------------

func TestParseBashHistory_SingleLineCommands(t *testing.T) {
	lines := []string{
		"#1773948684",
		"juju controllers",
		"#1773949015",
		"ll",
		"#1773949028",
		"history",
	}
	got := ParseBashHistory(lines)
	if len(got) != 3 {
		t.Fatalf("expected 3 records, got %d", len(got))
	}

	cases := []struct {
		ts  string
		cmd string
	}{
		{"#1773948684", "juju controllers"},
		{"#1773949015", "ll"},
		{"#1773949028", "history"},
	}
	for i, c := range cases {
		r := got[i]
		if r.Lines[0] != c.ts {
			t.Errorf("record %d: Lines[0] = %q, want %q", i, r.Lines[0], c.ts)
		}
		if len(r.Lines) != 2 {
			t.Errorf("record %d: len(Lines) = %d, want 2", i, len(r.Lines))
			continue
		}
		if r.Lines[1] != c.cmd {
			t.Errorf("record %d: Lines[1] = %q, want %q", i, r.Lines[1], c.cmd)
		}
		if r.CompareLine != c.cmd {
			t.Errorf("record %d: CompareLine = %q, want %q", i, r.CompareLine, c.cmd)
		}
	}
}

func TestParseBashHistory_MultiLineCommand(t *testing.T) {
	lines := []string{
		"#1773949100",
		"echo hello",
		"world",
		"#1773949200",
		"ls",
	}
	got := ParseBashHistory(lines)
	if len(got) != 2 {
		t.Fatalf("expected 2 records, got %d", len(got))
	}

	// First record: multi-line command
	r0 := got[0]
	if r0.Lines[0] != "#1773949100" {
		t.Errorf("record 0 ts: got %q, want %q", r0.Lines[0], "#1773949100")
	}
	if len(r0.Lines) != 3 {
		t.Errorf("record 0: len(Lines) = %d, want 3", len(r0.Lines))
	} else {
		if r0.Lines[1] != "echo hello" {
			t.Errorf("record 0 Lines[1]: got %q, want %q", r0.Lines[1], "echo hello")
		}
		if r0.Lines[2] != "world" {
			t.Errorf("record 0 Lines[2]: got %q, want %q", r0.Lines[2], "world")
		}
	}
	wantCompare := "echo hello\nworld"
	if r0.CompareLine != wantCompare {
		t.Errorf("record 0 CompareLine: got %q, want %q", r0.CompareLine, wantCompare)
	}

	// Second record: single-line command
	r1 := got[1]
	if r1.Lines[0] != "#1773949200" {
		t.Errorf("record 1 ts: got %q, want %q", r1.Lines[0], "#1773949200")
	}
	if r1.CompareLine != "ls" {
		t.Errorf("record 1 CompareLine: got %q, want %q", r1.CompareLine, "ls")
	}
}

func TestParseBashHistory_DropRecordsWithNoCommandLines(t *testing.T) {
	// Two consecutive timestamp lines — the first record has no command lines
	// and must be silently dropped.
	lines := []string{
		"#1000000001",
		// no command lines
		"#1000000002",
		"ls",
		// timestamp at EOF with no following commands
		"#1000000003",
	}
	got := ParseBashHistory(lines)
	if len(got) != 1 {
		t.Fatalf("expected 1 record (the ls one), got %d: %v", len(got), got)
	}
	if got[0].Lines[0] != "#1000000002" {
		t.Errorf("expected timestamp #1000000002, got %q", got[0].Lines[0])
	}
	if got[0].CompareLine != "ls" {
		t.Errorf("expected CompareLine %q, got %q", "ls", got[0].CompareLine)
	}
}

func TestParseBashHistory_SkipLeadingNonTimestampLines(t *testing.T) {
	// Lines before the first timestamp must be silently skipped.
	lines := []string{
		"some garbage line",
		"another garbage line",
		"#1773948684",
		"git status",
	}
	got := ParseBashHistory(lines)
	if len(got) != 1 {
		t.Fatalf("expected 1 record, got %d", len(got))
	}
	if got[0].Lines[0] != "#1773948684" {
		t.Errorf("expected timestamp #1773948684, got %q", got[0].Lines[0])
	}
	if got[0].CompareLine != "git status" {
		t.Errorf("expected CompareLine %q, got %q", "git status", got[0].CompareLine)
	}
}

func TestParseBashHistory_Empty(t *testing.T) {
	got := ParseBashHistory([]string{})
	if len(got) != 0 {
		t.Errorf("expected 0 records for empty input, got %d", len(got))
	}
}

func TestParseBashHistory_OnlyTimestampLines(t *testing.T) {
	// Every record has no command lines — all should be dropped.
	lines := []string{
		"#1000000001",
		"#1000000002",
		"#1000000003",
	}
	got := ParseBashHistory(lines)
	if len(got) != 0 {
		t.Errorf("expected 0 records, got %d", len(got))
	}
}

// ---------------------------------------------------------------------------
// DeduplicateBashHistory
// ---------------------------------------------------------------------------

func TestDeduplicateBashHistory_KeepsLatestTimestamp(t *testing.T) {
	// "juju controllers" appears twice; the newer timestamp (1773949035) must win.
	input := ContentType{
		{Lines: []string{"#1773948684", "juju controllers"}, CompareLine: "juju controllers"},
		{Lines: []string{"#1773949015", "ll"}, CompareLine: "ll"},
		{Lines: []string{"#1773949035", "juju controllers"}, CompareLine: "juju controllers"},
	}
	got := DeduplicateBashHistory(input)
	if len(got) != 2 {
		t.Fatalf("expected 2 records after dedup, got %d", len(got))
	}
	// Find the "juju controllers" entry and verify it carries the latest ts.
	for _, r := range got {
		if r.CompareLine == "juju controllers" {
			if r.Lines[0] != "#1773949035" {
				t.Errorf("juju controllers: expected ts #1773949035, got %q", r.Lines[0])
			}
		}
	}
}

func TestDeduplicateBashHistory_MultiLineCommandDuplicate(t *testing.T) {
	// Multi-line command block duplicated; keep latest timestamp.
	input := ContentType{
		{
			Lines:       []string{"#1000000010", "echo hello", "world"},
			CompareLine: "echo hello\nworld",
		},
		{
			Lines:       []string{"#1000000020", "echo hello", "world"},
			CompareLine: "echo hello\nworld",
		},
	}
	got := DeduplicateBashHistory(input)
	if len(got) != 1 {
		t.Fatalf("expected 1 record after dedup, got %d", len(got))
	}
	if got[0].Lines[0] != "#1000000020" {
		t.Errorf("expected ts #1000000020, got %q", got[0].Lines[0])
	}
	if got[0].CompareLine != "echo hello\nworld" {
		t.Errorf("expected CompareLine %q, got %q", "echo hello\nworld", got[0].CompareLine)
	}
}

func TestDeduplicateBashHistory_NewerTimestampFirst(t *testing.T) {
	// The newer timestamp appears BEFORE the older one in the file.
	// Dedup must still keep the newer one.
	input := ContentType{
		{Lines: []string{"#1773949035", "juju controllers"}, CompareLine: "juju controllers"},
		{Lines: []string{"#1773948684", "juju controllers"}, CompareLine: "juju controllers"},
	}
	got := DeduplicateBashHistory(input)
	if len(got) != 1 {
		t.Fatalf("expected 1 record, got %d", len(got))
	}
	if got[0].Lines[0] != "#1773949035" {
		t.Errorf("expected ts #1773949035, got %q", got[0].Lines[0])
	}
}

func TestDeduplicateBashHistory_AllUnique(t *testing.T) {
	// No duplicates — output must equal input (same length).
	input := ContentType{
		{Lines: []string{"#1000000001", "ls"}, CompareLine: "ls"},
		{Lines: []string{"#1000000002", "cd"}, CompareLine: "cd"},
		{Lines: []string{"#1000000003", "pwd"}, CompareLine: "pwd"},
	}
	got := DeduplicateBashHistory(input)
	if len(got) != 3 {
		t.Fatalf("expected 3 records, got %d", len(got))
	}
}

func TestDeduplicateBashHistory_Empty(t *testing.T) {
	got := DeduplicateBashHistory(ContentType{})
	if len(got) != 0 {
		t.Errorf("expected 0 records, got %d", len(got))
	}
}

// ---------------------------------------------------------------------------
// SortBashHistory
// ---------------------------------------------------------------------------

func TestSortBashHistory_Ascending(t *testing.T) {
	input := ContentType{
		{Lines: []string{"#1773949035", "juju controllers"}, CompareLine: "juju controllers"},
		{Lines: []string{"#1773948684", "juju controllers"}, CompareLine: "juju controllers"},
		{Lines: []string{"#1773949015", "ll"}, CompareLine: "ll"},
	}
	got := SortBashHistory(input)
	if len(got) != 3 {
		t.Fatalf("expected 3 records, got %d", len(got))
	}
	expectedTS := []string{"#1773948684", "#1773949015", "#1773949035"}
	for i, ts := range expectedTS {
		if got[i].Lines[0] != ts {
			t.Errorf("position %d: expected ts %q, got %q", i, ts, got[i].Lines[0])
		}
	}
}

func TestSortBashHistory_AlreadySorted(t *testing.T) {
	input := ContentType{
		{Lines: []string{"#1000000001", "a"}, CompareLine: "a"},
		{Lines: []string{"#1000000002", "b"}, CompareLine: "b"},
		{Lines: []string{"#1000000003", "c"}, CompareLine: "c"},
	}
	got := SortBashHistory(input)
	expectedTS := []string{"#1000000001", "#1000000002", "#1000000003"}
	for i, ts := range expectedTS {
		if got[i].Lines[0] != ts {
			t.Errorf("position %d: expected ts %q, got %q", i, ts, got[i].Lines[0])
		}
	}
}

func TestSortBashHistory_SingleElement(t *testing.T) {
	input := ContentType{
		{Lines: []string{"#1000000001", "ls"}, CompareLine: "ls"},
	}
	got := SortBashHistory(input)
	if len(got) != 1 {
		t.Fatalf("expected 1 record, got %d", len(got))
	}
	if got[0].Lines[0] != "#1000000001" {
		t.Errorf("expected ts #1000000001, got %q", got[0].Lines[0])
	}
}

func TestSortBashHistory_Empty(t *testing.T) {
	got := SortBashHistory(ContentType{})
	if len(got) != 0 {
		t.Errorf("expected 0 records, got %d", len(got))
	}
}

func TestSortBashHistory_DoesNotModifyInput(t *testing.T) {
	// SortBashHistory must return a new slice and not mutate the original.
	input := ContentType{
		{Lines: []string{"#1773949035", "z"}, CompareLine: "z"},
		{Lines: []string{"#1773948684", "a"}, CompareLine: "a"},
	}
	origFirst := input[0].Lines[0]
	SortBashHistory(input)
	if input[0].Lines[0] != origFirst {
		t.Errorf("SortBashHistory mutated input: input[0].Lines[0] changed from %q to %q", origFirst, input[0].Lines[0])
	}
}

// ---------------------------------------------------------------------------
// RunBashHistoryMode — full pipeline via temp file
// ---------------------------------------------------------------------------

func TestRunBashHistoryMode_FullPipeline(t *testing.T) {
	// Input from the specification:
	//   #1773948684  juju controllers   <- duplicate, older
	//   #1773949015  ll
	//   #1773949026  cd
	//   #1773949028  history
	//   #1773949035  juju controllers   <- duplicate, newer  (must win)
	//
	// Expected after dedup + sort:
	//   #1773949015  ll
	//   #1773949026  cd
	//   #1773949028  history
	//   #1773949035  juju controllers
	content := "#1773948684\njuju controllers\n#1773949015\nll\n#1773949026\ncd\n#1773949028\nhistory\n#1773949035\njuju controllers\n"

	tmpDir := t.TempDir()
	inputFile := filepath.Join(tmpDir, "bash_history")
	if err := os.WriteFile(inputFile, []byte(content), 0644); err != nil {
		t.Fatalf("failed to write temp file: %v", err)
	}

	got := RunBashHistoryMode([]string{inputFile})

	expected := []struct {
		ts  string
		cmd string
	}{
		{"#1773949015", "ll"},
		{"#1773949026", "cd"},
		{"#1773949028", "history"},
		{"#1773949035", "juju controllers"},
	}

	if len(got) != len(expected) {
		t.Fatalf("expected %d records, got %d", len(expected), len(got))
	}
	for i, e := range expected {
		if got[i].Lines[0] != e.ts {
			t.Errorf("record %d ts: expected %q, got %q", i, e.ts, got[i].Lines[0])
		}
		if len(got[i].Lines) < 2 {
			t.Errorf("record %d: no command lines", i)
			continue
		}
		if got[i].Lines[1] != e.cmd {
			t.Errorf("record %d cmd: expected %q, got %q", i, e.cmd, got[i].Lines[1])
		}
	}
}

func TestRunBashHistoryMode_MultiLineCommandDuplicateOlderFirst(t *testing.T) {
	// Multi-line command block: older timestamp comes first in the file.
	// After dedup+sort only the newer entry (ts=20) must survive.
	content := "#1000000010\necho hello\nworld\n#1000000020\necho hello\nworld\n#1000000015\nls\n"

	tmpDir := t.TempDir()
	inputFile := filepath.Join(tmpDir, "bash_history")
	if err := os.WriteFile(inputFile, []byte(content), 0644); err != nil {
		t.Fatalf("failed to write temp file: %v", err)
	}

	got := RunBashHistoryMode([]string{inputFile})

	// Two unique command blocks: "echo hello\nworld" and "ls".
	// Sorted by ts ascending: ls(15) < echo hello\nworld(20).
	if len(got) != 2 {
		t.Fatalf("expected 2 records, got %d", len(got))
	}

	if got[0].Lines[0] != "#1000000015" {
		t.Errorf("record 0 ts: expected #1000000015, got %q", got[0].Lines[0])
	}
	if got[0].CompareLine != "ls" {
		t.Errorf("record 0 CompareLine: expected %q, got %q", "ls", got[0].CompareLine)
	}

	if got[1].Lines[0] != "#1000000020" {
		t.Errorf("record 1 ts: expected #1000000020, got %q", got[1].Lines[0])
	}
	if got[1].CompareLine != "echo hello\nworld" {
		t.Errorf("record 1 CompareLine: expected %q, got %q", "echo hello\nworld", got[1].CompareLine)
	}
}

func TestRunBashHistoryMode_NewerTimestampBeforeOlderInFile(t *testing.T) {
	// The newer timestamp of "git pull" appears BEFORE the older one in the file.
	// Dedup must keep the newer one regardless of file order.
	content := "#1773949100\ngit pull\n#1773948900\ngit pull\n#1773949050\ngit status\n"

	tmpDir := t.TempDir()
	inputFile := filepath.Join(tmpDir, "bash_history")
	if err := os.WriteFile(inputFile, []byte(content), 0644); err != nil {
		t.Fatalf("failed to write temp file: %v", err)
	}

	got := RunBashHistoryMode([]string{inputFile})

	// Two unique commands: "git pull" (keep ts=1773949100) and "git status" (ts=1773949050).
	// Sorted ascending: git status (50) < git pull (100).
	if len(got) != 2 {
		t.Fatalf("expected 2 records, got %d: %v", len(got), got)
	}

	if got[0].Lines[0] != "#1773949050" {
		t.Errorf("record 0 ts: expected #1773949050, got %q", got[0].Lines[0])
	}
	if got[0].CompareLine != "git status" {
		t.Errorf("record 0 cmd: expected %q, got %q", "git status", got[0].CompareLine)
	}

	if got[1].Lines[0] != "#1773949100" {
		t.Errorf("record 1 ts: expected #1773949100, got %q", got[1].Lines[0])
	}
	if got[1].CompareLine != "git pull" {
		t.Errorf("record 1 cmd: expected %q, got %q", "git pull", got[1].CompareLine)
	}
}

func TestRunBashHistoryMode_EmptyFile(t *testing.T) {
	tmpDir := t.TempDir()
	inputFile := filepath.Join(tmpDir, "bash_history")
	if err := os.WriteFile(inputFile, []byte(""), 0644); err != nil {
		t.Fatalf("failed to write temp file: %v", err)
	}

	got := RunBashHistoryMode([]string{inputFile})
	if len(got) != 0 {
		t.Errorf("expected 0 records for empty file, got %d", len(got))
	}
}
