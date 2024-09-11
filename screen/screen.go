package screen

import (
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

	// key_ch is the key channel.
	key_ch chan *tcell.EventKey

	// width is the width of the screen.
	width int

	// height is the height of the screen.
	height int

	// dt is the display table.
	dt *DtTable

	// to_display is the table to display.
	to_display *rws.Safe[*DtTable]

	// pos_y is the position of the cursor.
	pos_y *rws.Safe[int]

	// pos_x is the position of the cursor.
	pos_x *rws.Safe[int]
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
		key_ch:     make(chan *tcell.EventKey),
		width:      80,
		height:     25,
		pos_x:      rws.NewSafe(0),
		pos_y:      rws.NewSafe(0),
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

			s.width, s.height = ev.Size()

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

	table := s.to_display.Get()

	if table == nil {
		s.screen.Show()

		return
	}

	y := 0
	x := 0

	pos_x := s.pos_x.Get()
	pos_y := s.pos_y.Get()

	if pos_y >= table.Height() {
		s.screen.Show()

		return
	} else if pos_y < 0 {
		y -= pos_y
		pos_y = 0
	}

	underlying_table := table.cells[pos_y:]

	for _, row := range underlying_table {
		if y >= s.height-3 {
			break
		}

		j := 0
		for _, c := range row {
			if j >= s.width {
				j = 0
				y++
				x = pos_x // Change this to indent the next line
			}

			if c == nil {
				c = NewDtCell(' ', BgStyle)
			}

			s.screen.SetContent(x, y, c.char, nil, c.style)

			j++
			x++
		}

		y++
		x = pos_x
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

// Table returns the table.
//
// Returns:
//   - *DtTable: The table. Nil only if the receiver is nil.
func (s *Screen) Table() *DtTable {
	if s == nil {
		return nil
	}

	return s.dt
}

// Display displays the screen.
//
// Parameters:
//   - drawer: The drawer to draw. Can be nil.
//
// Returns:
//   - error: An error if the screen could not be displayed.
func (s *Screen) Display(x, y int, drawer Drawer) error {
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
	s.pos_x.Set(x)
	s.pos_y.Set(y)

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
