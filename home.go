package main

import (
	c "github.com/dudubtw/osu-radio-native/components"
	"github.com/dudubtw/osu-radio-native/lib"
	"github.com/dudubtw/osu-radio-native/theme"
	rl "github.com/gen2brain/raylib-go/raylib"
)

var scrollSpeed = 40

func HomePage() {
	// Feed music into audio card
	if music.Selected != nil {
		rl.UpdateMusicStream(*music.Selected)
	}

	// Check if song ended
	if !music.IsSkeekMode && music.HasEnded() {
		music.Next(&songTable)
		UpdateSong()
		music.Play()
	}

	// Space pause
	if rl.IsKeyPressed(rl.KeySpace) && ui.FocusedId == "" {
		music.Toggle()
	}

	// Shuffle
	if rl.IsKeyPressed(rl.KeyF2) {
		Shuffle()
	}

	textures.ProcessPendingTextures()
	if songTable.HasSelectedSong() && textures.SelectedSong != nil {
		BackgroundImage()
	}

	// Components
	padding := lib.Padding{}
	var home = lib.NewConstrainedLayout(lib.ContrainedLayout{
		Contrains: rl.Rectangle{Width: float32(ui.ScreenW), Height: float32(ui.ScreenH), X: 0, Y: 0},
		Direction: lib.DIRECTION_ROW,
		Padding:   padding,
		ChildrenSize: []lib.ChildSize{
			{SizeType: lib.SIZE_ABSOLUTE, Value: 520},
			{SizeType: lib.SIZE_WEIGHT, Value: 1},
		},
	})

	home.Render(Panel)
	home.Render(SongDetails(typographyMap))
}

func BackgroundImage() {
	rl.BeginShaderMode(shaders.Blur.Shader)
	lib.DrawFitTexture(*textures.SelectedSong, rl.NewRectangle(0, 0, float32(ui.ScreenW), float32(ui.ScreenH)), rl.Gray)
	rl.EndShaderMode()
	rl.DrawRectangle(0, 0, ui.ScreenW, ui.ScreenH, rl.NewColor(18, 18, 18, 209))
}

func Panel(rect rl.Rectangle) {
	rl.DrawRectanglePro(rect, rl.Vector2{}, 0, rl.Fade(rl.Black, 0.42))
	padding := lib.Padding{}
	padding.Axis(20, 24)
	padding.Bottom(0)
	var container = lib.NewConstrainedLayout(lib.ContrainedLayout{
		Direction: lib.DIRECTION_COLUMN,
		Padding:   padding,
		Gap:       20,
		Contrains: rect,
		ChildrenSize: []lib.ChildSize{
			{SizeType: lib.SIZE_ABSOLUTE, Value: 42},
			{SizeType: lib.SIZE_ABSOLUTE, Value: 40},
			{SizeType: lib.SIZE_WEIGHT, Value: 1},
		},
	})

	container.Render(UpperPart)
	container.Render(Filters)

	switch ui.SelectedPanelPage {
	case lib.PANEL_PAGE_SONGS:
		container.Render(SongList())

	}
}

func UpperPart(rect rl.Rectangle) {
	var container = lib.NewConstrainedLayout(lib.ContrainedLayout{
		Direction: lib.DIRECTION_ROW,
		Gap:       16,
		Contrains: rect,
		ChildrenSize: []lib.ChildSize{
			{SizeType: lib.SIZE_ABSOLUTE, Value: c.ICON_BUTTON_SIZE_GHOST},
			{SizeType: lib.SIZE_WEIGHT, Value: 1},
			{SizeType: lib.SIZE_ABSOLUTE, Value: c.ICON_BUTTON_SIZE_GHOST},
		},
	})

	container.Render(PanelSidebarButton)
	container.Render(UpperPartTabs())
	container.Render(PanelSettingsButton)
}

func Filters(rect rl.Rectangle) {
	layout := lib.NewConstrainedLayout(lib.ContrainedLayout{
		Direction: lib.DIRECTION_ROW,
		Contrains: rect,
		Gap:       12,
		ChildrenSize: []lib.ChildSize{
			{SizeType: lib.SIZE_WEIGHT, Value: 1},
			{SizeType: lib.SIZE_ABSOLUTE, Value: c.ICON_BUTTON_SIZE_GHOST},
		},
	})

	layout.Render(SearchInput)
	layout.Render(FiltersExpandButton)
}

