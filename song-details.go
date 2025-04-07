package main

import (
	"strconv"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func SongDetails(rect rl.Rectangle) {
	originalRect := rect.Width
	rect.Width = Min(1000, rect.Width)

	if originalRect != rect.Width {
		rect.X += (originalRect - rect.Width) / 2
	}

	padding := Padding{}
	padding.All(94)
	var container = NewConstrainedLayout(ContrainedLayout{
		Direction: DIRECTION_COLUMN,
		Padding:   padding,
		Contrains: rect,
		Gap:       4,
		ChildrenSize: []ChildSize{
			{SizeType: SIZE_WEIGHT, Value: 1},
			{SizeType: SIZE_ABSOLUTE, Value: 184},
		},
	})

	container.Render(SongMiniature)
	container.Render(SongControls)
}

func SongMiniature(rect rl.Rectangle) {
	if !UI.HasSelectedSong() {
		return
	}

	miniature := UI.SelectedSongMiniature()
	if miniature == nil {
		return
	}

	rl.BeginShaderMode(Shaders.Shadow)
	var size float32 = 350
	x := int32(rect.X + ((rect.Width - size) / 2))
	y := int32(rect.Y + ((rect.Height - size) / 2))
	rl.DrawTexture(*UI.SelectedSongMiniature(), x, y, rl.White)
	rl.EndShaderMode()
}

func SongControls(rect rl.Rectangle) {
	var container = NewLayout(Layout{
		Direction: DIRECTION_COLUMN,
		Gap:       16,
	}, rect)

	container.Render(SongControlTexts)
	container.Render(SongProgress)
	container.Render(SongControlButton)
}

func durationFromSeconds(seconds int) string {
	minutes := strconv.Itoa(seconds / 60)

	resultSeconds := strconv.Itoa(seconds % 60)
	for len(resultSeconds) < 2 {
		resultSeconds = "0" + resultSeconds
	}

	return minutes + ":" + resultSeconds
}

func SongProgress(position Position, next Next) {
	if !UI.HasSelectedSong() {
		return
	}

	rect := position.ToRect(position.Contrains.Width, 28)
	rectInt32 := rect.ToInt32()

	music := UI.Music
	sliderProps := SliderProps{
		Value:        UI.Progress(),
		Rect:         rect,
		Padding:      4,
		BorderRadius: ROUNDED,
		Thumb: Thumb{
			Size: Size{
				Width:  6,
				Height: rect.Height,
			},
			Offset: rl.Vector2{
				X: 0,
				Y: 4,
			},
			Roundness: 3,
		},
		Color: rl.Pink,
		Id:    "progress-slider",
	}
	newSliderValue := Slider(sliderProps)

	if IsItemActivated() {
		UI.BeginSeekMode()
	}

	if IsItemDeactivatedAfterEdit() {
		UI.ExitSeekMode()
	}

	if IsActive() {
		UI.Seek(newSliderValue)
	}

	thumbRect := SliderThumbRect(sliderProps)
	progressSecongs := int(rl.GetMusicTimeLength(music) * sliderProps.Value)
	thumbRectInt := thumbRect.ToInt32()
	progressText := durationFromSeconds(progressSecongs)
	var progressFontSize int32 = 14
	var progressTextPaddingHorizontal int32 = 12
	var progressTextY int32 = rectInt32.Y + (rectInt32.Height-progressFontSize)/2

	rl.DrawText(progressText,
		MaxInt32(thumbRectInt.X-thumbRectInt.Width-rl.MeasureText(progressText, progressFontSize), rectInt32.X+progressTextPaddingHorizontal),
		progressTextY,
		progressFontSize,
		rl.White,
	)

	totalText := durationFromSeconds(int(rl.GetMusicTimeLength(music)))
	var totalTextX int32 = rectInt32.X + rectInt32.Width - rl.MeasureText(totalText, progressFontSize) - progressTextPaddingHorizontal
	const hideTOtalTextOffset int32 = 4
	if totalTextX > thumbRectInt.X+thumbRectInt.Width+hideTOtalTextOffset {
		rl.DrawText(totalText,
			totalTextX,
			progressTextY,
			progressFontSize,
			rl.White,
		)
	}

	next(rect)
}

func SongControlButton(position Position, next Next) {
	rect := position.ToRect(position.Contrains.Width, 52)

	container := NewConstrainedLayout(ContrainedLayout{
		Direction: DIRECTION_ROW,
		Gap:       12,
		ChildrenSize: []ChildSize{
			{SizeType: SIZE_ABSOLUTE, Value: ICON_BUTTON_SIZE_GHOST},
			{SizeType: SIZE_ABSOLUTE, Value: ICON_BUTTON_SIZE_PRIMARY},
			{SizeType: SIZE_WEIGHT, Value: 1.0},
			{SizeType: SIZE_ABSOLUTE, Value: ICON_BUTTON_SIZE_GHOST},
			{SizeType: SIZE_ABSOLUTE, Value: ICON_BUTTON_SIZE_GHOST},
			{SizeType: SIZE_ABSOLUTE, Value: ICON_BUTTON_SIZE_GHOST},
		},
		Contrains: rect,
	})

	offsetTop := (rect.Height - ICON_BUTTON_SIZE_GHOST) / 2

	isPlaying := rl.IsMusicStreamPlaying(UI.Music)

	var ghostWithOffset = func(iconRect rl.Rectangle) rl.Rectangle {
		iconRect.Y += offsetTop
		return iconRect
	}

	if IconButton("c-play2", ICON_PLAYER_PREVIOUS, ICON_BUTTON_GHOST, ghostWithOffset(container.Render(nil))) {
		UI.Previous()
	}

	var iconName IconName = ICON_PLAY
	if isPlaying {
		iconName = ICON_PAUSE
	}
	if IconButton("c-play", iconName, ICON_BUTTON_PRIMARY, container.Render(nil)) {
		if isPlaying {
			rl.PauseMusicStream(UI.Music)
		} else {
			rl.PlayMusicStream(UI.Music)
		}
	}
	if IconButton("c-play3", ICON_PLAYER_NEXT, ICON_BUTTON_GHOST, ghostWithOffset(container.Render(nil))) {
		UI.Next()
	}

	IconButton("repeat", ICON_REPEAT, ICON_BUTTON_GHOST, ghostWithOffset(container.Render(nil)))
	IconButton("shuffle", ICON_SHUFFLE, ICON_BUTTON_GHOST, ghostWithOffset(container.Render(nil)))
	IconButton("add-playlist", ICON_ADD_CIRCLE, ICON_BUTTON_GHOST, ghostWithOffset(container.Render(nil)))

	next(rect)
}

func SongControlTexts(position Position, next Next) {
	var container = NewLayout(Layout{
		Direction: DIRECTION_COLUMN,
		Gap:       4,
	}, position.ToRect(position.Contrains.Width, position.Contrains.Height))

	// FIXME - Line height
	container.Render(SongControlText(UI.SelectedSong().Title, 28))
	container.Render(SongControlText(UI.SelectedSong().Artist, 16))

	// FIXME - Getting the height of children after, is this ok? prob not right, at least it is weird rn
	next(position.ToRect(position.Contrains.Width, container.Size.Height))
}

func SongControlText(text string, fontSize float32) Component {
	return func(position Position, next Next) {
		font := rl.GetFontDefault()
		textHeight := DrawTextByWidth(font, text, rl.NewVector2(position.X, position.Y), position.Contrains.Width, fontSize, 2, rl.White)
		next(position.ToRect(position.Contrains.Width, textHeight))
	}
}
