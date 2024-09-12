package format

import (
	"fmt"
	"slices"
	"sync"

	dtb "github.com/PlayerR9/display/table"
	gcers "github.com/PlayerR9/go-commons/errors"
)

// Format is a type that represents a format.
type Format[T gcers.Enumer] struct {
	// m is the map of elements.
	m map[T]dtb.Displayer

	// order is the order of the elements.
	order []T

	// mu is the mutex of the format.
	mu sync.RWMutex
}

// Draw implements the drawtable.Displayer interface.
func (f *Format[T]) Draw(table *dtb.Table, x, y *int) error {
	order := f.GetOrder()

	xCoord := *x
	yCoord := *y

	for _, order := range order {
		elem, ok := f.GetElement(order)
		if !ok {
			continue
		}

		err := elem.Draw(table, &xCoord, &yCoord)
		if err != nil {
			return fmt.Errorf("error drawing element %s: %w", order.String(), err)
		}

		xCoord = *x
		yCoord += 2 // Skip a line
	}

	*x = xCoord
	*y = yCoord

	return nil
}

// NewFormat returns a new format.
//
// Returns:
//   - *Format: The new format.
func NewFormat[T gcers.Enumer]() *Format[T] {
	return &Format[T]{
		m: make(map[T]dtb.Displayer),
	}
}

// AddElement adds an element to the format.
//
// Parameters:
//   - key: The key of the element.
//   - elem: The element to add.
func (f *Format[T]) AddElement(key T, elem dtb.Displayer) {
	f.mu.Lock()
	defer f.mu.Unlock()

	f.m[key] = elem
	f.order = append(f.order, key)
}

// GetOrder returns the order of the elements.
//
// Returns:
//   - []T: The order of the elements.
func (f *Format[T]) GetOrder() []T {
	f.mu.RLock()
	defer f.mu.RUnlock()

	order := make([]T, 0, len(f.order))
	order = append(order, f.order...)

	return order
}

// GetElement returns an element from the format.
//
// Parameters:
//   - key: The key of the element.
//
// Returns:
//   - Displayer: The element.
//   - bool: True if the element exists, false otherwise.
func (f *Format[T]) GetElement(key T) (dtb.Displayer, bool) {
	f.mu.RLock()
	defer f.mu.RUnlock()

	elem, ok := f.m[key]
	return elem, ok
}

// ReplaceElement replaces an element in the format.
//
// Parameters:
//   - key: The key of the element.
//   - elem: The element to replace the existing element with.
func (f *Format[T]) ReplaceElement(key T, elem dtb.Displayer) {
	f.mu.Lock()
	defer f.mu.Unlock()

	f.m[key] = elem
}

// RemoveElement removes an element from the format.
//
// Parameters:
//   - key: The key of the element.
func (f *Format[T]) RemoveElement(key T) {
	f.mu.Lock()
	defer f.mu.Unlock()

	delete(f.m, key)

	for i, order := range f.order {
		if order == key {
			f.order = slices.Delete(f.order, i, i+1)
			break
		}
	}
}
