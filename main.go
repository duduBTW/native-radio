package main

import (
	"fmt"
	"net/http"
	"os"

	_ "net/http/pprof" // Import pprof

	rl "github.com/gen2brain/raylib-go/raylib"
)

type Page string

const (
	PAGE_HOME         Page = "home"
	PAGE_SETUP_WIZARD Page = "setup-wizard"
)

type PanelPage string

const (
	PANEL_PAGE_SONGS     PanelPage = "songs"
	PANEL_PAGE_PLAYLISTS PanelPage = "playlists"
	PANEL_PAGE_SETTINGS  PanelPage = "settings"
)

type UITextures struct {
	Icons rl.Texture2D
	Songs map[string]rl.Texture2D
}

type UIStruct struct {
	selectedSongindex  int
	Music              rl.Music
	SidePanelScrollTop float32
	ScreenW            int32
	ScreenH            int32

	SelectedPanelPage PanelPage

	ActiveId string
	HotId    string

	Songs []Song

	seekModeProgress     float32
	isSkeekMode          bool
	wasPlayingBeforeSeek bool
}

func (ui *UIStruct) BeginSeekMode() {
	if ui.isSkeekMode {
		fmt.Println("uh oh")
		return
	}

	ui.wasPlayingBeforeSeek = rl.IsMusicStreamPlaying(ui.Music)
	rl.PauseMusicStream(ui.Music)
	ui.seekModeProgress = ui.Progress()
	ui.isSkeekMode = true
}
func (ui *UIStruct) ExitSeekMode() {
	if !ui.isSkeekMode {
		fmt.Println("uh oh")
		return
	}

	ui.isSkeekMode = false
	rl.SeekMusicStream(UI.Music, rl.GetMusicTimeLength(ui.Music)*ui.seekModeProgress)

	if ui.wasPlayingBeforeSeek {
		rl.PlayMusicStream(ui.Music)
	}
}

func (ui *UIStruct) Play() {
	rl.PlayMusicStream(ui.Music)
}

func (ui *UIStruct) Next() {
	ui.SelectSong(MinInt(ui.selectedSongindex+1, len(ui.Songs)-1))
	ui.Play()
}

func (ui *UIStruct) Progress() float32 {
	if ui.isSkeekMode {
		return ui.seekModeProgress
	}

	return rl.GetMusicTimePlayed(ui.Music) / rl.GetMusicTimeLength(ui.Music)
}

func (ui *UIStruct) Seek(to float32) {
	ui.seekModeProgress = to
}

func (ui *UIStruct) HasEnded() bool {
	if !rl.IsMusicValid(ui.Music) {
		return false
	}

	// TODO deal with ms, since we are doing int() it only compare seconds
	return int(rl.GetMusicTimePlayed(ui.Music)) == int(rl.GetMusicTimeLength(ui.Music))
}

func (ui *UIStruct) Previous() {
	ui.SelectSong(MaxInt(ui.selectedSongindex-1, 0))
	ui.Play()
}

func (ui UIStruct) SelectedSong() Song {
	return ui.Songs[ui.selectedSongindex]
}

func (ui *UIStruct) SelectSong(songIndex int) {
	ui.selectedSongindex = songIndex
	rl.UnloadMusicStream(UI.Music)
	ui.Music = rl.LoadMusicStream(ui.SelectedSong().FileName)
	ui.Music.Looping = false
}

func (ui UIStruct) SelectedSongTexture() rl.Texture2D {
	return Textures.Songs[ui.SelectedSong().FileName]
}

func (ui UIStruct) HasSelectedSong() bool {
	return ui.SelectedSong().FileName != ""
}

var UI = UIStruct{}
var Textures = UITextures{
	Songs: make(map[string]rl.Texture2D),
}

func LoadTextures(songs SongTable) {
	for _, song := range songs.Songs {
		Textures.Songs[song.FileName] = rl.LoadTexture(song.Background)
	}
	Textures.Icons = rl.LoadTexture("/Users/carloseduardoalvesdonascimento/Personal/osu-song-native/sprites/icon-sprite.png")
}

func InitEverything() {
	// Songs
	var songTable = SongTable{}
	songTable.FromJson("/Users/carloseduardoalvesdonascimento/Personal/osu-song-native/songs.json")
	UI.Songs = songTable.Songs
	UI.SelectedPanelPage = PANEL_PAGE_SONGS

	// Textures
	LoadTextures(songTable)

	UI.SelectSong(0)
}

func main() {
	go func() {
		http.ListenAndServe("localhost:6060", nil)
	}()

	fmt.Println("Process ID:", os.Getpid())

	rl.SetConfigFlags(rl.FlagWindowResizable)
	rl.SetConfigFlags(rl.FlagBorderlessWindowedMode)
	rl.InitWindow(800, 600, "raylib [core] example - basic window")
	defer rl.CloseWindow()

	rl.InitAudioDevice()
	rl.SetTargetFPS(60)

	InitEverything()

	var currentPage Page = PAGE_HOME

	for !rl.WindowShouldClose() {
		rl.UpdateMusicStream(UI.Music)

		if !UI.isSkeekMode && UI.HasEnded() {
			UI.Next()
		}

		UI.ScreenW = int32(rl.GetScreenWidth())
		UI.ScreenH = int32(rl.GetScreenHeight())

		rl.BeginDrawing()
		rl.ClearBackground(rl.RayWhite)

		switch currentPage {
		case PAGE_HOME:
			HomePage()
		case PAGE_SETUP_WIZARD:
			SetupWizardPage()
		}

		// if rl.IsMouseButtonReleased(rl.MouseLeftButton) {
		// 	UI.ActiveId = ""
		// }

		rl.DrawFPS(10, 10)
		rl.EndDrawing()
	}
}

func SetupWizardPage() {
	rl.DrawText("Setup wizard", 190, 200, 20, rl.LightGray)
}
