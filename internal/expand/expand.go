package expand

import (
	"fmt"
	"os/exec"
	"regexp"
	"strings"
	"sync"

	"github.com/hagatasdelus/babbr/internal/config"
)

var (
	commandSeparators = []string{"&&", "||", ";", "|", "(", "{"}
	regexCache        = make(map[string]*regexp.Regexp)
	regexCacheMutex   sync.RWMutex
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

func compileRegexCached(pattern string) (*regexp.Regexp, error) {
	regexCacheMutex.RLock()
	if cached, exists := regexCache[pattern]; exists {
		regexCacheMutex.RUnlock()
		return cached, nil
	}
	regexCacheMutex.RUnlock()

	compiled, err := regexp.Compile(pattern)
	if err != nil {
		return nil, err
	}

	regexCacheMutex.Lock()
	regexCache[pattern] = compiled
	regexCacheMutex.Unlock()

	return compiled, nil
}

func isCommandSeparatorMatch(word, separator string) bool {
	return word == separator || strings.HasSuffix(word, separator)
}

func skipWhitespace(text string, start int) int {
	for start < len(text) && text[start] == ' ' {
		start++
	}
	return start
}

func (e *Expander) extractCurrentCommand(commandLine, word string) string {
	beforeWord := commandLine[:len(commandLine)-len(word)]
	commandStart := e.findCommandStart(beforeWord)
	return strings.TrimSpace(commandLine[commandStart:])
}

func (e *Expander) findCommandStartPosition(words []string) int {
	for i := len(words) - 1; i >= 0; i-- {
		word := words[i]
		for _, sep := range commandSeparators {
			if isCommandSeparatorMatch(word, sep) {
				return i + 1
			}
		}
	}
	return 0
}

func (e *Expander) Expand(req ExpandRequest) (*ExpandResult, error) {
	word, _ := e.extractWordBeforeCursor(req.LeftBuffer)
	if word == "" {
		return &ExpandResult{
			NewLeftBuffer:  req.LeftBuffer,
			NewRightBuffer: req.RightBuffer,
			CursorOffset:   len(req.LeftBuffer),
			HasExpansion:   false,
		}, nil
	}

	abbr, matchStart, _ := e.findMatchingAbbreviationWithPosition(word, req.LeftBuffer)
	if abbr == nil {
		return &ExpandResult{
			NewLeftBuffer:  req.LeftBuffer,
			NewRightBuffer: req.RightBuffer,
			CursorOffset:   len(req.LeftBuffer),
			HasExpansion:   false,
		}, nil
	}

	expansion, err := e.processSnippet(abbr, word, req.LeftBuffer)
	if err != nil {
		return nil, fmt.Errorf("failed to process snippet: %w", err)
	}

	newLeftBuffer := req.LeftBuffer[:matchStart] + expansion
	cursorPos := len(newLeftBuffer)
	setCursor := false

	if abbr.Options != nil && abbr.Options.SetCursor {
		if pos := strings.Index(expansion, "%"); pos != -1 {
			expansion = strings.Replace(expansion, "%", "", 1)
			newLeftBuffer = req.LeftBuffer[:matchStart] + expansion
			cursorPos = matchStart + pos
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

	// Optimize by finding the last word without allocating new strings
	end := len(leftBuffer)
	for end > 0 && leftBuffer[end-1] == ' ' {
		end--
	}

	if end == 0 {
		return "", 0
	}

	start := end
	for start > 0 && leftBuffer[start-1] != ' ' {
		start--
	}

	return leftBuffer[start:end], start
}
func (e *Expander) findMatchingAbbreviationWithPosition(word, commandLine string) (*config.Abbreviation, int, int) {
	for _, abbr := range e.config.Abbreviations {
		if e.matchesAbbreviation(&abbr, word, commandLine) {
			if abbr.Options != nil && abbr.Options.Condition != "" {
				if !e.checkCondition(abbr.Options.Condition) {
					continue
				}
			}

			matchStart, matchEnd := e.calculateMatchPosition(&abbr, word, commandLine)
			return &abbr, matchStart, matchEnd
		}
	}
	return nil, 0, 0
}

func (e *Expander) calculateMatchPosition(abbr *config.Abbreviation, word, commandLine string) (int, int) {
	// For regex-only abbreviations with ^ pattern: compute match position in the current command
	if abbr.Abbr == "" && abbr.Options != nil && abbr.Options.Regex != "" && strings.HasPrefix(abbr.Options.Regex, "^") {
		re, err := compileRegexCached(abbr.Options.Regex)
		if err != nil {
			wordStart := strings.LastIndex(commandLine, word)
			return wordStart, wordStart + len(word)
		}

		beforeWord := commandLine[:len(commandLine)-len(word)]
		commandStart := e.findCommandStart(beforeWord)
		currentCommand := strings.TrimSpace(commandLine[commandStart:])

		if matches := re.FindStringIndex(currentCommand); matches != nil {
			return commandStart + matches[0], commandStart + matches[1]
		}
	}

	// For regular abbreviations, match the word position
	wordStart := strings.LastIndex(commandLine, word)
	return wordStart, wordStart + len(word)
}

func (e *Expander) matchesAbbreviation(abbr *config.Abbreviation, word, commandLine string) bool {
	if abbr.Options != nil && abbr.Options.Regex != "" {
		re, err := compileRegexCached(abbr.Options.Regex)
		if err != nil {
			return false
		}

		if abbr.Abbr == "" {
			if strings.HasPrefix(abbr.Options.Regex, "^") {
				currentCommand := e.extractCurrentCommand(commandLine, word)
				return re.MatchString(currentCommand)
			}
			return re.MatchString(word)
		}

		if !re.MatchString(word) {
			return false
		}

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
	// Find the text before the last word
	lastSpaceIndex := strings.LastIndex(commandLine, " ")
	if lastSpaceIndex == -1 {
		return expectedCmd == ""
	}

	beforeWord := strings.TrimSpace(commandLine[:lastSpaceIndex])
	if beforeWord == "" {
		return expectedCmd == ""
	}

	words := strings.Fields(beforeWord)
	if len(words) == 0 {
		return expectedCmd == ""
	}

	// Handle compound commands
	expectedWords := strings.Fields(expectedCmd)
	if len(expectedWords) == 0 {
		return false
	}

	commandStart := e.findCommandStartPosition(words)
	currentCommandWords := words[commandStart:]

	// Check if current command matches expected command
	if len(currentCommandWords) < len(expectedWords) {
		return false
	}

	for i, expectedWord := range expectedWords {
		if i >= len(currentCommandWords) || currentCommandWords[i] != expectedWord {
			return false
		}
	}

	return true
}

func (e *Expander) findCommandStart(text string) int {
	lastSeparatorPos := -1
	lastSeparatorLen := 0

	for _, sep := range commandSeparators {
		pos := strings.LastIndex(text, sep)
		if pos > lastSeparatorPos {
			lastSeparatorPos = pos
			lastSeparatorLen = len(sep)
		}
	}

	if lastSeparatorPos >= 0 {
		return skipWhitespace(text, lastSeparatorPos+lastSeparatorLen)
	}

	return skipWhitespace(text, 0)
}

func (e *Expander) isValidPosition(abbr *config.Abbreviation, commandLine, word string) bool {
	if abbr.Options != nil && abbr.Options.Position == "anywhere" {
		return true
	}

	beforeWord := commandLine[:len(commandLine)-len(word)]
	beforeWord = strings.TrimSpace(beforeWord)
	// If nothing before the word, it's at command position
	if beforeWord == "" {
		return true
	}

	for _, sep := range commandSeparators {
		if isCommandSeparatorMatch(beforeWord, sep) {
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

func (e *Expander) processSnippet(abbr *config.Abbreviation, word string, commandLine string) (string, error) {
	snippet := abbr.Snippet
	// Perform regex capture-group substitution when regex is configured
	if abbr.Options != nil && abbr.Options.Regex != "" {
		re, err := compileRegexCached(abbr.Options.Regex)
		if err != nil {
			return "", fmt.Errorf("invalid regex: %w", err)
		}

		var matches []string
		if abbr.Abbr == "" && strings.HasPrefix(abbr.Options.Regex, "^") {
			currentCommand := e.extractCurrentCommand(commandLine, word)
			matches = re.FindStringSubmatch(currentCommand)
		} else {
			matches = re.FindStringSubmatch(word)
		}

		if len(matches) > 1 {
			for i, name := range re.SubexpNames() {
				if i > 0 && i < len(matches) && name != "" {
					placeholder := "$" + name
					snippet = strings.ReplaceAll(snippet, placeholder, matches[i])
				}
			}
		}
	}
	// Evaluate snippet only when explicitly requested
	if abbr.Options != nil && abbr.Options.Evaluate {
		escapedSnippet := strings.ReplaceAll(snippet, "\"", "\\\"")
		cmd := exec.Command("bash", "-c", "printf '%s' \""+escapedSnippet+"\"")
		output, err := cmd.Output()
		if err != nil {
			return snippet, nil
		}
		return string(output), nil
	}

	return snippet, nil
}
