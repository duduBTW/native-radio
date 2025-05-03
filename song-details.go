package main

import (
	"strconv"

	c "github.com/dudubtw/osu-radio-native/components"
	"github.com/dudubtw/osu-radio-native/lib"
	"github.com/dudubtw/osu-radio-native/theme"
	rl "github.com/gen2brain/raylib-go/raylib"
)

func SongDetails(typographyMap lib.TypographyMap) func(rect rl.Rectangle) {
	return func(rect rl.Rectangle) {
		isVolumeSliderActive := rl.CheckCollisionPointRec(mousePoint, rect)
		originalRect := rect.Width
		rect.Width = lib.Min(1000, rect.Width)

		if originalRect != rect.Width {
			rect.X += (originalRect - rect.Width) / 2
		}

		padding := lib.Padding{}
		padding.All(94)
		var container = lib.NewConstrainedLayout(lib.ContrainedLayout{
			Direction: lib.DIRECTION_COLUMN,
			Padding:   padding,
			Contrains: rect,
			Gap:       4,
			ChildrenSize: []lib.ChildSize{
				{SizeType: lib.SIZE_WEIGHT, Value: 1},
				{SizeType: lib.SIZE_ABSOLUTE, Value: 184},
			},
		})

		container.Render(SongMiniature)
		container.Render(SongControls(typographyMap))

		SetVolume(comp.VolumeSlider(ui.Volume, isVolumeSliderActive))
	}
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

func SongControls(typographyMap lib.TypographyMap) func(rect rl.Rectangle) {
	return func(rect rl.Rectangle) {
		var container = lib.NewLayout(lib.Layout{
			Direction: lib.DIRECTION_COLUMN,
			Gap:       16,
		}, rect)

		container.Render(SongControlTexts)
		container.Render(SongProgress(typographyMap))
		container.Render(SongControlButton)
	}
}

func durationFromSeconds(seconds int) string {
	minutes := strconv.Itoa(seconds / 60)

	resultSeconds := strconv.Itoa(seconds % 60)
	for len(resultSeconds) < 2 {
		resultSeconds = "0" + resultSeconds
	}

	return minutes + ":" + resultSeconds
}

func SongProgress(typographyMap lib.TypographyMap) func(avaliablePosition lib.Position, next lib.Next) {
	return func(position lib.Position, next lib.Next) {
		if !songTable.HasSelectedSong() {
			return
		}

		rect := position.ToRect(position.Contrains.Width, 28)

		sliderProps := SliderProps{
			Value:        music.Progress(),
			Rect:         rect,
			Padding:      4,
			BorderRadius: c.ROUNDED,
			Thumb: Thumb{
				Size: lib.Size{
					Width:  6,
					Height: rect.Height,
				},
				Offset: rl.Vector2{
					X: 0,
					Y: 4,
				},
				Roundness: 3,
			},
			Color:      rl.Pink,
			TrackColor: rl.Fade(rl.Black, 0.8),
			Id:         "progress-slider",
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

		thumbRect := HorizontalSliderThumbRect(sliderProps)
		progressSecongs := int(rl.GetMusicTimeLength(*music.Selected) * sliderProps.Value)
		progressText := durationFromSeconds(progressSecongs)
		var progressFontSize = theme.FontSize.Regular
		var progressFontWeight = theme.FontWeight.Bold
		var progressTextPaddingHorizontal float32 = 12
		var progressTextY float32 = rect.Y + (rect.Height-float32(progressFontSize))/2
		var textWidth = float32(lib.MeasureText(progressText, progressFontSize, progressFontWeight, typographyMap))

		progressPos := rl.NewVector2(
			lib.Max(
				thumbRect.X-thumbRect.Width-textWidth,
				rect.X+progressTextPaddingHorizontal,
			),
			progressTextY,
		)

		lib.Typography(progressText, progressPos, progressFontSize, progressFontWeight, theme.Colors.Text, typographyMap)

		totalText := durationFromSeconds(int(rl.GetMusicTimeLength(*music.Selected)))
		totalFontWeight := theme.FontWeight.Regular
		var totalTextWidth = float32(lib.MeasureText(totalText, progressFontSize, totalFontWeight, typographyMap))
		var totalTextX float32 = rect.X + rect.Width - totalTextWidth - progressTextPaddingHorizontal
		const hideTOtalTextOffset float32 = 4
		if totalTextX > thumbRect.X+thumbRect.Width+hideTOtalTextOffset {
			totalPos := rl.NewVector2(totalTextX, progressTextY)
			lib.Typography(totalText, totalPos, progressFontSize, totalFontWeight, theme.Colors.SubText, typographyMap)
		}

		next(rect)
	}
}

func SongControlButton(position lib.Position, next lib.Next) {
	rect := position.ToRect(position.Contrains.Width, 52)

	container := lib.NewConstrainedLayout(lib.ContrainedLayout{
		Direction: lib.DIRECTION_ROW,
		Gap:       12,
		ChildrenSize: []lib.ChildSize{
			{SizeType: lib.SIZE_ABSOLUTE, Value: c.ICON_BUTTON_SIZE_GHOST},
			{SizeType: lib.SIZE_ABSOLUTE, Value: c.ICON_BUTTON_SIZE_PRIMARY},
			{SizeType: lib.SIZE_WEIGHT, Value: 1.0},
			{SizeType: lib.SIZE_ABSOLUTE, Value: c.ICON_BUTTON_SIZE_GHOST},
			{SizeType: lib.SIZE_ABSOLUTE, Value: c.ICON_BUTTON_SIZE_GHOST},
			{SizeType: lib.SIZE_ABSOLUTE, Value: c.ICON_BUTTON_SIZE_GHOST},
		},
		Contrains: rect,
	})

	offsetTop := (rect.Height - c.ICON_BUTTON_SIZE_GHOST) / 2

	isPlaying := rl.IsMusicStreamPlaying(*music.Selected)

	var ghostWithOffset = func(iconRect rl.Rectangle) rl.Rectangle {
		iconRect.Y += offsetTop
		return iconRect
	}

	if comp.IconButton("c-play2", c.ICON_PLAYER_PREVIOUS, c.ICON_BUTTON_GHOST, ghostWithOffset(container.Render(nil))) {
		music.Previous(&songTable)
		UpdateSong()
	}

	var iconName c.IconName = c.ICON_PLAY
	if isPlaying {
		iconName = c.ICON_PAUSE
	}
	if comp.IconButton("c-play", iconName, c.ICON_BUTTON_PRIMARY, container.Render(nil)) {
		if isPlaying {
			music.Pause()
		} else {
			music.Play()
		}
	}
	if comp.IconButton("c-play3", c.ICON_PLAYER_NEXT, c.ICON_BUTTON_GHOST, ghostWithOffset(container.Render(nil))) {
		music.Next(&songTable)
		UpdateSong()
	}

	comp.IconButton("repeat", c.ICON_REPEAT, c.ICON_BUTTON_GHOST, ghostWithOffset(container.Render(nil)))
	if comp.IconButton("shuffle", c.ICON_SHUFFLE, c.ICON_BUTTON_GHOST, ghostWithOffset(container.Render(nil))) {
		Shuffle()
	}

	comp.IconButton("add-playlist", c.ICON_ADD_CIRCLE, c.ICON_BUTTON_GHOST, ghostWithOffset(container.Render(nil)))

	next(rect)
}

func SongControlTexts(position lib.Position, next lib.Next) {
	var container = lib.NewLayout(lib.Layout{
		Direction: lib.DIRECTION_COLUMN,
		Gap:       4,
	}, position.ToRect(position.Contrains.Width, position.Contrains.Height))

	// FIXME - Line height
	container.Render(SongControlText(songTable.SelectedSong().Title, theme.FontSize.ExtraLarge, theme.FontWeight.Bold, theme.Colors.Text))
	container.Render(SongControlText(songTable.SelectedSong().Artist, theme.FontSize.Regular, theme.FontWeight.Regular, theme.Colors.SubText))

	// FIXME - Getting the height of children after, is this ok? prob not right, at least it is weird rn
	next(position.ToRect(position.Contrains.Width, container.Size.Height))
}

func SongControlText(text string, fontSize theme.Text, weight theme.Weight, color rl.Color) lib.Component {
	return func(position lib.Position, next lib.Next) {
		pos := rl.NewVector2(position.X, position.Y)
		lib.Typography(text, pos, fontSize, weight, color, typographyMap)
		next(position.ToRect(position.Contrains.Width, float32(fontSize)))
	}
}
