package lib

import (
	"strconv"

	"github.com/dudubtw/osu-radio-native/theme"
	rl "github.com/gen2brain/raylib-go/raylib"
)

type mapKey string

type TypographyMap = map[mapKey]rl.Font

const fontName = "Torus"
const fontFormat = ".otf"

func newKey(size theme.Text, weight theme.Weight) mapKey {
	return mapKey(strconv.Itoa(int(size)) + string(weight))
}

func InitTypography() TypographyMap {
	sizes := []theme.Text{theme.FontSize.ExtraSuperLarge, theme.FontSize.ExtraLarge, theme.FontSize.Large, theme.FontSize.Regular, theme.FontSize.Small}
	weights := []theme.Weight{theme.FontWeight.Bold, theme.FontWeight.Regular, theme.FontWeight.Light}

	typographyMap := make(TypographyMap)

	for _, size := range sizes {
		for _, weight := range weights {
			key := newKey(size, weight)
			fileName := "fonts/" + fontName + "_" + string(weight) + fontFormat
			typographyMap[key] = rl.LoadFontEx(fileName, int32(size), nil, 0)
		}
	}

	return typographyMap
}

func MeasureText(text string, size theme.Text, weight theme.Weight, typographyMap TypographyMap) int32 {
	font := typographyMap[newKey(size, weight)]
	return int32(rl.MeasureTextEx(font, text, float32(size), 1).X)
}

func Typography(value string, pos rl.Vector2, size theme.Text, weight theme.Weight, color rl.Color, typographyMap TypographyMap) {
	font := typographyMap[newKey(size, weight)]
	rl.DrawTextEx(font, value, pos, float32(size), 1, color)
}
