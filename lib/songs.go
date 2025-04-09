package lib

import (
	"encoding/json"
	"os"
)

type Song struct {
	Audio      string `json:"audio"`
	Title      string `json:"title"`
	Artist     string `json:"artist"`
	Background string `json:"bg"`
	FileName   string
}

func (s *Song) Id() string {
	return s.Title + s.Artist + s.FileName
}

type SongTable struct {
	Songs             []Song
	selectedSongindex int
}

func (table *SongTable) HasSelectedSong() bool {
	return table.SelectedSong().FileName != ""
}

func (table *SongTable) SelectSong(songIndex int) {
	table.selectedSongindex = songIndex
}

func (table *SongTable) SelectedSong() Song {
	return table.Songs[table.selectedSongindex]
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
		songs[i].FileName = fileName
		i++
	}

	songTable.Songs = songs
	return songTable, nil
}
