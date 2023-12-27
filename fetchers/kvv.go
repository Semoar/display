package fetchers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"
)

type Departure struct {
	Time      time.Time
	Line      string
	Direction string
	Status    string
}

func (d Departure) String() string {
	// TODO use 'in x min' for x <= 9
	t := d.Time.Format("15:04")
	fmtString := "%s  %s  %s"
	if d.Status == "TRIP_CANCELLED" {
		fmtString = "~~%s  %s  %s~~"
	}
	return fmt.Sprintf(fmtString, t, d.Line, d.Direction)
}

type KVVResponse struct {
	DepartureList []KVVDepartures `json:"departureList"`
}
type KVVDepartures struct {
	DateTime       KVVDateTime    `json:"dateTime"`
	RealDateTime   KVVDateTime    `json:"realDateTime"`
	ServingLine    KVVServingLine `json:"servingLine"`
	RealtimeStatus string         `json:"realtimeStatus"`
}
type KVVServingLine struct {
	Symbol    string
	Direction string
}
type KVVDateTime struct {
	Year   string
	Month  string
	Day    string
	Hour   string
	Minute string
}

func KVV() []Departure {
	// Station "Tivoli". Find the ID by searching with `curl 'https://www.kvv.de/tunnelEfaDirect.php?action=XSLT_STOPFINDER_REQUEST&name_sf=tivoli&outputFormat=JSON&type_sf=any' | jq .`
	stationID := 7000084
	url := fmt.Sprintf("https://projekte.kvv-efa.de/sl3-alone/XSLT_DM_REQUEST?outputFormat=JSON&coordOutputFormat=WGS84[dd.ddddd]&depType=stopEvents&locationServerActive=1&mode=direct&name_dm=%d&type_dm=stop&useOnlyStops=1&useRealtime=1&limit=5", stationID)
	resp, err := http.Get(url)
	if err != nil {
		return []Departure{}
	}
	defer resp.Body.Close()

	byteValue, err := io.ReadAll(resp.Body)
	if err != nil {
		return []Departure{}
	}

	var data KVVResponse
	err = json.Unmarshal(byteValue, &data)
	if err != nil {
		fmt.Printf("could not unmarshal json: %s\n", err)
		return []Departure{}
	}

	fmt.Printf("json map: %v\n", data)
	d, err := data.toGolang()
	return d
}

func (k KVVResponse) toGolang() ([]Departure, error) {
	result := []Departure{}
	for _, d := range k.DepartureList {
		t, err := d.RealDateTime.toGolang()
		if err != nil {
			// Sometime only planned times are available
			t, _ = d.DateTime.toGolang()
		}
		result = append(result, Departure{
			Time:      t,
			Line:      d.ServingLine.Symbol,
			Direction: d.ServingLine.Direction,
			Status:    d.RealtimeStatus,
		})
	}
	fmt.Printf("result: %v\n", result)
	return result, nil
}

func (k KVVDateTime) toGolang() (time.Time, error) {
	year, err := strconv.Atoi(k.Year)
	if err != nil {
		return time.Time{}, err
	}
	month, err := strconv.Atoi(k.Month)
	if err != nil {
		return time.Time{}, err
	}
	day, err := strconv.Atoi(k.Day)
	if err != nil {
		return time.Time{}, err
	}
	hour, err := strconv.Atoi(k.Hour)
	if err != nil {
		return time.Time{}, err
	}
	minute, err := strconv.Atoi(k.Minute)
	if err != nil {
		return time.Time{}, err
	}
	return time.Date(year, time.Month(month), day, hour, minute, 0, 0, time.Local), nil
}
