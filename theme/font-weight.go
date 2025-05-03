package theme

type Weight string

type FontWeightStruct struct {
	Bold    Weight
	Light   Weight
	Regular Weight
}

var FontWeight = FontWeightStruct{
	Bold:    "Bold",
	Regular: "Regular",
	Light:   "Light",
}
