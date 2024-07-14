package Section

// Render is a type that represents a render of a section.
type Render struct {
	// runeTable is the table of runes of the render.
	runeTable [][]rune
}

// NewRender creates a new render.
//
// Parameters:
//   - runeTable: The table of runes of the render.
//
// Returns:
//   - *Render: The new render.
func NewRender(runeTable [][]rune) *Render {
	return &Render{
		runeTable: runeTable,
	}
}
