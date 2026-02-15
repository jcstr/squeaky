package cleaner

import "testing"

func TestParsePaccacheOutput(t *testing.T) {
	tests := []struct {
		name      string
		output    string
		wantItems int
		wantMin   int64
		wantMax   int64
	}{
		{
			name:      "GiB output",
			output:    "==> finished dry run: 800 candidates (disk space saved: 8.94 GiB)",
			wantItems: 800,
			wantMin:   9500000000,
			wantMax:   9700000000,
		},
		{
			name:      "MiB output",
			output:    "==> finished dry run: 12 candidates (disk space saved: 456.78 MiB)",
			wantItems: 12,
			wantMin:   478000000,
			wantMax:   480000000,
		},
		{
			name:      "no candidates",
			output:    "==> no candidate packages found for pruning",
			wantItems: 0,
			wantMin:   0,
			wantMax:   0,
		},
		{
			name:      "empty output",
			output:    "",
			wantItems: 0,
			wantMin:   0,
			wantMax:   0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := parsePaccacheOutput("test", tt.output)
			if result.ItemsFound != tt.wantItems {
				t.Errorf("ItemsFound = %d, want %d", result.ItemsFound, tt.wantItems)
			}
			if result.SpaceSaved < tt.wantMin || result.SpaceSaved > tt.wantMax {
				t.Errorf("SpaceSaved = %d, want between %d and %d",
					result.SpaceSaved, tt.wantMin, tt.wantMax)
			}
		})
	}
}

func TestPacmanAnalyze(t *testing.T) {
	orig := RunCommand
	defer func() { RunCommand = orig }()

	RunCommand = func(name string, args ...string) (string, error) {
		return "==> finished dry run: 5 candidates (disk space saved: 100.00 MiB)", nil
	}

	pc := NewPacmanCleaner(2)
	result, err := pc.Analyze()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.ItemsFound != 5 {
		t.Errorf("ItemsFound = %d, want 5", result.ItemsFound)
	}
	if result.SpaceSaved == 0 {
		t.Error("SpaceSaved should not be 0")
	}
}
