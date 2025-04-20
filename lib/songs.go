package lib

import (
	"encoding/json"
	"os"
)

type Song struct {
	ID             int
	Audio          string  `json:"audio"`
	Bg             string  `json:"bg,omitempty"`
	OsuFile        string  `json:"osuFile"`
	Path           string  `json:"path"`
	DateAdded      string  `json:"dateAdded"`
	Title          string  `json:"title"`
	Artist         string  `json:"artist"`
	Creator        string  `json:"creator"`
	Duration       float64 `json:"duration"`
	BeatmapSetID   *int    `json:"beatmapSetID,omitempty"`
	Mode           *int    `json:"mode,omitempty"`
	TitleUnicode   *string `json:"titleUnicode,omitempty"`
	ArtistUnicode  *string `json:"artistUnicode,omitempty"`
	PrimaryColor   *string `json:"primaryColor,omitempty"`
	SecondaryColor *string `json:"secondaryColor,omitempty"`
}

type SongTable struct {
	Songs             []Song
	selectedSongindex int
	selectedSong      Song
}

func (table *SongTable) SetSongs(songs []Song) {
	table.Songs = songs
}

func (table *SongTable) HasSelectedSong() bool {
	return table.SelectedSong().Path != ""
}

func (table *SongTable) SelectSong(songIndex int) {
	table.selectedSongindex = songIndex
	table.selectedSong = table.Songs[songIndex]
}

func (table *SongTable) SelectedSong() Song {
	return table.selectedSong
}

func (table *SongTable) Next() {
	table.SelectSong(MinInt(table.selectedSongindex+1, len(table.Songs)-1))
}

func (table *SongTable) Previous() {
	table.SelectSong(MaxInt(table.selectedSongindex-1, 0))
}

func NewSongTableFromJson(filename string) (*SongTable, error) {
	var songTable = &SongTable{}

	file, err := os.ReadFile(filename)
	if err != nil {
		return songTable, err
	}

	var songMap map[string]Song
	err = json.Unmarshal(file, &songMap)
	if err != nil {
		return songTable, err
	}

	var songs = make([]Song, len(songMap))
	var i = 0
	for fileName, song := range songMap {
		songs[i] = song
		songs[i].Path = fileName
		i++
	}

	songTable.Songs = songs
	return songTable, nil
}