func SearchInput(rect rl.Rectangle) {
	newSearchValue := comp.Input(c.InputProps{
		X:           rect.X,
		Y:           rect.Y,
		Width:       rect.Width,
		Placeholder: "Search in your songs...",
		Id:          "search",
		Value:       ui.SearchValue,
	})

	if newSearchValue != ui.SearchValue {
		ui.SearchValue = newSearchValue
		UpdateSongList()
	}
}

func FiltersExpandButton(rect rl.Rectangle) {
	rect.Y += (rect.Height - c.ICON_BUTTON_SIZE_GHOST) / 2
	comp.IconButton("filter-expand", c.ICON_FILTER, c.ICON_BUTTON_GHOST, rect)
}

func UpperPartTabs() lib.ContrainedComponent {
	return func(rect rl.Rectangle) {
		value := comp.Tabs(c.TabsProps{
			Items: []c.TabsItemProps{
				{
					Icon:  c.ICON_MUSIC,
					Label: "Songs",
					Value: "songs",
				},
				{
					Icon:  c.ICON_PLAYLIST,
					Label: "Playlists",
					Value: "playlists",
				},
			},
			Rect:  rect,
			Value: string(ui.SelectedPanelPage),
		})
		ui.SelectedPanelPage = lib.PanelPage(value)
	}
}

func PanelSidebarButton(rect rl.Rectangle) {
	comp.IconButton("sidebar-collapse", c.ICON_SIDEBAR, c.ICON_BUTTON_GHOST, rl.NewRectangle(rect.X, rect.Y+5, rect.Width, rect.Height))
}

func PanelSettingsButton(rect rl.Rectangle) {
	var variant c.IconButtonVariant = c.ICON_BUTTON_GHOST
	if ui.SelectedPanelPage == lib.PANEL_PAGE_SETTINGS {
		variant = c.ICON_BUTTON_SECONDARY
	}

	if comp.IconButton("settings-button", c.ICON_SETTINGS, variant, rl.NewRectangle(rect.X, rect.Y+5, rect.Width, rect.Height)) {
		ui.SelectedPanelPage = lib.PANEL_PAGE_SETTINGS
	}
}

func SongList() lib.ContrainedComponent {
	return func(rect rl.Rectangle) {
		iRect := rect.ToInt32()
		rl.BeginScissorMode(iRect.X-4, iRect.Y-4, iRect.Width+8, iRect.Height+8)

		rectWithOffset := rl.NewRectangle(rect.X, rect.Y, rect.Width, rect.Height)

		container := lib.NewLayout(lib.Layout{
			Direction: lib.DIRECTION_COLUMN,
			Gap:       12,
		}, rectWithOffset)

		for index, song := range songTable.Songs {
			container.Render(SongCard(song, index, rect))
		}

		scrollTo := func(dest float32) {
			scrollHeight := container.Size.Height - rect.Height
			ui.SidePanelScrollTop = lib.Clamp(dest, -(scrollHeight), 0)
		}

		if rl.IsKeyPressed(rl.KeyPageDown) || rl.IsKeyDown(rl.KeyPageDown) {
			scrollTo(ui.SidePanelScrollTop - float32(ui.ScreenH))
		} else if rl.IsKeyPressed(rl.KeyPageUp) || rl.IsKeyDown(rl.KeyPageUp) {
			scrollTo(ui.SidePanelScrollTop + float32(ui.ScreenH))
		} else if rl.CheckCollisionPointRec(mousePoint, rect) {
			newScroll := ui.SidePanelScrollTop + rl.GetMouseWheelMove()*float32(scrollSpeed)
			scrollHeight := container.Size.Height - rect.Height
			if scrollHeight > newScroll {
				scrollTo(newScroll)
			}
		}

		rl.EndScissorMode()

		h := float32(rect.Height) - 10
		if container.Size.Height > h && (IsActive() || rl.CheckCollisionPointRec(mousePoint, rl.NewRectangle(rect.X-12, rect.Y, rect.Width+12, rect.Height))) {
			scrollbarThumbRatio := h / container.Size.Height
			thumbHeight := lib.Clamp(h*scrollbarThumbRatio, 28, 400)

			maxScroll := container.Size.Height - h
			scrollProgress := rl.Clamp((ui.SidePanelScrollTop*-1)/maxScroll, 0, 1)
			movableSpace := h - thumbHeight

			thumbY := rect.Y + scrollProgress*movableSpace
			thumbRect := rl.NewRectangle(rect.X-12, thumbY, 4, thumbHeight)

			sliderProps := SliderProps{
				Value:      scrollProgress,
				Rect:       rl.NewRectangle(thumbRect.X, rect.Y, thumbRect.Width, rect.Height),
				Color:      rl.Blank,
				TrackColor: rl.Blank,
				Thumb: Thumb{
					Size: lib.Size{
						Width:  5,
						Height: thumbHeight,
					},
					Roundness: 4,
				},
				Direction: SLIDER_DIRECTION_VERTICAL,
				Id:        "scroll-bar",
			}
			newSliderValue := Slider(sliderProps)
			if IsActive() {
				scrollTo(-(maxScroll * newSliderValue))
			}
		}
	}
}

