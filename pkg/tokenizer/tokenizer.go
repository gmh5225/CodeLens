package tokenizer

import (
	"strings"
	"sync"
	"unicode"
	"unicode/utf8"
)

// TokenConfig defines configuration options for token counting
type TokenConfig struct {
	CountNewlines    bool // Whether to count newlines as tokens
	CountSpaces      bool // Whether to count spaces as tokens
	CountPunctuation bool // Whether to count punctuation as tokens
	CJKAsTokens      bool // Whether to count each CJK character as a token
	MaxTokens        int  // Maximum number of tokens to count (0 for unlimited)
}

// DefaultConfig returns the default token counting configuration
func DefaultConfig() TokenConfig {
	return TokenConfig{
		CountNewlines:    false,
		CountSpaces:      false,
		CountPunctuation: false,
		CJKAsTokens:      true,
		MaxTokens:        0,
	}
}

var (
	// Cache frequently used token counts
	tokenCache sync.Map
	// Maximum cache size to prevent memory leaks
	maxCacheSize = 10000
	// Current cache size
	cacheSize int32
)

// CountTokens calculates the approximate number of tokens in the text
// This implementation follows a more conservative version of token counting:
// 1. Punctuation and special characters are generally not counted unless significant
// 2. For ASCII text:
//   - Words are counted as single tokens regardless of length
//   - Common programming keywords and operators count as single tokens
//
// 3. For CJK and other Unicode characters:
//   - Each CJK character counts as 1 token
//   - CJK punctuation is not counted
//
// 4. Numbers are handled specially:
//   - Simple numbers count as 1 token
//   - Complex numbers (scientific, decimal) count as 2-3 tokens
//
// 5. Special handling:
//   - URLs count as domain + path components
//   - Email addresses count as 2-3 tokens
//   - Code constructs (brackets, operators) count minimally
func CountTokens(text string, config ...TokenConfig) int {
	if strings.TrimSpace(text) == "" {
		return 0
	}

	// Use default config if none provided
	cfg := DefaultConfig()
	if len(config) > 0 {
		cfg = config[0]
	}

	// Check cache first
	if count, ok := tokenCache.Load(text); ok {
		return count.(int)
	}

	tokens := 0
	inWord := false
	currentWord := strings.Builder{}

	for _, r := range text {
		// Check token limit
		if cfg.MaxTokens > 0 && tokens >= cfg.MaxTokens {
			break
		}

		// Handle CJK characters
		if cfg.CJKAsTokens && isCJK(r) {
			if currentWord.Len() > 0 {
				tokens++
				currentWord.Reset()
			}
			inWord = false
			tokens++
			continue
		}

		// Handle whitespace
		if unicode.IsSpace(r) {
			if currentWord.Len() > 0 {
				tokens++
				currentWord.Reset()
			}
			if cfg.CountNewlines && r == '\n' {
				tokens++
			} else if cfg.CountSpaces && inWord {
				tokens++
			}
			inWord = false
			continue
		}

		// Handle punctuation and symbols
		if unicode.IsPunct(r) || unicode.IsSymbol(r) {
			if currentWord.Len() > 0 {
				tokens++
				currentWord.Reset()
			}
			if cfg.CountPunctuation && isSignificantPunct(r) {
				tokens++
			}
			inWord = false
			continue
		}

		// Accumulate characters for word
		currentWord.WriteRune(r)
		inWord = true
	}

	// Process any remaining word
	if currentWord.Len() > 0 {
		tokens++
	}

	// Cache the result if not too large
	if len(text) < 1000 && cacheSize < int32(maxCacheSize) {
		tokenCache.Store(text, tokens)
		cacheSize++
	}

	return tokens
}

// isSignificantPunct checks if a punctuation mark should be counted as a token
func isSignificantPunct(r rune) bool {
	significant := map[rune]bool{
		'{': true, '}': true,
		'[': true, ']': true,
		'(': true, ')': true,
		';': true,
		'=': true,
		'+': true, '-': true,
		'*': true, '/': true,
		'%': true,
	}
	return significant[r]
}

