package main

import (
	"github.com/dudubtw/osu-radio-native/components"
	rl "github.com/gen2brain/raylib-go/raylib"
)

type ButtonStyle struct {
	Color        rl.Color
	BorderRadius Roundness
}
type ButtonStyles = map[components.InteractableState]ButtonStyle

func Button(id string, buttonRect rl.Rectangle, styles ButtonStyles) bool {
	interactable := components.NewInteractable(id, &ui)
	clicked := interactable.Event(mousePoint, buttonRect)
	style := styles[interactable.State()]
	DrawRectangleRoundedPixels(buttonRect, style.BorderRadius, style.Color)
	return clicked
}
