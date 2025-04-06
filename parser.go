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
	Songs [30]Song
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

	var songs [30]Song
	var i = 0
	for fileName, song := range songMap {
		if i > 29 {
			break
		}

		songs[i] = song
		songs[i].FileName = fileName
		i++
	}

	songTable.Songs = songs

	return songTable, nil
}

// This is kinda dumb but the easy way I found for raylib to detect the files with a sha256 name.
func ReadEncriptedFile(tempFile string, originalFilePath string) (*string, *func(), error) {
	fileData, err := os.ReadFile(originalFilePath)
	if err != nil {
		return nil, nil, err
	}

	tempFilePath := "D:\\Peronal\\native-radio\\temp\\" + tempFile
	if err := os.WriteFile(tempFilePath, fileData, 0644); err != nil {
		return nil, nil, err
	}
	cleanUp := func() {
		os.Remove(tempFilePath)
	}
	return &tempFilePath, &cleanUp, nil
}
