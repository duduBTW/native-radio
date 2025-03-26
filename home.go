package main

import (
	"fmt"

	rl "github.com/gen2brain/raylib-go/raylib"
)

var mousePoint = rl.Vector2{} // FIX-ME REMOVE DAMN GLOBAL

func HomePage(songTable SongTable) {
	mousePoint = rl.GetMousePosition()

	if UI.HasSelectedSong() {
		selectedSongTexture := UI.SelectedSongTexture()
		rl.DrawTexturePro(
			selectedSongTexture,
			rl.NewRectangle(0, 0, float32(selectedSongTexture.Width), float32(selectedSongTexture.Height)),
			rl.NewRectangle(0, 0, float32(UI.ScreenW), float32(UI.ScreenH)),
			rl.NewVector2(0, 0),
			0,
			rl.Gray,
		)
	}

	padding := Padding{}
	var home = NewConstrainedLayout(ContrainedLayout{
		Contrains: rl.Rectangle{Width: float32(UI.ScreenW), Height: float32(UI.ScreenH), X: 0, Y: 0},
		Direction: DIRECTION_ROW,
		Padding:   padding,
		ChildrenSize: []ChildSize{
			{SizeType: SIZE_ABSOLUTE, Value: 440},
			{SizeType: SIZE_WEIGHT, Value: 1},
		},
	})

	home.Render(Panel(songTable))
	home.Render(SongDetails)
}

func Panel(songTable SongTable) ContrainedComponent {
	return func(rect rl.Rectangle) {
		rl.DrawRectanglePro(rect, rl.Vector2{}, 0, rl.Fade(rl.Black, 0.42))
		padding := Padding{}
		padding.Axis(20, 24)
		var container = NewLayout(Layout{
			Direction: DIRECTION_COLUMN,
			Padding:   padding,
			Gap:       24,
		}, rect)

		container.Render(Tabs)
		// container.Render(Filters)
		container.Render(SongList(songTable))
	}
}

func Tabs(position Position, next Next) {
	rect := position.ToRect(position.Contrains.Width, 42)
	var container = NewConstrainedLayout(ContrainedLayout{
		Direction: DIRECTION_COLUMN,
		Gap:       20,
		Contrains: rect,
		ChildrenSize: []ChildSize{
			{SizeType: SIZE_ABSOLUTE, Value: 24},
			{SizeType: SIZE_WEIGHT, Value: 1},
			{SizeType: SIZE_ABSOLUTE, Value: 40},
		},
	})

	container.Render(PanelRightIcon)

	// rl.DrawRectanglePro(rect, rl.Vector2{}, 0, rl.Fade(rl.Blue, 0.42))
	next(rect)
}

func PanelRightIcon(rect rl.Rectangle) {
	intRect := rect.ToInt32()
	rl.DrawTexture(Textures.PanelRightIcon, intRect.X, intRect.Y, rl.White)
}

func Filters(position Position, next Next) {
	rect := position.ToRect(position.Contrains.Width, 40)
	rl.DrawRectanglePro(rect, rl.Vector2{}, 0, rl.Fade(rl.Blue, 0.42))
	next(rect)
}

func SongList(songTable SongTable) Component {
	return func(position Position, next Next) {
		rect := position.ToRect(position.Contrains.Width, position.Contrains.Height)
		iRect := rect.ToInt32()
		rl.BeginScissorMode(iRect.X-2, iRect.Y-2, iRect.Width+4, iRect.Height+4)
		rect.Y += UI.SidePanelScrollTop

		container := NewLayout(Layout{
			Direction: DIRECTION_COLUMN,
			Gap:       12,
		}, rect)

		for _, song := range songTable.Songs {
			container.Render(SongCard(song))
		}

		// if rl.CheckCollisionPointRec(mousePoint, rect) { // FIX-ME hover detection
		UI.SidePanelScrollTop = clamp(UI.SidePanelScrollTop+rl.GetMouseWheelMove()*float32(scrollSpeed), -container.Size.Height, 0) // FIXME: container.Size.Height IS WRONG
		// }/

		rl.EndScissorMode()
		next(container.Position.ToRect(container.Size.Width, container.Size.Height))
	}
}

