package Section

// Sectioner is an interface that represents a section.
type Sectioner interface {
	// FromTextBlock creates a section from words.
	//
	// Parameters:
	//   - lines: The words to create the section from.
	//
	// Returns:
	//   - error: An error if the section could not be created.
	FromTextBlock(lines [][]string) error

	// ApplyRender applies the render to the section.
	//
	// Parameters:
	//   - width: The width of the render.
	//   - height: The height of the render.
	//
	// Returns:
	//   - []*Render: The render of the section.
	//   - error: An error if the render could not be applied.
	ApplyRender(width, height int) ([]*Render, error)

	// GetRawContent returns the raw content of the section.
	//
	// Returns:
	//   - [][]string: The raw content of the section.
	GetRawContent() [][]string
}
