package icanhazwordz

import (
	"testing"
)

func TestSearcherFindExact(t *testing.T) {
	// Use a simple filter for testing
	searcher := NewSearcher(Filter{ExactLength: 5})

	tests := []struct {
		name       string
		text       string
		numMatches int
	}{
		{
			name:       "exact filter",
			text:       "hello world",
			numMatches: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := searcher.Find(tt.text)
			if result.WordCount != tt.numMatches {
				t.Errorf("expected %d matches, got %d. Words: %v", tt.numMatches, result.WordCount, result.Matches)
			}
		})
	}
}

func TestSearcherFindDefault(t *testing.T) {
	// Use a simple filter for testing
	searcher := NewSearcher(DefaultFilter)
	tests := []struct {
		name       string
		text       string
		numMatches int
	}{
		{
			name:       "default filter",
			text:       "hello world",
			numMatches: 9,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := searcher.Find(tt.text)
			if result.WordCount != tt.numMatches {
				t.Errorf("expected %d matches, got %d. Words: %v", tt.numMatches, result.WordCount, result.Matches)
			}
		})
	}
}

func TestSearcherFindWithFilter(t *testing.T) {
	filter := Filter{MinLength: 3, MaxLength: 4}
	searcher := NewSearcher(filter)

	tests := []struct {
		name       string
		text       string
		numMatches int
	}{
		{
			name:       "test min and max",
			text:       "hello world",
			numMatches: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := searcher.Find(tt.text)
			if result.WordCount != tt.numMatches {
				t.Errorf("expected %d matches, got %d. Words: %v", tt.numMatches, result.WordCount, result.Matches)
			}
		})
	}
}
