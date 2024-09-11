package strings

import (
	"strconv"
	"strings"
)

// ErrLongerSuffix is a struct that represents an error when the suffix is
// longer than the string.
type ErrLongerSuffix struct {
	// Str is the string that is shorter than the suffix.
	Str string

	// Suffix is the Suffix that is longer than the string.
	Suffix string
}

// Error implements the error interface.
//
// Message: "suffix {Suffix} is longer than the string {Str}"
func (e ErrLongerSuffix) Error() string {
	values := []string{
		"suffix",
		strconv.Quote(e.Suffix),
		"is longer than the string",
		strconv.Quote(e.Str),
	}

	msg := strings.Join(values, " ")

	return msg
}

// NewErrLongerSuffix is a constructor of ErrLongerSuffix.
//
// Parameters:
//   - str: The string that is shorter than the suffix.
//   - suffix: The suffix that is longer than the string.
//
// Returns:
//   - *ErrLongerSuffix: A pointer to the newly created error. Never returns nil.
func NewErrLongerSuffix(str, suffix string) *ErrLongerSuffix {
	e := &ErrLongerSuffix{
		Str:    str,
		Suffix: suffix,
	}
	return e
}

// ErrLinesGreaterThanWords is an error type that is returned when the
// number of lines in a text is greater than the number of words.
type ErrLinesGreaterThanWords struct {
	// NumberOfLines is the number of lines in the text.
	NumberOfLines int

	// NumberOfWords is the number of words in the text.
	NumberOfWords int
}

// Error implements the error interface.
//
// Message: "number of lines ({NumberOfLines}) is greater than the number of words ({NumberOfWords})"
func (e ErrLinesGreaterThanWords) Error() string {
	values := []string{
		"number of lines",
		"(",
		strconv.Itoa(e.NumberOfLines),
		")",
		"is greater than the number of words",
		"(",
		strconv.Itoa(e.NumberOfWords),
		")",
	}

	msg := strings.Join(values, " ")

	return msg
}

// NewErrLinesGreaterThanWords is a constructor of ErrLinesGreaterThanWords.
//
// Parameters:
//   - numberOfLines: The number of lines in the text.
//   - numberOfWords: The number of words in the text.
//
// Returns:
//   - *ErrLinesGreaterThanWords: A pointer to the newly created error. Never returns nil.
func NewErrLinesGreaterThanWords(numberOfLines, numberOfWords int) *ErrLinesGreaterThanWords {
	e := &ErrLinesGreaterThanWords{
		NumberOfLines: numberOfLines,
		NumberOfWords: numberOfWords,
	}
	return e
}
