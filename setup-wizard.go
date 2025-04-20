package main

import (
	c "github.com/dudubtw/osu-radio-native/components"
	"github.com/dudubtw/osu-radio-native/lib"
	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/ncruces/zenity"
)

var selectedFolder = "AA"

func drawBg() {
	rl.BeginShaderMode(shaders.BgBlur.Shader)
	shaders.BgBlur.SetTime(float32(rl.GetTime()))
	shaders.BgBlur.SetRes([2]float32{float32(ui.ScreenW), float32(ui.ScreenH)})
	rl.DrawRectangle(0, 0, ui.ScreenW, ui.ScreenH, rl.White)
	rl.EndShaderMode()
}

func SetupWizardPage() {
	drawBg()

	containerWidth := lib.MinInt32(ui.ScreenW, 860)
	containerHeight := ui.ScreenH - 160
	containerX := (ui.ScreenW - containerWidth) / 2
	containerY := 80
	rect := rl.NewRectangle(float32(containerX), float32(containerY), float32(containerWidth), float32(containerHeight))
	c.DrawRectangleRoundedPixels(rect, c.ROUNDED, rl.NewColor(39, 39, 42, 255))

	padding := Padding{}
	padding.Axis(26, 0)
	padding.Top(82)
	padding.Bottom(32)
	layout := NewConstrainedLayout(ContrainedLayout{
		Direction: DIRECTION_COLUMN,
		Padding:   padding,
		Contrains: rect,
		Gap:       24,
		ChildrenSize: []ChildSize{
			{SizeType: SIZE_WEIGHT, Value: 0.5},
			{SizeType: SIZE_ABSOLUTE, Value: 72},
			{SizeType: SIZE_WEIGHT, Value: 0.5},
		},
	})

	layout.Render(SetupTitle)
	layout.Render(SetupFolderSelector)
	layout.Render(SetupSubmitButton)
}

func SetupTitle(rect rl.Rectangle) {
	rectInt32 := rect.ToInt32()
	text := "Welcome to osu! radio"
	var fontSize int32 = 24
	x := rectInt32.X + (int32(rect.Width)-rl.MeasureText(text, fontSize))/2

	rl.DrawText(text, x, rectInt32.Y, fontSize, rl.White)
}

func SetupFolderSelector(rect rl.Rectangle) {
	layout := NewLayout(Layout{
		Direction: DIRECTION_COLUMN,
		Gap:       6,
	}, rect)

	layout.Render(SetupFolderSelectorTitle)
	layout.Render(SetupFolderSelectorInput)
}

func SetupFolderSelectorTitle(pos Position, next Next) {
	text := "Your osu! Songs folder"
	var fontSize int32 = 16
	rl.DrawText(text, int32(pos.X), int32(pos.Y), fontSize, rl.White)
	next(pos.ToRect(float32(rl.MeasureText(text, fontSize)), float32(fontSize)))
}

func SetupFolderSelectorInput(pos Position, next Next) {
	rect := pos.ToRect(pos.Contrains.Width, 46)
	padding := Padding{}
	padding.All(4)
	padding.Start(16)
	layout := NewConstrainedLayout(ContrainedLayout{
		Direction: DIRECTION_ROW,
		Padding:   padding,
		Contrains: rect,
		Gap:       12,
		ChildrenSize: []ChildSize{
			{SizeType: SIZE_WEIGHT, Value: 1},
			{SizeType: SIZE_ABSOLUTE, Value: 128},
		},
	})

	c.DrawRectangleRoundedPixels(rect, c.ROUNDED, rl.Fade(rl.Black, 0.22))
	layout.Render(SetupFolderSelectorText)
	layout.Render(SetupFolderSelectorButton)
	next(rect)
}

func SetupFolderSelectorText(rect rl.Rectangle) {
	rectInt32 := rect.ToInt32()
	var fontSize int32 = 16
	var y int32 = rectInt32.Y + 11
	rl.DrawText(selectedFolder, rectInt32.X, y, fontSize, rl.White)
}

func SetupFolderSelectorButton(rect rl.Rectangle) {
	getButtonStyle := func(color rl.Color) ButtonStyle {
		return ButtonStyle{Color: color, BorderRadius: c.ROUNDED}
	}

	buttonStyle := ButtonStyles{
		c.STATE_INITIAL: getButtonStyle(rl.Fade(rl.White, 0.15)),
		c.STATE_HOT:     getButtonStyle(rl.Fade(rl.White, 0.2)),
		c.STATE_ACTIVE:  getButtonStyle(rl.Fade(rl.White, 0.4)),
	}

	if Button("select-folder", rect, buttonStyle) {
		dir, err := zenity.SelectFile(
			zenity.Title("Select the osu! lazer folder"),
			zenity.Directory(), // This makes it folder-only
		)

		if err == nil {
			selectedFolder = dir
		}
	}

	rectInt32 := rect.ToInt32()
	text := "Select Folder"
	var fontSize int32 = 16
	y := rectInt32.Y + (int32(rect.Height)-fontSize)/2
	x := rectInt32.X + (int32(rect.Width)-rl.MeasureText(text, fontSize))/2
	rl.DrawText(text, x, y, int32(fontSize), rl.White)
}

func SetupSubmitButton(container rl.Rectangle) {
	var width float32 = 110
	var height float32 = 42
	btnX := container.X + container.Width - width
	btnY := container.Y + container.Height - height
	rect := rl.NewRectangle(btnX, btnY, width, height)
	getButtonStyle := func(color rl.Color) ButtonStyle {
		return ButtonStyle{Color: color, BorderRadius: c.ROUNDED}
	}

	buttonStyle := ButtonStyles{
		c.STATE_INITIAL: getButtonStyle(rl.Fade(rl.White, 1)),
		c.STATE_HOT:     getButtonStyle(rl.Fade(rl.White, 0.9)),
		c.STATE_ACTIVE:  getButtonStyle(rl.Fade(rl.White, 0.8)),
	}

	if Button("submit", rect, buttonStyle) {
		ExecTrash(selectedFolder)
	}

	rectInt32 := rect.ToInt32()
	text := "Submit"
	var fontSize int32 = 16
	y := rectInt32.Y + (int32(rect.Height)-fontSize)/2
	x := rectInt32.X + (int32(rect.Width)-rl.MeasureText(text, fontSize))/2
	rl.DrawText(text, x, y, int32(fontSize), rl.Black)
}
