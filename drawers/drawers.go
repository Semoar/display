package drawers

import (
	"image"
	"image/draw"
	"log"
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
}
