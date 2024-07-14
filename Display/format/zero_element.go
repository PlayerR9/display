package format

import (
	ddt "github.com/PlayerR9/MyGoLib/Display/drawtable"
)

// ZeroElement is a type that represents an element that takes up a certain amount of space.
type ZeroElement struct {
	// width is the width of the element.
	width int

	// height is the height of the element.
	height int
}

// Draw implements the drawtable.Displayer interface.
func (ze *ZeroElement) Draw(table *ddt.DrawTable, x, y *int) error {
	*x += ze.width
	*y += ze.height

	return nil
}

// NewZeroElement creates a new ZeroElement with the given width and height.
//
// Parameters:
//   - width: The width of the element.
//   - height: The height of the element.
//
// Returns:
//   - *ZeroElement: The new ZeroElement.
func NewZeroElement(width, height int) *ZeroElement {
	return &ZeroElement{
		width:  width,
		height: height,
	}
}
