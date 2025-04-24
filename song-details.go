package main

import (
	"strconv"

	"github.com/dudubtw/osu-radio-native/components"
	c "github.com/dudubtw/osu-radio-native/components"
	rl "github.com/gen2brain/raylib-go/raylib"
)

func SongDetails(rect rl.Rectangle) {
	isVolumeSliderActive := rl.CheckCollisionPointRec(mousePoint, rect)
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

	SetVolume(components.VolumeSlider(ui.Volume, isVolumeSliderActive, &ui, &textures, mousePoint, db))
}

func SongMiniature(rect rl.Rectangle) {
	if !songTable.HasSelectedSong() {
		return
	}

	miniature := textures.Miniature
	if miniature == nil {
		return
	}

	rl.BeginShaderMode(shaders.Shadow)
	var size float32 = 350
	x := int32(rect.X + ((rect.Width - size) / 2))
	y := int32(rect.Y + ((rect.Height - size) / 2))
	rl.DrawTexture(*miniature, x, y, rl.White)
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
	if !songTable.HasSelectedSong() {
		return
	}

	rect := position.ToRect(position.Contrains.Width, 28)
	rectInt32 := rect.ToInt32()

	sliderProps := SliderProps{
		Value:        music.Progress(),
		Rect:         rect,
		Padding:      4,
		BorderRadius: c.ROUNDED,
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
		music.BeginSeekMode()
	}

	if IsItemDeactivatedAfterEdit() {
		music.ExitSeekMode()
	}

	if IsActive() {
		music.Seek(newSliderValue)
	}

	thumbRect := SliderThumbRect(sliderProps)
	progressSecongs := int(rl.GetMusicTimeLength(*music.Selected) * sliderProps.Value)
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

	totalText := durationFromSeconds(int(rl.GetMusicTimeLength(*music.Selected)))
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
			{SizeType: SIZE_ABSOLUTE, Value: c.ICON_BUTTON_SIZE_GHOST},
			{SizeType: SIZE_ABSOLUTE, Value: c.ICON_BUTTON_SIZE_PRIMARY},
			{SizeType: SIZE_WEIGHT, Value: 1.0},
			{SizeType: SIZE_ABSOLUTE, Value: c.ICON_BUTTON_SIZE_GHOST},
			{SizeType: SIZE_ABSOLUTE, Value: c.ICON_BUTTON_SIZE_GHOST},
			{SizeType: SIZE_ABSOLUTE, Value: c.ICON_BUTTON_SIZE_GHOST},
		},
		Contrains: rect,
	})

	offsetTop := (rect.Height - c.ICON_BUTTON_SIZE_GHOST) / 2

	isPlaying := rl.IsMusicStreamPlaying(*music.Selected)

	var ghostWithOffset = func(iconRect rl.Rectangle) rl.Rectangle {
		iconRect.Y += offsetTop
		return iconRect
	}

	if c.IconButton("c-play2", c.ICON_PLAYER_PREVIOUS, c.ICON_BUTTON_GHOST, ghostWithOffset(container.Render(nil)), &ui, &textures, mousePoint) {
		music.Previous(&songTable)
		UpdateSong()
	}

	var iconName c.IconName = c.ICON_PLAY
	if isPlaying {
		iconName = c.ICON_PAUSE
	}
	if c.IconButton("c-play", iconName, c.ICON_BUTTON_PRIMARY, container.Render(nil), &ui, &textures, mousePoint) {
		if isPlaying {
			music.Pause()
		} else {
			music.Play()
		}
	}
	if c.IconButton("c-play3", c.ICON_PLAYER_NEXT, c.ICON_BUTTON_GHOST, ghostWithOffset(container.Render(nil)), &ui, &textures, mousePoint) {
		music.Next(&songTable)
		UpdateSong()
	}

	c.IconButton("repeat", c.ICON_REPEAT, c.ICON_BUTTON_GHOST, ghostWithOffset(container.Render(nil)), &ui, &textures, mousePoint)
	if c.IconButton("shuffle", c.ICON_SHUFFLE, c.ICON_BUTTON_GHOST, ghostWithOffset(container.Render(nil)), &ui, &textures, mousePoint) {
		Shuffle()
	}

	c.IconButton("add-playlist", c.ICON_ADD_CIRCLE, c.ICON_BUTTON_GHOST, ghostWithOffset(container.Render(nil)), &ui, &textures, mousePoint)

	next(rect)
}

func SongControlTexts(position Position, next Next) {
	var container = NewLayout(Layout{
		Direction: DIRECTION_COLUMN,
		Gap:       4,
	}, position.ToRect(position.Contrains.Width, position.Contrains.Height))

	// FIXME - Line height
	container.Render(SongControlText(songTable.SelectedSong().Title, 28))
	container.Render(SongControlText(songTable.SelectedSong().Artist, 16))

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
