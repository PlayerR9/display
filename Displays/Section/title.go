package Section

import (
	"fmt"
	"strings"
	"unicode/utf8"

	gcfs "github.com/PlayerR9/display/util/Formatting/strings"
	gcers "github.com/PlayerR9/go-commons/errors"
)

const (
	Asterisks string = "***"

	AsterisksLen int = len(Asterisks)

	TitleMinWidth int = 2 * (AsterisksLen + 1)
)

// Title represents the header of a process or application.
// It contains information about the title, current process, counters, message buffer,
// channels for receiving messages and errors, synchronization primitives, and the width
// of the header.
type Title struct {
	title    string
	subtitle string // empty string if no subtitle
}

// ForceDraw draws the title to the draw table.
//
// Parameters:
//   - table: The draw table.
//   - x: The x coordinate to draw the title at.
//   - y: The y coordinate to draw the title at.
//
// Returns:
//   - error: An error if the title could not be drawn.
func (t *Title) ApplyRender(width, height int) ([]*Render, error) {
	// 1. Generate the full title
	var fullTitle string

	if t.subtitle == "" {
		fullTitle = t.title
	} else {
		var builder strings.Builder

		builder.WriteString(t.title)
		builder.WriteRune(' ')
		builder.WriteRune('-')
		builder.WriteRune(' ')
		builder.WriteString(t.subtitle)

		fullTitle = builder.String()
	}

	// 2. Generate the lines
	lines, err := t.tryToFitLines(fullTitle, width)
	if err != nil {
		return nil, fmt.Errorf("failed to generate lines: %s", err.Error())
	}

	// 3. Write the lines with centered alignment
	var renders []*Render

	var runeTable [][]rune
	yCoord := 0

	for i := 0; i < len(lines); i++ {
		startPos := (width - utf8.RuneCountInString(lines[i])) / 2

		if yCoord < height {
			row := make([]rune, width)
			copy(row[startPos:], []rune(lines[i]))

			runeTable = append(runeTable, row)
		} else {
			renders = append(renders, NewRender(runeTable))

			runeTable = nil
			yCoord = 0
		}
	}

	if runeTable != nil {
		renders = append(renders, NewRender(runeTable))
	}

	return renders, nil
}

// GetRawContent returns the raw content of the Title.
//
// Returns:
//   - [][]string: The raw content of the Title.
func (t *Title) GetRawContent() [][]string {
	// 1. Generate the full title
	var fullTitle string

	if t.subtitle == "" {
		fullTitle = t.title
	} else {
		var builder strings.Builder

		builder.WriteString(t.title)
		builder.WriteRune(' ')
		builder.WriteRune('-')
		builder.WriteRune(' ')
		builder.WriteString(t.subtitle)

		fullTitle = builder.String()
	}

	return [][]string{{fullTitle}}
}

// NewTitle creates a new Title with the given title and a style.
//
// Parameters:
//   - title: The title of the new Title.
//   - style: The style of the new Title.
//
// Returns:
//   - *Title: The new Title.
func NewTitle(title string) *Title {
	return &Title{
		title:    title,
		subtitle: "",
	}
}

// SetSubtitle sets the subtitle of the Title.
//
// Parameters:
//   - subtitle: The new subtitle.
//
// Behaviors:
//   - If the subtitle is an empty string, the subtitle is removed.
func (t *Title) SetSubtitle(subtitle string) {
	t.subtitle = subtitle
}

// tryToFitLines is a helper method that tries to fit the full title in the draw table.
//
// Parameters:
//   - table: The draw table.
//   - fullTitle: The full title.
//
// Returns:
//   - []string: The lines of the title.
//   - error: An error if the full title could not be split in lines.
func (t *Title) tryToFitLines(fullTitle string, width int) ([]string, error) {
	lines, err := generateLines(fullTitle, width)
	if err == nil {
		return lines, nil
	}

	fullTitle = gcfs.FitString(fullTitle, width-TitleMinWidth)

	fullTitle, ok := gcfs.ReplaceSuffix(fullTitle, Hellip)
	if !ok {
		return nil, gcfs.NewErrLongerSuffix(fullTitle, Hellip)
	}

	var builder strings.Builder

	builder.WriteString(Asterisks)
	builder.WriteString(Space)
	builder.WriteString(fullTitle)
	builder.WriteString(Space)
	builder.WriteString(Asterisks)

	return []string{builder.String()}, nil
}

// generateLines is a helper method that generates the lines of the title.
//
// Parameters:
//   - fullTitle: The full title.
//   - width: The width of the lines.
//
// Returns:
//   - []string: The lines of the title.
//   - error: An error if the full title could not be split in lines.
func generateLines(fullTitle string, width int) ([]string, error) {
	contents := strings.Fields(fullTitle) // FIXME: Use a better method to split the text

	numberOfLines, err := gcfs.CalculateNumberOfLines(contents, width-TitleMinWidth)
	if err != nil {
		ok := gcers.Is[*gcfs.ErrLinesGreaterThanWords](err)
		if !ok {
			return nil, fmt.Errorf("could not calculate number of lines: %s", err.Error())
		}
	}

	ts, err := gcfs.SplitInEqualSizedLines(contents, width-TitleMinWidth, numberOfLines)
	if err != nil {
		return nil, fmt.Errorf("could not split text in equal sized lines: %s", err.Error())
	}

	lines := ts.Lines()
	var builder strings.Builder

	for i := 0; i < len(lines); i++ {
		builder.WriteString(Asterisks)
		builder.WriteString(Space)
		builder.WriteString(lines[i])
		builder.WriteString(Space)
		builder.WriteString(Asterisks)

		lines[i] = builder.String()
		builder.Reset()
	}

	return lines, nil
}

// GetTitle returns the title of the Title.
//
// Returns:
//   - string: The title of the Title.
func (t *Title) GetTitle() string {
	return t.title
}

// GetSubtitle returns the subtitle of the Title.
//
// Returns:
//   - string: The subtitle of the Title.
func (t *Title) GetSubtitle() string {
	return t.subtitle
}
