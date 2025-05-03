package theme

import (
	"log"
	"strconv"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type ColorsStruct struct {
	Text    rl.Color
	SubText rl.Color
}

var Colors = ColorsStruct{
	Text:    HexToRaylibColor("#EFF1F5"),
	SubText: HexToRaylibColor("#EFF1F5"),
}

func HexToRaylibColor(hex string) rl.Color {
	if hex[0] == '#' {
		hex = hex[1:]
	}

	var r, g, b, a uint8 = 0, 0, 0, 255 // default alpha = 255

	switch len(hex) {
	case 6: // RRGGBB
		ri, _ := strconv.ParseUint(hex[0:2], 16, 8)
		gi, _ := strconv.ParseUint(hex[2:4], 16, 8)
		bi, _ := strconv.ParseUint(hex[4:6], 16, 8)
		r, g, b = uint8(ri), uint8(gi), uint8(bi)
	case 8: // RRGGBBAA
		ri, _ := strconv.ParseUint(hex[0:2], 16, 8)
		gi, _ := strconv.ParseUint(hex[2:4], 16, 8)
		bi, _ := strconv.ParseUint(hex[4:6], 16, 8)
		ai, _ := strconv.ParseUint(hex[6:8], 16, 8)
		r, g, b, a = uint8(ri), uint8(gi), uint8(bi), uint8(ai)
	default:
		log.Fatalf("invalid hex color: %s", hex)
	}

	return rl.NewColor(r, g, b, a)
}
