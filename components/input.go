package components

import (
	"fmt"

	"github.com/dudubtw/osu-radio-native/lib"
	"github.com/dudubtw/osu-radio-native/theme"
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
}

const INPUT_HEIGHT float32 = 40

func (c *Components) Input(props InputProps) string {
	rect := rl.NewRectangle(props.X, props.Y, props.Width, INPUT_HEIGHT)
	rectInt32 := rect.ToInt32()
	InputEvent(rect, props, c)
	state := InputState(props.Id, c)
	borderColor := rl.White

	switch state {
	case STATE_HOT:
		borderColor = rl.Fade(rl.White, 0.5)
	case STATE_INITIAL:
		borderColor = rl.Fade(rl.White, 0.3)
	}

	isEmpty := props.Value == ""
	var themeFontSize = theme.FontSize.Regular
	var fontSize int32 = int32(themeFontSize)
	var textY int32 = rectInt32.Y + (int32(rect.Height)-fontSize)/2
	var textX int32 = rectInt32.X + 12
	var textPos = rl.NewVector2(float32(textX), float32(textY))
	if isEmpty && state != STATE_ACTIVE {
		lib.Typography(props.Placeholder, textPos, themeFontSize, theme.FontWeight.Regular, theme.Colors.SubText, c.typographyMap)
	} else if !isEmpty {
		lib.Typography(props.Value, textPos, themeFontSize, theme.FontWeight.Regular, theme.Colors.Text, c.typographyMap)
	}

	DrawRectangleRoundedLinePixels(rect, ROUNDED, 1, borderColor)
	newValue := InputValueChange(props, state, c)

	if rl.IsMouseButtonDown(rl.MouseButtonLeft) {
		UpdateClickedCursorPosition(newValue, textX, themeFontSize, c)
	}

	if state == STATE_ACTIVE && c.ui.InputCursorStart == c.ui.InputCursorEnd {
		DrawCusor(c.ui.InputCursorStart, newValue, textX, textY, themeFontSize, c)
	}

	return newValue
}

func UpdateClickedCursorPosition(value string, textX int32, fontSize theme.Text, c *Components) {
	mousePoint := c.mousePoint
	ui := c.ui
	totalTextWidth := lib.MeasureText(value, fontSize, theme.FontWeight.Regular, c.typographyMap)
	if mousePoint.X >= float32(textX+totalTextWidth) {
		ui.SetCursors(len(value))
		return
	}

	if mousePoint.X <= float32(textX) {
		ui.SetCursors(0)
		return
	}

	var lastPos int32 = textX
	index := 0
	for index <= len(value)-1 {
		char := value[index : index+1]
		charSize := lib.MeasureText(char, fontSize, theme.FontWeight.Regular, c.typographyMap)

		if mousePoint.X >= float32(lastPos) && mousePoint.X <= float32(lastPos+charSize) {
			if mousePoint.X > float32(lastPos+(charSize/2)) {
				ui.SetCursors(index + 1)
			} else {
				ui.SetCursors(index)
			}
			return
		}
		//                    letter spacing
		lastPos += charSize + 1
		index++
	}

	fmt.Println("uh oh")
}

func DrawCusor(position int, value string, textX, textY int32, fontSize theme.Text, c *Components) {
	color := rl.Fade(rl.White, 0.6)
	if ShouldBlink(c.ui) {
		color = rl.White
	}

	x := textX + lib.MeasureText(value[:position], fontSize, theme.FontWeight.Regular, c.typographyMap) + 1
	y := textY + 8
	cursorHeight := int32(fontSize) - 6
	rl.DrawLine(x, y-cursorHeight, x, y+cursorHeight, color)
}

const blinkInterval = 0.8
const blinkTotal = 0.5

func ShouldBlink(ui *lib.UIStruct) bool {
	// is blinking
	if ui.BlinkingTimer > 0 {
		return ShouldStayBlinked(ui)
	}

	ui.BlinkTimer += rl.GetFrameTime()
	if ui.BlinkTimer > blinkInterval {
		ui.BlinkTimer = 0
		ui.BlinkingTimer += 0.001
		return true
	}
	return false
}

func ShouldStayBlinked(ui *lib.UIStruct) bool {
	ui.BlinkingTimer += rl.GetFrameTime()
	if ui.BlinkingTimer > blinkInterval {
		ui.BlinkingTimer = 0
		return false
	}
	return true
}

func InputValueChange(props InputProps, state InteractableState, c *Components) string {
	value := props.Value
	ui := c.ui
	if state != STATE_ACTIVE {
		return value
	}

	key := rl.GetCharPressed()
	for key > 0 {
		if len(value) == 0 {
			value += string(key)
		} else {
			fmt.Println(ui.InputCursorStart)
			value = string(value[:ui.InputCursorStart]) + string(key) + string(value[ui.InputCursorStart:])
		}
		key = rl.GetCharPressed()
		ui.IncrementCursor()
	}

	if (rl.IsKeyPressed(rl.KeyBackspace) || rl.IsKeyPressedRepeat(rl.KeyBackspace)) && ui.InputCursorStart > 0 {
		value = string(value[:ui.InputCursorStart-1]) + string(value[ui.InputCursorStart:])
		ui.DecrementCursor()
	}

	if (rl.IsKeyPressed(rl.KeyLeft) || rl.IsKeyPressedRepeat(rl.KeyLeft)) && ui.InputCursorStart > 0 {
		if rl.IsKeyDown(rl.KeyLeftControl) {
			ui.SetCursors(0)
		} else {
			ui.DecrementCursor()
		}
	}

	if (rl.IsKeyPressed(rl.KeyRight) || rl.IsKeyPressedRepeat(rl.KeyRight)) && ui.InputCursorStart < len(props.Value) {
		if rl.IsKeyDown(rl.KeyLeftControl) {
			ui.SetCursors(len(props.Value))
		} else {
			ui.IncrementCursor()
		}
	}

	return value
}

func InputEvent(rect rl.Rectangle, props InputProps, c *Components) bool {
	ui := c.ui
	id := props.Id
	mousePoint := c.mousePoint
	isFocused := id == ui.FocusedId
	isInside := rl.CheckCollisionPointRec(mousePoint, rect)

	if !isInside {
		rl.SetMouseCursor(rl.MouseCursorDefault)
	} else {
		rl.SetMouseCursor(rl.MouseCursorIBeam)
	}

	if isFocused && rl.IsMouseButtonDown(rl.MouseButtonLeft) && !isInside {
		ui.FocusedId = ""
		return false
	}

	if isFocused {
		return false
	}

	if ui.HotId == id && !isInside {
		ui.HotId = ""
		return false
	}

	// Other element is being interacted with
	if (ui.ActiveId != "" && !isFocused) ||
		(ui.HotId != "" && !isInside) {
		return false
	}

	// clicked
	if isInside && rl.IsMouseButtonDown(rl.MouseButtonLeft) {
		ui.FocusedId = id
		ui.HotId = ""
		return true
	}

	if !isFocused && isInside {
		ui.HotId = id
	}

	return false
}

func InputState(id string, c *Components) InteractableState {
	ui := c.ui
	switch id {
	case ui.FocusedId:
		return STATE_ACTIVE
	case ui.HotId:
		return STATE_HOT
	default:
		return STATE_INITIAL
	}
}
