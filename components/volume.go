package components

import (
	"strconv"
	"time"

	"github.com/dudubtw/osu-radio-native/lib"
	"github.com/dudubtw/osu-radio-native/theme"
	rl "github.com/gen2brain/raylib-go/raylib"
)

var volumeSliderSize float32 = 100
var volumeSliderProgressWidth float32 = 5
var volumeSliderProgressPadding float32 = 12
var startAngle float32 = 90
var endAngle float32 = 360
var muteButtonBottomPadding float32 = 8
var scrollSpeed = 0.01
var hideDuration time.Duration = 1700 * time.Millisecond

func getMuteButtonIcon(ui *lib.UIStruct) IconName {
	if ui.IsMuted {
		return ICON_VOLUME_MUTED
	}

	return ICON_VOLUME_HIGH
}

func DrawVolumeSlider(volume float32, enabled bool, c *Components) (bool, bool) {
	// BG
	y := (c.ui.ScreenH - int32(volumeSliderSize)) / 2
	x := (c.ui.ScreenW - int32(volumeSliderSize)) - 20
	rl.DrawCircle(x, y, volumeSliderSize, rl.Fade(rl.Black, 0.8))

	// Text
	textContent := strconv.Itoa(int(volume * 100))
	var themeFontSize = theme.FontSize.ExtraSuperLarge
	var fontSize int32 = int32(themeFontSize)
	textX := x - lib.MeasureText(textContent, themeFontSize, theme.FontWeight.Bold, c.typographyMap)/2
	textY := y - fontSize/2
	pos := rl.NewVector2(float32(textX), float32(textY))
	lib.Typography(textContent, pos, themeFontSize, theme.FontWeight.Bold, theme.Colors.Text, c.typographyMap)

	// Progress ring
	center := rl.NewVector2(float32(x), float32(y))
	progressInnerRadius := volumeSliderSize - volumeSliderProgressWidth - volumeSliderProgressPadding
	progressOuterRadius := volumeSliderSize - volumeSliderProgressPadding
	scale := lib.NewLinearScale([]float32{0, 1}, []float32{startAngle, endAngle})
	rl.DrawRing(center, progressInnerRadius, progressOuterRadius, startAngle, scale(volume), 0, rl.Pink)

	// Mute button
	var muteX float32 = float32(x) - ICON_BUTTON_SIZE_GHOST/2
	var muteY float32 = float32(y) + volumeSliderSize - ICON_BUTTON_SIZE_GHOST - volumeSliderProgressPadding - volumeSliderProgressWidth - muteButtonBottomPadding
	isClicked := c.IconButton("volume-mute", getMuteButtonIcon(c.ui), ICON_BUTTON_GHOST, rl.NewRectangle(muteX, muteY, 0, 0))
	isMouseInside := lib.CheckCollisionPointCircle(x, y, volumeSliderSize, c.mousePoint)
	return isClicked, isMouseInside
}

func (c *Components) VolumeSlider(volume float32, enabled bool) (float32, bool, bool) {
	isMuteClicked := false
	isMouseInside := false
	now := time.Now()
	if now.Sub(c.ui.LastTimeScrolled) < hideDuration {
		isMuteClicked, isMouseInside = DrawVolumeSlider(volume, enabled, c)
	}

	if !enabled {
		return volume, isMuteClicked, false
	}

	// Calculate new volume
	mouseMovment := rl.GetMouseWheelMove()
	hasChangedVolume := mouseMovment != 0
	if hasChangedVolume || isMuteClicked || isMouseInside {
		c.ui.LastTimeScrolled = time.Now()
	}

	newVolume := lib.Clamp(volume+mouseMovment*float32(scrollSpeed), 0, 1)

	return newVolume, isMuteClicked, hasChangedVolume
}
