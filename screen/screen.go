package screen

import (
	"context"

	dtb "github.com/PlayerR9/display/table"
	gcers "github.com/PlayerR9/go-commons/errors"
	gda "github.com/PlayerR9/go-debug/assert"
	"github.com/gdamore/tcell"
)

type DisplayKey string

// Screen is a screen.
type Screen struct {
	// bg_style is the background style.
	bg_style tcell.Style

	// screen is the tcell screen.
	screen tcell.Screen

	// event_ch is the event channel.
	event_ch chan tcell.Event

	// key_ch is the key channel.
	key_ch chan *tcell.EventKey

	// dt is the draw table.
	dt *Display
}

// NewScreen creates a new screen.
//
// Parameters:
//   - bg_style: The background style of the screen.
//
// Returns:
//   - *Screen: The new screen.
//   - error: The error if any.
func NewScreen(bg_style tcell.Style) (*Screen, error) {
	screen, err := tcell.NewScreen()
	if err != nil {
		return nil, err
	}

	buffer, err := dtb.NewTable(80, 25)
	gda.AssertErr(err, "table.NewTable(80, 25)")

	frame, err := dtb.NewTable(80, 25)
	gda.AssertErr(err, "table.NewTable(80, 25)")

	dt := &Display{
		buffer: buffer,
		frame:  frame,
	}

	return &Screen{
		bg_style: bg_style,
		screen:   screen,
		event_ch: make(chan tcell.Event, 1),
		key_ch:   make(chan *tcell.EventKey),
		dt:       dt,
	}, nil
}

// event_listener is a helper function that listens for events.
func (s *Screen) event_listener() {
	for {
		ev := s.screen.PollEvent()
		if ev == nil {
			break
		}

		s.event_ch <- ev
	}
}

// Start starts the screen.
//
// Returns:
//   - error: The error if any.
func (s *Screen) Start() (context.Context, error) {
	k := DisplayKey("display")

	ctx := context.WithValue(context.Background(), k, dt)

	if s == nil {
		return gcers.NilReceiver
	}

	err := s.screen.Init()
	if err != nil {
		return err
	}

	s.screen.SetStyle(s.bg_style)

	s.screen.EnableMouse()

	s.screen.Clear()

	width, height := s.screen.Size()

	err = s.dt.resize(width, height)
	if err != nil {
		return err
	}

	go s.event_listener()

	go s.run()

	return nil
}

// Close closes the screen.
func (s *Screen) Close() {
	if s == nil {
		return
	}

	s.screen.Fini()

	if s.event_ch != nil {
		close(s.event_ch)
		s.event_ch = nil
	}

	if s.key_ch != nil {
		close(s.key_ch)
		s.key_ch = nil
	}
}

// run runs the screen.
func (s *Screen) run() {
	for ev := range s.event_ch {
		switch ev := ev.(type) {
		case *tcell.EventKey:
			select {
			case s.key_ch <- ev:
			}
		case *tcell.EventResize:
			// s.screen.Sync()

			width, height := ev.Size()

			err := s.dt.resize(width, height)
			gda.AssertErr(err, "s.dt.ResizeWidth(%d, %d)", width, height)

			s.show_display()
			// case *tcell.EventMouse:
			// 	button := ev.Buttons()

			// 	if button == tcell.WheelUp {
			// 		val := s.pos_y.Get()
			// 		s.pos_y.Set(val - 1)
			// 	} else if button == tcell.WheelDown {
			// 		val := s.pos_y.Get()
			// 		s.pos_y.Set(val + 1)
			// 	}

			// 	s.show_display()
		}
	}
}

// show_display is a helper function that shows the display.
func (s *Screen) show_display() {
	s.screen.Clear()

	y := 0

	for row := range s.dt.Row() {
		for x := 0; x < len(row); x++ {
			cell := row[x]

			if cell == nil {
				s.screen.SetContent(x, y, ' ', nil, s.bg_style)
			} else {
				s.screen.SetContent(x, y, cell.Char, nil, cell.Style)
			}
		}
	}

	s.screen.Show()
}

// SetCell is a helper function that sets a cell.
//
// Parameters:
//   - x: The x position of the cell.
//   - y: The y position of the cell.
//   - c: The character of the cell.
//   - style: The style of the cell.
func (s *Screen) DrawCell(x, y int, c rune, style tcell.Style) {
	if s == nil {
		return
	}

	s.screen.SetContent(x, y, c, nil, style)
}

// display_label is a helper function that displays a label.
//
// Parameters:
//   - x: The x position of the label.
//   - y: The y position of the label.
//   - style: The style of the label.
//   - text: The text of the label.
func (s *Screen) display_label(x, y int, style tcell.Style, text string) {
	for _, c := range []rune(text) {
		s.screen.SetContent(x, y, c, nil, style)

		x++
	}
}

// Show shows the screen.
//
// Parameters:
//   - elem: The element to draw. Can be nil.
//   - x: The x position of the element.
//   - y: The y position of the element.
//
// Returns:
//   - int: The new x position.
//   - int: The new y position.
//   - error: An error if the screen could not be shown.
func (s *Screen) Show(elem Drawer, x, y int) (int, int, error) {
	if s == nil {
		return x, y, nil
	}

	s.screen.Clear()
	defer s.screen.Show()

	var err error

	if elem != nil {
		err = elem.Draw(s, &x, &y)
	}

	return x, y, err
}

// Clear clears the screen.
func (s *Screen) Clear() {
	if s == nil {
		return
	}

	s.screen.Clear()
	s.screen.Show()
}

// ListenForKey listens for a key press event on the screen.
//
// Parameters:
//   - None.
//
// Returns:
//   - *tcell.EventKey: The key press event.
//   - bool: Whether the channel is still open.
func (s *Screen) ListenForKey() (*tcell.EventKey, bool) {
	if s == nil || s.key_ch == nil {
		return nil, false
	}

	ev, ok := <-s.key_ch
	if !ok {
		return nil, false
	}

	return ev, true
}

// Height returns the height of the screen.
//
// Returns:
//   - int: The height of the screen.
func (s *Screen) Height() int {
	if s == nil || s.dt == nil {
		return 0
	}

	return s.dt.Height()
}

// Width returns the width of the screen.
//
// Returns:
//   - int: The width of the screen.
func (s *Screen) Width() int {
	if s == nil {
		return 0
	}

	return s.dt.Width()
}

// BgStyle returns the background style.
//
// Returns:
//   - tcell.Style: The background style.
func (d *Screen) BgStyle() tcell.Style {
	return d.bg_style
}

func (s *Screen) Table() *Display {
	return s.dt
}
