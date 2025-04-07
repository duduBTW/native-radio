package main

import rl "github.com/gen2brain/raylib-go/raylib"

type InteractableState = int8

const (
	STATE_INITIAL InteractableState = 0
	STATE_HOT     InteractableState = 1
	STATE_ACTIVE  InteractableState = 2
)

type ButtonStyle struct {
	Color        rl.Color
	BorderRadius Roundness
}
type ButtonStyles = map[InteractableState]ButtonStyle

func Button(id string, buttonRect rl.Rectangle, styles ButtonStyles) bool {
	clicked := ButtonEventHandler(id, buttonRect)
	state := ButtonStateHandler(id, buttonRect)

	style := styles[state]
	DrawRectangleRoundedPixels(buttonRect, style.BorderRadius, style.Color)
	return clicked
}

func ButtonStateHandler(id string, rect rl.Rectangle) InteractableState {
	if id == UI.ActiveId {
		return ICON_BUTTON_STATE_ACTIVE
	}
	if id == UI.HotId {
		return ICON_BUTTON_STATE_HOT
	}

	return ICON_BUTTON_STATE_INITIAL
}

func ButtonEventHandler(id string, rect rl.Rectangle) bool {
	isActive := id == UI.ActiveId
	isInside := rl.CheckCollisionPointRec(MousePoint, rect)

	if isActive && rl.IsMouseButtonUp(rl.MouseButtonLeft) {
		UI.ActiveId = ""
		return isInside
	}

	if isActive {
		return false
	}

	if UI.HotId == id && !isInside {
		UI.HotId = ""
	}

	if !isInside {
		return false
	}

	if rl.IsMouseButtonDown(rl.MouseButtonLeft) {
		UI.ActiveId = id
		return false
	}

	if UI.HotId == "" {
		UI.HotId = id
	}

	return false
}
