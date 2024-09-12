package Table

import dtb "github.com/PlayerR9/display/table"

const (
	// EmptyRuneCell is a rune that represents an empty cell.
	EmptyRuneCell rune = '\000'
)

type Colorer interface {
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
	Runes(width, height int) ([][]rune, error)
}

type TableDrawer interface {
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
	Draw(table *dtb.Table, x, y *int) error

	Colorer
}
