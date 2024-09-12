package screen

import "github.com/gdamore/tcell"

// Drawable is an interface that can be drawn.
type Drawable interface {
	// DrawCell draws a cell.
	//
	// Parameters:
	//   - x: The x position of the cell.
	//   - y: The y position of the cell.
	//   - char: The character of the cell.
	//   - style: The style of the cell.
	DrawCell(x, y int, char rune, style tcell.Style)

	// BgStyle returns the background style.
	//
	// Returns:
	//   - tcell.Style: The background style.
	BgStyle() tcell.Style
}
