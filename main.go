package main

import (
	"fmt"
	"image"
	"image/draw"
	"image/png"
	"log"
	"os"
	"os/exec"
	"time"

	"git.ff02.de/display/fetchers"
	"golang.org/x/image/font"
	"golang.org/x/image/font/gofont/goregular"
	"golang.org/x/image/font/opentype"
	"golang.org/x/image/math/fixed"
)

func main() {
	trains := fetchers.KVV()
	weather := fetchers.DWD()
	time := fetchers.WordClock(time.Now())

	const width, height = 600, 800

	f, err := opentype.Parse(goregular.TTF)
	if err != nil {
		log.Fatalf("Could not parse TTF, %s", err)
	}
	face, err := opentype.NewFace(f, &opentype.FaceOptions{
		Size:    28,
		DPI:     72,
		Hinting: font.HintingNone,
	})
	if err != nil {
		log.Fatalf("Could not initialize font face, %s", err)
	}

	img := image.NewGray(image.Rect(0, 0, width, height))
	// Draw white background
	draw.Draw(img, image.Rect(0, 0, width, height), image.White, image.Point{0, 0}, draw.Over)

	marginLeft := 32
	marginTop := 25
	lineSpacing := 36
	// TODO move this drawing into fetchers and only pass them where to draw?
	currentLineStart := marginTop + lineSpacing
	d := font.Drawer{
		Dst:  img,
		Src:  image.Black,
		Face: face,
		Dot:  fixed.P(marginLeft, currentLineStart),
	}
	for _, train := range trains {
		d.DrawString(train.String())
		currentLineStart += lineSpacing
		d.Dot = fixed.P(marginLeft, currentLineStart)
	}

	marginLeft = 150
	currentLineStart = 300
	d.Dot = fixed.P(marginLeft, currentLineStart)
	for _, line := range time {
		d.DrawString(line)
		currentLineStart += lineSpacing
		d.Dot = fixed.P(marginLeft, currentLineStart)
	}

	// Draw weather
	marginLeft = 100
	currentLineStart = 500
	d.Dot = fixed.P(marginLeft, currentLineStart)
	// TODO draw today a bit bigger, fancier etc
	today := weather[0]
	d.DrawString(today.String())

	// upcoming days
	marginLeft = 80
	currentLineStart = 650
	for i, w := range weather[1:4] {
		d.Dot = fixed.P(marginLeft*(2*i+1), currentLineStart)
		d.DrawString(w.Date.Format("2.1."))
	}
	currentLineStart += lineSpacing
	for i, w := range weather[1:4] {
		d.Dot = fixed.P(marginLeft*(2*i+1), currentLineStart)
		d.DrawString(fmt.Sprintf("%.1f° | %.1f°", w.TempMin, w.TempMax))
	}
	currentLineStart += lineSpacing
	for i, w := range weather[1:4] {
		d.Dot = fixed.P(marginLeft*(2*i+1), currentLineStart)
		d.DrawString(fmt.Sprintf("%.1fmm", w.Rain))
	}

	// Write to PNG
	fileName := "example.png"
	file, err := os.Create(fileName)
	if err != nil {
		log.Fatalf("Could not create file %s, %s", fileName, err)
	}
	if err := png.Encode(file, img); err != nil {
		file.Close()
		log.Fatalf("Could not write PNG to file %s, %s", fileName, err)
	}

	if err := file.Close(); err != nil {
		log.Fatalf("Could not close file %s, %s", fileName, err)
	}

	clearCmd := exec.Command("eips", "-c")
	clearCmd.Run()
	drawCmd := exec.Command("eips", "-g", "example.png")
	drawCmd.Run()
}
