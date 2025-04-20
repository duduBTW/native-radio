package main

import (
	c "github.com/dudubtw/osu-radio-native/components"
	"github.com/dudubtw/osu-radio-native/lib"
	rl "github.com/gen2brain/raylib-go/raylib"
)

var scrollSpeed = 40

func HomePage() {
	if music.Selected != nil {
		rl.UpdateMusicStream(*music.Selected)
	}

	if !music.IsSkeekMode && music.HasEnded() {
		music.Next(&songTable)
		UpdateSong()
		music.Play()
	}

	if rl.IsKeyPressed(rl.KeySpace) && ui.FocusedId == "" {
		music.Toggle()
	}

	if rl.IsKeyPressed(rl.KeyF2) {
		SelectSong(lib.RandomRange(0, len(songTable.Songs)))
	}

	textures.ProcessPendingTextures()

	if songTable.HasSelectedSong() && textures.SelectedSong != nil {
		rl.BeginShaderMode(shaders.Blur.Shader)
		lib.DrawFitTexture(*textures.SelectedSong, rl.NewRectangle(0, 0, float32(ui.ScreenW), float32(ui.ScreenH)), rl.Gray)
		rl.EndShaderMode()
		rl.DrawRectangle(0, 0, ui.ScreenW, ui.ScreenH, rl.NewColor(18, 18, 18, 209))
	}

	padding := Padding{}
	var home = NewConstrainedLayout(ContrainedLayout{
		Contrains: rl.Rectangle{Width: float32(ui.ScreenW), Height: float32(ui.ScreenH), X: 0, Y: 0},
		Direction: DIRECTION_ROW,
		Padding:   padding,
		ChildrenSize: []ChildSize{
			{SizeType: SIZE_ABSOLUTE, Value: 520},
			{SizeType: SIZE_WEIGHT, Value: 1},
		},
	})

	home.Render(Panel())
	home.Render(SongDetails)
}

func Panel() ContrainedComponent {
	return func(rect rl.Rectangle) {
		rl.DrawRectanglePro(rect, rl.Vector2{}, 0, rl.Fade(rl.Black, 0.42))
		padding := Padding{}
		padding.Axis(20, 24)
		padding.Bottom(0)
		var container = NewConstrainedLayout(ContrainedLayout{
			Direction: DIRECTION_COLUMN,
			Padding:   padding,
			Gap:       20,
			Contrains: rect,
			ChildrenSize: []ChildSize{
				{SizeType: SIZE_ABSOLUTE, Value: 42},
				{SizeType: SIZE_ABSOLUTE, Value: 40},
				{SizeType: SIZE_WEIGHT, Value: 1},
			},
		})

		container.Render(UpperPart)
		container.Render(Filters)

		switch ui.SelectedPanelPage {
		case lib.PANEL_PAGE_SONGS:
			container.Render(SongList())

		}
	}
}

func UpperPart(rect rl.Rectangle) {
	var container = NewConstrainedLayout(ContrainedLayout{
		Direction: DIRECTION_ROW,
		Gap:       16,
		Contrains: rect,
		ChildrenSize: []ChildSize{
			{SizeType: SIZE_ABSOLUTE, Value: c.ICON_BUTTON_SIZE_GHOST},
			{SizeType: SIZE_WEIGHT, Value: 1},
			{SizeType: SIZE_ABSOLUTE, Value: c.ICON_BUTTON_SIZE_GHOST},
		},
	})

	container.Render(PanelSidebarButton)
	container.Render(UpperPartTabs())
	container.Render(PanelSettingsButton)
}

func Filters(rect rl.Rectangle) {
	newSearchValue := c.Input(c.InputProps{
		X:           rect.X,
		Y:           rect.Y,
		Width:       rect.Width,
		Placeholder: "Search in your songs...",
		Id:          "search",
		Ui:          &ui,
		MousePoint:  mousePoint,
		Value:       ui.SearchValue,
	})

	if newSearchValue != ui.SearchValue {
		ui.SearchValue = newSearchValue
		UpdateSongList()
	}
}

