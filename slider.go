package main

import (
	c "github.com/dudubtw/osu-radio-native/components"
	"github.com/dudubtw/osu-radio-native/lib"
	rl "github.com/gen2brain/raylib-go/raylib"
)

type SliderDirection int

const (
	SLIDER_DIRECTION_HORIZONTAL SliderDirection = 0
	SLIDER_DIRECTION_VERTICAL   SliderDirection = 1
)

type SliderMode int

const (
	SLIDER_STATE_INITIAL SliderMode = 0
	SLIDER_STATE_ACTIVE  SliderMode = 1
	SLIDER_STATE_HOT     SliderMode = 2
)

type Thumb struct {
	Size      lib.Size
	Offset    rl.Vector2
	Roundness c.Roundness
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

func SliderWidth(slider SliderProps) (float32, float32) {
	var totalProgressRect = slider.Rect.Width - (slider.Padding * 2)
	var valueWidth = totalProgressRect * slider.Value
	return totalProgressRect, valueWidth
}

func HorizontalSliderThumbRect(slider SliderProps) rl.Rectangle {
	_, valueWidth := SliderWidth(slider)
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

func SliderHeight(slider SliderProps) (float32, float32) {
	var totalProgressRect = slider.Rect.Height - (slider.Padding * 2)
	var valueWidth = totalProgressRect * slider.Value
	return totalProgressRect, valueWidth
}

func VerticalSliderThumbRect(slider SliderProps) rl.Rectangle {
	total, valueHeight := SliderHeight(slider)
	rect := slider.Rect
	thumb := slider.Thumb

	offsetScale := lib.NewLinearScale([]float32{0, total}, []float32{0, -thumb.Size.Height})

	// CHECK: X OFFSET MAY BE BROKENGE DIDNT TEST IT
	return rl.Rectangle{
		X:      rect.X - thumb.Offset.X,
		Y:      rect.Y + valueHeight - thumb.Offset.Y + offsetScale(valueHeight),
		Width:  thumb.Size.Width + thumb.Offset.X,
		Height: thumb.Size.Height + (thumb.Offset.Y * 2),
	}
}

type SliderProps struct {
	Id           string
	Value        float32
	Rect         rl.Rectangle
	Padding      float32
	BorderRadius c.Roundness
	Thumb        Thumb
	Color        rl.Color
	Direction    SliderDirection
	TrackColor   rl.Color
}

func Slider(slider SliderProps) float32 {
	SliderEventHandler(slider)
	state := SliderStateHandler(slider)
	c.DrawRectangleRoundedPixels(slider.Rect, slider.BorderRadius, slider.TrackColor)

	totalProgressRect, valueWidth := SliderWidth(slider)
	progressRect := rl.Rectangle{
		X:      slider.Rect.X + slider.Padding,
		Y:      slider.Rect.Y + slider.Padding,
		Width:  totalProgressRect,
		Height: slider.Rect.Height - (slider.Padding * 2),
	}

	var rectInt32 = slider.Rect.ToInt32()
	rl.BeginScissorMode(rectInt32.X+3, rectInt32.Y, int32(valueWidth), int32(slider.Rect.Height))
	c.DrawRectangleRoundedPixels(progressRect, slider.BorderRadius-slider.Padding, slider.Color)
	rl.EndScissorMode()

	thumbColor := rl.DarkGray
	switch state {
	case SLIDER_STATE_HOT:
		thumbColor = rl.Gray
	case SLIDER_STATE_ACTIVE:
		thumbColor = rl.LightGray
	}

	var thumbRect rl.Rectangle
	switch slider.Direction {
	case SLIDER_DIRECTION_HORIZONTAL:
		thumbRect = HorizontalSliderThumbRect(slider)
	case SLIDER_DIRECTION_VERTICAL:
		thumbRect = VerticalSliderThumbRect(slider)
	}

	c.DrawRectangleRoundedPixels(thumbRect, slider.Thumb.Roundness, thumbColor)
	newValue := SliderValueHanlder(slider)
	return newValue
}

func SliderValueHanlder(slider SliderProps) float32 {
	if slider.Id != ui.ActiveId {
		return slider.Value
	}

	var start float32
	switch slider.Direction {
	case SLIDER_DIRECTION_HORIZONTAL:
		start = slider.Rect.X
	case SLIDER_DIRECTION_VERTICAL:
		start = slider.Rect.Y
	}

	var end float32
	switch slider.Direction {
	case SLIDER_DIRECTION_HORIZONTAL:
		end = slider.Rect.Width
	case SLIDER_DIRECTION_VERTICAL:
		end = slider.Rect.Height
	}

	var mousePos float32
	switch slider.Direction {
	case SLIDER_DIRECTION_HORIZONTAL:
		mousePos = mousePoint.X - start
	case SLIDER_DIRECTION_VERTICAL:
		mousePos = mousePoint.Y - start
	}

	if mousePos < 0 {
		return 0.001
	}

	if mousePos > end {
		return 1
	}

	return mousePos / end
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

	if slider.Id != ui.ActiveId && isInside && rl.IsMouseButtonPressed(rl.MouseButtonLeft) {
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
