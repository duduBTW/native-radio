package main

import rl "github.com/gen2brain/raylib-go/raylib"

// Helper function to find the minimum of two float32 values
func minF(a, b float32) float32 {
	if a < b {
		return a
	}
	return b
}

// Helper function to find the maximum of two float32 values
func maxF(a, b float32) float32 {
	if a > b {
		return a
	}
	return b
}

type Roundness = float32

const (
	ROUNDED = 8
)

func DrawRectangleRoundedPixels(rec rl.Rectangle, radiusPixels Roundness, color rl.Color) {
	// Ensure radius is not negative
	radiusPixels = maxF(0, radiusPixels)

	// Prevent division by zero or weirdness with non-positive dimensions
	if rec.Width <= 0 || rec.Height <= 0 {
		// Optionally draw a regular rectangle or just return
		// rl.DrawRectangleRec(rec, color)
		return
	}

	// Find the smaller dimension
	minDimension := minF(rec.Width, rec.Height)

	// Calculate roundness based on pixel radius
	// roundness = (radius * 2) / minDimension
	roundness := (radiusPixels * 2) / minDimension

	// Clamp roundness to the valid range [0.0, 1.0]
	// If requested radiusPixels * 2 > minDimension, it means the radius
	// is too large, so we cap at full roundness (1.0).
	roundness = maxF(0.0, minF(roundness, 1.0))

	// Call the original raylib function with the calculated roundness
	rl.DrawRectangleRounded(rec, roundness, 0, color)
}
