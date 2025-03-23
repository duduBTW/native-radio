package main

import (
	"fmt"

	rl "github.com/gen2brain/raylib-go/raylib"
)

// TODO -
// 1. Position elements.
func HomePage(songTable SongTable) {
	var home = Layout{
		Contrains: rl.Rectangle{Width: float32(UI.ScreenW), Height: float32(UI.ScreenH), X: 0, Y: 0},
		Direction: DIRECTION_ROW,
		ChildrenSize: []ChildSize{
			{SizeType: SIZE_ABSOLUTE, Value: 400},
			{SizeType: SIZE_WEIGHT, Value: 1},
		},
	}
	home.ComputeChildren() // FIXME: ERROR HANDLING

	Panel(home.Current(), home.Next, songTable)
	SongDetails(home.Current(), home.Next)
	// SongDetails(home.Next())
}

func SongDetails(rect rl.Rectangle, next func(position Position)) {
	rl.DrawRectangle(int32(rect.X), int32(rect.Y), int32(rect.Width), int32(rect.Height), rl.Fade(rl.Red, 0.5))
	next(Position{X: int(rect.X + rect.Width), Y: int(rect.Y)})

	var container = Layout{
		Direction:  DIRECTION_COLUMN,
		Padding:    Padding{},
		Contrains:  rect,
		currentPos: Position{X: int(rect.X), Y: int(rect.Y)},
		ChildrenSize: []ChildSize{
			{SizeType: SIZE_WEIGHT, Value: 1},
			{SizeType: SIZE_ABSOLUTE, Value: 184},
		},
	}
	container.ComputeChildren()

	SongImage(container.Current(), container.Next)
	SongControls(container.Current(), container.Next)
}

func SongImage(rect rl.Rectangle, next func(position Position)) {
	rl.DrawRectangle(int32(rect.X), int32(rect.Y), int32(rect.Width), int32(rect.Height), rl.Fade(rl.Blue, 0.5))
	next(Position{Y: int(rect.Y + rect.Height), X: int(rect.X)})
}
func SongControls(rect rl.Rectangle, next func(position Position)) {
	rl.DrawRectangle(int32(rect.X), int32(rect.Y), int32(rect.Width), int32(rect.Height), rl.Fade(rl.Pink, 0.5))
	next(Position{Y: int(rect.Y + rect.Height), X: int(rect.X)})
}

func Panel(rect rl.Rectangle, next func(position Position), songTable SongTable) {
	rl.DrawRectangle(int32(rect.X), int32(rect.Y), int32(rect.Width), int32(rect.Height), rl.Fade(rl.LightGray, 0.5))
	next(Position{X: int(rect.X + rect.Width), Y: int(rect.Y)})

	var container = Layout{
		Direction:  DIRECTION_COLUMN,
		Padding:    Padding{},
		Gap:        24,
		Contrains:  rect,
		currentPos: Position{X: int(rect.X), Y: int(rect.Y)},
	}
	Tabs(container.Current(), container.Next)
	Filters(container.Current(), container.Next)
	SongList(container.Current(), container.Next, songTable)
}
func Tabs(rect rl.Rectangle, next func(position Position)) {
	container := Layout{
		Direction: DIRECTION_ROW,
		Gap:       20,
		Contrains: rect,
		ChildrenSize: []ChildSize{
			{SizeType: SIZE_ABSOLUTE, Value: 20},
			{SizeType: SIZE_WEIGHT, Value: 1},
			{SizeType: SIZE_ABSOLUTE, Value: 40},
		},
	}
	container.ComputeChildren()

	CollapseSidebarButton(container.Current(), container.Next)
	TabsContent(container.Current(), container.Next)
	Settings(container.Current(), container.Next)

	rl.DrawRectangle(int32(rect.X), int32(rect.Y), int32(rect.Width), int32(container.MaxH), rl.Fade(rl.Purple, 0.5))
	next(Position{Y: int(rect.Y) + container.MaxH, X: int(rect.X)})
}

func CollapseSidebarButton(rect rl.Rectangle, next func(position Position)) {
	rl.DrawRectangle(int32(rect.X), int32(rect.Y), int32(rect.Width), int32(rect.Width), rl.Fade(rl.Red, 0.5))
	next(Position{X: int(rect.X + rect.Width), Y: int(rect.Y), H: int(rect.Width)})
}

func TabsContent(rect rl.Rectangle, next func(position Position)) {
	rl.DrawRectangle(int32(rect.X), int32(rect.Y), int32(rect.Width), 70, rl.Fade(rl.Blue, 0.5))
	next(Position{X: int(rect.X + rect.Width), Y: int(rect.Y), H: 70})
}

func Settings(rect rl.Rectangle, next func(position Position)) {
	rl.DrawRectangle(int32(rect.X), int32(rect.Y), int32(rect.Width), int32(rect.Width), rl.Fade(rl.Green, 0.5))
	next(Position{X: int(rect.X + rect.Width), Y: int(rect.Y), H: int(rect.Width)})
}

func Filters(rect rl.Rectangle, next func(position Position)) {
	rl.DrawRectangle(int32(rect.X), int32(rect.Y), int32(rect.Width), 60, rl.Fade(rl.DarkBrown, 0.5))
	next(Position{Y: int(rect.Y + 60), X: int(rect.X)})
}

func SongList(rect rl.Rectangle, next func(position Position), songTable SongTable) {
	container := Layout{
		Direction:  DIRECTION_COLUMN,
		Contrains:  rect,
		currentPos: Position{X: int(rect.X), Y: int(rect.Y)},
		Gap:        16,
	}

	for _, song := range songTable.Songs {
		var rect = container.Current()
		var itemHeight = 70
		rl.DrawRectangle(int32(rect.X), int32(rect.Y), int32(rect.Width), int32(itemHeight), rl.Fade(rl.Green, 0.5))
		rl.DrawText(song.FileName, 8, int32(rect.Y)+(int32(itemHeight)/2), 8, rl.Black)
		container.Next(Position{Y: int(rect.Y) + itemHeight, X: int(rect.X)})
	}

	next(Position{})
}

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
	var x float32 = ((float32(UI.ScreenW) - 200 - size) / 2) + 200
	var y float32 = (float32(UI.ScreenH) - size) / 2
	var imgW float32 = float32(UI.SelectedSongTexture.Width)
	var imgH float32 = float32(UI.SelectedSongTexture.Height)

	// background
	rl.DrawTexturePro(
		UI.SelectedSongTexture,
		rl.NewRectangle(0, 0, imgW, imgH),
		rl.NewRectangle(0, 0, float32(UI.ScreenW), float32(UI.ScreenH)),
		rl.NewVector2(0, 0),
		0,
		rl.Gray,
	)

	// middle img
	rl.DrawTexturePro(
		UI.SelectedSongTexture,
		rl.NewRectangle(0, 0, imgW, imgH),
		rl.NewRectangle(x, y, size, size),
		rl.NewVector2(0, 0),
		0,
		rl.Gray,
	)
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
			rl.UnloadTexture(UI.SelectedSongTexture)
			UI.SelectedSongTexture = rl.LoadTexture(song.Background)
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
