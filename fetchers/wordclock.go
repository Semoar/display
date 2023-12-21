package fetchers

import (
	"time"
	"fmt"

	// embed timezones as the Kindle doesn't have timezone infos
	_ "time/tzdata"
)

func wordsMinutes(m int) string {
	// TODO add "kurz vor" and "kurz nach"
	if m >= 8 && m <= 22 {
		return "Viertel "
	} else if m >= 23 && m <= 37 {
		return "Halb "
	} else if m >= 38 && m <= 52 {
		return "Dreiviertel "
	} else {
		return ""
	}
}

func wordsHour(t time.Time) string {
	h := t.Hour()
	if t.Minute() >= 8 {
		h += 1
	}
	h = h % 12
	names := []string{"Zwölf", "Eins", "Zwei", "Drei", "Vier", "Fünf", "Sechs", "Sieben", "Acht", "Neun", "Zehn", "Elf"}
	return names[h]
}

var wochentage = []string{"Sonntag", "Montag", "Dienstag", "Mittwoch", "Donnerstag", "Freitag", "Samstag"}

func WordClock() []string {
	// TODO doesn't work on Kindle (no timedatectl etc)
	// Without location, it's UTC and I would need to add 1 or two depending on summer /winter time
	loc, err := time.LoadLocation("Europe/Berlin")
	if err != nil {
		fmt.Printf("Could not load location, %s", err)
	}
	now := time.Now().In(loc)
	fmt.Printf("%v", now)
	return []string{
		fmt.Sprintf("%s%s", wordsMinutes(now.Minute()), wordsHour(now)),
		wochentage[now.Weekday()],
		// TODO translate to german
		now.Format("2. Jan 2006"),
	}
}