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
	Icons      rl.Texture2D
	Songs      map[string]rl.Texture2D
	Miniatures map[string](*rl.Texture2D)
}

type BlurShader struct {
	Shader       rl.Shader
	TexResLoc    int32
	ScreenResLoc int32
	MouseLoc     int32
}

type UIShaders struct {
	Blur   BlurShader
	Shadow rl.Shader
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

	Songs [30]Song

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

	texture := UI.SelectedSongTexture()
	rl.SetShaderValue(Shaders.Blur.Shader, Shaders.Blur.TexResLoc, []float32{float32(texture.Width), float32(texture.Height)}, rl.ShaderUniformVec2)

	musicPath, err := ReadEncriptedMusic(ui.SelectedSong().FileName)
	if err != nil {
		panic(1)
	}

	rl.UnloadMusicStream(UI.Music)
	ui.Music = *musicPath
	ui.Music.Looping = false

}

func (ui UIStruct) SelectedSongTexture() rl.Texture2D {
	return Textures.Songs[ui.SelectedSong().FileName]
}
func (ui UIStruct) SelectedSongMiniature() *rl.Texture2D {
	return Textures.Miniatures[ui.SelectedSong().FileName]
}

func (ui UIStruct) HasSelectedSong() bool {
	return ui.SelectedSong().FileName != ""
}

var UI = UIStruct{}
var Textures = UITextures{}
var Shaders = UIShaders{}

func LoadTextures(songs SongTable) {
	for _, song := range songs.Songs {
		texture, err := ReadEncriptedTexture(song.Background)
		if err != nil {
			continue
		}
		Textures.Songs[song.FileName] = *texture

		m, err := GenerateImage(song.Background, "D:\\Peronal\\native-radio\\masks\\miniature.png", Size{Width: 350, Height: 350}, rl.White)
		if err != nil {
			panic(1)
		}

		Textures.Miniatures[song.FileName] = m
	}
	Textures.Icons = rl.LoadTexture("D:\\Peronal\\native-radio\\sprites\\icon-sprite.png")
}

func LoadShaders() {
	// Shadow
	Shaders.Shadow = rl.LoadShader("", "D:\\Peronal\\native-radio\\shaders\\shadow.fs")

	// Blur
	shader := rl.LoadShader("", "D:\\Peronal\\native-radio\\shaders\\blur.fs")
	blur := BlurShader{
		Shader:       shader,
		TexResLoc:    rl.GetShaderLocation(shader, "textureResolution"),
		ScreenResLoc: rl.GetShaderLocation(shader, "resolution"),
		MouseLoc:     rl.GetShaderLocation(shader, "mouse"),
	}

	blurRadiusLoc := rl.GetShaderLocation(shader, "blurRadius")
	var blurRadius float32 = 150
	rl.SetShaderValue(shader, blurRadiusLoc, []float32{blurRadius}, rl.ShaderUniformFloat)
	Shaders.Blur = blur
}

func InitEverything() {
	// Songs
	var songTable = SongTable{}
	songTable.FromJson(("C:\\Users\\carlo\\AppData\\Roaming\\osu-radio\\storage\\songs.json"))
	UI.Songs = songTable.Songs
	UI.SelectedPanelPage = PANEL_PAGE_SONGS

	Textures.Songs = make(map[string]rl.Texture2D, len(songTable.Songs))
	Textures.Miniatures = make(map[string]*rl.Texture2D, len(songTable.Songs))

	LoadTextures(songTable)
	LoadShaders()

	UI.SelectSong(0)

}

var MousePoint = rl.Vector2{} // FIX-ME REMOVE DAMN GLOBAL

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
	rl.SetTargetFPS(165)

	InitEverything()

	var currentPage Page = PAGE_HOME

	for !rl.WindowShouldClose() {
		rl.UpdateMusicStream(UI.Music)

		if !UI.isSkeekMode && UI.HasEnded() {
			UI.Next()
		}

		UI.ScreenW = int32(rl.GetScreenWidth())
		UI.ScreenH = int32(rl.GetScreenHeight())
		MousePoint = rl.GetMousePosition()

		rl.SetShaderValue(Shaders.Blur.Shader, Shaders.Blur.ScreenResLoc, []float32{float32(UI.ScreenW), float32(UI.ScreenH)}, rl.ShaderUniformVec2)
		rl.SetShaderValue(Shaders.Blur.Shader, Shaders.Blur.MouseLoc, []float32{MousePoint.X, float32(UI.ScreenH) - MousePoint.Y}, rl.ShaderUniformVec2)

		rl.BeginDrawing()
		rl.ClearBackground(rl.RayWhite)

		switch currentPage {
		case PAGE_HOME:
			HomePage()
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
