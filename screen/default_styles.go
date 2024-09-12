package screen

import "github.com/gdamore/tcell"

var (
	// LightModeStyle is a default style for light mode.
	LightModeStyle tcell.Style

	// DarkModeStyle is a default style for dark mode.
	DarkModeStyle tcell.Style

	// GreenStyle is a default style for green text.
	GreenStyle tcell.Style
)

func init() {
	LightModeStyle = tcell.StyleDefault.Background(tcell.ColorGhostWhite).Foreground(tcell.ColorBlack)

	DarkModeStyle = tcell.StyleDefault.Background(tcell.ColorDarkGray).Foreground(tcell.ColorWhiteSmoke)

	GreenStyle = tcell.StyleDefault.Background(tcell.ColorSpringGreen).Foreground(tcell.ColorDarkGrey)
}
