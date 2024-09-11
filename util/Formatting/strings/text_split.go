package strings

import (
	"strings"
	"unicode/utf8"

	gcers "github.com/PlayerR9/go-commons/errors"
)

// TextSplit represents a split text with a maximum width and height.
type TextSplit struct {
	// lines is the lines of the split text.
	lines []*line_of_splitter

	// max_width is the maximum length of a line.
	max_width int

	// max_height is the maximum number of lines.
	max_height int
}

// Copy is a method that creates a copy of the TextSplit.
//
// Returns:
//   - *TextSplit: A copy of the TextSplit. Never returns nil.
func (ts TextSplit) Copy() *TextSplit {
	lines_copy := make([]*line_of_splitter, 0, len(ts.lines))

	for _, line := range ts.lines {
		lines_copy = append(lines_copy, line.Copy())
	}

	return &TextSplit{
		max_width:  ts.max_width,
		lines:      lines_copy,
		max_height: ts.max_height,
	}
}

// Runes is a method of TextSplit that returns the runes of the TextSplit.
//
// Returns:
//   - [][]rune: A slice of runes representing the words in the TextSplit.
//
// Behaviors:
//   - It is always a slice of runes with one line.
func (ts TextSplit) Runes() [][]rune {
	if len(ts.lines) == 0 {
		return [][]rune{{}}
	}

	runeTable := make([][]rune, 0, len(ts.lines))

	for _, line := range ts.lines {
		row := line.Runes()
		runeTable = append(runeTable, row[0])
	}

	return runeTable
}

// NewTextSplit creates a new TextSplit with the given maximum width and height.
//
// Parameters:
//   - max_width: The maximum length of a line.
//   - max_height: The maximum number of lines.
//
// Returns:
//   - *TextSplit: A pointer to the newly created TextSplit.
//   - error: An error of type *errors.ErrInvalidParameter if the maxWidth or
//     maxHeight is less than 0.
func NewTextSplit(max_width, max_height int) (*TextSplit, error) {
	if max_width < 0 {
		return nil, gcers.NewErrInvalidParameter(
			"max_width",
			gcers.NewErrGTE(0),
		)
	}

	if max_height < 0 {
		return nil, gcers.NewErrInvalidParameter(
			"max_height",
			gcers.NewErrGTE(0),
		)
	}

	return &TextSplit{
		max_width:  max_width,
		lines:      make([]*line_of_splitter, 0, max_height),
		max_height: max_height,
	}, nil
}

// can_insert_word_at is a helper method that checks if a given word can be inserted
// into a specific line without exceeding the width of the TextSplit.
//
// Parameters:
//   - word: The word to check.
//   - line_index: The index of the line to check.
//
// Returns:
//   - bool: True if the word can be inserted into the line at line_index without
//     exceeding the width, and false otherwise.
func (ts *TextSplit) can_insert_word_at(word string, line_index int) bool {
	if ts == nil || line_index < 0 || line_index >= len(ts.lines) {
		return false
	}

	wordLen := utf8.RuneCountInString(word)
	totalLen := ts.lines[line_index].len + wordLen

	return totalLen+1 <= ts.max_width
}

// InsertWord is a method that attempts to insert a given word into
// the TextSplit.
//
// Parameters:
//   - word: The word to insert.
//
// Returns:
//   - bool: True if the word was successfully inserted, and false if the word is
//     too long to fit within the width of the TextSplit.
func (ts *TextSplit) InsertWord(word string) bool {
	if ts == nil {
		return false
	}

	if len(ts.lines) < ts.max_height {
		word_len := utf8.RuneCountInString(word)

		if word_len > ts.max_width {
			return false
		}

		los := new_line_of_splitter(word)

		ts.lines = append(ts.lines, los)

		return true
	}

	last_line_index := ts.max_height - 1

	// Check if adding the next word to the last line exceeds the width.
	// If it does, we shift the words of the last line to the left.
	ok := ts.can_insert_word_at(word, last_line_index)

	for !ok && last_line_index >= 0 {
		last_line := ts.lines[last_line_index]

		first_word, _ := last_line.shift_left()
		last_line.insert_word(word)

		word = first_word

		last_line_index--
	}

	ok = ts.can_insert_word_at(word, last_line_index)
	if !ok {
		return false
	}

	last_line := ts.lines[last_line_index]
	last_line.insert_word(word)

	return true
}

// InsertWords is a method that attempts to insert multiple words into the TextSplit.
//
// Parameters:
//   - words: The words to insert.
//
// Returns:
//   - int: The index of the first word that could not be inserted, or -1 if all words were inserted.
func (ts *TextSplit) InsertWords(words []string) int {
	if ts == nil {
		return 0
	}

	for i, word := range words {
		ok := ts.InsertWord(word)
		if !ok {
			return i
		}
	}

	return -1
}

// can_shift_up is an helper method that checks if the first word of a given line
// can be shifted up to the previous line without exceeding the width.
//
// Parameters:
//   - line_index: The index of the line to check.
//
// Returns:
//   - bool: True if the first word of the line at line_index can be shifted up to the
//     previous line without exceeding the width, and false otherwise.
func (ts TextSplit) can_shift_up(line_index int) bool {
	ok := ts.can_insert_word_at(ts.lines[line_index].line[0], line_index-1)

	return ok
}

// shift_up is an helper method that shifts the first word of a given line up to
// the previous line.
//
// Parameters:
//   - line_index: The index of the line to shift up.
//
// Returns:
//   - bool: True if the line was successfully shifted up, and false otherwise.
func (ts *TextSplit) shift_up(line_index int) bool {
	if ts == nil || line_index < 0 || line_index >= len(ts.lines) {
		return false
	}

	lastLine := ts.lines[line_index]
	firstWord, _ := lastLine.shift_left()

	secondLastLine := ts.lines[line_index-1]
	secondLastLine.insert_word(firstWord)

	return true
}

// Height is a method that returns the height of the TextSplit.
//
// Returns:
//   - int: The height of the TextSplit.
func (ts TextSplit) Height() int {
	height := len(ts.lines)
	return height
}

// Lines is a method that returns the lines of the TextSplit.
//
// Returns:
//   - []string: The lines of the TextSplit.
func (ts TextSplit) Lines() []string {
	if len(ts.lines) == 0 {
		return nil
	}

	lines := make([]string, 0, len(ts.lines))

	for _, line := range ts.lines {
		str := strings.Join(line.line, " ")
		lines = append(lines, str)
	}

	return lines
}

// FirstLine is a method that returns the first line of the TextSplit.
//
// Returns:
//   - []string: The first line of the TextSplit, or nil if the TextSplit is empty.
//
// Behaviors:
//   - If the TextSplit is empty, the method returns nil.
func (ts TextSplit) FirstLine() []string {
	if len(ts.lines) == 0 {
		return nil
	}

	return ts.lines[0].line
}

// FurthestRightEdge is a method that returns the number of characters in
// the longest line of the TextSplit.
//
// Returns:
//   - int: The number of characters in the longest line.
//   - bool: True if the TextSplit is not empty, and false otherwise.
func (ts TextSplit) FurthestRightEdge() (int, bool) {
	if len(ts.lines) == 0 {
		return 0, false
	}

	max := ts.lines[0].len

	for _, line := range ts.lines[1:] {
		if line.len > max {
			max = line.len
		}
	}

	return max, true
}
