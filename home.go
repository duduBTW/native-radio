package main

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

var scrollSpeed = 40

func HomePage() {
	if UI.HasSelectedSong() {
		rl.BeginShaderMode(Shaders.Blur.Shader)
		DrawFitImage(UI.SelectedSongTexture(), rl.NewRectangle(0, 0, float32(UI.ScreenW), float32(UI.ScreenH)), rl.Gray)
		rl.EndShaderMode()

		rl.DrawRectangle(0, 0, UI.ScreenW, UI.ScreenH, rl.NewColor(18, 18, 18, 209))
	}

	padding := Padding{}
	var home = NewConstrainedLayout(ContrainedLayout{
		Contrains: rl.Rectangle{Width: float32(UI.ScreenW), Height: float32(UI.ScreenH), X: 0, Y: 0},
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
				{SizeType: SIZE_WEIGHT, Value: 1},
			},
		})

		container.Render(UpperPart)

		switch UI.SelectedPanelPage {
		case PANEL_PAGE_SONGS:
			container.Render(SongList())

		}
		// container.Render(Filters)
	}
}

func UpperPart(rect rl.Rectangle) {
	var container = NewConstrainedLayout(ContrainedLayout{
		Direction: DIRECTION_ROW,
		Gap:       16,
		Contrains: rect,
		ChildrenSize: []ChildSize{
			{SizeType: SIZE_ABSOLUTE, Value: ICON_BUTTON_SIZE_GHOST},
			{SizeType: SIZE_WEIGHT, Value: 1},
			{SizeType: SIZE_ABSOLUTE, Value: ICON_BUTTON_SIZE_GHOST},
		},
	})

	container.Render(PanelSidebarButton)
	container.Render(UpperPartTabs())
	container.Render(PanelSettingsButton)
}

func UpperPartTabs() ContrainedComponent {
	return func(rect rl.Rectangle) {
		value := Tabs(TabsProps{
			Items: []TabsItemProps{
				{
					Icon:  ICON_REPEAT,
					Label: "Songs",
					Value: "songs",
				},
				{
					Icon:  ICON_VOLUME_MUTED,
					Label: "Playlists",
					Value: "playlists",
				},
			},
			Rect:  rect,
			Value: string(UI.SelectedPanelPage),
		})
		UI.SelectedPanelPage = PanelPage(value)
	}
}

func PanelSidebarButton(rect rl.Rectangle) {
	IconButton("sidebar-collapse", ICON_SIDEBAR, ICON_BUTTON_GHOST, rl.NewRectangle(rect.X, rect.Y+5, rect.Width, rect.Height))
}

func PanelSettingsButton(rect rl.Rectangle) {
	var variant IconButtonVariant = ICON_BUTTON_GHOST
	if UI.SelectedPanelPage == PANEL_PAGE_SETTINGS {
		variant = ICON_BUTTON_SECONDARY
	}

	if IconButton("settings-button", ICON_SETTINGS, variant, rl.NewRectangle(rect.X, rect.Y+5, rect.Width, rect.Height)) {
		UI.SelectedPanelPage = PANEL_PAGE_SETTINGS
	}
}

func Filters(position Position, next Next) {
	rect := position.ToRect(position.Contrains.Width, 40)
	rl.DrawRectanglePro(rect, rl.Vector2{}, 0, rl.Fade(rl.Blue, 0.42))
	next(rect)
}

func SongList() ContrainedComponent {
	return func(rect rl.Rectangle) {
		iRect := rect.ToInt32()
		rl.BeginScissorMode(iRect.X-2, iRect.Y-2, iRect.Width+4, iRect.Height+4)

		rectWithOffset := rl.NewRectangle(rect.X, rect.Y+UI.SidePanelScrollTop, rect.Width, rect.Height)

		container := NewLayout(Layout{
			Direction: DIRECTION_COLUMN,
			Gap:       12,
		}, rectWithOffset)

		for index, song := range UI.Songs {
			container.Render(SongCard(song, index))
		}

		if rl.CheckCollisionPointRec(MousePoint, rect) { // FIX-ME hover detection
			UI.SidePanelScrollTop = clamp(UI.SidePanelScrollTop+rl.GetMouseWheelMove()*float32(scrollSpeed), -(container.Size.Height - rect.Height), 0) // FIXME: container.Size.Height IS WRONG
		}

		rl.EndScissorMode()
	}
}

var titleHeight float32 = 30
var artistHeight float32 = 21

func SongCard(song Song, index int) Component {
	return func(position Position, next Next) {
		var height float32 = 72
		var rect = position.ToRect(position.Contrains.Width, height)

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

		var id = "song-card" + song.FileName

		// -- BURN WITH FIRE
		if rl.IsMouseButtonReleased(rl.MouseButtonLeft) {
			// UI.ActiveId = ""

			if rl.CheckCollisionPointRec(MousePoint, rect) { // HOW DO I IMUI
				UI.SelectSong(index)
			}
		} else if rl.CheckCollisionPointRec(MousePoint, rect) {
			if rl.IsMouseButtonPressed(rl.MouseButtonLeft) {
				UI.ActiveId = id
			} else {
				UI.HotId = id
			}
		} else if UI.HotId == id {
			UI.HotId = ""
		}

		isSelected := UI.SelectedSong().FileName == song.FileName
		buttonColor := rl.Fade(rl.Black, 0.42)
		if isSelected {
			buttonColor = rl.Fade(rl.Black, 0.1)
		} else if UI.ActiveId == id {
			buttonColor = rl.Fade(rl.Black, 0.1)
		} else if UI.HotId == id {
			buttonColor = rl.Fade(rl.Black, 0.2)
		}
		// ----------------
		// FIXME
		// DrawFitImage(Textures.Songs[song.FileName], rect, rl.White)
		rl.DrawRectanglePro(rect, rl.Vector2{}, 0, buttonColor)

		if isSelected {
			rl.DrawRectangleRoundedLinesEx(rect, 0, 0, 2, rl.White)
		}

		cardContent.Render(SongCardText(song.Title, 20))
		cardContent.Render(SongCardText(song.Artist, 14))

		next(rect)
	}
}

func SongCardText(text string, fontSize float32) ContrainedComponent {
	return func(rect rl.Rectangle) {
		font := rl.GetFontDefault()
		rl.DrawTextEx(font, text, rl.NewVector2(rect.X, rect.Y), fontSize, 2, rl.White)
	}
}

func clamp(value, min, max float32) float32 {
	if value > max {
		return max
	}

	if value < min {
		return min
	}

	return value
}

// func RenderList(songTable SongTable) {
// 	for i, song := range songTable.Songs {
// 		y := int32(40 * i)
// 		var fontSize int32 = 20
// 		rl.DrawRectangle(0, y, rl.MeasureText(song.FileName, fontSize)+16, 40, rl.Fade(rl.LightGray, 0.5))
// 		rl.DrawText(song.FileName, 8, y+10, fontSize, rl.Black)
// 		i++
// 	}
// }
