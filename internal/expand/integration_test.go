package expand

import (
	"testing"

	"github.com/hagatasdelus/babbr/internal/config"
)

func TestRealConfigExpansion(t *testing.T) {
	cfg, err := config.LoadConfig()
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	expander := NewExpander(cfg)

	tests := []struct {
		name         string
		leftBuffer   string
		rightBuffer  string
		wantExpanded bool
		wantContains string
	}{
		{
			name:         "basic l expansion",
			leftBuffer:   "l",
			rightBuffer:  "",
			wantExpanded: true,
			wantContains: "ls -l",
		},
		{
			name:         "git status with s",
			leftBuffer:   "git s",
			rightBuffer:  "",
			wantExpanded: true,
			wantContains: "git status",
		},
		{
			name:         "git commit with cursor",
			leftBuffer:   "git cm",
			rightBuffer:  "",
			wantExpanded: true,
			wantContains: "commit -m",
		},
		{
			name:         "anywhere position L",
			leftBuffer:   "cat file.txt L",
			rightBuffer:  "",
			wantExpanded: true,
			wantContains: "| less",
		},
		{
			name:         "null redirection",
			leftBuffer:   "command >null",
			rightBuffer:  "",
			wantExpanded: true,
			wantContains: ">/dev/null 2>&1",
		},
		{
			name:         "python file suffix",
			leftBuffer:   "test.py",
			rightBuffer:  "",
			wantExpanded: true,
			wantContains: "python3 test.py",
		},
		{
			name:         "shell script suffix",
			leftBuffer:   "script.sh",
			rightBuffer:  "",
			wantExpanded: true,
			wantContains: "bash script.sh",
		},
		{
			name:         "tar.gz suffix",
			leftBuffer:   "archive.tar.gz",
			rightBuffer:  "",
			wantExpanded: true,
			wantContains: "tar -xzvf archive.tar.gz",
		},
		{
			name:         "no expansion for non-matching",
			leftBuffer:   "randomcommand",
			rightBuffer:  "",
			wantExpanded: false,
			wantContains: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := ExpandRequest{
				LeftBuffer:  tt.leftBuffer,
				RightBuffer: tt.rightBuffer,
			}

			result, err := expander.Expand(req)
			if err != nil {
				t.Fatalf("Expand() error = %v", err)
			}

			if result.HasExpansion != tt.wantExpanded {
				t.Errorf("Expand() HasExpansion = %v, want %v", result.HasExpansion, tt.wantExpanded)
			}

			if tt.wantExpanded && tt.wantContains != "" {
				fullResult := result.NewLeftBuffer + result.NewRightBuffer
				if !contains(fullResult, tt.wantContains) {
					t.Errorf("Expand() result %q does not contain %q", fullResult, tt.wantContains)
				}
			}
		})
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || findSubstring(s, substr))
}

func findSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

func TestConditionalExpansion(t *testing.T) {
	cfg, err := config.LoadConfig()
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	expander := NewExpander(cfg)

	req := ExpandRequest{
		LeftBuffer:  "del",
		RightBuffer: "",
	}

	result, err := expander.Expand(req)
	if err != nil {
		t.Fatalf("Expand() error = %v", err)
	}

	if result.HasExpansion {
		fullResult := result.NewLeftBuffer + result.NewRightBuffer
		if fullResult != "trash" && fullResult != "rm -i" {
			t.Errorf("Expected either 'trash' or 'rm -i', got %q", fullResult)
		}
		t.Logf("Conditional expansion resulted in: %q", fullResult)
	} else {
		t.Error("Expected conditional expansion to occur")
	}
}
