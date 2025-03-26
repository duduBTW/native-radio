package main

import (
	"strings"

	rl "github.com/gen2brain/raylib-go/raylib"
)

// https://github.com/raysan5/raylib/blob/master/examples/text/text_rectangle_bounds.c
func DrawTextByWidth(font rl.Font, text string, pos rl.Vector2, width, fontSize, spacing float32, tint rl.Color) float32 {
	words := strings.Fields(text)
	line := ""
	y := pos.Y

	// Loop over words and build lines until the measured width exceeds the constraint.
	for i, word := range words {
		var testLine string
		if line == "" {
			testLine = word
		} else {
			testLine = line + " " + word
		}
		// Measure the width of the test line using the given font, size and spacing.
		size := rl.MeasureTextEx(font, testLine, fontSize, spacing)
		if size.X > width && line != "" {
			// Draw the current line since adding the new word would exceed the width.
			rl.DrawTextEx(font, line, rl.Vector2{X: pos.X, Y: y}, fontSize, spacing, tint)
			y += size.Y // Increase y by the line height.
			line = word // Start a new line with the current word.
		} else {
			line = testLine
		}
		// If this is the last word, draw the remaining text.
		if i == len(words)-1 {
			size = rl.MeasureTextEx(font, line, fontSize, spacing)
			rl.DrawTextEx(font, line, rl.Vector2{X: pos.X, Y: y}, fontSize, spacing, tint)
			y += size.Y
		}
	}
	return y - pos.Y // Final height = total drawn height minus the starting y position.
}
