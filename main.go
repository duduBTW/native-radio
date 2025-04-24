package main

import (
	"database/sql"
	"fmt"
	"os"
	"os/exec"

	dbManager "github.com/dudubtw/osu-radio-native/db-manager"
	"github.com/dudubtw/osu-radio-native/lib"
	rl "github.com/gen2brain/raylib-go/raylib"
)

var songTable lib.SongTable
var music lib.Music
var textures lib.Textures
var ui lib.UIStruct
var shaders lib.Shaders
var mousePoint rl.Vector2
var db *sql.DB

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
	newVolume := float32(dbManager.GetUserVolume(db)) / 100
	rl.SetMusicVolume(*music.Selected, newVolume)
	ui.Volume = newVolume
	textures.LoadSelectedSong(songTable, shaders)
}
func SelectSong(songIndex int) {
	songTable.SelectSong(songIndex)
	UpdateSong()
	go func() {
		dbManager.UpdateSelectedIndex(db, songIndex)
	}()
}
func Shuffle() {
	SelectSong(lib.RandomRange(0, len(songTable.Songs)))
	ui.ScrollToIndex(songTable.SelectedSongindex)
	music.Play()
}
func ExecTrash(lazerFilePath string) {
	go func() {
		cmd := exec.Command("node", "D:\\Peronal\\native-radio\\trash\\index.js", "--lazer="+lazerFilePath)
		stdout, err := cmd.Output()
		if err != nil {
			fmt.Println(err)
			return
		}

		songMap, err := lib.ParseNodeOutput(stdout)
		if err != nil {
			fmt.Println(err)
			return
		}

		for _, song := range songMap {
			dbManager.InsertSong(db, &song)
		}

		songs, _, err := dbManager.SelectAllSongs(db, ui.SearchValue, 0)
		if err != nil {
			fmt.Println(err)
			return
		}

		songTable.SetSongs(songs)
		textures.SyncWithTable(&songTable)
		SelectSong(0)
		ui.SelectedPage = lib.PAGE_HOME
	}()
}

func UpdateSongList() {
	go func() {
		songs, containsIndex, err := dbManager.SelectAllSongs(db, ui.SearchValue, songTable.SelectedSong().ID)
		if err != nil {
			fmt.Println(err)
			return
		}

		songTable.SetSongs(songs)
		if containsIndex != -1 {
			ui.ScrollToIndex(containsIndex)
		} else {
			ui.SidePanelScrollTop = 0
		}
	}()
}

func InitEverything() {
	dbPath := "user-data.db"
	newDb, err := dbManager.InitDB(dbPath)
	if err != nil {
		fmt.Println(err)
		panic(1)
	}
	db = newDb

	if err := dbManager.SetupUser(db); err != nil {
		fmt.Println(err)
		panic(1)
	}

	table, err := dbManager.NewSongTableFromDb(db)
	if err != nil {
		fmt.Println("Could not load songs!")
		panic(1)
	}

	songTable = *table
	ui = lib.NewUi(table)
	textures = lib.NewTexture(table)
	shaders = lib.NewShaders()

	fmt.Print("Page", ui.SelectedPage)

	if len(table.Songs) == 0 {
		return
	}
	selectedIndex := dbManager.GetUserSelectedIndex(db)
	ui.ScrollToIndex(selectedIndex)
	SelectSong(selectedIndex)
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