var titleHeight float32 = 30
var artistHeight float32 = 21

func isSongCardaHidden(rect rl.Rectangle) bool {
	return rect.Y < -200 || rect.Y > float32(ui.ScreenH)+200
}

func SongCard(song lib.Song, index int, containerRect rl.Rectangle) lib.Component {
	return func(position lib.Position, next lib.Next) {
		var height float32 = 72
		var rect = position.ToRect(position.Contrains.Width, height)
		rect.Y += ui.SidePanelScrollTop

		if isSongCardaHidden(rect) {
			textures.UnloadSongCard(song)
			next(rect)
			return
		}

		textures.LoadSongCard(song, rect)
		padding := lib.Padding{}
		padding.Axis(20, 16)
		cardContent := lib.NewConstrainedLayout(lib.ContrainedLayout{
			Direction: lib.DIRECTION_COLUMN,
			Padding:   padding,
			Contrains: rect,
			ChildrenSize: []lib.ChildSize{
				{SizeType: lib.SIZE_ABSOLUTE, Value: titleHeight},
				{SizeType: lib.SIZE_ABSOLUTE, Value: artistHeight},
			},
		})

		id := "song-card" + song.Path
		interactable := c.NewInteractable(id, &ui)
		if interactable.Event(mousePoint, rect) && rl.CheckCollisionPointRec(mousePoint, containerRect) {
			SelectSong(index)
		}

		buttonColor := rl.Fade(rl.Black, 0.42)
		isSelected := songTable.SelectedSong().Path == song.Path
		if isSelected {
			buttonColor = rl.Fade(rl.Black, 0.3)
		} else {
			state := interactable.State()
			switch state {
			case c.STATE_HOT:
				buttonColor = rl.Fade(rl.Black, 0.5)
			case c.STATE_ACTIVE:
				buttonColor = rl.Fade(rl.Black, 0.4)
			}
		}

		texture := textures.GetSong(song)
		if texture != nil {
			rl.DrawTexture(*texture, rect.ToInt32().X, rect.ToInt32().Y, rl.White)
		}
		c.DrawRectangleRoundedPixels(rect, c.ROUNDED-1, buttonColor)

		if isSelected {
			c.DrawRectangleRoundedLinePixels(rect, c.ROUNDED+2, 4, rl.White)
		}

		cardContent.Render(SongCardText(song.Title, theme.FontSize.Large, theme.FontWeight.Bold, rl.White))
		cardContent.Render(SongCardText(song.Artist, theme.FontSize.Small, theme.FontWeight.Regular, theme.Colors.Text))

		next(rect)
	}
}

func SongCardText(text string, fontSize theme.Text, weight theme.Weight, color rl.Color) lib.ContrainedComponent {
	return func(rect rl.Rectangle) {
		lib.Typography(text, rl.NewVector2(rect.X-1, rect.Y-1), fontSize, weight, rl.Fade(rl.Black, 0.6), typographyMap)
		lib.Typography(text, rl.NewVector2(rect.X+1, rect.Y+1), fontSize, weight, rl.Fade(rl.Black, 0.6), typographyMap)
		lib.Typography(text, rl.NewVector2(rect.X, rect.Y), fontSize, weight, color, typographyMap)
	}
}
