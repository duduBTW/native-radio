package main

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

type SongTable struct {
	Songs []Song
}

func (songTable *SongTable) FromJson(filename string) (*SongTable, error) {
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
