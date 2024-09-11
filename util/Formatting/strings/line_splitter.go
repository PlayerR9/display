package strings

import (
	"strings"
	"unicode/utf8"
)

// line_of_splitter is a helper struct used in the SplitTextInEqualSizedLines function.
// It represents a line of text.
type line_of_splitter struct {
	// The line field is a slice of strings, each representing a word in the line.
	line []string

	// The len field is an integer representing the total length of the line,
	// including spaces between words.
	len int
}

// Copy is a method that creates a copy of the line of splitter.
//
// Returns:
//   - *lineOfSplitter: A copy of the line of splitter. Never returns nil.
func (sl line_of_splitter) Copy() *line_of_splitter {
	lines := make([]string, len(sl.line))
	copy(lines, sl.line)

	return &line_of_splitter{
		line: lines,
		len:  sl.len,
	}
}

// Runes is a method of SpltLine that returns the runes of the line.
//
// Always returns a slice of runes with one line.
//
// Returns:
//   - [][]rune: A slice of runes representing the words in the line.
//
// Behaviors:
//   - It is always a slice of runes with one line.
func (sl line_of_splitter) Runes() [][]rune {
	if len(sl.line) == 0 {
		return [][]rune{{}}
	}

	str := strings.Join(sl.line, " ")

	return [][]rune{[]rune(str)}
}

// new_line_of_splitter is a helper function that creates a new line of
// splitter with the given word.
//
// Parameters:
//   - word: The initial word to add to the line.
//
// Returns:
//   - *line_of_splitter: A pointer to the newly created line of splitter. Never returns nil.
func new_line_of_splitter(word string) *line_of_splitter {
	len := utf8.RuneCountInString(word)

	return &line_of_splitter{
		line: []string{word},
		len:  len,
	}
}

// shift_left is an helper method of SpltLine that removes the first word of the line.
//
// Returns:
//   - string: The word that was removed.
//   - bool: True if the receiver is not nil, false otherwise.
func (sl *line_of_splitter) shift_left() (string, bool) {
	if sl == nil {
		return "", false
	}

	first_word := sl.line[0]

	sl.line = sl.line[1:]
	sl.len -= utf8.RuneCountInString(first_word)
	sl.len-- // Remove the extra space

	return first_word, true
}

// InsertWord is a method of SpltLine that adds a given word to the end of the line.
//
// Parameters:
//   - word: The word to add to the line.
//
// Behaviors:
//   - If the word is an empty string, it is ignored.
func (sl *line_of_splitter) insert_word(word string) {
	if sl == nil || word == "" {
		return
	}

	sl.line = append(sl.line, word)
	sl.len += utf8.RuneCountInString(word)
	sl.len++ // Add the extra space
}
