package fetchers

import (
	"testing"
	"time"
)

func TestWordclock(t *testing.T) {
	type testCase struct {
		name   string
		input  time.Time
		result []string
	}
	tcs := []testCase{
		{
			name:   "Add 1 hour in winter",
			input:  time.Date(2009, time.November, 10, 9, 0, 0, 0, time.UTC),
			result: []string{"zehn", "Dienstag", "10. Nov 2009"},
		},
		{
			name:   "Add 2 hours in summer",
			input:  time.Date(2009, time.July, 10, 9, 0, 0, 0, time.UTC),
			result: []string{"elf", "Freitag", "10. Jul 2009"},
		},
		{
			name:   "Add another hour if we are at 'Viertel' or later",
			input:  time.Date(2009, time.July, 10, 9, 15, 0, 0, time.UTC),
			result: []string{"viertel zwölf", "Freitag", "10. Jul 2009"},
		},
		{
			name:   "Use 'kurz nach'",
			input:  time.Date(2009, time.July, 10, 9, 7, 0, 0, time.UTC),
			result: []string{"kurz nach elf", "Freitag", "10. Jul 2009"},
		},
		{
			name:   "Use 'kurz vor'",
			input:  time.Date(2009, time.July, 10, 9, 8, 0, 0, time.UTC),
			result: []string{"kurz vor viertel zwölf", "Freitag", "10. Jul 2009"},
		},
	}

	/*
		// To print all texts for one hour
		for  i := 0; i < 60; i++ {
			x:= time.Date(2009, time.July, 10, 9, i, 0, 0, time.UTC)
			got := WordClock(x)
			fmt.Printf("Timestamp %v resulted in %v\n", x, got)
		}
	*/

	for _, test := range tcs {
		got := WordClock(test.input)
		if !stringSlicesMatch(test.result, got) {
			t.Fatalf(`Case %q failed: got %v, want match for %#q`, test.name, got, test.result)
		}
	}
}

func stringSlicesMatch(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
