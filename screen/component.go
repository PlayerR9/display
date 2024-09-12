package screen

import (
	dtb "github.com/PlayerR9/display/table"
	"github.com/gdamore/tcell"
)

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

// Drawer is a table drawer.
type Drawer interface {
	// DrawTable draws the table.
	//
	// Parameters:
	//   - bg_style: The background style of the table.
	//
	// Returns:
	//   - *DtTable: The table that was drawn.
	//   - error: An error if the table could not be drawn.
	DrawTable(bg_style tcell.Style) (*dtb.Table, error)
}
