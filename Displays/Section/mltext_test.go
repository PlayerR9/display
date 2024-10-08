package Section

import (
	"testing"

	ddt "github.com/PlayerR9/display/Displays/Table"
	dtb "github.com/PlayerR9/display/table"
	"github.com/gdamore/tcell"
)

func TestWriteLines_ShortLines(t *testing.T) {
	mlt := new(MultilineText)

	err := mlt.FromTextBlock([][]string{{"Hello", "World"}})
	if err != nil {
		t.Fatalf("Expected no error, but got %s", err.Error())
	}

	table, err := dtb.NewTable(18, 2)
	if err != nil {
		t.Fatalf("Expected no error, but got %s", err.Error())
	}

	cell := ddt.NewColoredElement(&MockSection{}, tcell.StyleDefault)

	x, y := 0, 0

	err = cell.Draw(table, &x, &y)
	if err != nil {
		t.Fatalf("Expected no error, but got %s", err.Error())
	}

	lines := table.GetLines()

	if lines[0] != "Hello World       " {
		t.Errorf("Expected first line to be 'Hello World       ', but got '%s'", lines[0])
	}
}

func TestWriteLines_LongLine(t *testing.T) {
	mlt := new(MultilineText)

	err := mlt.FromTextBlock([][]string{{
		"This is really a very long line that should be truncated and end with an ellipsis",
	}})
	if err != nil {
		t.Fatalf("Expected no error, but got %s", err.Error())
	}

	table, err := dtb.NewTable(18, 1)
	if err != nil {
		t.Fatalf("Expected no error, but got %s", err.Error())
	}

	cell := ddt.NewColoredElement(&MockSection{}, tcell.StyleDefault)

	x, y := 0, 0

	err = cell.Draw(table, &x, &y)
	if err != nil {
		t.Fatalf("Expected no error, but got %s", err.Error())
	}

	lines := table.GetLines()

	if lines[0] != "This is really... " {
		t.Fatalf("Expected first line to be 'This is really... ', but got '%s'", lines[0])
	}
}
