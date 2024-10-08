package Section

import (
	"fmt"
	"strings"
	"unicode/utf8"

	gcstr "github.com/PlayerR9/display/util/Formatting/strings"
	gda "github.com/PlayerR9/go-debug/assert"
)

const (
	// Hellip defines the string to be used as an ellipsis when the content
	// of the MultilineText is truncated.
	// It is set to "...", which is the standard representation of an ellipsis
	// in text.
	Hellip string = "..."

	// HellipLen defines the length of the Hellip string.
	// It is set to 3, which is the number of characters in the Hellip string.
	HellipLen int = len(Hellip)

	// Space defines the string to be used as a space when writing content
	// into the MultilineText.
	// It is set to " ", which is the standard representation of a space in
	// text.
	Space string = " "

	// FieldSpacing defines the number of spaces between each field (word)
	// when they are written into the MultilineText.
	// It is set to 1, meaning there will be one spaces between each field.
	FieldSpacing int = 1

	// IndentLevel defines the number of spaces used for indentation when
	// writing content into the MultilineText.
	// It is set to 2, meaning there will be two spaces at the start of each
	// new line of content.
	IndentLevel int = 3
)

// MultilineText represents a box that contains content.
type MultilineText struct {
	// lines is a two-dimensional slice of strings representing the content
	// of the box.
	lines [][]string
}

// FromTextBlock sets the content of the MultilineText to the specified
// two-dimensional slice of strings.
//
// Parameters:
//   - lines - a two-dimensional slice of strings representing the content
//     of the MultilineText.
//
// Returns:
//   - error - an error if the content could not be set.
func (mlt *MultilineText) FromTextBlock(lines [][]string) error {
	if len(lines) == 0 {
		mlt.lines = [][]string{{}}

		return nil
	}

	mlt.lines = lines

	return nil
}

// GetRawContent returns the raw content of the MultilineText.
//
// Returns:
//   - [][]string - the raw content of the MultilineText.
func (mlt *MultilineText) GetRawContent() [][]string {
	return mlt.lines
}

// Runes returns the content of the unit as a 2D slice of runes
// given the size of the table.
//
// Parameters:
//   - width: The width of the table.
//   - height: The height of the table.
//
// Returns:
//   - [][]rune: The content of the unit as a 2D slice of runes.
//   - error: An error if the content could not be converted to runes.
//
// Behaviors:
//   - Always assume that the width and height are greater than 0. No need to check for
//     this.
//   - Errors are only for critical issues, such as the content not being able to be
//     converted to runes. However, out of bounds or other issues should not error.
//     Instead, the content should be drawn as much as possible before unable to be
//     drawn.
func (mlt *MultilineText) ApplyRender(width, height int) ([]*Render, error) {
	tss, err := mlt.forceApply(width, height)
	if err != nil {
		return nil, err
	}

	totalHeight := 0
	tableHeight := height

	for _, ts := range tss {
		totalHeight += ts.Height()

		if totalHeight > tableHeight {
			break
		}
	}

	var renders []*Render

	var runeTable [][]rune
	yCoord := 0

	for _, ts := range tss {
		currentHeight := ts.Height()

		canRenderMore := currentHeight+yCoord <= tableHeight
		if !canRenderMore {
			renders = append(renders, NewRender(runeTable))

			runeTable = nil
			yCoord = 0
		} else {
			runeTable = append(runeTable, ts.Runes()...)

			yCoord += currentHeight
		}
	}

	if runeTable != nil {
		renders = append(renders, NewRender(runeTable))
	}

	return renders, nil
}

