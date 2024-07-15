package Section

import (
	"testing"

	cdd "github.com/PlayerR9/display/Displays/Table"
	"github.com/gdamore/tcell"
)

type MockSection struct {
	text string
}

func (s *MockSection) Runes(width, height int) ([][]rune, error) {
	return nil, nil
}

func TestTitle(t *testing.T) {
	const (
		Title        string = "Test Title"
		ExpectedLine string = " *** Test Title *** "
	)

	table := cdd.NewDrawTable(20, 1)

	// title := NewTitle(Title)

	cell := cdd.NewColoredElement(&MockSection{}, tcell.StyleDefault)

	x, y := 0, 0

	err := cell.Draw(table, &x, &y)
	if err != nil {
		t.Fatalf("Expected no error, but got %s", err.Error())
	}

	lines := table.GetLines()

	if lines[0] != ExpectedLine {
		t.Fatalf("Expected line to be '%s', but got '%s'", ExpectedLine, lines[0])
	}
}

func TestMiddleSplit(t *testing.T) {
	type titleTest struct {
		title         string
		width         int
		height        int
		expectedLines []string
	}

	tests := []titleTest{
		{
			title:  "This is a very long title",
			width:  13,
			height: 5,
			expectedLines: []string{
				"*** This *** ",
				"*** is a *** ",
				"*** very *** ",
				"*** long *** ",
				"*** title ***",
			},
		},
		{
			title:  "Hello world, this is a test",
			width:  19,
			height: 3,
			expectedLines: []string{
				"   *** Hello ***   ",
				"*** world, this ***",
				" *** is a test *** ",
			},
		},
		{
			title:  "Hi You They",
			width:  14,
			height: 2,
			expectedLines: []string{
				"*** Hi You ***",
				" *** They *** ",
			},
		},
	}

	for i, test := range tests {
		// title := NewTitle(test.title)

		table := cdd.NewDrawTable(test.width, test.height)

		cell := cdd.NewColoredElement(&MockSection{}, tcell.StyleDefault)

		x, y := 0, 0

		err := cell.Draw(table, &x, &y)
		if err != nil {
			t.Fatalf("At test %d, expected no error, but got %s", i, err.Error())
		}

		lines := table.GetLines()

		if len(lines) != len(test.expectedLines) {
			t.Errorf("At test %d, expected %d lines, but got %d", i, len(test.expectedLines), len(lines))
		}

		for j := 0; j < len(lines); j++ {
			if lines[j] != test.expectedLines[j] {
				t.Errorf("At test %d, expected line %d to be '%s', but got '%s'", i, j, test.expectedLines[j], lines[j])
			}
		}
	}
}

func TestTitleTruncation(t *testing.T) {
	const (
		Title        string = "Thisisaverylongtitle"
		ExpectedLine string = "*** Th... ***"
	)

	// title := NewTitle(Title)

	table := cdd.NewDrawTable(13, 1)

	cell := cdd.NewColoredElement(&MockSection{}, tcell.StyleDefault)

	x, y := 0, 0

	err := cell.Draw(table, &x, &y)
	if err != nil {
		t.Fatalf("Expected no error, but got %s", err.Error())
	}

	lines := table.GetLines()

	if lines[0] != ExpectedLine {
		t.Errorf("Expected line to be '%s', but got '%s'", ExpectedLine, lines[0])
	}
}