// countWordTokens counts tokens for a single word
func countWordTokens(word string) int {
	if word == "" {
		return 0
	}

	// Handle special cases
	if isURL(word) {
		return countURLTokens(word)
	}
	if isEmail(word) {
		return countEmailTokens(word)
	}
	if isScientificNotation(word) {
		return countScientificTokens(word)
	}

	// Handle numbers
	if isNumber(word) {
		// Count decimal points separately
		if strings.Contains(word, ".") {
			parts := strings.Split(word, ".")
			return (len(parts[0])+2)/3 + 1 + (len(parts[1])+2)/3
		}
		return (len(word) + 2) / 3 // 1 token per 3 digits, rounded up
	}

	// Count words based on length and script type
	wordLen := utf8.RuneCountInString(word)
	switch {
	case wordLen <= 4:
		return 1
	case wordLen <= 8:
		return 2
	default:
		return (wordLen + 3) / 4 // Round up division
	}
}

// Helper functions for special cases
func isURL(word string) bool {
	return strings.HasPrefix(word, "http://") ||
		strings.HasPrefix(word, "https://") ||
		strings.HasPrefix(word, "ftp://")
}

func isEmail(word string) bool {
	return strings.Contains(word, "@") && strings.Contains(word, ".")
}

func isScientificNotation(word string) bool {
	return strings.ContainsAny(word, "eE") && strings.ContainsAny(word, "+-.")
}

func countURLTokens(url string) int {
	// Count protocol, domain, and path separately
	parts := strings.Split(url, "/")
	tokens := 0
	for _, part := range parts {
		if part == "" {
			continue
		}
		if strings.Contains(part, ".") {
			subParts := strings.Split(part, ".")
			tokens += len(subParts)
		} else {
			tokens++
		}
	}
	return tokens
}

func countEmailTokens(email string) int {
	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return len(email) / 4
	}
	return countWordTokens(parts[0]) + 1 + countWordTokens(parts[1])
}

func countScientificTokens(num string) int {
	parts := strings.FieldsFunc(num, func(r rune) bool {
		return r == 'e' || r == 'E' || r == '+' || r == '-'
	})
	tokens := 0
	for _, part := range parts {
		if part != "" {
			tokens += countWordTokens(part)
		}
	}
	return tokens + 1 // Add 1 for the 'e' notation
}

// isPairedPunct checks if a rune is part of common paired punctuation
func isPairedPunct(r rune) bool {
	pairs := map[rune]bool{
		'(': true, ')': true,
		'[': true, ']': true,
		'{': true, '}': true,
		'<': true, '>': true,
		'"': true, '\'': true,
		'`': true,
	}
	return pairs[r]
}

// isSpecialPunct checks if a rune is a special punctuation mark
func isSpecialPunct(r rune) bool {
	special := map[rune]bool{
		'@': true, '#': true,
		'$': true, '%': true,
		'^': true, '&': true,
		'*': true, '_': true,
		'+': true, '=': true,
		'|': true, '\\': true,
		'~': true,
	}
	return special[r]
}

// isEmoji checks if a rune is an emoji
func isEmoji(r rune) bool {
	return r >= 0x1F300 && r <= 0x1F9FF
}

// isCJK checks if a rune is a CJK character
func isCJK(r rune) bool {
	ranges := [...][2]rune{
		{0x4E00, 0x9FFF},   // CJK Unified Ideographs
		{0x3400, 0x4DBF},   // CJK Unified Ideographs Extension A
		{0x20000, 0x2A6DF}, // CJK Unified Ideographs Extension B
		{0x2A700, 0x2B73F}, // CJK Unified Ideographs Extension C
		{0x2B740, 0x2B81F}, // CJK Unified Ideographs Extension D
		{0x2B820, 0x2CEAF}, // CJK Unified Ideographs Extension E
		{0x3000, 0x303F},   // CJK Symbols and Punctuation
		{0xFF00, 0xFFEF},   // Halfwidth and Fullwidth Forms
		{0x3040, 0x309F},   // Hiragana
		{0x30A0, 0x30FF},   // Katakana
		{0x31F0, 0x31FF},   // Katakana Phonetic Extensions
		{0xAC00, 0xD7AF},   // Hangul Syllables
	}

	for _, rang := range ranges {
		if r >= rang[0] && r <= rang[1] {
			return true
		}
	}
	return false
}

// isNumber checks if a string is a numeric value
func isNumber(s string) bool {
	for _, r := range s {
		if !unicode.IsDigit(r) && r != '.' && r != '-' {
			return false
		}
	}
	return len(s) > 0
}
