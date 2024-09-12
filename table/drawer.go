package table

type Displayer interface {
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
	Draw(table *Table, x, y *int) error
}
