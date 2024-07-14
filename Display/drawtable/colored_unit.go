package drawtable

import (
	"github.com/gdamore/tcell"
)

// ColoredUnit represents a unit that contains a single color.
type ColoredUnit struct {
	// content is the content of the color.
	content rune

	// style is the style of the color.
	style tcell.Style
}

// Runes returns the content of the color as a 2D slice of runes given the size of the table.
//
// Parameters:
//   - width: The width of the table.
//   - height: The height of the table.
//
// Returns:
//   - [][]rune: The content of the color as a 2D slice of runes.
//   - error: An error if the content could not be converted to runes.
//
// Behaviors:
//   - Always assume that the width and height are greater than 0. No need to check for
func (cu *ColoredUnit) Runes(width, height int) ([][]rune, error) {
	return [][]rune{{cu.content}}, nil
}

// NewColoredUnit creates a new ColorUnit with the given content and style.
//
// Parameters:
//   - content: The content of the color.
//   - style: The style of the color.
//
// Returns:
//   - *ColoredUnit: The new ColoredUnit.
func NewColoredUnit(content rune, style tcell.Style) *ColoredUnit {
	return &ColoredUnit{
		content: content,
		style:   style,
	}
}

// GetContent returns the content of the color.
//
// Returns:
//   - rune: The content of the color.
func (cu *ColoredUnit) GetContent() rune {
	return cu.content
}

// GetStyle returns the style of the color.
//
// Returns:
//   - tcell.Style: The style of the color.
func (cu *ColoredUnit) GetStyle() tcell.Style {
	return cu.style
}

// LineToCells converts a string to a slice of ColoredUnits with the given style.
//
// Parameters:
//   - line: The string to convert.
//   - style: The style to use.
//
// Returns:
//   - []*ColoredUnit: The slice of ColoredUnits.
func LineToCells(line string, style tcell.Style) []*ColoredUnit {
	cells := make([]*ColoredUnit, 0, len(line))

	for _, char := range line {
		cells = append(cells, NewColoredUnit(char, style))
	}

	return cells
}
