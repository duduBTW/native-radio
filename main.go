package main

import (
	"fmt"
	"os"

	"github.com/dudubtw/osu-radio-native/lib"
	rl "github.com/gen2brain/raylib-go/raylib"
)

var songTable lib.SongTable
var music lib.Music
var textures lib.Textures
var ui lib.UIStruct
var shaders lib.Shaders
var mousePoint rl.Vector2

func SetVolume(newVolume float32, wasMuteClicked bool) {
	ui.Volume = newVolume
	if wasMuteClicked {
		ui.IsMuted = !ui.IsMuted
	}

	if ui.IsMuted {
		music.SetVolume(0)
	} else {
		music.SetVolume(newVolume)
	}
}
func UpdateSong() {
	music.LoadMusic(&songTable)
	textures.LoadSelectedSong(songTable, shaders)
}
func SelectSong(songIndex int) {
	songTable.SelectSong(songIndex)
	UpdateSong()
}

func InitEverything() {
	table, err := lib.NewSongTableFromJson("C:\\Users\\carlo\\AppData\\Roaming\\osu-radio\\storage\\songs.json")
	if err != nil {
		fmt.Println("Could not load songs!")
		panic(1)
	}

	songTable = *table
	ui = lib.NewUi()
	textures = lib.NewTexture(table)
	shaders = lib.NewShaders()

	SelectSong(0)
}

func main() {
	fmt.Println("Process ID:", os.Getpid())

	rl.SetConfigFlags(rl.FlagWindowResizable)
	rl.SetConfigFlags(rl.FlagBorderlessWindowedMode)
	rl.InitWindow(1400, 1000, "osu! radio")
	defer rl.CloseWindow()

	rl.InitAudioDevice()
	rl.SetTargetFPS(165)

	InitEverything()

	for !rl.WindowShouldClose() {
		if music.Selected != nil {
			rl.UpdateMusicStream(*music.Selected)
		}

		if !music.IsSkeekMode && music.HasEnded() {
			music.Next(&songTable)
			UpdateSong()
			music.Play()
		}

		ui.ScreenW = int32(rl.GetScreenWidth())
		ui.ScreenH = int32(rl.GetScreenHeight())
		mousePoint = rl.GetMousePosition()

		shaders.Blur.Update(mousePoint, &ui)

		rl.BeginDrawing()
		rl.ClearBackground(rl.RayWhite)

		switch ui.SelectedPage {
		case lib.PAGE_HOME:
			HomePage()
		case lib.PAGE_SETUP_WIZARD:
			SetupWizardPage()
		}

		rl.DrawFPS(10, 10)
		rl.EndDrawing()
	}
}

func SetupWizardPage() {
	rl.DrawText("Setup wizard", 190, 200, 20, rl.LightGray)
}
