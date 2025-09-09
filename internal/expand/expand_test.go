package expand

import (
	"testing"

	"github.com/hagatasdelus/babbr/internal/config"
)

func TestExpandBasicAbbreviation(t *testing.T) {
	cfg := &config.Config{
		Abbreviations: []config.Abbreviation{
			{
				Name:    "list files",
				Abbr:    "l",
				Snippet: "ls -l",
			},
		},
	}

	expander := NewExpander(cfg)

	tests := []struct {
		name         string
		leftBuffer   string
		rightBuffer  string
		wantLeft     string
		wantRight    string
		wantExpanded bool
	}{
		{
			name:         "basic abbreviation expansion",
			leftBuffer:   "l",
			rightBuffer:  "",
			wantLeft:     "ls -l",
			wantRight:    "",
			wantExpanded: true,
		},
		{
			name:         "no expansion for partial match",
			leftBuffer:   "ls",
			rightBuffer:  "",
			wantLeft:     "ls",
			wantRight:    "",
			wantExpanded: false,
		},
		{
			name:         "expansion in middle of command",
			leftBuffer:   "echo hello && l",
			rightBuffer:  " && echo world",
			wantLeft:     "echo hello && ls -l",
			wantRight:    " && echo world",
			wantExpanded: true,
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

			if result.NewLeftBuffer != tt.wantLeft {
				t.Errorf("Expand() NewLeftBuffer = %q, want %q", result.NewLeftBuffer, tt.wantLeft)
			}

			if result.NewRightBuffer != tt.wantRight {
				t.Errorf("Expand() NewRightBuffer = %q, want %q", result.NewRightBuffer, tt.wantRight)
			}
		})
	}
}

func TestExpandWithPosition(t *testing.T) {
	cfg := &config.Config{
		Abbreviations: []config.Abbreviation{
			{
				Name:    "pipe to less",
				Abbr:    "L",
				Snippet: "| less",
				Options: &config.AbbreviationOptions{
					Position: "anywhere",
				},
			},
			{
				Name:    "git status",
				Abbr:    "gst",
				Snippet: "git status",
			},
		},
	}

	expander := NewExpander(cfg)

	tests := []struct {
		name         string
		leftBuffer   string
		wantLeft     string
		wantExpanded bool
	}{
		{
			name:         "anywhere position works",
			leftBuffer:   "cat file.txt L",
			wantLeft:     "cat file.txt | less",
			wantExpanded: true,
		},
		{
			name:         "command position at start",
			leftBuffer:   "gst",
			wantLeft:     "git status",
			wantExpanded: true,
		},
		{
			name:         "command position not in middle",
			leftBuffer:   "echo gst",
			wantLeft:     "echo gst",
			wantExpanded: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := ExpandRequest{
				LeftBuffer:  tt.leftBuffer,
				RightBuffer: "",
			}

			result, err := expander.Expand(req)
			if err != nil {
				t.Fatalf("Expand() error = %v", err)
			}

			if result.HasExpansion != tt.wantExpanded {
				t.Errorf("Expand() HasExpansion = %v, want %v", result.HasExpansion, tt.wantExpanded)
			}

			if result.NewLeftBuffer != tt.wantLeft {
				t.Errorf("Expand() NewLeftBuffer = %q, want %q", result.NewLeftBuffer, tt.wantLeft)
			}
		})
	}
}

func TestExpandWithCommand(t *testing.T) {
	cfg := &config.Config{
		Abbreviations: []config.Abbreviation{
			{
				Name:    "git status",
				Abbr:    "s",
				Snippet: "status",
				Options: &config.AbbreviationOptions{
					Position: "anywhere",
					Command:  "git",
				},
			},
		},
	}

	expander := NewExpander(cfg)

	tests := []struct {
		name         string
		leftBuffer   string
		wantLeft     string
		wantExpanded bool
	}{
		{
			name:         "command specific expansion works",
			leftBuffer:   "git s",
			wantLeft:     "git status",
			wantExpanded: true,
		},
		{
			name:         "command specific expansion fails for wrong command",
			leftBuffer:   "hg s",
			wantLeft:     "hg s",
			wantExpanded: false,
		},
		{
			name:         "command specific expansion works with complex command",
			leftBuffer:   "git add . && git s",
			wantLeft:     "git add . && git status",
			wantExpanded: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := ExpandRequest{
				LeftBuffer:  tt.leftBuffer,
				RightBuffer: "",
			}

			result, err := expander.Expand(req)
			if err != nil {
				t.Fatalf("Expand() error = %v", err)
			}

			if result.HasExpansion != tt.wantExpanded {
				t.Errorf("Expand() HasExpansion = %v, want %v", result.HasExpansion, tt.wantExpanded)
			}

			if result.NewLeftBuffer != tt.wantLeft {
				t.Errorf("Expand() NewLeftBuffer = %q, want %q", result.NewLeftBuffer, tt.wantLeft)
			}
		})
	}
}

