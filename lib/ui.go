package lib

import (
	"time"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type UIStruct struct {
	SidePanelScrollTop float32
	ScreenW            int32
	ScreenH            int32

	SelectedPanelPage PanelPage
	SelectedPage      Page

	FocusedId string
	ActiveId  string
	HotId     string

	Volume           float32
	IsMuted          bool
	LastTimeScrolled time.Time

	SearchValue      string
	InputCursorStart int
	InputCursorEnd   int
	// Time between blinks
	BlinkTimer float32
	// Time a blink stayed active
	BlinkingTimer float32
}

func (ui *UIStruct) SetCursors(pos int) {
	ui.InputCursorStart = pos
	ui.InputCursorEnd = pos
}
func (ui *UIStruct) IncrementCursor() {
	ui.InputCursorStart += 1
	ui.InputCursorEnd += 1
}
func (ui *UIStruct) DecrementCursor() {
	ui.InputCursorStart -= 1
	ui.InputCursorEnd -= 1
}
func (ui *UIStruct) ScrollToIndex(index int) {
	ui.SidePanelScrollTop = -float32(index * (72 + 12))
}

func NewUi(table *SongTable) UIStruct {
	var selectedPage Page = PAGE_SETUP_WIZARD
	if len(table.Songs) > 0 {
		selectedPage = PAGE_HOME
	}

	ui := UIStruct{
		SelectedPanelPage: PANEL_PAGE_SONGS,
		SelectedPage:      selectedPage,
		Volume:            0.5,
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

// func DrawFitImage(texture rl.Texture2D, rect rl.Rectangle, color rl.Color) {
// 	origin := rl.NewVector2(float32(texture.Width), float32(texture.Height))
// 	target := rl.NewVector2(float32(rect.Width), float32(rect.Height))

// 	rl.DrawTexturePro(
// 		texture,
// 		ImageFitCordinates(origin, target),
// 		rect,
// 		rl.NewVector2(0, 0),
// 		0,
// 		color,
// 	)
// }

func DrawFitTexture(tex rl.Texture2D, dest rl.Rectangle, tint rl.Color) {
	// use the texture’s fields, not some stale origin you pulled earlier
	origin := rl.NewVector2(float32(tex.Width), float32(tex.Height))
	target := rl.NewVector2(dest.Width, dest.Height)

	src := ImageFitCordinates(origin, target)
	// now src.X/Y/W/H line up exactly with GPU‐side pixels

	rl.DrawTexturePro(
		tex,
		src,
		dest,
		rl.NewVector2(0, 0),
		0,
		tint,
	)
}
