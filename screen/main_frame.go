package screen

import (
	"fmt"

	gcers "github.com/PlayerR9/go-commons/errors"
)

type PartitionTyper interface {
	~int

	// String returns the string representation of the partition.
	//
	// Returns:
	//   - string: The string representation of the partition.
	String() string
}

type PartitionOption func(*PartitionInfo)

func WithPartitionX(x int, from_right bool) PartitionOption {
	if x < 0 {
		x *= -1
	}

	return func(info *PartitionInfo) {
		info.x = x
		info.x_from_right = from_right
	}
}

func WithPartitionY(y int, from_bottom bool) PartitionOption {
	if y < 0 {
		y *= -1
	}

	return func(info *PartitionInfo) {
		info.y = y
		info.y_from_bottom = from_bottom
	}
}

func WithWidth(width int, from_right bool) PartitionOption {
	if width < 0 {
		width *= -1
	}

	return func(info *PartitionInfo) {
		info.width = width
		info.width_from_right = from_right
	}
}

func WithHeight(height int, from_bottom bool) PartitionOption {
	if height < 0 {
		height *= -1
	}

	return func(info *PartitionInfo) {
		info.height = height
		info.height_from_bottom = from_bottom
	}
}

type PartitionInfo struct {
	x            int
	x_from_right bool

	y             int
	y_from_bottom bool

	width            int
	width_from_right bool

	height             int
	height_from_bottom bool
}

func NewPartitionInfo(opts ...PartitionOption) PartitionInfo {
	pi := PartitionInfo{
		x:      0,
		y:      0,
		width:  -1,
		height: -1,
	}

	for _, opt := range opts {
		opt(&pi)
	}

	return pi
}

type MainFrame[T PartitionTyper] struct {
	table       map[T]PartitionInfo
	assoc_table map[T]Component
}

func NewMainFrame[T PartitionTyper]() MainFrame[T] {
	return MainFrame[T]{
		table:       make(map[T]PartitionInfo),
		assoc_table: make(map[T]Component),
	}
}

func (mf *MainFrame[T]) Register(id T, info PartitionInfo) error {
	if mf == nil {
		return gcers.NilReceiver
	}

	_, ok := mf.table[id]
	if ok {
		return fmt.Errorf("partition %v is already registered", id)
	}

	mf.table[id] = info
	mf.assoc_table[id] = nil

	return nil
}

func (mf *MainFrame[T]) Associate(id T, cmp Component) error {
	if mf == nil {
		return gcers.NilReceiver
	}

	prev, ok := mf.assoc_table[id]
	if !ok {
		return fmt.Errorf("partition %v is not registered", id)
	} else if prev != nil {
		return fmt.Errorf("partition %v is already associated", id)
	}

	mf.assoc_table[id] = cmp

	return nil
}

func (mf MainFrame[T]) Apply(screen *Screen) (map[T]*VirtualTable, error) {
	if screen == nil {
		return nil, gcers.NewErrNilParameter("screen")
	}

	table, ok := screen.Table()
	if !ok {
		return nil, fmt.Errorf("screen does not have a table")
	}

	partition_table := make(map[T]*VirtualTable, len(mf.table)+1)
	partition_table[T(0)] = table

	for id, info := range mf.table {
		alloc, err := table.Allocate(info)
		if err != nil {
			return partition_table, fmt.Errorf("could not allocate table: %w", err)
		}

		cmp, ok := mf.assoc_table[id]
		if !ok || cmp == nil {
			return partition_table, fmt.Errorf("partition %v is not associated", id)
		}

		err = alloc.Associate(mf.assoc_table[id])
		if err != nil {
			return partition_table, fmt.Errorf("could not associate table: %w", err)
		}

		partition_table[id] = alloc
	}

	return partition_table, nil
}
