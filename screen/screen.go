package screen

import (
	"sync"

	gcers "github.com/PlayerR9/go-commons/errors"
	rws "github.com/PlayerR9/safe/rw_safe"
	"github.com/gdamore/tcell"
)

var (
	// BgStyle is the background style.
	BgStyle tcell.Style
)

func init() {
	BgStyle = tcell.StyleDefault.Background(tcell.ColorGhostWhite).Foreground(tcell.ColorBlack)
}

// Screen is a screen.
type Screen struct {
	// screen is the tcell screen.
	screen tcell.Screen

	// event_ch is the event channel.
	event_ch chan tcell.Event

	// width is the width of the screen.
	width int

	// height is the height of the screen.
	height int

	// to_display is the table to display.
	to_display *rws.Safe[*DtTable]

	// pos is the position of the cursor.
	pos *rws.Safe[int]

	// wg is the wait group.
	wg sync.WaitGroup
}

// NewScreen creates a new screen.
//
// Returns:
//   - *Screen: The new screen.
//   - error: The error if any.
func NewScreen() (*Screen, error) {
	screen, err := tcell.NewScreen()
	if err != nil {
		return nil, err
	}

	return &Screen{
		screen:     screen,
		event_ch:   make(chan tcell.Event, 1),
		width:      80,
		height:     25,
		pos:        rws.NewSafe(0),
		to_display: rws.NewSafe[*DtTable](nil),
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
func (s *Screen) Start() error {
	if s == nil {
		return gcers.NilReceiver
	}

	err := s.screen.Init()
	if err != nil {
		return err
	}

	s.screen.SetStyle(BgStyle)

	s.screen.EnableMouse()

	s.screen.Clear()

	s.width, s.height = s.screen.Size()

	go s.event_listener()

	s.wg.Add(1)

	go s.run()

	return nil
}

// Close closes the screen.
func (s *Screen) Close() {
	if s == nil {
		return
	}

	s.wg.Wait()

	s.screen.Fini()

	close(s.event_ch)
}

// handle_event handles an event.
//
// Parameters:
//   - ev: The event to handle.
//
// Returns:
//   - bool: True if the screen should be closed, false otherwise.
func (s *Screen) handle_event(ev tcell.Event) bool {
	switch ev := ev.(type) {
	case *tcell.EventKey:
		switch ev.Key() {
		case tcell.KeyEnter:
			return true
		case tcell.KeyUp:
			val := s.pos.Get()
			s.pos.Set(val - 1)

			s.show_display()
		case tcell.KeyDown:
			val := s.pos.Get()
			s.pos.Set(val + 1)

			s.show_display()
		}
	case *tcell.EventResize:
		// s.screen.Sync()

		s.width, s.height = ev.Size()

		s.show_display()
	case *tcell.EventMouse:
		button := ev.Buttons()

		if button == tcell.WheelUp {
			val := s.pos.Get()
			s.pos.Set(val - 1)
		} else if button == tcell.WheelDown {
			val := s.pos.Get()
			s.pos.Set(val + 1)
		}

		s.show_display()
	}

	return false
}

// run runs the screen.
func (s *Screen) run() {
	defer s.wg.Done()

	for {
		select {
		case ev := <-s.event_ch:
			should_close := s.handle_event(ev)
			if should_close {
				return
			}
		}
	}
}

// show_display is a helper function that shows the display.
func (s *Screen) show_display() {
	s.screen.Clear()

	table := s.to_display.Get()

	if table == nil {
		s.screen.Show()

		return
	}

	y := 0
	x := 0

	pos := s.pos.Get()

	if pos >= table.Height() {
		s.screen.Show()

		return
	} else if pos < 0 {
		y -= pos
		pos = 0
	}

	underlying_table := table.cells[pos:]

	for _, row := range underlying_table {
		if y >= s.height-3 {
			break
		}

		j := 0
		for _, c := range row {
			if j >= s.width {
				j = 0
				y++
				x = 0 // Change this to indent the next line
			}

			if c == nil {
				c = NewDtCell(' ', BgStyle)
			}

			s.screen.SetContent(x, y, c.char, nil, c.style)

			j++
			x++
		}

		y++
		x = 0
	}

	style := tcell.StyleDefault.Background(tcell.ColorCornflowerBlue).Foreground(tcell.ColorWhite)

	s.display_label(0, s.height-2, style, "Press UP/DOWN to scroll")

	s.display_label(0, s.height-1, style, "Press ENTER to exit")

	s.screen.Show()
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

// Display displays the screen.
//
// Parameters:
//   - drawer: The drawer to draw. Can be nil.
//
// Returns:
//   - error: An error if the screen could not be displayed.
func (s *Screen) Display(drawer Drawer) error {
	if s == nil {
		return gcers.NilReceiver
	}

	var table *DtTable

	if drawer != nil {
		tmp, err := drawer.DrawTable(BgStyle)
		if err != nil {
			return err
		}

		table = tmp
	}

	s.to_display.Set(table)

	s.show_display()

	return nil
}
