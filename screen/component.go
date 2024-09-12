package screen

import (
	"context"
	"sync"

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
	//   - screen: The screen to draw on.
	//   - x_coord: The x coordinate of the top-left corner of the area to draw in. Assumed not nil.
	//   - y_coord: The y coordinate of the top-left corner of the area to draw in. Assumed not nil.
	//
	// Returns:
	//   - error: An error if the table could not be drawn.
	//
	// NOTE: Out of bounds draw or inability to draw should not be considered an error; more specifically,
	// only panic-level of errors should be returned.
	Draw(screen *dtb.Table, x_coord, y_coord *int) error
}

type Display struct {
	buffer *dtb.Table
	frame  *dtb.Table

	mu sync.RWMutex
}

func (d *Display) resize(new_width, new_height int) error {
	d.mu.Lock()
	defer d.mu.Unlock()

	err := d.buffer.ResizeWidth(new_width)
	if err != nil {
		return err
	}

	err = d.buffer.ResizeHeight(new_height)
	if err != nil {
		return err
	}

	return nil
}

func Draw(ctx context.Context, elem Drawer, x, y int) (int, int, error) {
	k := DisplayKey("display")

	display := ctx.Value(k).(*Display)

	display.mu.Lock()
	defer display.mu.Unlock()

	err := elem.Draw(display.buffer, &x, &y)
	if err != nil {
		return x, y, err
	}

	return x, y, nil
}
