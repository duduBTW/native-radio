package main

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

func HomePage(songTable SongTable) {
	mousePoint := rl.GetMousePosition()
	RenderList(songTable, mousePoint)
	RenderSelectedSong()

	if rl.IsKeyPressed(rl.KeySpace) {
		rl.PlayMusicStream(UI.Music)
	}
}

func RenderSelectedSong() {
	if UI.SelectedSong.FileName == "" {
		return
	}

	rl.DrawText(UI.SelectedSong.FileName, 220, 200, 14, rl.Black)
}

var scrollSpeed = 4

func RenderList(songTable SongTable, mousePoint rl.Vector2) {
	if rl.CheckCollisionPointRec(mousePoint, rl.Rectangle{X: 0, Y: 0, Width: 200, Height: 450}) {
		UI.SidePanelScrollTop -= (rl.GetMouseWheelMove() * float32(scrollSpeed) * -1)
	}

	for i, song := range songTable.Songs {
		rowGap := 8
		rect := rl.Rectangle{X: 4, Y: float32(4+(40*i)+(rowGap*i)) + UI.SidePanelScrollTop, Width: 200, Height: 40}
		if rl.CheckCollisionPointRec(mousePoint, rect) && rl.IsMouseButtonPressed(rl.MouseButtonLeft) {
			UI.SelectedSong = song
			UI.Music = rl.LoadMusicStream(song.FileName)
		}

		rectInt32 := rect.ToInt32()
		rl.DrawRectangle(rectInt32.X, rectInt32.Y, rectInt32.Width, rectInt32.Height, rl.Fade(rl.LightGray, 0.5))
		i++
	}
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
