package lib

import (
	"fmt"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type Textures struct {
	Icons        *rl.Texture2D
	SelectedSong *rl.Texture2D
	Miniature    *rl.Texture2D

	Songs map[string](*rl.Texture2D)

	SongsLoading bool
	SongsERROR   map[string]bool
	IsLoadingSmt bool
}

var pendingTextures = make(chan struct {
	song  Song
	image *rl.Image
}, 32)

func (textures *Textures) GetSong(song Song) *rl.Texture2D {
	return textures.Songs[song.Id()]
}
func (textures *Textures) SetSong(song Song, texture *rl.Texture2D) {
	textures.Songs[song.Id()] = texture
}
func (textures *Textures) SetSongError(song Song, hasErrror bool) {
	textures.SongsERROR[song.Id()] = hasErrror
}
func (textures *Textures) GetSongError(song Song) bool {
	return textures.SongsERROR[song.Id()]
}
func (textures *Textures) GetSongLoading(song Song) bool {
	return textures.SongsLoading
}
func (textures *Textures) SetSongLoading(isLoading bool) {
	textures.SongsLoading = isLoading
}

func (textures *Textures) LoadIcon() {
	iconTexture := rl.LoadTexture("D:\\Peronal\\native-radio\\sprites\\icon-sprite.png")
	textures.Icons = &iconTexture
}

func GenerateImage(filePath string, maskPath string, size rl.Vector2, color rl.Color) (*rl.Image, error) {
	image, err := ReadEncriptedImage(filePath)
	if err != nil {
		return nil, err
	}

	var mask *rl.Image = nil
	if maskPath != "" {
		mask = rl.LoadImage(maskPath)
	}

	rl.ImageCrop(image, ImageFitCordinates(rl.NewVector2(float32(image.Width), float32(image.Height)), size))
	rl.ImageResize(image, int32(size.X), int32(size.Y))

	if mask != nil {
		rl.ImageAlphaMask(image, mask)
	}

	return image, nil
}
func GenerateTexture(filePath string, maskPath string, size rl.Vector2, color rl.Color) (*rl.Texture2D, error) {
	image, err := GenerateImage(filePath, maskPath, size, color)
	if err != nil {
		return nil, err
	}

	texture := rl.LoadTextureFromImage(image)
	rl.UnloadImage(image)
	return &texture, nil
}

// --- Song
func (textures *Textures) LoadSongCard(song Song, rect rl.Rectangle) {
	if textures.GetSong(song) != nil || textures.GetSongLoading(song) || textures.GetSongError(song) {
		return
	}

	textures.IsLoadingSmt = true
	textures.SetSongLoading(true)

	go func() {
		image, err := GenerateImage(song.Background, "D:\\Peronal\\native-radio\\masks\\card.png", rl.NewVector2(rect.Width, rect.Height), rl.White)
		if err != nil {
			fmt.Println("ERROR LOADING SONG TEXTURE!", err)
			textures.SetSongLoading(false)
			textures.SetSongError(song, true)
			return
		}

		pendingTextures <- struct {
			song  Song
			image *rl.Image
		}{song: song, image: image}
	}()
}
func (textures *Textures) UnloadSongCard(song Song) {
	texture := textures.GetSong(song)
	if texture == nil {
		return
	}

	textures.SetSong(song, nil)
	rl.UnloadTexture(*texture)
}

func (textures *Textures) LoadSelectedSong(table SongTable, shaders Shaders) {
	selectedSong := table.SelectedSong()

	textures.LoadSongBackground(selectedSong, shaders)
	textures.LoadSongMiniature(selectedSong)
}

func (textures *Textures) LoadSongBackground(song Song, shaders Shaders) {
	previusTexture := textures.SelectedSong
	if previusTexture != nil {
		rl.UnloadTexture(*previusTexture)
	}
	texture, err := ReadEncriptedTexture(song.Background)
	if err != nil {
		panic(1)
	}

	rl.SetShaderValue(shaders.Blur.Shader, shaders.Blur.TexResLoc, []float32{float32(texture.Width), float32(texture.Height)}, rl.ShaderUniformVec2)
	textures.SelectedSong = texture
}

func (textures *Textures) LoadSongMiniature(song Song) {
	previusMiniatureTexture := textures.Miniature
	if previusMiniatureTexture != nil {
		rl.UnloadTexture(*previusMiniatureTexture)
	}

	texture, err := GenerateTexture(song.Background, "D:\\Peronal\\native-radio\\masks\\miniature.png", rl.NewVector2(STYLE_MINIATURE_SIZE, STYLE_MINIATURE_SIZE), rl.White)
	if err != nil {
		fmt.Println(err)
		panic(1)
	}

	textures.Miniature = texture
}

func (textures *Textures) ProcessPendingTextures() {
	for {
		select {
		case job := <-pendingTextures:
			texture := rl.LoadTextureFromImage(job.image)
			rl.UnloadImage(job.image)

			textures.SetSong(job.song, &texture)
			textures.SetSongLoading(false)
		default:
			return
		}
	}
}

func NewTexture(table *SongTable) Textures {
	textures := Textures{
		Songs:      make(map[string]*rl.Texture2D, len(table.Songs)),
		SongsERROR: make(map[string]bool, len(table.Songs)),
	}
	textures.LoadIcon()

	return textures
}
