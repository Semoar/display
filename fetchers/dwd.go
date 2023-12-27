package fetchers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type Weather struct {
	Date          time.Time
	TempMin       float32
	TempMax       float32
	Rain          float32
	RainHourly    []float32
	WindSpeed     float32
	WindDirection float32
	Sunshine      float32
}

func (w Weather) String() string {
	d := w.Date.Format("2.1.")
	return fmt.Sprintf("%s\t%.1f° | %.1f°\t%.1fmm", d, w.TempMin, w.TempMax, w.Rain)
}

type DWDResponse struct {
	Foo DWDFoo `json:"10727"`
}
type DWDFoo struct {
	Forecast1 DWDForecast `json:"forecast1"`
	Days      []DWDDay    `json:"days"`
}
type DWDForecast struct {
	//Start int //actually a time from epoch
	TimeStep           int
	PrecipitationTotal []int
}
type DWDDay struct {
	Date           string `json:"dayDate"`
	TemperatureMin int
	TemperatureMax int
	Precipitation  int
	WindSpeed      int
	WindDirection  int
	Sunshine       int
}

func DWD() []Weather {
	// Forecast for Karlsruhe
	stationID := 10727
	//stationID := 10731 // Rheinstetten, also need to adjust Foo above
	url := fmt.Sprintf("https://dwd.api.proxy.bund.dev/v30/stationOverviewExtended?stationIds=%d", stationID)
	resp, err := http.Get(url)
	if err != nil {
		return []Weather{}
	}
	defer resp.Body.Close()

	byteValue, err := io.ReadAll(resp.Body)
	if err != nil {
		return []Weather{}
	}

	var data DWDResponse
	err = json.Unmarshal(byteValue, &data)
	if err != nil {
		fmt.Printf("could not unmarshal json: %s\n", err)
		return []Weather{}
	}

	fmt.Printf("json map: %v\n", data)
	d, err := data.toGolang()
	return d
}

func (k DWDResponse) toGolang() ([]Weather, error) {
	result := []Weather{}
	for _, d := range k.Foo.Days {
		t, err := time.Parse("2006-01-02", d.Date)
		if err != nil {
			return []Weather{}, err
		}
		result = append(result, Weather{
			Date:          t,
			TempMin:       float32(d.TemperatureMin) / 10,
			TempMax:       float32(d.TemperatureMax) / 10,
			Rain:          float32(d.Precipitation) / 10,
			RainHourly:    make([]float32, 24),
			WindSpeed:     float32(d.WindSpeed) / 10,
			WindDirection: float32(d.WindDirection) / 10,
			Sunshine:      float32(d.Sunshine) / 600, // I guess its tenth minutes
		})
	}
	// TODO parse time, align ... For now assume that it always starts at 0:00 of the first day
	for i, r := range k.Foo.Forecast1.PrecipitationTotal {
		// TODO validate time step to be 3600?
		day := i / 24
		hour := i % 24
		// 32767 is used for values in the past; ignore those
		if r == 32767 {
			continue
		}
		result[day].RainHourly[hour] = float32(r) / 10
	}
	fmt.Printf("result: %v\n", result)
	return result, nil
}
