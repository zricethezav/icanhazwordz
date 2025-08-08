// Have you ever asked yourself can i haz wordz? More specifically, does a string contain any english words?
// This package provides functionality to detect English words in a given text using, you guessed it, our favorite friend, ahocorasick.
package icanhazwordz

import (
	"strings"

	ahocorasick "github.com/BobuSumisu/aho-corasick"
)

// Default filter ignores single character words.
// If you want to include them, set MinLength to 1.
var DefaultFilter = Filter{MinLength: 2, MaxLength: 0}

// Match represents a found word match
type Match struct {
	Word     string `json:"word"`
	StartPos int    `json:"start_pos"`
	EndPos   int    `json:"end_pos"`
}

// Result represents the result of English word detection
type Result struct {
	WordCount   int      `json:"word_count"`
	UniqueWords []string `json:"unique_words"`
	AllMatches  []Match  `json:"all_matches"`
}

// Filter defines criteria for filtering words
type Filter struct {
	MinLength   int `json:"min_length"`   // Minimum word length (inclusive)
	MaxLength   int `json:"max_length"`   // Maximum word length (inclusive, 0 means no limit)
	ExactLength int `json:"exact_length"` // Exact word length (0 means no exact match required)
}

// Searcher provides nltk English word detection with configurable filters
type Searcher struct {
	trie   *ahocorasick.Trie
	words  []string
	filter Filter
}

// NewSearcher creates a new searcher with the specified word filters
func NewSearcher(filter Filter) *Searcher {
	words := getFilteredWords(filter)

	return &Searcher{
		trie:   ahocorasick.NewTrieBuilder().AddStrings(words).Build(),
		words:  words,
		filter: filter,
	}
}

// GetWordCount returns the number of words in this searcher's dictionary
func (s *Searcher) GetWordCount() int {
	return len(s.words)
}

// GetWords returns all words in this searcher's dictionary
func (s *Searcher) GetWords() []string {
	return append([]string(nil), s.words...)
}

// Find analyzes text using this searcher's word list
func (s *Searcher) Find(text string) Result {
	if text == "" {
		return Result{
			WordCount:   0,
			UniqueWords: []string{},
			AllMatches:  []Match{},
		}
	}

	textLower := strings.ToLower(text)
	matches := s.trie.MatchString(textLower)

	var allMatches []Match
	uniqueWordsMap := make(map[string]bool)

	for _, match := range matches {
		word := match.MatchString()
		startPos := int(match.Pos())
		endPos := startPos + len(word)

		allMatches = append(allMatches, Match{
			Word:     word,
			StartPos: startPos,
			EndPos:   endPos,
		})

		uniqueWordsMap[word] = true

	}

	var uniqueWords []string
	for word := range uniqueWordsMap {
		uniqueWords = append(uniqueWords, word)
	}

	return Result{
		WordCount:   len(matches),
		UniqueWords: uniqueWords,
		AllMatches:  allMatches,
	}
}

// getFilteredWords returns words filtered by the given criteria
func getFilteredWords(filter Filter) []string {
	var filtered []string

	for _, word := range nltkWords {
		wordLen := len(word)

		// Check exact length filter first
		if filter.ExactLength > 0 {
			if wordLen == filter.ExactLength {
				filtered = append(filtered, word)
			}
			continue
		}

		// Check min length
		if filter.MinLength > 0 && wordLen < filter.MinLength {
			continue
		}

		// Check max length (0 means no limit)
		if filter.MaxLength > 0 && wordLen > filter.MaxLength {
			continue
		}

		filtered = append(filtered, word)
	}

	return filtered
}
