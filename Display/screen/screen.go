package screen

import (
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

	dtb "github.com/PlayerR9/display/table"
	rws "github.com/PlayerR9/safe/rw_safe"
	"github.com/gdamore/tcell"
)

// Display represents a display that can draw elements to the screen.
type Display struct {
	// screen is the screen of the display.
	screen tcell.Screen

	// width is the width of the display.
	width int

	// height is the height of the display.
	height int

	// evChan is the channel of events.
	evChan chan tcell.Event

	// wg is the wait group of the display.
	wg sync.WaitGroup

	// table is the draw table of the display.
	table *dtb.Table

	// element is the element to draw.
	element *rws.Safe[dtb.Displayer]

	// errChan is the channel of errors.
	errChan chan error

	// keyChan is the channel of key events.
	keyChan chan tcell.EventKey

	// shouldClose is the subject of whether the display should close.
	shouldClose *rws.Safe[bool]

	// bgStyle is the background style of the display.
	bgStyle tcell.Style
}

// NewDisplay creates a new display with the given background style.
//
// Parameters:
//   - bgStyle: The background style of the display.
//
// Returns:
//   - *Display: The new display.
//   - error: An error if the display could not be created.
func NewDisplay(bgStyle tcell.Style) (*Display, error) {
	screen, err := tcell.NewScreen()
	if err != nil {
		return nil, err
	}

	err = screen.Init()
	if err != nil {
		return nil, err
	}

	screen.SetStyle(bgStyle)
	screen.Clear()

	width, height := screen.Size()

	table, err := dtb.NewTable(width, height)
	if err != nil {
		panic(err)
	}

	return &Display{
		screen:  screen,
		width:   width,
		height:  height,
		table:   table,
		bgStyle: bgStyle,
	}, nil
}

// Start starts the display.
func (d *Display) Start() {
	d.evChan = make(chan tcell.Event)
	d.errChan = make(chan error, 1)
	d.keyChan = make(chan tcell.EventKey)

	d.shouldClose = rws.NewSafe[bool](false)
	d.element = rws.NewSafe[dtb.Displayer](nil)

	d.screen.EnableMouse()

	go d.eventListener()

	d.wg.Add(1)

	go d.mainListener()
}

// Close closes the display.
func (d *Display) Close() {
	d.shouldClose.Set(true)

	d.wg.Wait()

	close(d.errChan) // Check this
	d.errChan = nil

	d.screen.Fini()

	close(d.evChan)
	d.evChan = nil

	close(d.keyChan)
	d.keyChan = nil
}

// ReceiveErr receives an error from the display.
//
// Returns:
//   - error: The error.
//   - bool: True if the error was received, false otherwise.
func (d *Display) ReceiveErr() (error, bool) {
	err, ok := <-d.errChan
	if !ok {
		return nil, false
	}

	return err, true
}

// Draw draws an element to the display.
//
// Parameters:
//   - elem: The element to draw.
func (d *Display) Draw(elem dtb.Displayer) {
	d.element.Set(elem)

	d.drawScreen()
}

// eventListener is a helper method that listens for events.
func (d *Display) eventListener() {
	for {
		ev := d.screen.PollEvent()
		if ev == nil {
			break
		}

		d.evChan <- ev
	}
}

// mainListener is a helper method that listens for events.
func (d *Display) mainListener() {
	defer d.wg.Done()

	for {
		select {
		case <-time.After(time.Microsecond * 100):
			if d.shouldClose.Get() {
				return
			}
		case ev := <-d.evChan:
			switch ev := ev.(type) {
			case *tcell.EventResize:
				d.resizeEvent()
			case *tcell.EventKey:
				d.keyChan <- *ev
			}
		}
	}
}

// ListenForKey listens for a key event.
//
// Returns:
//   - rune: The key.
//   - bool: True if the key was received, false otherwise.
func (d *Display) ListenForKey() (rune, bool) {
	key, ok := <-d.keyChan
	if !ok {
		return 0, false
	}

	switch key.Key() {
	case tcell.KeyRune:
		return key.Rune(), true
	case tcell.KeyEnter:
		return '\n', true
	case tcell.KeyBackspace, tcell.KeyBackspace2:
		return '\b', true
	}

	return 0, false
}

// resizeEvent is a helper method that handles a resize event.
func (d *Display) resizeEvent() {
	d.width, d.height = d.screen.Size()

	tmp, err := dtb.NewTable(d.width, d.height)
	if err != nil {
		panic(err)
	}

	d.table = tmp
}

// drawScreen is a helper method that draws the screen.
func (d *Display) drawScreen() {
	d.screen.Clear()

	elem := d.element.Get()

	if elem != nil {
		xCoord := 2
		yCoord := 2

		err := elem.Draw(d.table, &xCoord, &yCoord)
		if err != nil {
			d.errChan <- fmt.Errorf("error drawing element: %w", err)
		}

		for row := range d.table.Row() {
			for j := 0; j < len(row); j++ {
				cell := row[j]

				if cell == nil {
					continue
				}

				d.screen.SetContent(j, yCoord, cell.Char, nil, cell.Style)
			}

			yCoord++
		}
	}

	d.screen.Show()
	time.Sleep(time.Millisecond * 100)
}

// ListenForNumber listens for a number.
//
// Returns:
//   - int: The number.
//   - error: An error if the number could not be received.
func (d *Display) ListenForNumber() (int, error) {
	var builder strings.Builder

	for {
		key, ok := <-d.keyChan
		if !ok {
			break
		}

		kk := key.Key()

		if kk == tcell.KeyEnter {
			break
		}

		if kk == tcell.KeyBackspace || kk == tcell.KeyBackspace2 {
			if builder.Len() > 0 {
				str := builder.String()
				builder.Reset()

				builder.WriteString(str[:len(str)-1])
			}

			continue
		}

		r := key.Rune()

		if r < '0' || r > '9' {
			continue
		}

		builder.WriteRune(r)
	}

	num, err := strconv.Atoi(builder.String())
	if err != nil {
		return 0, err
	}

	return num, nil
}
