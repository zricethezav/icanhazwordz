package icanhazwordz

import (
	"sort"
	"strings"

	ahocorasick "github.com/BobuSumisu/aho-corasick"
)

// Default filter ignores single character words.
// If you want to include them, set MinLength to 1.
var DefaultFilter = Filter{MinLength: 2, MaxLength: 0, PreferLongestNonOverlapping: false}

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
	Matches     []Match  `json:"matches"`
}

// Filter defines criteria for filtering words
type Filter struct {
	MinLength                   int  // Minimum word length (inclusive)
	MaxLength                   int  // Maximum word length (inclusive, 0 means no limit)
	ExactLength                 int  // Exact word length (0 means no exact match required)
	PreferLongestNonOverlapping bool // Prefer longest non-overlapping matches. e.g. "hello" vs "hell"
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
			Matches:     []Match{},
		}
	}

	textLower := strings.ToLower(text)
	matches := s.trie.MatchString(textLower)

	var allMatches []Match
	for _, match := range matches {
		word := match.MatchString()
		startPos := int(match.Pos())
		endPos := startPos + len(word)

		allMatches = append(allMatches, Match{
			Word:     word,
			StartPos: startPos,
			EndPos:   endPos,
		})
	}

	// Conditionally filter out overlapping matches
	finalMatches := allMatches
	if s.filter.PreferLongestNonOverlapping {
		finalMatches = filterOverlappingMatches(allMatches)
	}

	// Build unique words map from final matches
	uniqueWordsMap := make(map[string]bool)
	for _, match := range finalMatches {
		uniqueWordsMap[match.Word] = true
	}

	var uniqueWords []string
	for word := range uniqueWordsMap {
		uniqueWords = append(uniqueWords, word)
	}

	return Result{
		WordCount:   len(finalMatches),
		UniqueWords: uniqueWords,
		Matches:     finalMatches,
	}
}

// filterOverlappingMatches removes overlapping matches, keeping only the longest non-overlapping ones
func filterOverlappingMatches(matches []Match) []Match {
	if len(matches) == 0 {
		return matches
	}

	// Sort matches by start position, then by length (longest first for same start position)
	sort.Slice(matches, func(i, j int) bool {
		if matches[i].StartPos == matches[j].StartPos {
			return len(matches[i].Word) > len(matches[j].Word) // Longer words first
		}
		return matches[i].StartPos < matches[j].StartPos
	})

	var result []Match
	lastEndPos := -1

	for _, match := range matches {
		// If this match doesn't overlap with the last accepted match
		if match.StartPos >= lastEndPos {
			result = append(result, match)
			lastEndPos = match.EndPos
		}
		// If it overlaps, we skip it (since we sorted by length, the previous one was longer)
	}

	return result
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