func SongCard(song Song) Component {
	return func(position Position, next Next) {
		padding := Padding{}
		padding.Axis(20, 16)
		cardContent := NewLayout(Layout{
			Direction: DIRECTION_COLUMN,
			Padding:   padding,
		}, position.ToRect(position.Contrains.Width, position.Contrains.Height))
		cardContent.Render(SongCardText(song.Title, 20))
		cardContent.Render(SongCardText(song.Artist, 14))

		var rect = position.ToRect(position.Contrains.Width, cardContent.Size.Height)
		var id = "song-card" + song.FileName

		// -- BURN WITH FIRE
		if rl.IsMouseButtonReleased(rl.MouseButtonLeft) {
			UI.ActiveId = ""

			if rl.CheckCollisionPointRec(mousePoint, rect) { // HOW DO I IMUI
				UI.SelectedSong = song
				rl.UnloadMusicStream(UI.Music)
				UI.Music = rl.LoadMusicStream(song.FileName)
			}
		} else if rl.CheckCollisionPointRec(mousePoint, rect) {
			if rl.IsMouseButtonPressed(rl.MouseButtonLeft) {
				UI.ActiveId = id
			} else {
				UI.HotId = id
			}
		} else if UI.HotId == id {
			UI.HotId = ""
		}

		isSelected := UI.SelectedSong.FileName == song.FileName
		buttonColor := rl.Fade(rl.Black, 0.3)
		if isSelected {
			buttonColor = rl.Fade(rl.Black, 0.1)
		} else if UI.ActiveId == id {
			buttonColor = rl.Fade(rl.Black, 0.1)
		} else if UI.HotId == id {
			buttonColor = rl.Fade(rl.Black, 0.2)
		}
		// ----------------

		texture := Textures.Songs[song.FileName]
		rl.DrawTexturePro(
			texture,
			rl.NewRectangle(0, 0, float32(texture.Width), float32(texture.Height)),
			rect,
			rl.NewVector2(0, 0),
			0,
			rl.White,
		)
		rl.DrawRectanglePro(rect, rl.Vector2{}, 0, buttonColor)

		if isSelected {
			rl.DrawRectangleRoundedLinesEx(rect, 0, 0, 2, rl.White)
		}

		// FIXME: FIGURE OUT A BETTER WAY TO DO THIS
		// WE ARE RENDERING AGAIN ON TOP OF THE RECT WITH THE SAME POSITION NOW
		cardContent2 := NewLayout(Layout{
			Direction: DIRECTION_COLUMN,
			Padding:   padding,
		}, position.ToRect(position.Contrains.Width, position.Contrains.Height))
		cardContent2.Render(SongCardText(song.Title, 20))
		cardContent2.Render(SongCardText(song.Artist, 14))

		next(rect)
	}
}

func SongCardText(text string, fontSize float32) Component {
	return func(position Position, next Next) {
		font := rl.GetFontDefault()
		textHeight := DrawTextByWidth(font, text, rl.NewVector2(position.X, position.Y), position.Contrains.Width, fontSize, 2, rl.White)
		next(position.ToRect(position.Contrains.Width, textHeight))
	}
}

func SongDetails(rect rl.Rectangle) {
	// rl.DrawRectanglePro(rect, rl.Vector2{}, 0, rl.Fade(rl.Black, 0.42))

	padding := Padding{}
	padding.Axis(94, 126)
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
	container.Render(SongDetail)
}

