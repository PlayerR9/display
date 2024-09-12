package screen

import (
	gcers "github.com/PlayerR9/go-commons/errors"
	gcch "github.com/PlayerR9/go-commons/runes"
	"github.com/gdamore/tcell"
)

type TextBox struct {
	chars []rune
	style tcell.Style
}

func (tb TextBox) Draw(screen Drawable, x_coord, y_coord *int) error {
	var x int

	if x_coord == nil {
		return gcers.NewErrNilParameter("x_coord")
	} else {
		x = *x_coord
	}

	var y int

	if y_coord == nil {
		return gcers.NewErrNilParameter("y_coord")
	} else {
		y = *y_coord
	}

	if screen == nil {
		return nil
	}

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

func NewTextBox() *TextBox {
	return &TextBox{
		style: tcell.StyleDefault,
	}
}

func (tb *TextBox) ChangeText(text string) error {
	if tb == nil {
		return gcers.NilReceiver
	}

	chars, err := gcch.StringToUtf8(text)
	if err != nil {
		return err
	}

	tb.chars = chars

	return nil
}

func (tb *TextBox) ChangeStyle(style tcell.Style) error {
	if tb == nil {
		return gcers.NilReceiver
	}

	tb.style = style

	return nil
}
