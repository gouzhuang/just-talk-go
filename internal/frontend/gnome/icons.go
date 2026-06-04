//go:build gnome
package gnome

import (
	"bytes"
	"image"
	"image/color"
	"image/png"
	"math"
)

var iconCache = map[string][]byte{}

func getIcon(state string) []byte {
	if b, ok := iconCache[state]; ok {
		return b
	}
	b := generateIcon(state)
	iconCache[state] = b
	return b
}

func generateIcon(state string) []byte {
	const size = 22
	img := image.NewRGBA(image.Rect(0, 0, size, size))

	// Background color based on state
	var bg color.Color
	switch state {
	case "connecting":
		bg = color.RGBA{0xF5, 0xBE, 0x46, 0xFF}
	case "recording":
		bg = color.RGBA{0xFF, 0x41, 0x41, 0xFF}
	case "stopping", "stopping_delayed":
		bg = color.RGBA{0xFF, 0x8C, 0x3C, 0xFF}
	case "error":
		bg = color.RGBA{0xFF, 0x41, 0x41, 0xFF}
	default:
		bg = color.RGBA{0x91, 0x91, 0x91, 0xFF}
	}

	white := color.RGBA{0xFF, 0xFF, 0xFF, 0xFF}

	// Draw circular background
	drawFilledCircle(img, size/2, size/2, size/2-1, bg)

	// Draw microphone symbol
	drawMicrophone(img, white)

	// Recording indicator dot (top-right)
	if state == "recording" {
		drawFilledCircle(img, 17, 5, 2, white)
	}

	// Error exclamation mark (center)
	if state == "error" {
		drawExclamation(img, white)
	}

	var buf bytes.Buffer
	_ = png.Encode(&buf, img)
	return buf.Bytes()
}

func drawFilledCircle(img *image.RGBA, cx, cy, r int, c color.Color) {
	for y := -r; y <= r; y++ {
		for x := -r; x <= r; x++ {
			if x*x+y*y <= r*r+r/2 {
				img.Set(cx+x, cy+y, c)
			}
		}
	}
}

func drawMicrophone(img *image.RGBA, c color.Color) {
	// Microphone head: circle at (11, 6) with radius 4
	drawFilledCircle(img, 11, 6, 4, c)

	// Microphone body/stem: vertical line from (11, 10) to (11, 15)
	drawLine(img, 11, 10, 11, 15, c)

	// Microphone stand arc: semi-circle at bottom
	// Draw arc from angle 0 to PI (bottom half)
	for angle := 0.0; angle <= math.Pi; angle += 0.05 {
		x := 11 + int(4*math.Cos(angle))
		y := 15 + int(4*math.Sin(angle))
		img.Set(x, y, c)
		img.Set(x, y+1, c) // Make arc 2px thick
	}

	// Microphone stand vertical connector
	drawLine(img, 11, 15, 11, 17, c)

	// Microphone stand base (horizontal line)
	for x := 7; x <= 15; x++ {
		img.Set(x, 17, c)
		img.Set(x, 18, c) // Make base 2px thick
	}
}

func drawExclamation(img *image.RGBA, c color.Color) {
	// Main stem of exclamation mark
	for y := 4; y <= 12; y++ {
		img.Set(10, y, c)
		img.Set(11, y, c)
		img.Set(12, y, c)
	}
	// Dot at bottom
	drawFilledCircle(img, 11, 16, 2, c)
}

func drawLine(img *image.RGBA, x0, y0, x1, y1 int, c color.Color) {
	// Bresenham's line algorithm
	dx := abs(x1 - x0)
	dy := abs(y1 - y0)
	sx := 1
	if x0 > x1 {
		sx = -1
	}
	sy := 1
	if y0 > y1 {
		sy = -1
	}
	err := dx - dy

	for {
		img.Set(x0, y0, c)
		if x0 == x1 && y0 == y1 {
			break
		}
		e2 := 2 * err
		if e2 > -dy {
			err -= dy
			x0 += sx
		}
		if e2 < dx {
			err += dx
			y0 += sy
		}
	}
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
