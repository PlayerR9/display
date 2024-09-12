package table

import (
	"iter"
	"strings"
	"sync"

	gcers "github.com/PlayerR9/go-commons/errors"
	"github.com/PlayerR9/go-commons/ints"
	"github.com/gdamore/tcell"
)

// Table represents a table of cells that can be drawn to the screen.
type Table struct {
	// table is the table of cells.
	table [][]*Cell

	// width is the width of the table.
	width int

	// height is the height of the table.
	height int

	// mu is the table mutex.
	mu sync.RWMutex
}

// NewTable creates a new table of type DrawCell with the given width and height.
// Negative parameters are treated as absolute values.
//
// Parameters:
//   - width: The width of the table.
//   - height: The height of the table.
//
// Returns:
//   - *Table: The new table.
//   - error: An error if the table could not be created.
//
// Errors:
//   - *gcers.ErrInvalidParameter: If the width or height is less than 0.
func NewTable(width, height int) (*Table, error) {
	if width < 0 {
		return nil, gcers.NewErrInvalidParameter("width", gcers.NewErrGTE(0))
	} else if height < 0 {
		return nil, gcers.NewErrInvalidParameter("height", gcers.NewErrGTE(0))
	}

	table := make([][]*Cell, 0, height)
	for i := 0; i < height; i++ {
		table = append(table, make([]*Cell, width))
	}

	return &Table{
		table:  table,
		width:  width,
		height: height,
	}, nil
}

// Cell returns an iterator that is a pull-model iterator that scans the table row by
// row as it was an array of elements of type DrawCell.
//
// Example:
//
//	[ a b c ]
//	[ d e f ]
//
//	Cell() -> [ a ] -> [ b ] -> [ c ] -> [ d ] -> [ e ] -> [ f ]
func (t *Table) Cell() iter.Seq[*Cell] {
	if t == nil {
		return func(yield func(*Cell) bool) {}
	}

	t.mu.RLock()
	defer t.mu.RUnlock()

	fn := func(yield func(*Cell) bool) {
		for i := 0; i < t.height; i++ {
			for j := 0; j < t.width; j++ {
				if !yield(t.table[i][j]) {
					return
				}
			}
		}
	}

	return fn
}

// Row returns an iterator that is a pull-model iterator that scans the table row by
// row as it was an array of elements of type DrawCell.
//
// Example:
//
//	[ a b c ]
//	[ d e f ]
//
//	Row(0) -> [ a b c ]
//	Row(1) -> [ d e f ]
func (t *Table) Row() iter.Seq[[]*Cell] {
	if t == nil {
		return func(yield func([]*Cell) bool) {}
	}

	t.mu.RLock()
	defer t.mu.RUnlock()

	fn := func(yield func([]*Cell) bool) {
		for i := 0; i < t.height; i++ {
			if !yield(t.table[i]) {
				return
			}
		}
	}

	return fn
}

// Cleanup is a method that cleans up the table.
//
// It sets all cells in the table to the zero value of type int.
func (t *Table) Cleanup() {
	if t == nil {
		return
	}

	t.mu.Lock()
	defer t.mu.Unlock()

	for i := 0; i < t.height; i++ {
		for j := 0; j < t.width; j++ {
			t.table[i][j] = nil
		}
	}
}

// Width returns the width of the table.
//
// Returns:
//   - int: The width of the table. Never negative.
func (t *Table) Width() int {
	if t == nil {
		return 0
	}

	t.mu.RLock()
	defer t.mu.RUnlock()

	return t.width
}

// Height returns the height of the table.
//
// Returns:
//   - int: The height of the table. Never negative.
func (t *Table) Height() int {
	if t == nil {
		return 0
	}

	t.mu.RLock()
	defer t.mu.RUnlock()

	return t.height
}

// WriteAt writes a cell to the table at the given coordinates. However, out-of-bounds
// coordinates do nothing.
//
// Parameters:
//   - x: The x-coordinate of the cell.
//   - y: The y-coordinate of the cell.
//   - cell: The cell to write to the table.
func (t *Table) WriteAt(x, y int, cell *Cell) {
	if t == nil {
		return
	}

	t.mu.Lock()
	defer t.mu.Unlock()

	if x < 0 || x >= t.width || y < 0 || y >= t.height {
		return
	}

	t.table[y][x] = cell
}

