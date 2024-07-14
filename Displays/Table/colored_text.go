package Table

import "github.com/gdamore/tcell"

// ColoredText is a unit that represents a colored text.
type ColoredText struct {
	// lines are the lines of the colored text.
	lines []TableDrawer

	// bgStyle is the background color of the colored text.
	bgStyle tcell.Style
}

// Draw is a method of cdd.TableDrawer that draws the unit to the table at the given x and y
// coordinates.
//
// Parameters:
//   - table: The table to draw the unit to.
//   - x: The x coordinate to draw the unit at.
//   - y: The y coordinate to draw the unit at.
//
// Returns:
//   - error: An error of type *ers.ErrInvalidParameter if the table is nil.
//
// Behaviors:
//   - Any value that would be drawn outside of the table is not drawn.
//   - Assumes that the table is not nil.
func (ct *ColoredText) Draw(table *DrawTable, x, y *int) error {
	height := table.GetHeight()
	X, Y := *x, *y

	var offsetX int
	offsetY := Y

	for _, line := range ct.lines {
		offsetX = X

		err := line.Draw(table, &offsetX, &offsetY)
		if err != nil {
			return err
		}

		offsetY++

		if offsetY >= height {
			break
		}
	}

	*x = offsetX
	*y = offsetY

	return nil
}

// NewColoredText creates a new colored text with the given background color.
//
// Parameters:
//   - bgStyle: The background color of the colored text.
//
// Returns:
//   - *ColoredText: The new colored text.
func NewColoredText(bgStyle tcell.Style) *ColoredText {
	return &ColoredText{
		lines:   make([]TableDrawer, 0),
		bgStyle: bgStyle,
	}
}

// Append appends the given text to the colored text with the given style.
//
// Parameters:
//   - elem: The text to append.
//   - style: The style of the text.
func (ct *ColoredText) Append(elem TableDrawer) {
	if elem == nil {
		return
	}

	ct.lines = append(ct.lines, elem)
}
