package main

import (
	"errors"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type Padding struct {
	top    int
	bottom int
	start  int
	end    int
}

func (p *Padding) Axis(horizontal, vertical int) {
	p.top = vertical
	p.bottom = vertical
	p.start = horizontal
	p.end = horizontal
}
func (p *Padding) All(padding int) {
	p.top = padding
	p.bottom = padding
	p.start = padding
	p.end = padding
}
func (p *Padding) Top(top int) {
	p.top = top
}
func (p *Padding) Bottom(bottom int) {
	p.bottom = bottom
}
func (p *Padding) Start(start int) {
	p.start = start
}
func (p *Padding) End(end int) {
	p.end = end
}

type Alignment string

const (
	ALIGNMENT_START  Alignment = "start"
	ALIGNMENT_CENTER Alignment = "center"
	ALIGNMENT_END    Alignment = "end"
)

type Direction string

const (
	DIRECTION_ROW    Direction = "row"
	DIRECTION_COLUMN Direction = "column"
)

type SizeType string

const (
	SIZE_ABSOLUTE SizeType = "absolute"
	SIZE_WEIGHT   SizeType = "weight"
)

type Position struct {
	X int
	Y int
	W int
	H int
}
type ChildSize struct {
	SizeType
	Value float32
}
type Layout struct {
	Padding

	index int

	Direction             Direction
	Gap                   int
	VerticalAlignment     Alignment
	HorizontalAlighment   Alignment
	Contrains             rl.Rectangle
	ChildrenSize          []ChildSize
	ChildrenComputedSizes []float32

	MaxW int
	MaxH int

	currentPos Position
}

func (layout *Layout) ComputeChildren() error {
	if len(layout.ChildrenSize) == 0 {
		return errors.New("no children to compute")
	}

	type Index = int
	var remainingSize float32 = 0
	switch layout.Direction {
	case DIRECTION_ROW:
		remainingSize = layout.Contrains.Width
	case DIRECTION_COLUMN:
		remainingSize = layout.Contrains.Height
	}

	var weightSum = 0
	var weightSizes = make(map[Index]float32)
	var computedSizes = make([]float32, len(layout.ChildrenSize))
	for index, childSize := range layout.ChildrenSize {
		value := childSize.Value
		if childSize.SizeType == SIZE_WEIGHT {
			weightSizes[index] = value
			weightSum += int(value)
			continue
		}

		computedSizes[index] = value
		remainingSize -= value
	}

	if weightSum != 1 {
		return errors.New("weight sum not equal to 1")
	}

	for index, weight := range weightSizes {
		computedSizes[index] = float32(remainingSize) * weight
	}

	layout.ChildrenComputedSizes = computedSizes
	return nil
}

func (layout *Layout) Current() rl.Rectangle {
	currentRect := rl.Rectangle{X: float32(layout.currentPos.X), Y: float32(layout.currentPos.Y)}

	if len(layout.ChildrenComputedSizes) > 0 {
		switch layout.Direction {
		case DIRECTION_ROW:
			currentRect.Width = layout.ChildrenComputedSizes[layout.index]
			currentRect.Height = layout.Contrains.Height
		case DIRECTION_COLUMN:
			currentRect.Height = layout.ChildrenComputedSizes[layout.index]
			currentRect.Width = layout.Contrains.Width
		}
	} else {
		switch layout.Direction {
		case DIRECTION_ROW:
			currentRect.Height = layout.Contrains.Height
		case DIRECTION_COLUMN:
			currentRect.Width = layout.Contrains.Width
		}
	}

	return currentRect
}
func (layout *Layout) Next(position Position) {
	layout.currentPos = position
	layout.index++
	layout.MaxW = Max(layout.MaxW, position.W)
	layout.MaxH = Max(layout.MaxH, position.H)
}

func Max(value, max int) int {
	if value < max {
		return max
	}

	return value
}