// CellAt returns the cell at the given coordinates in the table. However, out-of-bounds
// coordinates return nil.
//
// Parameters:
//   - x: The x-coordinate of the cell.
//   - y: The y-coordinate of the cell.
//
// Returns:
//   - *DrawCell: The cell at the given coordinates.
func (t *Table) CellAt(x, y int) *Cell {
	if t == nil {
		return nil
	}

	t.mu.RLock()
	defer t.mu.RUnlock()

	if x < 0 || x >= t.width || y < 0 || y >= t.height {
		return nil
	}

	return t.table[y][x]
}

// WriteVerticalSequence is a function that writes the specified values to the table
// starting from the specified coordinates (top = 0, 0) and continuing down the
// table in the vertical direction until either the sequence is exhausted or
// the end of the table is reached; at which point any remaining values in the
// sequence are ignored.
//
// Due to implementation details, any value that would be written outside are ignored.
// As such, if x is out-of-bounds, the function does nothing and, if y is out-of-bounds,
// only out-of-bounds values are not written.
//
// Parameters:
//   - x: The x-coordinate of the starting cell. (Never changes)
//   - y: The y-coordinate of the starting cell.
//   - sequence: The sequence of cells to write to the table.
//
// At the end of the function, the y coordinate points to the cell right below the
// last cell in the sequence that was written.
//
// Example:
//
//	// [ a b c ]
//	// [ d e f ]
//	//
//	// seq := [ g h i ], x = 0, y = -1
//
//	WriteVerticalSequence(x, y, seq)
//
//	// [ h b c ]
//	// [ i e f ]
//	//
//	// x = 0, y = 2
//
// As you can see, the 'g' value was ignored as it would be out-of-bounds.
// Finally, if either x or y is nil, the function does nothing.
func (t *Table) WriteVerticalSequence(x, y *int, sequence []*Cell) {
	if t == nil || x == nil || y == nil || len(sequence) == 0 {
		return
	}

	t.mu.Lock()
	defer t.mu.Unlock()

	actualX, actualY := *x, *y

	if actualX < 0 || actualX >= t.width || actualY >= t.height {
		return
	}

	if actualY < 0 {
		sequence = sequence[-actualY:]

		*y = 0
	} else if actualY+len(sequence) > t.height {
		sequence = sequence[:t.height-actualY]
	}

	for i, cell := range sequence {
		t.table[actualY+i][actualX] = cell
	}

	*y += len(sequence)
}

// WriteHorizontalSequence is the equivalent of WriteVerticalSequence but for horizontal
// sequences.
//
// See WriteVerticalSequence for more information.
//
// Parameters:
//   - x: The x-coordinate of the starting cell.
//   - y: The y-coordinate of the starting cell.
//   - sequence: The sequence of cells to write to the table.
func (t *Table) WriteHorizontalSequence(x, y *int, sequence []*Cell) {
	if t == nil || x == nil || y == nil || len(sequence) == 0 {
		return
	}

	t.mu.Lock()
	defer t.mu.Unlock()

	actualX, actualY := *x, *y

	if actualY < 0 || actualY >= t.height || actualX >= t.width {
		return
	}

	if actualX < 0 {
		sequence = sequence[-actualX:]

		*x = 0
	} else if actualX+len(sequence) > t.width {
		sequence = sequence[:t.width-actualX]
	}

	copy(t.table[actualY][actualX:], sequence)

	*x = actualX + len(sequence)
}

// FullTable returns the full table as a 2D slice of elements of type DrawCell.
//
// Returns:
//   - [][]*DrawCell: The full table.
func (t *Table) FullTable() [][]*Cell {
	if t == nil {
		return nil
	}

	t.mu.RLock()
	defer t.mu.RUnlock()

	dt := make([][]*Cell, 0, t.height)

	for y := 0; y < t.height; y++ {
		row := make([]*Cell, 0, t.width)

		for x := 0; x < t.width; x++ {
			row = append(row, t.table[y][x])
		}

		dt = append(dt, row)
	}

	return dt
}

// IsXInBounds checks if the given x-coordinate is within the bounds of the table.
//
// Parameters:
//   - x: The x-coordinate to check.
//
// Returns:
//   - error: An error of type *ints.ErrOutOfBounds if the x-coordinate is out of bounds.
func (t *Table) IsXInBounds(x int) error {
	if t == nil {
		return nil
	}

	t.mu.RLock()
	defer t.mu.RUnlock()

	if x < 0 || x >= t.width {
		return ints.NewErrOutOfBounds(x, 0, t.width)
	} else {
		return nil
	}
}

