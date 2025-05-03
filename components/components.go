package components

import (
	"github.com/dudubtw/osu-radio-native/lib"
	rl "github.com/gen2brain/raylib-go/raylib"
)

type Components struct {
	ui            *lib.UIStruct
	textures      lib.Textures
	mousePoint    rl.Vector2
	typographyMap lib.TypographyMap
}

func (c *Components) Update(mousePoint rl.Vector2) {
	c.mousePoint = mousePoint
}

func NewComponents(ui *lib.UIStruct, textures lib.Textures, typographyMap lib.TypographyMap) Components {
	return Components{
		ui:            ui,
		textures:      textures,
		typographyMap: typographyMap,
	}
}
