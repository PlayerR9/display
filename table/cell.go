package table

import "github.com/gdamore/tcell"

// Cell is a table cell.
type Cell struct {
	// Char is the character of the cell.
	Char rune

	// Style is the Style of the cell.
	Style tcell.Style
}

// NewCell creates a new table cell.
//
// Parameters:
//   - char: The character of the cell.
//   - style: The style of the cell.
//
// Returns:
//   - *Cell: The new table cell. Never returns nil.
func NewCell(char rune, style tcell.Style) *Cell {
	return &Cell{
		Char:  char,
		Style: style,
	}
}

// LineToCells converts a string to a slice of DrawCells with the given style.
//
// Parameters:
//   - line: The string to convert.
//   - style: The style to use.
//
// Returns:
//   - []*DrawCell: The slice of DrawCells.
func LineToCells(line string, style tcell.Style) []*Cell {
	if line == "" {
		return nil
	}

	cells := make([]*Cell, 0, len(line))

	for _, char := range line {
		cells = append(cells, NewCell(char, style))
	}

	return cells
}
