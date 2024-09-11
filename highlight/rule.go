package highlight

import "github.com/gdamore/tcell"

// WriteDataFn is a function that writes data.
//
// Parameters:
//   - data: The data to write.
//
// Returns:
//   - []rune: The data as runes.
type WriteDataFn func(data string) []rune

// TokenRule is a token rule.
type TokenRule struct {
	// style is the style of the token.
	style tcell.Style

	// fn is the function that is applied to the token data.
	fn WriteDataFn
}

// NewTokenRule creates a new token rule.
//
// Parameters:
//   - style: The style of the token.
//   - fn: The function that is applied to the token data.
//
// Returns:
//   - *TokenRule: The new token rule. Never returns nil.
func NewTokenRule(style tcell.Style, fn WriteDataFn) *TokenRule {
	if fn == nil {
		fn = func(data string) []rune {
			return []rune(data)
		}
	}

	return &TokenRule{
		style: style,
		fn:    fn,
	}
}
