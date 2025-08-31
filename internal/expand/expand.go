package expand

import (
	"fmt"
	"os/exec"
	"regexp"
	"strings"

	"github.com/hagatasdelus/babbr/internal/config"
)

type ExpandRequest struct {
	LeftBuffer  string
	RightBuffer string
}

type ExpandResult struct {
	NewLeftBuffer  string
	NewRightBuffer string
	CursorOffset   int
	HasExpansion   bool
	SetCursor      bool
}

type Expander struct {
	config *config.Config
}

func NewExpander(cfg *config.Config) *Expander {
	return &Expander{config: cfg}
}

func (e *Expander) Expand(req ExpandRequest) (*ExpandResult, error) {
	word, wordStart := e.extractWordBeforeCursor(req.LeftBuffer)
	if word == "" {
		return &ExpandResult{
			NewLeftBuffer:  req.LeftBuffer,
			NewRightBuffer: req.RightBuffer,
			HasExpansion:   false,
		}, nil
	}

	abbr := e.findMatchingAbbreviation(word, req.LeftBuffer)
	if abbr == nil {
		return &ExpandResult{
			NewLeftBuffer:  req.LeftBuffer,
			NewRightBuffer: req.RightBuffer,
			HasExpansion:   false,
		}, nil
	}

	expansion, err := e.processSnippet(abbr, word)
	if err != nil {
		return nil, fmt.Errorf("failed to process snippet: %w", err)
	}

	newLeftBuffer := req.LeftBuffer[:wordStart] + expansion
	cursorPos := len(newLeftBuffer)
	setCursor := false

	if abbr.Options != nil && abbr.Options.SetCursor {
		if pos := strings.Index(expansion, "%"); pos != -1 {
			expansion = strings.Replace(expansion, "%", "", 1)
			newLeftBuffer = req.LeftBuffer[:wordStart] + expansion
			cursorPos = wordStart + pos
			setCursor = true
		}
	}

	return &ExpandResult{
		NewLeftBuffer:  newLeftBuffer,
		NewRightBuffer: req.RightBuffer,
		CursorOffset:   cursorPos,
		HasExpansion:   true,
		SetCursor:      setCursor,
	}, nil
}

func (e *Expander) extractWordBeforeCursor(leftBuffer string) (string, int) {
	if leftBuffer == "" {
		return "", 0
	}

	words := strings.Fields(leftBuffer)
	if len(words) == 0 {
		return "", 0
	}

	lastWord := words[len(words)-1]
	wordStart := strings.LastIndex(leftBuffer, lastWord)

	return lastWord, wordStart
}

func (e *Expander) findMatchingAbbreviation(word, commandLine string) *config.Abbreviation {
	for _, abbr := range e.config.Abbreviations {
		if e.matchesAbbreviation(&abbr, word, commandLine) {
			if abbr.Options != nil && abbr.Options.Condition != "" {
				if !e.checkCondition(abbr.Options.Condition) {
					continue
				}
			}
			return &abbr
		}
	}
	return nil
}

func (e *Expander) matchesAbbreviation(abbr *config.Abbreviation, word, commandLine string) bool {
	if abbr.Options != nil && abbr.Options.Regex != "" {
		re, err := regexp.Compile(abbr.Options.Regex)
		if err != nil {
			return false
		}
		if !re.MatchString(word) {
			return false
		}
		// For regex-based abbreviations, also check position
		return e.isValidPosition(abbr, commandLine, word)
	}

	if abbr.Abbr != word {
		return false
	}

	if abbr.Options != nil && abbr.Options.Command != "" {
		return e.matchesCommand(abbr.Options.Command, commandLine)
	}

	return e.isValidPosition(abbr, commandLine, word)
}

func (e *Expander) matchesCommand(expectedCmd, commandLine string) bool {
	beforeWord := commandLine[:strings.LastIndex(commandLine, " ")]
	beforeWord = strings.TrimSpace(beforeWord)

	if beforeWord == "" {
		return expectedCmd == ""
	}

	words := strings.Fields(beforeWord)
	if len(words) == 0 {
		return expectedCmd == ""
	}

	commandSeparators := []string{"&&", "||", ";", "|", "(", "{"}

	for i := len(words) - 1; i >= 0; i-- {
		word := words[i]
		isSeparator := false
		for _, sep := range commandSeparators {
			if word == sep || strings.HasSuffix(word, sep) {
				isSeparator = true
				break
			}
		}
		if isSeparator && i+1 < len(words) {
			return words[i+1] == expectedCmd
		}
	}

	return words[0] == expectedCmd
}

func (e *Expander) isValidPosition(abbr *config.Abbreviation, commandLine, word string) bool {
	if abbr.Options != nil && abbr.Options.Position == "anywhere" {
		return true
	}

	beforeWord := commandLine[:len(commandLine)-len(word)]
	beforeWord = strings.TrimSpace(beforeWord)

	if beforeWord == "" {
		return true
	}

	commandSeparators := []string{"&&", "||", ";", "|", "(", "{"}
	for _, sep := range commandSeparators {
		if strings.HasSuffix(beforeWord, sep) {
			return true
		}
		if strings.HasSuffix(beforeWord, sep+" ") {
			return true
		}
	}

	return false
}

func (e *Expander) checkCondition(condition string) bool {
	cmd := exec.Command("bash", "-c", condition)
	return cmd.Run() == nil
}

func (e *Expander) processSnippet(abbr *config.Abbreviation, word string) (string, error) {
	snippet := abbr.Snippet

	// Handle regex-based variable substitution
	if abbr.Options != nil && abbr.Options.Regex != "" && abbr.Options.Evaluate {
		re, err := regexp.Compile(abbr.Options.Regex)
		if err != nil {
			return "", fmt.Errorf("invalid regex: %w", err)
		}

		matches := re.FindStringSubmatch(word)
		if len(matches) > 1 {
			for i, name := range re.SubexpNames() {
				if i > 0 && i < len(matches) && name != "" {
					placeholder := "$" + name
					snippet = strings.ReplaceAll(snippet, placeholder, matches[i])
				}
			}
		}
	}

	// Handle command substitution ($(command)) when evaluate is true and no regex
	if abbr.Options != nil && abbr.Options.Evaluate && abbr.Options.Regex == "" {
		// Check if snippet contains command substitution pattern
		if strings.Contains(snippet, "$(") && strings.Contains(snippet, ")") {
			cmd := exec.Command("bash", "-c", snippet)
			output, err := cmd.Output()
			if err != nil {
				return snippet, nil
			}
			return strings.TrimSpace(string(output)), nil
		}
	}

	return snippet, nil
}
