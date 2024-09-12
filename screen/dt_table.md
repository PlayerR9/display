package screen

import (
	dtb "github.com/PlayerR9/display/table"
	gcch "github.com/PlayerR9/go-commons/runes"
	"github.com/gdamore/tcell"
)

// AppendRow appends a new row to the table.
//
// Parameters:
//   - row: The row to append.
//
// Returns:
//   - bool: True if the receiver is not nil, false otherwise.
func (dt *dtb.Table) AppendRow(row []*DtCell) {
	if dt == nil {
		return
	}

	width := len(row)

	dt.cells = append(dt.cells, row)

	if width > dt.width {
		dt.width = width
	}

	dt.height++
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
func NewTableFromBytes(data []byte, fg_style tcell.Style) (*dtb.Table, error) {
	chars, err := gcch.BytesToUtf8(data)
	if err != nil {
		return nil, err
	}

	table, err := NewDtTable(0, 0)
	if err != nil {
		panic(err)
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

/*
// DrawTable implements Drawer interface.
func (dt *dtb.Table) DrawTable(bg_style tcell.Style) (*dtb.Table, error) {
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
*/
