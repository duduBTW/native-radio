package components

import (
	"github.com/dudubtw/osu-radio-native/lib"
	rl "github.com/gen2brain/raylib-go/raylib"
)

type IconName int8

const (
	ICON_PLAY            IconName = 0
	ICON_PAUSE           IconName = 1
	ICON_PLAYER_PREVIOUS IconName = 2
	ICON_PLAYER_NEXT     IconName = 3
	ICON_REPEAT          IconName = 4
	ICON_SHUFFLE         IconName = 5
	ICON_ADD_CIRCLE      IconName = 6
	ICON_VOLUME_LOW      IconName = 7
	ICON_VOLUME_MEDIUM   IconName = 8
	ICON_VOLUME_HIGH     IconName = 9
	ICON_VOLUME_MUTED    IconName = 10
	ICON_SIDEBAR         IconName = 11
	ICON_MUSIC           IconName = 12
	ICON_PLAYLIST        IconName = 13
	ICON_CHEVRON_DOWN    IconName = 14
	ICON_SEARCH          IconName = 15
	ICON_FILTER          IconName = 16
	ICON_SETTINGS        IconName = 17
)

const ICON_SIZE float32 = 24

func DrawIcon(name IconName, position rl.Vector2, textures *lib.Textures) {
	posOnSprite := rl.Rectangle{
		Y:      0,
		X:      ICON_SIZE * float32(name),
		Width:  ICON_SIZE,
		Height: ICON_SIZE,
	}
	rl.DrawTexturePro(*textures.Icons, posOnSprite, rl.NewRectangle(position.X, position.Y, ICON_SIZE, ICON_SIZE), rl.Vector2{}, 0, rl.White)
}

type IconButtonVariant = int8

const (
	ICON_BUTTON_PRIMARY   = 0
	ICON_BUTTON_GHOST     = 1
	ICON_BUTTON_SECONDARY = 2
)

type IconButtonState = int8

const (
	ICON_BUTTON_STATE_INITIAL = 0
	ICON_BUTTON_STATE_HOT     = 1
	ICON_BUTTON_STATE_ACTIVE  = 2
)

const ICON_BUTTON_SIZE_PRIMARY = 52
const ICON_BUTTON_SIZE_SECONDARY = 32
const ICON_BUTTON_SIZE_GHOST = 32

func IconButtonPosition(name IconName, variant IconButtonVariant, position rl.Rectangle) rl.Rectangle {
	var buttonSize float32 = 0
	switch variant {
	case ICON_BUTTON_PRIMARY:
		buttonSize = ICON_BUTTON_SIZE_PRIMARY
	case ICON_BUTTON_GHOST:
		buttonSize = ICON_BUTTON_SIZE_GHOST
	case ICON_BUTTON_SECONDARY:
		buttonSize = ICON_BUTTON_SIZE_SECONDARY
	}

	return rl.NewRectangle(position.X, position.Y, buttonSize, buttonSize)
}

func (c *Components) IconButton(id string, name IconName, variant IconButtonVariant, position rl.Rectangle) bool {
	buttonRect := IconButtonPosition(name, variant, position)
	interactable := NewInteractable(id, c.ui)
	clicked := interactable.Event(c.mousePoint, buttonRect)
	state := interactable.State()

	var borderRadius float32 = 0
	var segments int32 = 0
	var color = rl.Fade(rl.Black, 0)

	switch variant {
	case ICON_BUTTON_PRIMARY:
		segments = 12
		borderRadius = 1

		switch state {
		case ICON_BUTTON_STATE_INITIAL:
			color = rl.Purple
		case ICON_BUTTON_STATE_HOT:
			color = rl.DarkPurple
		case ICON_BUTTON_STATE_ACTIVE:
			color = rl.Pink

		}
	case ICON_BUTTON_GHOST:
		borderRadius = 0.2
		switch state {
		case ICON_BUTTON_STATE_HOT:
			color = rl.Fade(rl.Black, 0.2)
		case ICON_BUTTON_STATE_ACTIVE:
			color = rl.Fade(rl.White, 0.1)
		}
	case ICON_BUTTON_SECONDARY:
		borderRadius = 0.2
		switch state {
		case ICON_BUTTON_STATE_INITIAL:
			color = rl.Fade(rl.White, 0.1)
		case ICON_BUTTON_STATE_HOT:
			color = rl.Fade(rl.White, 0.2)
		case ICON_BUTTON_STATE_ACTIVE:
			color = rl.Fade(rl.White, 0.3)
		}
	}

	rl.DrawRectangleRounded(buttonRect, borderRadius, segments, color)
	DrawIcon(name, rl.NewVector2(position.X+(buttonRect.Width-ICON_SIZE)/2, position.Y+(buttonRect.Height-ICON_SIZE)/2), &c.textures)

	return clicked
}
