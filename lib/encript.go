package lib

import (
	"net/http"
	"os"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func ReadEncriptedTexture(originalFilePath string) (*rl.Texture2D, error) {
	image, err := ReadEncriptedImage(originalFilePath)
	if err != nil {
		return nil, err
	}

	texture := rl.LoadTextureFromImage(image)
	rl.UnloadImage(image)
	return &texture, nil
}

func detectImageFormat(data []byte) string {
	contentType := http.DetectContentType(data)

	switch contentType {
	case "image/png":
		return ".png"
	case "image/jpeg":
		return ".jpg"
	case "image/gif":
		return ".gif"
	default:
		// Fallback to png
		return ".png"
	}
}

func ReadEncriptedImage(originalFilePath string) (*rl.Image, error) {
	imageBytes, err := os.ReadFile(originalFilePath)
	if err != nil {
		return nil, err
	}

	imageFormat := detectImageFormat(imageBytes)
	image := rl.LoadImageFromMemory(imageFormat, imageBytes, int32(len(imageBytes)))
	return image, nil
}

func detectAudioFormat(data []byte) string {
	contentType := http.DetectContentType(data)
	switch contentType {
	case "audio/mpeg":
		return ".mp3"
	case "audio/wav":
		return ".wav"
	case "audio/ogg":
		return ".ogg"
	default:
		// Fallback to mp3
		return ".mp3"
	}
}

// keep the bytes alive
var bytes []byte

func ReadEncriptedMusic(originalFilePath string) (*rl.Music, error) {
	musicBytes, err := os.ReadFile(originalFilePath)
	if err != nil {
		return nil, err
	}

	audioFormat := detectAudioFormat(musicBytes)
	bytes = musicBytes
	loadedMusic := rl.LoadMusicStreamFromMemory(audioFormat, musicBytes, int32(len(musicBytes)))
	return &loadedMusic, nil
}
