package main

import (
	"fmt"
	"image"
	"image/draw"
	"image/png"
	"log"
	"math"
	"os"
	"os/exec"
	"strings"
	"time"

	"git.ff02.de/display/drawers"
	"git.ff02.de/display/fetchers"
)

func main() {
	trains := fetchers.KVV()
	weather := fetchers.DWD()
	tim := fetchers.WordClock(time.Now())

	const width, height = 600, 800

	img := image.NewGray(image.Rect(0, 0, width, height))
	// Draw white background
	draw.Draw(img, image.Rect(0, 0, width, height), image.White, image.Point{0, 0}, draw.Over)

	// draw train departures
	q := ""
	for _, train := range trains {
		q += train.String() + "\n"
	}
	drawers.DrawText(img, 50, 35, q, 24)

	// draw wordclock
	drawers.DrawText(img, 150, 260, strings.Join(tim, "\n"), 32)

	// Draw weather
	today := weather[0]
	drawers.DrawText(img, 80, 480, fmt.Sprintf("%d°", round(today.TempMax)), 28)
	drawers.DrawText(img, 80, 560, fmt.Sprintf("%d°", round(today.TempMin)), 28)
	drawers.DrawImage(img, 130, 490, "assets/thermometer.png")
	// Draw rain (diagram) - bar chart
	nowHour := time.Now().UTC().Hour() // TODO check whether DWD starts at UTC 0 or german 0
	values := today.RainHourly[nowHour:]
	drawers.DrawBarChart(img, 210, 480, 120, 100, values, fmt.Sprintf("+%d h", len(values)))
	// upcoming days
	for i, w := range weather[1:4] {
		drawers.DrawText(img,
			80*(2*i+1),
			650,
			fmt.Sprintf("%s\n%d° | %d°\n%.1fmm",
				w.Date.Format("2.1."),
				round(w.TempMin),
				round(w.TempMax),
				w.Rain,
			),
			24)
	}
	// visualizing general weather using penguin
	drawers.DrawImage(img, width-256, 370, "assets/penguins/penguin.png")
	if (today.TempMin < 0 || today.TempMax < 5) && today.Rain > 1 {
		drawers.DrawImage(img, width-256, 370, "assets/penguins/snow.png")
	}
	if today.WindSpeed > 10 {
		drawers.DrawImage(img, width-256, 370, "assets/penguins/scarf.png")
	}
	if today.TempMin < 5 || today.TempMax < 10 {
		drawers.DrawImage(img, width-256, 370, "assets/penguins/hat.png")
	}
	if today.Sunshine > 2 {
		drawers.DrawImage(img, width-256, 370, "assets/penguins/sunglasses.png")
	}
	if today.Rain > 1 {
		drawers.DrawImage(img, width-256, 370, "assets/penguins/umbrella.png")
	}
	// TODO add wind (direction and speed). That is only available per-day via the API
	// TODO add sun shine /rise /set
	// TODO show warnings

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

	// Only clear on every 5th refresh to be less flashy
	if time.Now().Minute()%5 == 0 {
		clearCmd := exec.Command("eips", "-c")
		clearCmd.Run()
	}
	drawCmd := exec.Command("eips", "-g", "example.png")
	drawCmd.Run()
}

func round(f float32) int {
	return int(math.Round(float64(f)))
}
