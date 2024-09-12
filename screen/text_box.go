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

func (tb TextBox) Draw(vt *VirtualTable) {
	bg_style := vt.BgStyle()

	x, y := 0, 0

	for i := 0; i < len(tb.chars); i++ {
		switch tb.chars[i] {
		case '\n':
			x = 0
			y++
		case '\t':
			for j := 0; j < 3; j++ {
				vt.DrawCellAt(x, y, NewDtCell(' ', bg_style))
				x++
			}
		default:
			vt.DrawCellAt(x, y, NewDtCell(tb.chars[i], tb.style))
			x++
		}
	}
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
