package fetchers

import (
	"fmt"
	"time"

	// embed timezones as the Kindle doesn't have timezone infos
	_ "time/tzdata"
)

func wordsMinutes(m int) string {
	quarters := []string{"", "viertel ", "halb ", "dreiviertel "}
	quarter := (m + 7) / 15 % 4
	additions := []string{"", "kurz nach ", "kurz vor "}
	x := (m + 2) / 5 % 3
	return additions[x] + quarters[quarter]
}

func wordsHour(t time.Time) string {
	h := t.Hour()
	if t.Minute() >= 8 {
		h += 1
	}
	h = h % 12
	names := []string{"zwölf", "eins", "zwei", "drei", "vier", "fünf", "sechs", "sieben", "acht", "neun", "zehn", "elf"}
	return names[h]
}

var wochentage = []string{"Sonntag", "Montag", "Dienstag", "Mittwoch", "Donnerstag", "Freitag", "Samstag"}

func WordClock(t time.Time) []string {
	loc, err := time.LoadLocation("Europe/Berlin")
	if err != nil {
		fmt.Printf("Could not load location, %s", err)
	}
	t = t.In(loc)
	fmt.Printf("%v", t)
	return []string{
		fmt.Sprintf("%s%s", wordsMinutes(t.Minute()), wordsHour(t)),
		wochentage[t.Weekday()],
		// TODO translate to german
		t.Format("2. Jan 2006"),
	}
}
