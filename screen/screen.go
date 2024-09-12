package screen

import (
	"fmt"

	gcers "github.com/PlayerR9/go-commons/errors"
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

	// key_ch is the key channel.
	key_ch chan *tcell.EventKey

	// dt is the display table.
	dt *DtTable

	// vt is the virtual table.
	vt *VirtualTable
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

	dt, err := NewDtTable(80, 25)
	if err != nil {
		panic(fmt.Sprintf("could not create table: %v", err.Error()))
	}

	vt := &VirtualTable{
		actual_table: dt,
	}

	return &Screen{
		screen:   screen,
		event_ch: make(chan tcell.Event, 1),
		key_ch:   make(chan *tcell.EventKey),
		dt:       dt,
		vt:       vt,
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

	width, height := s.screen.Size()

	s.dt.ResizeHeight(height)
	s.dt.ResizeWidth(width)

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

// Table returns the display table.
//
// Returns:
//   - Tabler: The display table.
//   - bool: True if the table exists, false otherwise.
func (s *Screen) Table() (*VirtualTable, bool) {
	if s == nil {
		return nil, false
	}

	return s.vt, true
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

			s.dt.ResizeHeight(height)
			s.dt.ResizeWidth(width)

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

	s.vt.Refresh()

	x, y := 0, 0

	width := s.dt.Width()
	height := s.dt.Height()

	for i := 0; i < len(s.dt.cells) && y < height; i++ {
		for j := 0; j < len(s.dt.cells[i]) && x < width; j++ {
			if x >= s.dt.Width() {
				x = 0
				y++
			}

			s.screen.SetContent(x, i, s.dt.cells[i][j].char, nil, s.dt.cells[i][j].style)
			x++
		}

		x = 0
		y++
	}

	style := tcell.StyleDefault.Background(tcell.ColorCornflowerBlue).Foreground(tcell.ColorWhite)

	s.display_label(0, height-2, style, "Press UP/DOWN to scroll")

	s.display_label(0, height-1, style, "Press ENTER to exit")

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
func (s *Screen) Display() error {
	if s == nil {
		return gcers.NilReceiver
	}

	s.show_display()

	return nil
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
