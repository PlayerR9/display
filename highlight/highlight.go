package highlight

import (
	pkg "github.com/PlayerR9/display/screen"

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

// DrawTable implements the Highlighter interface.
func (h Highlight[E, T]) DrawTable(bg_style tcell.Style) (*pkg.DtTable, error) {
	dt := pkg.NewDtTable()

	var row []*pkg.DtCell

	last_pos := 0

	for _, tk := range h.tokens[:len(h.tokens)-1] {
		pos := tk.GetPos()

		if last_pos > pos {
			panic("tokens must be ordered by position")
		}

		if last_pos < pos {
			for i := last_pos; i < pos; i++ {
				switch h.data[i] {
				case '\n':
					dt.AppendRow(row)
					row = nil
				case '\t':
					for j := 0; j < 3; j++ {
						row = append(row, pkg.NewDtCell(' ', bg_style))
					}
				default:
					row = append(row, pkg.NewDtCell(' ', bg_style))
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
			row = append(row, pkg.NewDtCell(c, style))
		}

		last_pos = pos + len(chars)
	}

	dt.AppendRow(row)

	return dt, nil
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