// processLine processes a line of text represented by a slice of fields.
// It calculates the number of lines the text would occupy if split into
// lines of a specified width. If the text cannot be split into lines of
// the specified width, it replaces the suffix of the text with a hellip
// and adds the resulting line to the TextSplitter. If the text can be split
// into more than one line, it creates a new line with the first field and
// as many subsequent fields as can fit into the line width, adding a hellip
// if necessary. If the text can be split into exactly one line, it splits
// the text into equal-sized lines and adds the first line to the TextSplitter.
//
// Parameters:
//   - isFirst - a boolean indicating whether the line is the first line of text.
//   - maxWidth - the maximum width of the line.
//   - ts - the TextSplitter to add the line to.
//   - words - a slice of fields representing the line of text.
//
// Returns:
//   - *sext.TextSplit - the updated TextSplitter.
//   - bool - a boolean indicating whether the text was truncated.
//   - error - an error if the text could not be processed.
func (mlt *MultilineText) processLine(isFirst bool, maxWidth int, ts *gcstr.TextSplit, words []string) (*gcstr.TextSplit, bool, error) {
	if !isFirst {
		maxWidth -= IndentLevel
	}

	numberOfLines, err := gcstr.CalculateNumberOfLines(words, maxWidth)

	if err != nil {
		line := strings.Join(words, "")[:maxWidth]

		line, ok := gcstr.ReplaceSuffix(line, Hellip)
		if !ok {
			return nil, false, gcstr.NewErrLongerSuffix(line, Hellip)
		}

		ok = ts.InsertWord(line)
		gda.AssertOk(ok, "ts.InsertWord(%q)", line)

		return ts, true, nil
	}

	if numberOfLines > 1 {
		wordsProcessed := []string{words[0]}
		wpLen := utf8.RuneCountInString(words[0])

		var nextField string

		for i, currentField := range words[1 : len(words)-1] {
			nextField = words[i+1]

			totalLen := wpLen + 2 + utf8.RuneCountInString(currentField) +
				utf8.RuneCountInString(nextField)

			if totalLen+HellipLen > maxWidth {
				currentField += Hellip

				wordsProcessed = append(wordsProcessed, currentField)
				wpLen += utf8.RuneCountInString(currentField) + 1
				break
			}

			wordsProcessed = append(wordsProcessed, currentField)
			wpLen += utf8.RuneCountInString(currentField) + 1
		}

		if wpLen+1+utf8.RuneCountInString(nextField)+HellipLen <= maxWidth {
			wordsProcessed = append(wordsProcessed, nextField)
			wpLen += utf8.RuneCountInString(nextField) + 1
		}

		firstNotInserted := ts.InsertWords(wordsProcessed)
		gda.AssertF(firstNotInserted == -1, "ts.InsertWords(%q)", wordsProcessed)

		return ts, true, nil
	} else {
		halfTs, err := gcstr.SplitInEqualSizedLines(
			words, maxWidth, numberOfLines,
		)

		if err != nil {
			return nil, false, fmt.Errorf("could not split text: %s", err.Error())
		}

		wordsProcessed := halfTs.FirstLine()

		firstNotInserted := ts.InsertWords(wordsProcessed)
		gda.AssertF(firstNotInserted == -1, "ts.InsertWords(%q)", wordsProcessed)

		return ts, false, nil
	}
}

// createTextSplitter takes a two-dimensional slice of strings
// representing a list of fields and a width, and creates a
// TextSplitter that splits the fields into lines of the specified
// width. It processes the first line of fields separately from
// the other lines. If an error occurs while processing a line,
// it returns an error with a message indicating the line number
// and the original error.
//
// The function returns a pointer to the created TextSplitter
// and an error. If no errors occur during the creation of the
// TextSplitter, the error is nil.
func (mlt *MultilineText) createTextSplitter(lines [][]string, maxWidth, maxHeight int) (*gcstr.TextSplit, error) {
	ts, err := gcstr.NewTextSplit(maxWidth, maxHeight)
	if err != nil {
		return nil, fmt.Errorf("could not create TextSplitter: %s", err.Error())
	}

	possibleNewLine := false

	ts, possibleNewLine, err = mlt.processLine(true, maxWidth, ts, lines[0])
	if err != nil {
		return nil, err
	}

	for _, line := range lines[1:] {
		if possibleNewLine {
			for len(line) > 0 {
				ok := ts.InsertWord(line[0])
				if !ok {
					break
				}

				line = line[1:]
			}
		}

		if len(line) == 0 {
			continue
		}

		ts, possibleNewLine, err = mlt.processLine(false, maxWidth, ts, line)
		if err != nil {
			return nil, fmt.Errorf("could not process line: %s", err.Error())
		}
	}

	return ts, nil
}

// apply takes a maximum width and height, and applies the content of the MultilineText
// to the specified width and height. It splits the content into lines of the specified
// width, and optimizes the text if possible.
//
// Parameters:
//   - maxWidth - the maximum width of the content.
//   - maxHeight - the maximum height of the content.
//
// Returns:
//   - []*sext.TextSplit - a slice of TextSplit objects representing the optimized content.
//   - error - an error if the content could not be applied.
func (mlt *MultilineText) forceApply(maxWidth, maxHeight int) ([]*gcstr.TextSplit, error) {
	finalTs := make([]*gcstr.TextSplit, 0, len(mlt.lines))

	for _, line := range mlt.lines {
		sentences := [][]string{line}

		ts, err := mlt.createTextSplitter(sentences, maxWidth, maxHeight)
		if err != nil {
			return nil, err
		}

		// If it is possible to optimize the text, optimize it.
		// Otherwise, the unoptimized text is also fine.
		optimizedTs, err := gcstr.SplitInEqualSizedLines(ts.Lines(), maxWidth, -1)
		if err != nil {
			finalTs = append(finalTs, ts)
		} else {
			finalTs = append(finalTs, optimizedTs)
		}
	}

	return finalTs, nil
}
