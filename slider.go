package main

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

type SliderMode int

const (
	SLIDER_STATE_INITIAL SliderMode = 0
	SLIDER_STATE_ACTIVE  SliderMode = 1
	SLIDER_STATE_HOT     SliderMode = 2
)

type Thumb struct {
	Size      Size
	Offset    rl.Vector2
	Roundness Roundness
}

type SliderState struct {
	IsItemActivated            bool
	IsItemDeactivatedAfterEdit bool
	IsActive                   bool
}

var sliderState = SliderState{}

func IsItemActivated() bool {
	temp := sliderState.IsItemActivated
	sliderState.IsItemActivated = false
	return temp
}

func IsActive() bool {
	return sliderState.IsActive
}

func IsItemDeactivatedAfterEdit() bool {
	temp := sliderState.IsItemDeactivatedAfterEdit
	sliderState.IsItemDeactivatedAfterEdit = false
	return temp
}

func SliderValue(slider SliderProps) (float32, float32) {
	var totalProgressRect = slider.Rect.Width - (slider.Padding * 2)
	var valueWidth = totalProgressRect * slider.Value
	return totalProgressRect, valueWidth
}

func SliderThumbRect(slider SliderProps) rl.Rectangle {
	_, valueWidth := SliderValue(slider)
	rect := slider.Rect
	thumb := slider.Thumb

	// CHECK: X OFFSET MAY BE BROKENGE DIDNT TEST IT
	return rl.Rectangle{
		X:      rect.X + valueWidth - thumb.Offset.X,
		Y:      rect.Y - thumb.Offset.Y,
		Width:  thumb.Size.Width + thumb.Offset.X,
		Height: thumb.Size.Height + (thumb.Offset.Y * 2),
	}
}

type SliderProps struct {
	Id           string
	Value        float32
	Rect         rl.Rectangle
	Padding      float32
	BorderRadius Roundness
	Thumb        Thumb
	Color        rl.Color
}

func Slider(slider SliderProps) float32 {
	SliderEventHandler(slider)
	state := SliderStateHandler(slider)
	DrawRectangleRoundedPixels(slider.Rect, slider.BorderRadius, rl.Fade(rl.Black, 0.8))

	totalProgressRect, valueWidth := SliderValue(slider)
	progressRect := rl.Rectangle{
		X:      slider.Rect.X + slider.Padding,
		Y:      slider.Rect.Y + slider.Padding,
		Width:  totalProgressRect,
		Height: slider.Rect.Height - (slider.Padding * 2),
	}

	var rectInt32 = slider.Rect.ToInt32()
	rl.BeginScissorMode(rectInt32.X+3, rectInt32.Y, int32(valueWidth), int32(slider.Rect.Height))
	DrawRectangleRoundedPixels(progressRect, slider.BorderRadius-slider.Padding, slider.Color)
	rl.EndScissorMode()

	thumbColor := rl.DarkGray
	switch state {
	case SLIDER_STATE_HOT:
		thumbColor = rl.Gray
	case SLIDER_STATE_ACTIVE:
		thumbColor = rl.LightGray
	}

	thumbRect := SliderThumbRect(slider)
	DrawRectangleRoundedPixels(thumbRect, slider.Thumb.Roundness, thumbColor)

	newValue := SliderValueHanlder(slider)
	return newValue
}

func SliderValueHanlder(slider SliderProps) float32 {
	if slider.Id != ui.ActiveId {
		return slider.Value
	}

	start := slider.Rect.X
	end := slider.Rect.Width
	mouseX := mousePoint.X - start

	if mouseX < 0 {
		return 0.001
	}

	if mouseX > end {
		return 1
	}

	return mouseX / end
}

func SliderStateHandler(slider SliderProps) SliderMode {
	if slider.Id == ui.ActiveId {
		sliderState.IsActive = true
		return SLIDER_STATE_ACTIVE
	}

	if slider.Id == ui.HotId {
		return SLIDER_STATE_HOT
	}

	return SLIDER_STATE_INITIAL
}

func SliderEventHandler(slider SliderProps) {
	isActive := slider.Id == ui.ActiveId
	isInside := rl.CheckCollisionPointRec(mousePoint, slider.Rect)

	// Other element is being interacted with
	if ui.ActiveId != "" && !isActive {
		return
	}

	if isActive && rl.IsMouseButtonReleased(rl.MouseLeftButton) {
		ui.ActiveId = ""
		sliderState.IsItemDeactivatedAfterEdit = true
		sliderState.IsActive = false
		return
	}

	if isActive {
		return
	}

	if slider.Id != ui.ActiveId && isInside && rl.IsMouseButtonDown(rl.MouseButtonLeft) {
		ui.ActiveId = slider.Id
		sliderState.IsItemActivated = true
		return
	}

	if ui.HotId == slider.Id && !isInside {
		ui.HotId = ""
		return
	}

	if ui.HotId != slider.Id && isInside {
		ui.HotId = slider.Id
	}
}
