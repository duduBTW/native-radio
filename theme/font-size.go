package theme

type Text int32

type FontSizeStruct struct {
	ExtraSuperLarge Text
	ExtraLarge      Text
	Large           Text
	Regular         Text
	Small           Text
}

var FontSize = FontSizeStruct{
	ExtraSuperLarge: 40,
	ExtraLarge:      28,
	Large:           22,
	Regular:         16,
	Small:           14,
}