// IsYInBounds checks if the given y-coordinate is within the bounds of the table.
//
// Parameters:
//   - y: The y-coordinate to check.
//
// Returns:
//   - error: An error of type *ints.ErrOutOfBounds if the y-coordinate is out of bounds.
func (t *Table) IsYInBounds(y int) error {
	if t == nil {
		return nil
	}

	t.mu.RLock()
	defer t.mu.RUnlock()

	if y < 0 || y >= t.height {
		return ints.NewErrOutOfBounds(y, 0, t.height)
	} else {
		return nil
	}
}

// WriteTableAt is a convenience function that copies the values from the given
// table to the table starting at the given coordinates in a more efficient way
// than using any other methods.
//
// While it acts in the same way as both WriteVerticalSequence and WriteHorizontalSequence
// combined, it is more efficient than calling those two functions separately.
//
// See WriteVerticalSequence for more information.
//
// Parameters:
//   - table: The table to write to the table.
//   - x: The x-coordinate to write the table at.
//   - y: The y-coordinate to write the table at.
//
// If the table is nil, x or y are nil, nothing happens.
func (t *Table) WriteTableAt(table *Table, x, y *int) {
	if t == nil || table == nil || x == nil || y == nil {
		return
	}

	t.mu.Lock()
	defer t.mu.Unlock()

	offsetX, offsetY := 0, 0
	X, Y := *x, *y

	for offsetY < table.height && Y+offsetY < t.height {
		offsetX = 0

		for offsetX < table.width && X+offsetX < t.width {
			t.table[Y+offsetY][X+offsetX] = table.table[offsetY][offsetX]
			offsetX++
		}

		offsetY++
	}

	*x += offsetX
	*y += offsetY
}

// ResizeWidth resizes the table to the given width.
//
// Parameters:
//   - new_width: The new width of the table.
//
// Returns:
//   - error: An error if the table could not be resized.
//
// Errors:
//   - *gcers.ErrInvalidParameter: If the new width is less than 0.
//   - gcers.NilReceiver: If the table is nil.
func (t *Table) ResizeWidth(new_width int) error {
	if t == nil {
		return gcers.NilReceiver
	} else if new_width < 0 {
		return gcers.NewErrInvalidParameter("new_width", gcers.NewErrGTE(0))
	}

	t.mu.Lock()
	defer t.mu.Unlock()

	if new_width == t.width {
		return nil
	} else if new_width < t.width {
		for i := 0; i < t.height; i++ {
			t.table[i] = t.table[i][:new_width]
		}
	} else {
		for i := 0; i < t.height; i++ {
			t.table[i] = append(t.table[i], make([]*Cell, new_width-t.width)...)
		}
	}

	t.width = new_width

	return nil
}

// ResizeHeight resizes the table to the given height.
//
// Parameters:
//   - new_height: The new height of the table.
//
// Returns:
//   - error: An error if the table could not be resized.
//
// Errors:
//   - *gcers.ErrInvalidParameter: If the new height is less than 0.
//   - gcers.NilReceiver: If the table is nil.
func (t *Table) ResizeHeight(new_height int) error {
	if t == nil {
		return gcers.NilReceiver
	} else if new_height < 0 {
		return gcers.NewErrInvalidParameter("new_height", gcers.NewErrGTE(0))
	}

	t.mu.Lock()
	defer t.mu.Unlock()

	if new_height == t.height {
		return nil
	} else if new_height < t.height {
		t.table = t.table[:new_height]
	} else {
		t.table = append(t.table, make([][]*Cell, new_height-t.height)...)
	}

	t.height = new_height

	return nil
}

// GetLines returns each line of the drawTable as a string.
//
// Parameters:
//   - table: The drawTable.
//
// Returns:
//   - []string: The lines of the drawTable.
//
// Behaviors:
//   - Any nil cells in the drawTable are represented by a space character.
//   - The last line may not be full width.
func (t *Table) GetLines() []string {
	if t == nil {
		return nil
	}

	var lines []string
	var builder strings.Builder

	for i := 0; i < t.height; i++ {
		for j := 0; j < t.width; j++ {
			if t.table[i][j] == nil {
				builder.WriteRune(' ')
			} else {
				builder.WriteRune(t.table[i][j].Char)
			}
		}

		lines = append(lines, builder.String())
		builder.Reset()
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
func (t *Table) WriteLineAt(x, y *int, line string, style tcell.Style, isHorizontal bool) {
	if t == nil {
		return
	}

	runes := []rune(line)

	sequence := make([]*Cell, 0, len(runes))

	for _, r := range runes {
		sequence = append(sequence, NewCell(r, style))
	}

	if isHorizontal {
		t.WriteHorizontalSequence(x, y, sequence)
	} else {
		t.WriteVerticalSequence(x, y, sequence)
	}
}
