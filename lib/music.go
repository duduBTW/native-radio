package lib

import (
	"fmt"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type Music struct {
	Selected             *rl.Music
	IsSkeekMode          bool
	wasPlayingBeforeSeek bool
	seekModeProgress     float32
}

func (music *Music) LoadMusic(table *SongTable) error {
	if music.Selected != nil {
		rl.UnloadMusicStream(*music.Selected)
	}

	loadedMusic, err := ReadEncriptedMusic(table.SelectedSong().FileName)
	if err != nil {
		return err
	}

	music.Selected = loadedMusic
	music.Selected.Looping = false
	return nil
}

func (music *Music) BeginSeekMode() {
	if music.IsSkeekMode {
		fmt.Println("uh oh")
		return
	}

	music.wasPlayingBeforeSeek = rl.IsMusicStreamPlaying(*music.Selected)
	rl.PauseMusicStream(*music.Selected)
	music.seekModeProgress = music.Progress()
	music.IsSkeekMode = true
}

func (music *Music) ExitSeekMode() {
	if !music.IsSkeekMode {
		fmt.Println("uh oh")
		return
	}

	music.IsSkeekMode = false
	rl.SeekMusicStream(*music.Selected, rl.GetMusicTimeLength(*music.Selected)*music.seekModeProgress)

	if music.wasPlayingBeforeSeek {
		rl.PlayMusicStream(*music.Selected)
	}
}

func (music *Music) Play() {
	rl.PlayMusicStream(*music.Selected)
}
func (music *Music) Pause() {
	rl.PauseMusicStream(*music.Selected)
}

func (music *Music) Previous(table *SongTable) {
	table.Previous()
	music.Play()
}

func (music *Music) Next(table *SongTable) {
	table.Next()
	music.Play()
}

func (music *Music) Progress() float32 {
	if music.IsSkeekMode {
		return music.seekModeProgress
	}

	// fmt.Println(music.Selected, rl.GetMusicTimePlayed(*music.Selected), rl.GetMusicTimeLength(*music.Selected))
	return rl.GetMusicTimePlayed(*music.Selected) / rl.GetMusicTimeLength(*music.Selected)
}

func (music *Music) Seek(to float32) {
	music.seekModeProgress = to
}

func (music *Music) HasEnded() bool {
	if !rl.IsMusicValid(*music.Selected) {
		return false
	}

	// TODO deal with ms, since we are doing int() it only compare seconds
	return int(rl.GetMusicTimePlayed(*music.Selected)) == int(rl.GetMusicTimeLength(*music.Selected))
}

func (music *Music) SetVolume(newVolume float32) {
	if music.Selected == nil {
		return
	}

	rl.SetMusicVolume(*music.Selected, newVolume)
}
