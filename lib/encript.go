package lib

import (
	"fmt"
	"os"
	"strconv"
	"time"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func ReadEncriptedTexture(originalFilePath string) (*rl.Texture2D, error) {
	tempFilePath, cleanUp, err := ReadEncriptedFile("image.png", originalFilePath, true)
	if err != nil {
		return nil, err
	}

	texture := rl.LoadTexture(*tempFilePath)
	(*cleanUp)()
	return &texture, nil
}

func ReadEncriptedImage(originalFilePath string) (*rl.Image, error) {
	tempFilePath, cleanUp, err := ReadEncriptedFile("image.png", originalFilePath, true)
	if err != nil {
		return nil, err
	}

	fmt.Print("temp", tempFilePath)
	image := rl.LoadImage(*tempFilePath)
	(*cleanUp)()
	return image, nil
}

// This is kinda dumb but the easy way I found for raylib to detect the files with a sha256 name.
func ReadEncriptedFile(tempFile string, originalFilePath string, random bool) (*string, *func(), error) {
	fileData, err := os.ReadFile(originalFilePath)
	if err != nil {
		return nil, nil, err
	}

	tempFilePath := "D:\\Peronal\\native-radio\\temp\\"
	if random {
		tempFilePath += strconv.Itoa(int(time.Now().UnixMilli()))
	}
	tempFilePath += tempFile

	if err := os.WriteFile(tempFilePath, fileData, 0644); err != nil {
		return nil, nil, err
	}
	cleanUp := func() {
		fmt.Println("Cleaning up")
		os.Remove(tempFilePath)
	}
	return &tempFilePath, &cleanUp, nil
}