func SongMiniature(rect rl.Rectangle) {
	if !UI.HasSelectedSong() {
		return
	}

	selectedSongTexture := UI.SelectedSongTexture()
	var size float32 = 300
	rl.DrawTexturePro(
		selectedSongTexture,
		rl.NewRectangle(0, 0, float32(selectedSongTexture.Width), float32(selectedSongTexture.Height)),
		rl.NewRectangle(rect.X+((rect.Width-size)/2), rect.Y+((rect.Height-size)/2), size, size),
		rl.NewVector2(0, 0),
		0,
		rl.White,
	)
}

func SongDetail(rect rl.Rectangle) {
	rl.DrawRectanglePro(rect, rl.Vector2{}, 0, rl.Fade(rl.Blue, 0.42))
}

// func SongDetails(rect rl.Rectangle, next func(position Position)) {
// 	rl.DrawRectangle(int32(rect.X), int32(rect.Y), int32(rect.Width), int32(rect.Height), rl.Fade(rl.Red, 0.5))
// 	next(Position{X: int(rect.X + rect.Width), Y: int(rect.Y)})

// var container = Layout{
// 	Direction:  DIRECTION_COLUMN,
// 	Padding:    Padding{},
// 	Contrains:  rect,
// 	currentPos: Position{X: int(rect.X), Y: int(rect.Y)},
// 	ChildrenSize: []ChildSize{
// 		{SizeType: SIZE_WEIGHT, Value: 1},
// 		{SizeType: SIZE_ABSOLUTE, Value: 184},
// 	},
// }
// 	container.ComputeChildren()

// 	SongImage(container.Current(), container.Next)
// 	SongControls(container.Current(), container.Next)
// }

// func SongImage(rect rl.Rectangle, next func(position Position)) {
// 	rl.DrawRectangle(int32(rect.X), int32(rect.Y), int32(rect.Width), int32(rect.Height), rl.Fade(rl.Blue, 0.5))
// 	next(Position{Y: int(rect.Y + rect.Height), X: int(rect.X)})
// }
// func SongControls(rect rl.Rectangle, next func(position Position)) {
// 	rl.DrawRectangle(int32(rect.X), int32(rect.Y), int32(rect.Width), int32(rect.Height), rl.Fade(rl.Pink, 0.5))
// 	next(Position{Y: int(rect.Y + rect.Height), X: int(rect.X)})
// }

// func Tabs(rect rl.Rectangle, next func(position Position)) {
// 	container := NewConstrainedLayout(ContrainedLayout{
// 		Direction: DIRECTION_ROW,
// 		Gap:       20,
// 		Contrains: rect,
// 		ChildrenSize: []ChildSize{
// 			{SizeType: SIZE_ABSOLUTE, Value: 20},
// 			{SizeType: SIZE_WEIGHT, Value: 1},
// 			{SizeType: SIZE_ABSOLUTE, Value: 40},
// 		},
// 	})

// 	CollapseSidebarButton(container.Current(), container.Next)
// 	TabsContent(container.Current(), container.Next)
// 	Settings(container.Current(), container.Next)

// 	rl.DrawRectangle(int32(rect.X), int32(rect.Y), int32(rect.Width), int32(container.MaxH), rl.Fade(rl.Purple, 0.5))
// 	next(Position{Y: int(rect.Y) + container.MaxH, X: int(rect.X)})
// }

// func CollapseSidebarButton(rect rl.Rectangle, next func(position Position)) {
// 	rl.DrawRectangle(int32(rect.X), int32(rect.Y), int32(rect.Width), int32(rect.Width), rl.Fade(rl.Red, 0.5))
// 	next(Position{X: int(rect.X + rect.Width), Y: int(rect.Y), H: int(rect.Width)})
// }

// func TabsContent(rect rl.Rectangle, next func(position Position)) {
// 	rl.DrawRectangle(int32(rect.X), int32(rect.Y), int32(rect.Width), 70, rl.Fade(rl.Blue, 0.5))
// 	next(Position{X: int(rect.X + rect.Width), Y: int(rect.Y), H: 70})
// }

