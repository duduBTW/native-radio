package main

import (
	"fmt"
	"os"

	_ "net/http/pprof" // Import pprof

	"github.com/dudubtw/osu-radio-native/lib"
	rl "github.com/gen2brain/raylib-go/raylib"
)

var SongTable lib.SongTable
var music lib.Music
var Textures lib.Textures
var UI lib.UIStruct
var Shaders lib.Shaders
var MousePoint rl.Vector2

func UpdateSong() {
	music.LoadMusic(&SongTable)
	Textures.LoadSelectedSong(SongTable, Shaders)
}
func SelectSong(songIndex int) {
	SongTable.SelectSong(songIndex)
	UpdateSong()
}

func InitEverything() {
	// Songs
	table, err := lib.NewSongTableFromJson("C:\\Users\\carlo\\AppData\\Roaming\\osu-radio\\storage\\songs.json")
	if err != nil {
		fmt.Println("Could not load songs!")
		panic(1)
	}

	SongTable = *table
	UI = lib.NewUi()
	Textures = lib.NewTexture(table)
	Shaders = lib.NewShaders()

	SelectSong(0)
}

func main() {
	fmt.Println("Process ID:", os.Getpid())

	rl.SetConfigFlags(rl.FlagWindowResizable)
	rl.SetConfigFlags(rl.FlagBorderlessWindowedMode)
	rl.InitWindow(1280, 720, "raylib [core] example - basic window")
	defer rl.CloseWindow()

	rl.InitAudioDevice()
	rl.SetTargetFPS(165)

	InitEverything()

	for !rl.WindowShouldClose() {
		if music.Selected != nil {
			rl.UpdateMusicStream(*music.Selected)
		}

		if !music.IsSkeekMode && music.HasEnded() {
			music.Next(&SongTable)
			UpdateSong()
		}

		UI.ScreenW = int32(rl.GetScreenWidth())
		UI.ScreenH = int32(rl.GetScreenHeight())
		MousePoint = rl.GetMousePosition()

		rl.SetShaderValue(Shaders.Blur.Shader, Shaders.Blur.ScreenResLoc, []float32{float32(UI.ScreenW), float32(UI.ScreenH)}, rl.ShaderUniformVec2)
		rl.SetShaderValue(Shaders.Blur.Shader, Shaders.Blur.MouseLoc, []float32{MousePoint.X, float32(UI.ScreenH) - MousePoint.Y}, rl.ShaderUniformVec2)

		rl.BeginDrawing()
		rl.ClearBackground(rl.RayWhite)

		switch UI.SelectedPage {
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
