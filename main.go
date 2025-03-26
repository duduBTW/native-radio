package main

import (
	"fmt"
	"os"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type Page string

const (
	PAGE_HOME         Page = "home"
	PAGE_SETUP_WIZARD Page = "setup-wizard"
)

type UITextures struct {
	PanelRightIcon rl.Texture2D
	Songs          map[string]rl.Texture2D
}

type UIStruct struct {
	SelectedSong       Song
	Music              rl.Music
	SidePanelScrollTop float32
	ScreenW            int32
	ScreenH            int32

	ActiveId string
	HotId    string
}

func (ui UIStruct) SelectedSongTexture() rl.Texture2D {
	return Textures.Songs[ui.SelectedSong.FileName]
}

func (ui UIStruct) HasSelectedSong() bool {
	return ui.SelectedSong.FileName != ""
}

var UI = UIStruct{}
var Textures = UITextures{
	Songs: make(map[string]rl.Texture2D),
}

func LoadTextures(songs SongTable) {
	for _, song := range songs.Songs {
		Textures.Songs[song.FileName] = rl.LoadTexture(song.Background)
	}
	Textures.PanelRightIcon = rl.LoadTexture("/Users/carloseduardoalvesdonascimento/Personal/osu-song-native/icons/panel-right.svg")
}

func main() {
	fmt.Println("Process ID:", os.Getpid())

	rl.SetConfigFlags(rl.FlagWindowResizable)
	rl.InitWindow(800, 600, "raylib [core] example - basic window")
	defer rl.CloseWindow()

	rl.InitAudioDevice()
	rl.SetTargetFPS(60)

	var currentPage Page = PAGE_HOME
	var songs = SongTable{}
	songs.FromJson("/Users/carloseduardoalvesdonascimento/Personal/osu-song-native/songs.json")
	LoadTextures(songs)

	for !rl.WindowShouldClose() {
		rl.UpdateMusicStream(UI.Music)

		UI.ScreenW = int32(rl.GetScreenWidth())
		UI.ScreenH = int32(rl.GetScreenHeight())

		rl.BeginDrawing()
		rl.ClearBackground(rl.RayWhite)
		switch currentPage {
		case PAGE_HOME:
			HomePage(songs)
		case PAGE_SETUP_WIZARD:
			SetupWizardPage()
		}

		rl.DrawFPS(10, 10)
		rl.EndDrawing()
	}
}

func SetupWizardPage() {
	rl.DrawText("Setup wizard", 190, 200, 20, rl.LightGray)
}
