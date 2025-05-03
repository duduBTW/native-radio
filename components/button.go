package components

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

type ButtonStyle struct {
	Color        rl.Color
	BorderRadius Roundness
}

type ButtonStyles = map[InteractableState]ButtonStyle

func (c *Components) Button(id string, buttonRect rl.Rectangle, styles ButtonStyles) bool {
	interactable := NewInteractable(id, c.ui)
	clicked := interactable.Event(c.mousePoint, buttonRect)
	style := styles[interactable.State()]
	DrawRectangleRoundedPixels(buttonRect, style.BorderRadius, style.Color)
	return clicked
}
