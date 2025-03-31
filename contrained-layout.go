package main

import (
	"errors"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type ChildSize struct {
	SizeType
	Value float32
}

type ContrainedLayout struct {
	Padding

	index int

	Contrains             rl.Rectangle
	Direction             Direction
	Gap                   float32
	VerticalAlignment     Alignment
	HorizontalAlighment   Alignment
	ChildrenSize          []ChildSize
	ChildrenComputedSizes []float32

	rect rl.Rectangle
}

type ContrainedComponent func(current rl.Rectangle)

func NewConstrainedLayout(layout ContrainedLayout) ContrainedLayout {
	layout.Contrains.Width -= (layout.Padding.start + layout.Padding.end)
	layout.Contrains.Height -= (layout.Padding.top + layout.Padding.bottom)

	layout.ComputeChildren()
	return layout
}

func (layout *ContrainedLayout) Render(component ContrainedComponent) rl.Rectangle {
	targetRect := rl.Rectangle{X: layout.rect.X + layout.Padding.start + layout.Contrains.X, Y: layout.rect.Y + layout.Padding.top + layout.Contrains.Y}

	switch layout.Direction {
	case DIRECTION_ROW:
		targetRect.Width = layout.ChildrenComputedSizes[layout.index]
		targetRect.Height = layout.Contrains.Height
	case DIRECTION_COLUMN:
		targetRect.Height = layout.ChildrenComputedSizes[layout.index]
		targetRect.Width = layout.Contrains.Width
	}

	if component != nil {
		component(targetRect)
	}

	switch layout.Direction {
	case DIRECTION_ROW:
		layout.rect.X += targetRect.Width + layout.Gap
	case DIRECTION_COLUMN:
		layout.rect.Y += targetRect.Height + layout.Gap
	}

	layout.rect.Width = targetRect.Width
	layout.rect.Height = targetRect.Height
	layout.index++

	return targetRect
}

func (layout *ContrainedLayout) ComputeChildren() error {
	if len(layout.ChildrenSize) == 0 {
		return errors.New("no children to compute")
	}

	type Index = int
	var remainingSize float32 = 0
	switch layout.Direction {
	case DIRECTION_ROW:
		remainingSize = layout.Contrains.Width - layout.Gap*float32(len(layout.ChildrenSize)-1)
	case DIRECTION_COLUMN:
		remainingSize = layout.Contrains.Height - layout.Gap*float32(len(layout.ChildrenSize)-1)
	}

	var weightSum float32 = 0
	var weightSizes = make(map[Index]float32)
	var computedSizes = make([]float32, len(layout.ChildrenSize))
	for index, childSize := range layout.ChildrenSize {
		value := childSize.Value
		if childSize.SizeType == SIZE_WEIGHT {
			weightSizes[index] = value
			weightSum += value
			continue
		}

		computedSizes[index] = value
		remainingSize -= value
	}

	if weightSum > 1 {
		return errors.New("weight sum not equal to 1")
	}

	for index, weight := range weightSizes {
		computedSizes[index] = float32(remainingSize) * weight
	}

	layout.ChildrenComputedSizes = computedSizes
	return nil
}
