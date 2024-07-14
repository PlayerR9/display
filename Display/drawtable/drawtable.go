package drawtable

import (
	"strings"

	cdt "github.com/PlayerR9/MyGoLib/CustomData/Table"
	"github.com/gdamore/tcell"
)

// DrawTable represents a table of cells that can be drawn to the screen.
type DrawTable struct {
	*cdt.Table[*ColoredUnit]
}

// NewDrawTable creates a new DrawTable with the given width and height.
//
// Parameters:
//   - width: The width of the drawTable.
//   - height: The height of the drawTable.
//
// Returns:
//   - *DrawTable: The new drawTable.
//
// Behaviors:
//   - If the width or height is negative, the absolute value is used.
func NewDrawTable(width, height int) *DrawTable {
	table := cdt.NewTable[*ColoredUnit](width, height)

	return &DrawTable{table}
}

// GetLines returns each line of the drawTable as a string.
//
// Returns:
//   - []string: The lines of the drawTable.
//
// Behaviors:
//   - Any nil cells in the drawTable are represented by a space character.
//   - The last line may not be full width.
func (dt *DrawTable) GetLines() []string {
	width := dt.GetWidth()
	iter := dt.Iterator()

	var lines []string
	var builder strings.Builder

	for count := 0; ; count++ {
		if count == width {
			lines = append(lines, builder.String())
			builder.Reset()

			count = 0
		}

		unit, err := iter.Consume()
		if err != nil {
			break
		}

		if unit != nil {
			builder.WriteRune(unit.GetContent())
		} else {
			builder.WriteRune(' ')
		}
	}

	if builder.Len() > 0 {
		lines = append(lines, builder.String())
	}

	return lines
}

// WriteLineAt writes a string to the drawTable at the given coordinates.
//
// Parameters:
//   - x: The x-coordinate of the starting cell.
//   - y: The y-coordinate of the starting cell.
//   - line: The string to write to the drawTable.
//   - style: The style of the string.
//   - isHorizontal: A boolean that determines if the string should be written
//     horizontally or vertically.
//
// Behaviors:
//   - This is just a convenience function that converts the string to a sequence
//     of cells and calls WriteHorizontalSequence or WriteVerticalSequence.
//   - x and y are updated to the next available cell after the line is written.
func (dt *DrawTable) WriteLineAt(x, y *int, line string, style tcell.Style, isHorizontal bool) {
	runes := []rune(line)

	sequence := make([]*ColoredUnit, 0, len(runes))

	for _, r := range runes {
		sequence = append(sequence, NewColoredUnit(r, style))
	}

	if isHorizontal {
		dt.WriteHorizontalSequence(x, y, sequence)
	} else {
		dt.WriteVerticalSequence(x, y, sequence)
	}
}
