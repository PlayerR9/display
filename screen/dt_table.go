package screen

import (
	"fmt"

	"github.com/gdamore/tcell"

	gcers "github.com/PlayerR9/go-commons/errors"
	gcch "github.com/PlayerR9/go-commons/runes"
)

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
	DrawTable(bg_style tcell.Style) (*DtTable, error)
}

// DtCell is a table cell.
type DtCell struct {
	// char is the character of the cell.
	char rune

	// style is the style of the cell.
	style tcell.Style
}

// NewDtCell creates a new table cell.
//
// Parameters:
//   - char: The character of the cell.
//   - style: The style of the cell.
//
// Returns:
//   - *DtCell: The new table cell. Never returns nil.
func NewDtCell(char rune, style tcell.Style) *DtCell {
	return &DtCell{
		char:  char,
		style: style,
	}
}

// DtTable is a table.
type DtTable struct {
	// cells is the table cells.
	cells [][]*DtCell

	// width is the width of the table.
	width int

	// height is the height of the table.
	height int
}

// NewDtTable creates a new table.
//
// Parameters:
//   - width: The width of the table.
//   - height: The height of the table.
//
// Returns:
//   - *DtTable: The new table.
//   - error: An error if the table could not be created.
func NewDtTable(width, height int) (*DtTable, error) {
	if width < 0 {
		return nil, gcers.NewErrInvalidParameter("width", gcers.NewErrGTE(0))
	} else if height < 0 {
		return nil, gcers.NewErrInvalidParameter("height", gcers.NewErrGTE(0))
	}

	cells := make([][]*DtCell, 0, height)

	for i := 0; i < height; i++ {
		cells = append(cells, make([]*DtCell, width))
	}

	return &DtTable{
		cells:  cells,
		width:  width,
		height: height,
	}, nil
}

// ResizeHeight resizes the height of the table.
//
// Parameters:
//   - height: The new height of the table.
//
// Returns:
//   - error: An error if the table could not be resized.
func (dt *DtTable) ResizeHeight(height int) error {
	if dt == nil {
		return gcers.NilReceiver
	} else if height < 0 {
		return gcers.NewErrInvalidParameter("height", gcers.NewErrGTE(0))
	}

	dt.height = height

	return nil
}

// ResizeWidth resizes the width of the table.
//
// Parameters:
//   - width: The new width of the table.
//
// Returns:
//   - error: An error if the table could not be resized.
func (dt *DtTable) ResizeWidth(width int) error {
	if dt == nil {
		return gcers.NilReceiver
	} else if width < 0 {
		return gcers.NewErrInvalidParameter("width", gcers.NewErrGTE(0))
	}

	dt.width = width

	return nil
}

// DrawTable implements Drawer interface.
func (dt DtTable) DrawTable(bg_style tcell.Style) (*DtTable, error) {
	cells := make([][]*DtCell, 0, dt.height)

	for i := 0; i < dt.height; i++ {
		cells = append(cells, make([]*DtCell, dt.width))
	}

	for i, row := range dt.cells {
		for j, cell := range row {
			if cell == nil {
				cells[i][j] = NewDtCell(' ', bg_style)
			} else {
				cells[i][j] = cell
			}
		}
	}

	for i, row := range cells {
		if len(row) == dt.width {
			continue
		}

		for j := len(row); j < dt.width; j++ {
			cells[i] = append(cells[i], NewDtCell(' ', bg_style))
		}
	}

	return &DtTable{
		cells:  cells,
		width:  dt.width,
		height: dt.height,
	}, nil
}

// Width returns the width of the table.
//
// Returns:
//   - int: The width of the table.
func (dt DtTable) Width() int {
	return dt.width
}

// Height returns the height of the table.
//
// Returns:
//   - int: The height of the table.
func (dt DtTable) Height() int {
	return dt.height
}

// AppendRow appends a new row to the table.
//
// Parameters:
//   - row: The row to append.
//
// Returns:
//   - bool: True if the receiver is not nil, false otherwise.
func (dt *DtTable) AppendRow(row []*DtCell) bool {
	if dt == nil {
		return false
	}

	width := len(row)

	dt.cells = append(dt.cells, row)

	if width > dt.width {
		dt.width = width
	}

	dt.height++

	return true
}

// NewTableFromBytes creates a new table from a byte slice.
//
// Parameters:
//   - data: The byte slice to create the table from.
//   - fg_style: The foreground style of the table.
//
// Returns:
//   - *DtTable: The new table.
//   - error: An error if the table could not be created.
func NewTableFromBytes(data []byte, fg_style tcell.Style) (*DtTable, error) {
	chars, err := gcch.BytesToUtf8(data)
	if err != nil {
		return nil, err
	}

	table, err := NewDtTable(0, 0)
	if err != nil {
		panic(fmt.Sprintf("could not create table: %v", err.Error()))
	}

	var row []*DtCell

	for _, c := range chars {
		if c == '\n' {
			table.AppendRow(row)
			row = nil
		} else {
			row = append(row, NewDtCell(c, fg_style))
		}
	}

	if len(row) > 0 {
		table.AppendRow(row)
	}

	return table, nil
}

func (dt *DtTable) DrawCellAt(x, y int, cell *DtCell) {
	if dt == nil || x < 0 || x >= dt.width || y < 0 || y >= dt.height {
		return
	}

	dt.cells[y][x] = cell
}