func TestExpandWithSetCursor(t *testing.T) {
	cfg := &config.Config{
		Abbreviations: []config.Abbreviation{
			{
				Name:    "git commit with message",
				Abbr:    "cm",
				Snippet: "commit -m '%'",
				Options: &config.AbbreviationOptions{
					Position:  "anywhere",
					Command:   "git",
					SetCursor: true,
				},
			},
		},
	}

	expander := NewExpander(cfg)

	req := ExpandRequest{
		LeftBuffer:  "git cm",
		RightBuffer: "",
	}

	result, err := expander.Expand(req)
	if err != nil {
		t.Fatalf("Expand() error = %v", err)
	}

	if !result.HasExpansion {
		t.Error("Expected expansion to occur")
	}

	expectedLeft := "git commit -m ''"
	if result.NewLeftBuffer != expectedLeft {
		t.Errorf("NewLeftBuffer = %q, want %q", result.NewLeftBuffer, expectedLeft)
	}

	expectedCursorPos := len("git commit -m '")
	if result.CursorOffset != expectedCursorPos {
		t.Errorf("CursorOffset = %d, want %d", result.CursorOffset, expectedCursorPos)
	}
}

func TestExpandWithRegex(t *testing.T) {
	cfg := &config.Config{
		Abbreviations: []config.Abbreviation{
			{
				Name:    "Execute Python scripts",
				Snippet: "python3 $file",
				Options: &config.AbbreviationOptions{
					Regex:    `^(?P<file>\S+\.py)$`,
					Evaluate: true,
				},
			},
		},
	}

	expander := NewExpander(cfg)

	tests := []struct {
		name         string
		leftBuffer   string
		wantLeft     string
		wantExpanded bool
	}{
		{
			name:         "python file regex match",
			leftBuffer:   "script.py",
			wantLeft:     "python3 script.py",
			wantExpanded: true,
		},
		{
			name:         "python file with path regex match",
			leftBuffer:   "./myscript.py",
			wantLeft:     "python3 ./myscript.py",
			wantExpanded: true,
		},
		{
			name:         "non-python file no match",
			leftBuffer:   "script.txt",
			wantLeft:     "script.txt",
			wantExpanded: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := ExpandRequest{
				LeftBuffer:  tt.leftBuffer,
				RightBuffer: "",
			}

			result, err := expander.Expand(req)
			if err != nil {
				t.Fatalf("Expand() error = %v", err)
			}

			if result.HasExpansion != tt.wantExpanded {
				t.Errorf("Expand() HasExpansion = %v, want %v", result.HasExpansion, tt.wantExpanded)
			}

			if result.NewLeftBuffer != tt.wantLeft {
				t.Errorf("Expand() NewLeftBuffer = %q, want %q", result.NewLeftBuffer, tt.wantLeft)
			}
		})
	}
}

func TestExpandWithCondition(t *testing.T) {
	cfg := &config.Config{
		Abbreviations: []config.Abbreviation{
			{
				Name:    "Use echo command if available",
				Abbr:    "testcmd",
				Snippet: "echo available",
				Options: &config.AbbreviationOptions{
					Condition: "command -v echo >/dev/null 2>&1",
				},
			},
			{
				Name:    "Use non-existent command",
				Abbr:    "failcmd",
				Snippet: "will not expand",
				Options: &config.AbbreviationOptions{
					Condition: "command -v nonexistentcommand123 >/dev/null 2>&1",
				},
			},
		},
	}

	expander := NewExpander(cfg)

	tests := []struct {
		name         string
		leftBuffer   string
		wantExpanded bool
	}{
		{
			name:         "condition passes",
			leftBuffer:   "testcmd",
			wantExpanded: true,
		},
		{
			name:         "condition fails",
			leftBuffer:   "failcmd",
			wantExpanded: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := ExpandRequest{
				LeftBuffer:  tt.leftBuffer,
				RightBuffer: "",
			}

			result, err := expander.Expand(req)
			if err != nil {
				t.Fatalf("Expand() error = %v", err)
			}

			if result.HasExpansion != tt.wantExpanded {
				t.Errorf("Expand() HasExpansion = %v, want %v", result.HasExpansion, tt.wantExpanded)
			}
		})
	}
}

func TestExpandWithEvaluate(t *testing.T) {
	cfg := &config.Config{
		Abbreviations: []config.Abbreviation{
			{
				Name:    "Current date command",
				Abbr:    "now",
				Snippet: "echo $(date +%Y-%m-%d)",
				Options: &config.AbbreviationOptions{
					Evaluate: true,
				},
			},
		},
	}

	expander := NewExpander(cfg)

	req := ExpandRequest{
		LeftBuffer:  "now",
		RightBuffer: "",
	}

	result, err := expander.Expand(req)
	if err != nil {
		t.Fatalf("Expand() error = %v", err)
	}

	if !result.HasExpansion {
		t.Error("Expected expansion to occur")
	}

	if len(result.NewLeftBuffer) < 10 {
		t.Errorf("Expected expanded date string, got: %q", result.NewLeftBuffer)
	}
}
