package highlight

import (
	ds "github.com/PlayerR9/display/screen"
	gcers "github.com/PlayerR9/go-commons/errors"
	gda "github.com/PlayerR9/go-debug/assert"
	"github.com/gdamore/tcell"
)

// Highlight is a highlighter.
type Highlight[E interface {
	~int
}, T interface {
	GetPos() int
	GetData() string
	GetType() E
}] struct {
	// tokens is the tokens to highlight.
	tokens []T

	// table is the table of token rules.
	table map[E]*TokenRule

	// data is the data to highlight.
	data []byte
}

// Draw draws the highlighter on the screen.
//
// Parameters:
//   - screen: The screen to draw on.
//   - x_coord: The x coordinate of the top-left corner of the area to draw in.
//   - y_coord: The y coordinate of the top-left corner of the area to draw in.
//
// Returns:
//   - error: An error if a parameter is nil.
func (h Highlight[E, T]) Draw(screen ds.Drawable, x_coord, y_coord *int) error {
	if screen == nil {
		return nil
	}

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

	bg_style := screen.BgStyle()

	last_pos := 0

	for _, tk := range h.tokens[:len(h.tokens)-1] {
		pos := tk.GetPos()

		gda.Assert(pos <= last_pos, "tokens must be ordered by position")

		if last_pos < pos {
			for i := last_pos; i < pos; i++ {
				switch h.data[i] {
				case '\n':
					x = 0
					y++
				case '\t':
					for j := 0; j < 3; j++ {
						screen.DrawCell(x, y, ' ', bg_style)
						x++
					}
				default:
					screen.DrawCell(x, y, ' ', bg_style)
					x++
				}
			}

			last_pos = pos
		}

		var style tcell.Style
		var fn func(data string) []rune

		rule, ok := h.table[tk.GetType()]
		if !ok {
			style = bg_style
			fn = func(data string) []rune {
				return []rune(data)
			}
		} else {
			style = rule.style
			fn = rule.fn
		}

		chars := fn(tk.GetData())

		for _, c := range chars {
			screen.DrawCell(x, y, c, style)
			x++
		}

		last_pos = pos + len(chars)
	}

	if last_pos < len(h.data) {
		for i := last_pos; i < len(h.data); i++ {
			switch h.data[i] {
			case '\n':
				x = 0
				y++
			case '\t':
				for j := 0; j < 3; j++ {
					screen.DrawCell(x, y, ' ', bg_style)
					x++
				}
			default:
				screen.DrawCell(x, y, ' ', bg_style)
				x++
			}
		}
	}

	*x_coord = x
	*y_coord = y

	return nil
}

// NewHighlight creates a new highlighter.
//
// Returns:
//   - *Highlight: The new highlighter. Never returns nil.
func NewHighlight[E interface {
	~int
}, T interface {
	GetPos() int
	GetData() string
	GetType() E
}]() *Highlight[E, T] {
	return &Highlight[E, T]{
		table: make(map[E]*TokenRule),
	}
}

// Register registers a new token rule.
//
// Parameters:
//   - type_: The type of the token.
//   - style: The style of the token.
//   - fn: The function that is applied to the token data.
func (h *Highlight[E, T]) Register(type_ E, style tcell.Style, fn WriteDataFn) {
	if h == nil {
		return
	}

	h.table[type_] = NewTokenRule(style, fn)
}

// SetTokens sets the tokens to highlight.
//
// Parameters:
//   - data: The data to highlight.
//   - tokens: The tokens to highlight.
//
// Returns:
//   - *Highlight: The highlighter. Nil only if receiver is nil.
func (h *Highlight[E, T]) SetTokens(data []byte, tokens []T) *Highlight[E, T] {
	if h == nil {
		return nil
	}

	table := make(map[E]*TokenRule, len(h.table))

	for k, v := range h.table {
		table[k] = v
	}

	return &Highlight[E, T]{
		tokens: tokens,
		table:  table,
		data:   data,
	}
}
