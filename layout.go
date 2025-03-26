package main

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

type Padding struct {
	top    float32
	bottom float32
	start  float32
	end    float32
}

func (p *Padding) Axis(horizontal, vertical float32) {
	p.top = vertical
	p.bottom = vertical
	p.start = horizontal
	p.end = horizontal
}
func (p *Padding) All(padding float32) {
	p.top = padding
	p.bottom = padding
	p.start = padding
	p.end = padding
}
func (p *Padding) Top(top float32) {
	p.top = top
}
func (p *Padding) Bottom(bottom float32) {
	p.bottom = bottom
}
func (p *Padding) Start(start float32) {
	p.start = start
}
func (p *Padding) End(end float32) {
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
	X         float32
	Y         float32
	Contrains rl.Rectangle
}

func (position Position) ToRect(width, height float32) rl.Rectangle {
	return rl.NewRectangle(position.X, position.Y, width, height)
}

type Size struct {
	Width  float32
	Height float32
}

type Layout struct {
	Padding

	index int

	Direction           Direction
	Gap                 int
	VerticalAlignment   Alignment
	HorizontalAlighment Alignment

	Position Position
	Size     Size
}

type Next func(rl.Rectangle)

type Component func(avaliablePosition Position, next Next)

func NewLayout(layout Layout, contrains rl.Rectangle) Layout {
	layout.Position.X = contrains.X + layout.start
	layout.Position.Y = contrains.Y + layout.top
	layout.Position.Contrains = contrains
	layout.Position.Contrains.Width -= (layout.Padding.start + layout.Padding.end)
	layout.Position.Contrains.Height -= (layout.Padding.top + layout.Padding.bottom)

	layout.Size.Width += (layout.Padding.start + layout.Padding.end)
	layout.Size.Height += (layout.Padding.top + layout.Padding.bottom)
	return layout
}

func (layout *Layout) Render(component Component) {
	component(layout.Position, layout.Next)
}

func (layout *Layout) Next(component rl.Rectangle) {
	switch layout.Direction {
	case DIRECTION_ROW:
		// size
		layout.Size.Width = layout.Size.Width + component.Width
		layout.Size.Height = Max(layout.Size.Height, component.Height)

		// position
		layout.Position.X = layout.Position.X + component.Width + float32(layout.Gap)
	case DIRECTION_COLUMN:
		// size
		layout.Size.Height = layout.Size.Height + component.Height
		layout.Size.Width = Max(layout.Size.Width, component.Width)

		// positioN
		layout.Position.Y = layout.Position.Y + component.Height + float32(layout.Gap)
	}
}

func Max(value, max float32) float32 {
	if value < max {
		return max
	}

	return value
}

func Min(value, min float32) float32 {
	if value > min {
		return min
	}

	return value
}