func UpperPartTabs() ContrainedComponent {
	return func(rect rl.Rectangle) {
		value := Tabs(TabsProps{
			Items: []TabsItemProps{
				{
					Icon:  c.ICON_REPEAT,
					Label: "Songs",
					Value: "songs",
				},
				{
					Icon:  c.ICON_VOLUME_MUTED,
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
	c.IconButton("sidebar-collapse", c.ICON_SIDEBAR, c.ICON_BUTTON_GHOST, rl.NewRectangle(rect.X, rect.Y+5, rect.Width, rect.Height), &ui, &textures, mousePoint)
}

func PanelSettingsButton(rect rl.Rectangle) {
	var variant c.IconButtonVariant = c.ICON_BUTTON_GHOST
	if ui.SelectedPanelPage == lib.PANEL_PAGE_SETTINGS {
		variant = c.ICON_BUTTON_SECONDARY
	}

	if c.IconButton("settings-button", c.ICON_SETTINGS, variant, rl.NewRectangle(rect.X, rect.Y+5, rect.Width, rect.Height), &ui, &textures, mousePoint) {
		ui.SelectedPanelPage = lib.PANEL_PAGE_SETTINGS
	}
}

func SongList() ContrainedComponent {
	return func(rect rl.Rectangle) {
		iRect := rect.ToInt32()
		rl.BeginScissorMode(iRect.X-4, iRect.Y-4, iRect.Width+8, iRect.Height+8)

		rectWithOffset := rl.NewRectangle(rect.X, rect.Y, rect.Width, rect.Height)

		container := NewLayout(Layout{
			Direction: DIRECTION_COLUMN,
			Gap:       12,
		}, rectWithOffset)

		for index, song := range songTable.Songs {
			container.Render(SongCard(song, index, rect))
		}

		if rl.CheckCollisionPointRec(mousePoint, rect) {
			newScroll := ui.SidePanelScrollTop + rl.GetMouseWheelMove()*float32(scrollSpeed)
			scrollHeight := container.Size.Height - rect.Height
			if scrollHeight > newScroll {
				ui.SidePanelScrollTop = lib.Clamp(newScroll, -(scrollHeight), 0)
			}
		}

		rl.EndScissorMode()
	}
}

var titleHeight float32 = 30
var artistHeight float32 = 21

func isSongCardaHidden(rect rl.Rectangle) bool {
	return rect.Y < -200 || rect.Y > float32(ui.ScreenH)+200
}

func SongCard(song lib.Song, index int, containerRect rl.Rectangle) Component {
	return func(position Position, next Next) {
		var height float32 = 72
		var rect = position.ToRect(position.Contrains.Width, height)
		rect.Y += ui.SidePanelScrollTop

		if isSongCardaHidden(rect) {
			textures.UnloadSongCard(song)
			next(rect)
			return
		}

		textures.LoadSongCard(song, rect)
		padding := Padding{}
		padding.Axis(20, 16)
		cardContent := NewConstrainedLayout(ContrainedLayout{
			Direction: DIRECTION_COLUMN,
			Padding:   padding,
			Contrains: rect,
			ChildrenSize: []ChildSize{
				{SizeType: SIZE_ABSOLUTE, Value: titleHeight},
				{SizeType: SIZE_ABSOLUTE, Value: artistHeight},
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
			buttonColor = rl.Fade(rl.Black, 0.1)
		} else {
			state := interactable.State()
			switch state {
			case c.STATE_HOT:
				buttonColor = rl.Fade(rl.Black, 0.4)
			case c.STATE_ACTIVE:
				buttonColor = rl.Fade(rl.Black, 0.3)
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

		cardContent.Render(SongCardText(song.Title, 20))
		cardContent.Render(SongCardText(song.Artist, 14))

		next(rect)
	}
}

func SongCardText(text string, fontSize float32) ContrainedComponent {
	return func(rect rl.Rectangle) {
		font := rl.GetFontDefault()
		rl.DrawTextEx(font, text, rl.NewVector2(rect.X-1, rect.Y-1), fontSize, 2, rl.Fade(rl.Black, 0.6))
		rl.DrawTextEx(font, text, rl.NewVector2(rect.X+1, rect.Y+1), fontSize, 2, rl.Fade(rl.Black, 0.6))
		rl.DrawTextEx(font, text, rl.NewVector2(rect.X, rect.Y), fontSize, 2, rl.White)
	}
}
