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

type UIStruct struct {
	SelectedSong       Song
	Music              rl.Music
	SidePanelScrollTop float32
}

var UI = UIStruct{}

func main() {
	fmt.Println("Process ID:", os.Getpid())

	rl.InitWindow(1200, 450, "raylib [core] example - basic window")
	defer rl.CloseWindow()

	rl.InitAudioDevice()
	rl.SetTargetFPS(60)

	var currentPage Page = PAGE_HOME
	var songs = SongTable{}
	songs.FromJson("/Users/carloseduardoalvesdonascimento/Personal/osu-song-native/songs.json")

	for !rl.WindowShouldClose() {
		rl.UpdateMusicStream(UI.Music)

		rl.BeginDrawing()
		rl.ClearBackground(rl.RayWhite)
		switch currentPage {
		case PAGE_HOME:
			HomePage(songs)
		case PAGE_SETUP_WIZARD:
			SetupWizardPage()
		}

		rl.EndDrawing()
	}
}

func SetupWizardPage() {
	rl.DrawText("Setup wizard", 190, 200, 20, rl.LightGray)
}
