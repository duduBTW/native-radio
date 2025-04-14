package components

import (
	"strconv"
	"time"

	"github.com/dudubtw/osu-radio-native/lib"
	rl "github.com/gen2brain/raylib-go/raylib"
)

var volumeSliderSize float32 = 100
var volumeSliderProgressWidth float32 = 5
var volumeSliderProgressPadding float32 = 12
var startAngle float32 = 90
var endAngle float32 = 360
var muteButtonBottomPadding float32 = 8
var scrollSpeed = 0.01
var hideDuration time.Duration = 1400 * time.Millisecond

func getMuteButtonIcon(ui *lib.UIStruct) IconName {
	if ui.IsMuted {
		return ICON_VOLUME_MUTED
	}

	return ICON_VOLUME_HIGH
}

func DrawVolumeSlider(volume float32, enabled bool, ui *lib.UIStruct, textures *lib.Textures, mousePoint rl.Vector2) bool {
	// BG
	y := (ui.ScreenH - int32(volumeSliderSize)) / 2
	x := (ui.ScreenW - int32(volumeSliderSize)) - 20
	rl.DrawCircle(x, y, volumeSliderSize, rl.Fade(rl.Black, 0.8))

	// Text
	textContent := strconv.Itoa(int(volume * 100))
	var fontSize int32 = 40
	textX := x - rl.MeasureText(textContent, fontSize)/2
	textY := y - fontSize/2
	rl.DrawText(textContent, textX, textY, fontSize, rl.White)

	// Progress ring
	center := rl.NewVector2(float32(x), float32(y))
	progressInnerRadius := volumeSliderSize - volumeSliderProgressWidth - volumeSliderProgressPadding
	progressOuterRadius := volumeSliderSize - volumeSliderProgressPadding
	scale := lib.NewLinearScale([2]float32{0, 1}, [2]float32{startAngle, endAngle})
	rl.DrawRing(center, progressInnerRadius, progressOuterRadius, startAngle, scale(volume), 0, rl.Pink)

	// Mute button
	var muteX float32 = float32(x) - ICON_BUTTON_SIZE_GHOST/2
	var muteY float32 = float32(y) + volumeSliderSize - ICON_BUTTON_SIZE_GHOST - volumeSliderProgressPadding - volumeSliderProgressWidth - muteButtonBottomPadding
	return IconButton("volume-mute", getMuteButtonIcon(ui), ICON_BUTTON_GHOST, rl.NewRectangle(muteX, muteY, 0, 0), ui, textures, mousePoint)
}

func VolumeSlider(volume float32, enabled bool, ui *lib.UIStruct, textures *lib.Textures, mousePoint rl.Vector2) (float32, bool) {
	isMuteClicked := false
	now := time.Now()
	if now.Sub(ui.LastTimeScrolled) < hideDuration {
		isMuteClicked = DrawVolumeSlider(volume, enabled, ui, textures, mousePoint)
	}

	if !enabled {
		return volume, isMuteClicked
	}

	// Calculate new volume
	mouseMovment := rl.GetMouseWheelMove()
	if mouseMovment != 0 || isMuteClicked {
		ui.LastTimeScrolled = time.Now()
	}

	newVolume := volume + mouseMovment*float32(scrollSpeed)
	return lib.Clamp(newVolume, 0, 1), isMuteClicked
}
