package components

import (
	"github.com/dudubtw/osu-radio-native/lib"
	rl "github.com/gen2brain/raylib-go/raylib"
)

type InputProps struct {
	Id          string
	X           float32
	Y           float32
	Width       float32
	Placeholder string
	Value       string
	Icon        IconName
	Ui          *lib.UIStruct
	MousePoint  rl.Vector2
}

const INPUT_HEIGHT float32 = 40

func Input(props InputProps) string {
	rect := rl.NewRectangle(props.X, props.Y, props.Width, INPUT_HEIGHT)
	rectInt32 := rect.ToInt32()
	InputEvent(rect, props)
	state := InputState(props)
	borderColor := rl.White

	switch state {
	case STATE_HOT:
		borderColor = rl.Fade(rl.White, 0.5)
	case STATE_INITIAL:
		borderColor = rl.Fade(rl.White, 0.3)
	}

	var fontSize int32 = 16
	var textY int32 = rectInt32.Y + (int32(rect.Height)-fontSize)/2
	var textX int32 = rectInt32.X + 12
	if props.Value == "" {
		rl.DrawText(props.Placeholder, textX, textY, fontSize, rl.Fade(rl.White, 0.42))
	} else {
		rl.DrawText(props.Value, textX, textY, fontSize, rl.White)
	}

	DrawRectangleRoundedLinePixels(rect, ROUNDED, 1, borderColor)
	return InputValueChange(props, state)
}

func InputValueChange(props InputProps, state InteractableState) string {
	value := props.Value
	if state != STATE_ACTIVE {
		return value
	}

	key := rl.GetCharPressed()
	for key > 0 {
		value += string(key)
		key = rl.GetCharPressed()
	}

	if (rl.IsKeyPressed(rl.KeyBackspace) || rl.IsKeyPressedRepeat(rl.KeyBackspace)) && len(props.Value) > 0 {
		value = value[:len(props.Value)-1]
	}

	return value
}

func InputEvent(rect rl.Rectangle, props InputProps) {
	ui := props.Ui
	id := props.Id
	mousePoint := props.MousePoint
	isFocused := id == ui.FocusedId
	isInside := rl.CheckCollisionPointRec(mousePoint, rect)

	if !isInside {
		rl.SetMouseCursor(rl.MouseCursorDefault)
	} else {
		rl.SetMouseCursor(rl.MouseCursorIBeam)
	}

	if isFocused && rl.IsMouseButtonDown(rl.MouseButtonLeft) && !isInside {
		ui.FocusedId = ""
		return
	}

	if isFocused {
		return
	}

	if ui.HotId == id && !isInside {
		ui.HotId = ""
		return
	}

	// Other element is being interacted with
	if (ui.ActiveId != "" && !isFocused) ||
		(ui.HotId != "" && !isInside) {
		return
	}

	// clicked
	if isInside && rl.IsMouseButtonDown(rl.MouseButtonLeft) {
		ui.FocusedId = id
		ui.HotId = ""
		return
	}

	if !isFocused && isInside {
		ui.HotId = id
	}
}

func InputState(props InputProps) InteractableState {
	ui := props.Ui
	switch props.Id {
	case ui.FocusedId:
		return STATE_ACTIVE
	case ui.HotId:
		return STATE_HOT
	default:
		return STATE_INITIAL
	}
}
