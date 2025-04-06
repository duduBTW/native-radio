package main

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

func ReadEncriptedMusic(originalFilePath string) (*rl.Music, error) {
	tempFilePath, _, err := ReadEncriptedFile("music.mp3", originalFilePath)
	if err != nil {
		return nil, err
	}

	music := rl.LoadMusicStream(*tempFilePath)
	// (*cleanUp)()
	return &music, nil
}
