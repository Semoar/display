package fetchers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

type Weather struct {
	Date    time.Time
	TempMin float32
	TempMax float32
	Rain    float32
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

	byteValue, err := ioutil.ReadAll(resp.Body)
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
			Date:    t,
			TempMin: float32(d.TemperatureMin) / 10,
			TempMax: float32(d.TemperatureMax) / 10,
			Rain:    float32(d.Precipitation) / 10,
		})
	}
	fmt.Printf("result: %v\n", result)
	return result, nil
}
