package Table

import (
	cdt "github.com/PlayerR9/MyGoLib/CustomData/Table"
	ddt "github.com/PlayerR9/MyGoLib/Display/drawtable"
	"github.com/gdamore/tcell"
)

// ColoredElement represents an element that can be colored.
type ColoredElement[T Colorer] struct {
	// elem is the element of the color.
	elem T

	// style is the style of the color.
	style tcell.Style
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
func (ce *ColoredElement[T]) Draw(table *DrawTable, x, y *int) error {
	width, height := table.GetWidth(), table.GetHeight()

	runeTable, err := ce.elem.Runes(width, height)
	if err != nil {
		return err
	}

	// Fix the boundaries of the rune table
	runeTable = cdt.FixBoundaries(width, height, runeTable, x, y)
	if len(runeTable) == 0 {
		return nil
	}

	var offsetX int

	for i, row := range runeTable {
		offsetX = *x

		if len(row) == 0 {
			continue
		}

		sequence := make([]*ddt.ColoredUnit, 0, len(row))

		for _, r := range row {
			if r == EmptyRuneCell {
				sequence = append(sequence, nil)
			} else {
				sequence = append(sequence, ddt.NewColoredUnit(r, ce.style))
			}
		}

		offsetY := *y + i

		table.WriteHorizontalSequence(&offsetX, &offsetY, sequence)
	}

	*x = offsetX
	*y += len(runeTable)

	return nil
}

// NewColoredElement creates a new ColoredElement with the given element and style.
//
// Parameters:
//   - elem: The element of the color.
//   - style: The style of the color.
//
// Returns:
//   - *ColoredElement: The new ColoredElement.
func NewColoredElement[T Colorer](elem T, style tcell.Style) *ColoredElement[T] {
	return &ColoredElement[T]{
		elem:  elem,
		style: style,
	}
}

// Apply applies the color to the element.
//
// Parameters:
//   - width: The width of the element.
//   - height: The height of the element.
//
// Returns:
//   - [][]*ColoredUnit: The colored element.
//   - error: An error if the element could not be colored.
//
// Behaviors:
//   - Always assume that the width and height are greater than 0. No need to check for
//     this.
//   - Errors are only for critical issues, such as the element not being able to be
//     colored. However, out of bounds or other issues should not error. Instead, the
//     element should be colored as much as possible before unable to be colored.
func (ce *ColoredElement[T]) Apply(width, height int) ([][]*ddt.ColoredUnit, error) {
	runeTable, err := ce.elem.Runes(width, height)
	if err != nil {
		return nil, err
	}

	colorTable := make([][]*ddt.ColoredUnit, len(runeTable))

	for _, row := range runeTable {
		var colorRow []*ddt.ColoredUnit

		for _, r := range row {
			colorRow = append(colorRow, ddt.NewColoredUnit(r, ce.style))
		}

		colorTable = append(colorTable, colorRow)
	}

	return colorTable, nil
}
