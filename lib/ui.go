package lib

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

type UIStruct struct {
	SidePanelScrollTop float32
	ScreenW            int32
	ScreenH            int32

	SelectedPanelPage PanelPage
	SelectedPage      Page

	ActiveId string
	HotId    string
}

func NewUi() UIStruct {
	ui := UIStruct{
		SelectedPanelPage: PANEL_PAGE_SONGS,
		SelectedPage:      PAGE_HOME,
	}
	return ui
}

func ImageFitCordinates(origin rl.Vector2, target rl.Vector2) rl.Rectangle {
	textureW := origin.X
	textureH := origin.Y

	screenAspectRatio := target.X / target.Y
	textureAspectRatio := textureW / textureH

	var newWidth, newHeight float32
	var sourceX, sourceY float32

	if textureAspectRatio < screenAspectRatio {
		// Image is too tall, scale width to fit and crop top/bottom
		newWidth = textureW
		newHeight = textureW / screenAspectRatio
		sourceX = 0
		sourceY = float32(int32((textureH - newHeight) / 2)) // Crop vertically
	} else {
		// Image is too wide, scale height to fit and crop sides
		newHeight = textureH
		newWidth = textureH * screenAspectRatio
		sourceX = float32(int32((textureW - newWidth) / 2)) // Crop horizontally
		sourceY = 0
	}

	return rl.NewRectangle(sourceX, sourceY, newWidth, newHeight)
}

func DrawFitImage(texture rl.Texture2D, target rl.Rectangle, color rl.Color) {
	rl.DrawTexturePro(
		texture,
		ImageFitCordinates(rl.NewVector2(float32(texture.Width), float32(texture.Height)), rl.NewVector2(float32(target.Width), float32(target.Height))),
		rl.NewRectangle(target.X, target.Y, target.Width, target.Height),
		rl.NewVector2(0, 0),
		0,
		color,
	)
}
