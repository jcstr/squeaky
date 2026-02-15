package cleaner

import "testing"

func TestParseJournalSize(t *testing.T) {
	tests := []struct {
		name    string
		output  string
		wantMin int64
		wantMax int64
	}{
		{
			name:    "megabytes",
			output:  "Archived and active journals take up 237.5M in the file system.",
			wantMin: 248000000,
			wantMax: 250000000,
		},
		{
			name:    "gigabytes",
			output:  "Archived and active journals take up 1.2G in the file system.",
			wantMin: 1280000000,
			wantMax: 1300000000,
		},
		{
			name:    "empty",
			output:  "",
			wantMin: 0,
			wantMax: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := parseJournalSize(tt.output)
			if got < tt.wantMin || got > tt.wantMax {
				t.Errorf("parseJournalSize() = %d, want between %d and %d",
					got, tt.wantMin, tt.wantMax)
			}
		})
	}
}

func TestJournalAnalyze(t *testing.T) {
	orig := RunCommand
	defer func() { RunCommand = orig }()

	RunCommand = func(name string, args ...string) (string, error) {
		return "Archived and active journals take up 512.0M in the file system.", nil
	}

	j := NewJournalCleaner("2weeks")
	result, err := j.Analyze()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.SpaceSaved == 0 {
		t.Error("SpaceSaved should not be 0")
	}
}
