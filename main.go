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
	tim := fetchers.WordClock(time.Now())

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
	for _, line := range tim {
		d.DrawString(line)
		currentLineStart += lineSpacing
		d.Dot = fixed.P(marginLeft, currentLineStart)
	}

	// Draw weather
	marginLeft = 80
	currentLineStart = 480
	d.Dot = fixed.P(marginLeft, currentLineStart)
	// TODO draw today a bit bigger, fancier etc
	today := weather[0]
	d.DrawString(fmt.Sprintf("%.1f°", today.TempMax))
	currentLineStart += 2 * lineSpacing
	d.Dot = fixed.P(marginLeft, currentLineStart)
	d.DrawString(fmt.Sprintf("%.1f°", today.TempMin))
	// Draw rain (diagram) - bar chart
	fmt.Printf("Rain today: %v", today.RainHourly)
	nowHour := time.Now().UTC().Hour() // TODO check whether DWD starts at UTC 0 or german 0
	values := today.RainHourly[nowHour:]
	marginLeft = 280
	currentLineStart = 480
	w := 120
	barWidth := w / len(values)
	he := 50
	maxValue := float32(0.0)
	for _, v := range values {
		if v > maxValue {
			maxValue = v
		}
	}
	heightFactor := he / (int(maxValue) + 1)
	for i, h := range values {
		draw.Draw(img, image.Rect(
			marginLeft+i*barWidth,
			currentLineStart+he-int(h*float32(heightFactor)),
			marginLeft+(i+1)*barWidth-2,
			currentLineStart+he), image.Black, image.Point{0, 0}, draw.Over)
	}
	// print legend/axis/units
	// TODO right align
	d.Dot = fixed.P(marginLeft-20, currentLineStart)
	d.DrawString(fmt.Sprintf("%d", int(maxValue)+1))
	d.Dot = fixed.P(marginLeft-20, currentLineStart+he)
	d.DrawString("0")
	d.Dot = fixed.P(marginLeft+w, currentLineStart+he+lineSpacing)
	d.DrawString(fmt.Sprintf("+%d h", len(values)))

	// upcoming days
	// TODO print in smaller font
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

	// TODO only clear on every nth refresh to be less flashy
	clearCmd := exec.Command("eips", "-c")
	clearCmd.Run()
	drawCmd := exec.Command("eips", "-g", "example.png")
	drawCmd.Run()
}
