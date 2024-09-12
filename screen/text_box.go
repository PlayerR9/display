package screen

import (
	gcch "github.com/PlayerR9/go-commons/runes"
	"github.com/gdamore/tcell"
)

// TextBox is a text box.
type TextBox struct {
	// chars is the text of the text box.
	chars []rune

	// style is the style of the text box.
	style tcell.Style
}

// Draw implements the Drawable interface.
func (tb TextBox) Draw(screen Drawable, x_coord, y_coord *int) error {
	if screen == nil {
		return nil
	}

	x := *x_coord
	y := *y_coord

	bg_style := screen.BgStyle()

	for i := 0; i < len(tb.chars); i++ {
		switch tb.chars[i] {
		case '\n':
			x = 0
			y++
		case '\t':
			for j := 0; j < 3; j++ {
				screen.DrawCell(x, y, ' ', bg_style)
				x++
			}
		default:
			screen.DrawCell(x, y, tb.chars[i], tb.style)
			x++
		}
	}

	*x_coord = x
	*y_coord = y

	return nil
}

// NewTextBox creates a new text box.
//
// Returns:
//   - *TextBox: The new text box. Never returns nil.
func NewTextBox() *TextBox {
	return &TextBox{
		style: tcell.StyleDefault,
	}
}

// ChangeText changes the text of the text box. Does nothing with a nil receiver.
//
// Returns:
//   - error: An error if the text could not be changed.
//
// Errors:
//   - *runes.ErrInvalidUTF8Encoding: If the text is not valid UTF-8.
func (tb *TextBox) ChangeText(text string) error {
	if tb == nil {
		return nil
	}

	chars, err := gcch.StringToUtf8(text)
	if err != nil {
		return err
	}

	tb.chars = chars

	return nil
}

// ChangeStyle changes the style of the text box. Does nothing with a nil receiver.
//
// Returns:
//   - error: An error if the style could not be changed.
func (tb *TextBox) ChangeStyle(style tcell.Style) {
	if tb == nil {
		return
	}

	tb.style = style
}
