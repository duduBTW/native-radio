package main

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

func ReadEncriptedTexture(originalFilePath string) (*rl.Texture2D, error) {
	tempFilePath, cleanUp, err := ReadEncriptedFile("image.png", originalFilePath)
	if err != nil {
		return nil, err
	}

	texture := rl.LoadTexture(*tempFilePath)
	(*cleanUp)()
	return &texture, nil
}
func ReadEncriptedImage(originalFilePath string) (*rl.Image, error) {
	tempFilePath, cleanUp, err := ReadEncriptedFile("image.png", originalFilePath)
	if err != nil {
		return nil, err
	}

	image := rl.LoadImage(*tempFilePath)
	(*cleanUp)()
	return image, nil
}

func GenerateImage(filePath string, maskPath string, size Size, color rl.Color) (*rl.Texture2D, error) {
	image, err := ReadEncriptedImage(filePath)
	if err != nil {
		return nil, err
	}

	mask := rl.LoadImage(maskPath)

	rl.ImageCrop(image, ImageFitCordinates(image, size))
	rl.ImageResize(image, int32(size.Width), int32(size.Height))
	rl.ImageAlphaMask(image, mask)
	texture := rl.LoadTextureFromImage(image)
	rl.UnloadImage(image)
	return &texture, nil
}

func ImageFitCordinates(image *rl.Image, target Size) rl.Rectangle {
	textureW := float32(image.Width)
	textureH := float32(image.Height)

	screenAspectRatio := target.Width / target.Height
	textureAspectRatio := textureW / textureH

	var newWidth, newHeight float32
	var sourceX, sourceY float32

	if textureAspectRatio < screenAspectRatio {
		// Image is too tall, scale width to fit and crop top/bottom
		newWidth = textureW
		newHeight = textureW / screenAspectRatio
		sourceX = 0
		sourceY = (textureH - newHeight) / 2 // Crop vertically
	} else {
		// Image is too wide, scale height to fit and crop sides
		newHeight = textureH
		newWidth = textureH * screenAspectRatio
		sourceX = (textureW - newWidth) / 2 // Crop horizontally
		sourceY = 0
	}

	return rl.NewRectangle(sourceX, sourceY, newWidth, newHeight)
}

func DrawFitImage(texture rl.Texture2D, target rl.Rectangle, color rl.Color) {
	textureW := float32(texture.Width)
	textureH := float32(texture.Height)

	screenAspectRatio := target.Width / target.Height
	textureAspectRatio := textureW / textureH

	var newWidth, newHeight float32
	var sourceX, sourceY float32

	if textureAspectRatio < screenAspectRatio {
		// Image is too tall, scale width to fit and crop top/bottom
		newWidth = textureW
		newHeight = textureW / screenAspectRatio
		sourceX = 0
		sourceY = (textureH - newHeight) / 2 // Crop vertically
	} else {
		// Image is too wide, scale height to fit and crop sides
		newHeight = textureH
		newWidth = textureH * screenAspectRatio
		sourceX = (textureW - newWidth) / 2 // Crop horizontally
		sourceY = 0
	}

	// Draw the cropped and scaled texture
	rl.DrawTexturePro(
		texture,
		rl.NewRectangle(sourceX, sourceY, newWidth, newHeight),
		rl.NewRectangle(target.X, target.Y, target.Width, target.Height),
		rl.NewVector2(0, 0),
		0,
		color,
	)
}