// func Settings(rect rl.Rectangle, next func(position Position)) {
// 	rl.DrawRectangle(int32(rect.X), int32(rect.Y), int32(rect.Width), int32(rect.Width), rl.Fade(rl.Green, 0.5))
// 	next(Position{X: int(rect.X + rect.Width), Y: int(rect.Y), H: int(rect.Width)})
// }

// func Filters(rect rl.Rectangle, next func(position Position)) {
// 	rl.DrawRectangle(int32(rect.X), int32(rect.Y), int32(rect.Width), 60, rl.Fade(rl.DarkBrown, 0.5))
// 	next(Position{Y: int(rect.Y + 60), X: int(rect.X)})
// }

// --- OLD

// 	Tabs() // FIXME: We should pass a function as argument here, this function will update the size of the layout
// 	Filters(container.Next())
// 	SongList(container.Next())
// }

func HomePageOld(songTable SongTable) {
	mousePoint := rl.GetMousePosition()
	RenderSelectedSong()
	RenderList(songTable, mousePoint)

	if rl.IsKeyPressed(rl.KeySpace) {
		rl.PlayMusicStream(UI.Music)
	}
}

func RenderSelectedSong() {
	if UI.SelectedSong.FileName == "" {
		return
	}

	const size float32 = 350
	// var x float32 = ((float32(UI.ScreenW) - 200 - size) / 2) + 200
	// var y float32 = (float32(UI.ScreenH) - size) / 2
	// var imgW float32 = float32(UI.SelectedSongTexture.Width)
	// var imgH float32 = float32(UI.SelectedSongTexture.Height)

	// // background
	// rl.DrawTexturePro(
	// 	UI.SelectedSongTexture,
	// 	rl.NewRectangle(0, 0, imgW, imgH),
	// 	rl.NewRectangle(0, 0, float32(UI.ScreenW), float32(UI.ScreenH)),
	// 	rl.NewVector2(0, 0),
	// 	0,
	// 	rl.Gray,
	// )

	// // middle img
	// rl.DrawTexturePro(
	// 	UI.SelectedSongTexture,
	// 	rl.NewRectangle(0, 0, imgW, imgH),
	// 	rl.NewRectangle(x, y, size, size),
	// 	rl.NewVector2(0, 0),
	// 	0,
	// 	rl.Gray,
	// )
}

var scrollSpeed = 4

func getItemRectByIndex(i int, rowGap int) float32 {
	return float32(4 + (40 * i) + (rowGap * i))
}

func RenderList(songTable SongTable, mousePoint rl.Vector2) {
	rowGap := 8

	if rl.CheckCollisionPointRec(mousePoint, rl.Rectangle{X: 0, Y: 0, Width: 200, Height: 450}) {
		UI.SidePanelScrollTop += rl.GetMouseWheelMove() * float32(scrollSpeed)
		fmt.Println(UI.SidePanelScrollTop)
	}

	for i, song := range songTable.Songs {
		rect := rl.Rectangle{X: 4, Y: getItemRectByIndex(i, rowGap) + UI.SidePanelScrollTop, Width: 200, Height: 40}
		if rl.CheckCollisionPointRec(mousePoint, rect) && rl.IsMouseButtonPressed(rl.MouseButtonLeft) {
			// rl.UnloadTexture(UI.SelectedSongTexture)
			// UI.SelectedSongTexture = rl.LoadTexture(song.Background)
			UI.SelectedSong = song
			rl.UnloadMusicStream(UI.Music)
			UI.Music = rl.LoadMusicStream(song.FileName)
		}

		rectInt32 := rect.ToInt32()
		rl.DrawRectangle(rectInt32.X, rectInt32.Y, rectInt32.Width, rectInt32.Height, rl.Fade(rl.LightGray, 0.5))
		if UI.SelectedSong.FileName == song.FileName {
			rl.DrawRectangleRoundedLinesEx(rect, 0, 0, 1, rl.White)
		}

		i++
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
