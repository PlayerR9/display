package screen

import (
	gcers "github.com/PlayerR9/go-commons/errors"
	"github.com/gdamore/tcell"
)

type Component interface {
	Draw(vt *VirtualTable)
}

type VirtualTable struct {
	actual_table *DtTable

	x, y          int
	width, height int
	bg_style      tcell.Style

	data Component

	partitions []*VirtualTable
}

func (vt *VirtualTable) Allocate(info PartitionInfo) (*VirtualTable, error) {
	if vt == nil {
		return nil, gcers.NilReceiver
	}

	var x int

	if info.x_from_right {
		x = vt.x + vt.width - info.x
	} else {
		x = vt.x + info.x
	}

	var y int

	if info.y_from_bottom {
		y = vt.y + vt.height - info.y
	} else {
		y = vt.y + info.y
	}

	var width int

	if info.width < 0 {
		width = vt.width - x
	} else if info.width_from_right {
		width = vt.width - x - info.width
	} else {
		width = info.width
	}

	var height int

	if info.height < 0 {
		height = vt.height - y
	} else if info.height_from_bottom {
		height = vt.height - y - info.height
	} else {
		height = info.height
	}

	sub_vt := &VirtualTable{
		actual_table: vt.actual_table,
		x:            x,
		y:            y,
		width:        width,
		height:       height,
		bg_style:     vt.bg_style,
	}

	// TODO: Check for overlapping partitions.

	sub_vt.partitions = append(sub_vt.partitions, vt.partitions...)

	return sub_vt, nil
}

func (vt *VirtualTable) Associate(cmp Component) error {
	if vt == nil {
		return gcers.NilReceiver
	} else if cmp == nil {
		return gcers.NewErrNilParameter("cmp")
	}

	vt.data = cmp

	return nil
}

func (vt *VirtualTable) ChangeBgStyle(style tcell.Style) {
	if vt == nil {
		return
	}

	vt.bg_style = style

	for _, partition := range vt.partitions {
		partition.ChangeBgStyle(style)
	}
}

func (vt VirtualTable) Height() int {
	return vt.height
}

func (vt VirtualTable) Width() int {
	return vt.width
}

func (vt *VirtualTable) DrawCellAt(x, y int, cell *DtCell) {
	if vt == nil {
		return
	}

	actual_x := vt.x + x
	actual_y := vt.y + y

	vt.actual_table.DrawCellAt(actual_x, actual_y, cell)
}

func (vt *VirtualTable) ShiftHorizontal(x int) bool {
	if vt == nil {
		return false
	}

	vt.x += x

	for _, partition := range vt.partitions {
		_ = partition.ShiftHorizontal(x)
	}

	return true
}

func (vt *VirtualTable) ShiftVertical(y int) bool {
	if vt == nil {
		return false
	}

	vt.y += y

	for _, partition := range vt.partitions {
		_ = partition.ShiftVertical(y)
	}

	return true
}

func (vt VirtualTable) BgStyle() tcell.Style {
	return vt.bg_style
}

func (vt *VirtualTable) Refresh() {
	if vt == nil {
		return
	}

	vt.data.Draw(vt)

	for _, partition := range vt.partitions {
		partition.Refresh()
	}
}
