package drawers

import (
	"fmt"
	"image"
	"image/draw"
	"log"
	"math"
	"strings"

	"golang.org/x/image/font"
	"golang.org/x/image/font/gofont/goregular"
	"golang.org/x/image/font/opentype"
	"golang.org/x/image/math/fixed"
)

func getFace(size int) font.Face {
	f, err := opentype.Parse(goregular.TTF)
	if err != nil {
		log.Fatalf("Could not parse TTF, %s", err)
	}
	face, err := opentype.NewFace(f, &opentype.FaceOptions{
		Size:    float64(size),
		DPI:     72,
		Hinting: font.HintingNone,
	})
	if err != nil {
		log.Fatalf("Could not initialize font face, %s", err)
	}
	return face
}

// DrawText renders multiline text into the image. Upper left corner is approximately at (x,y).
func DrawText(img draw.Image, x int, y int, text string, fontsize int) {
	// TODO add option for horizontal alignment
	face := getFace(fontsize)
	lineSpacing := int(float32(fontsize) * 1.3)
	currentLineStart := y + lineSpacing
	d := font.Drawer{
		Dst:  img,
		Src:  image.Black,
		Face: face,
		Dot:  fixed.P(x, currentLineStart),
	}
	lines := strings.Split(text, "\n")
	for _, line := range lines {
		d.DrawString(line)
		currentLineStart += lineSpacing
		d.Dot = fixed.P(x, currentLineStart)
	}
	// TODO maybe return lower right corner to help layouting other elements
}

// DrawBarChart draws a bar chart approximately at (x,y) in the upper lest corner
// and with width w and height h. Automatically adds 0 and max value to y-axis.
// The xLabel is printed in the bottom right corner.
func DrawBarChart(img draw.Image, x, y, w, h int, values []float32, xLabel string) {
	barWidth := w / len(values)
	maxValue := float32(0.0)
	for _, v := range values {
		if v > maxValue {
			maxValue = v
		}
	}
	maxCeiled := int(math.Ceil(float64(maxValue)))
	heightFactor := float32(h) / float32(maxCeiled)
	for i, v := range values {
		draw.Draw(img, image.Rect(
			x+i*barWidth,
			y+h-int(v*heightFactor),
			x+(i+1)*barWidth-2,
			y+h),
			image.Black, image.Point{0, 0}, draw.Over)
	}
	// print legend/axis/units
	// TODO right align
	DrawText(img, x-20, y, fmt.Sprintf("%d", maxCeiled), 16)
	DrawText(img, x-20, y+h, fmt.Sprintf("0"), 16)
	DrawText(img, x+w, y+h, xLabel, 16)
}
