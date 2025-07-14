package maprun

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

// Result represents a single runner's result from the MapRun API.
type Result struct {
	Score    int
	TimeSecs int
	Name     string
	Date     time.Time
}

type apiResponse struct {
	Results []apiResult `json:"results"`
}

type apiResult struct {
	NetScore              int    `json:"NetScore"`
	TotalTimeSecs         int    `json:"TotalTimeSecs"`
	Firstname             string `json:"Firstname"`
	Surname               string `json:"Surname"`
	TrackStartDateTimeUTC string `json:"TrackStartDateTimeUTC"`
}

// GetResults fetches all results for a given event name.
func GetResults(eventName string) ([]Result, error) {
	apiURL := fmt.Sprintf("https://p.fne.com.au:8886/resultsGetPublicForEventv2?eventName=%s", url.QueryEscape(eventName))
	resp, err := http.Get(apiURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var apiResp apiResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return nil, err
	}

	var results []Result
	for _, r := range apiResp.Results {
		date, _ := time.Parse(time.RFC3339, r.TrackStartDateTimeUTC)
		results = append(results, Result{
			Score: r.NetScore, TimeSecs: r.TotalTimeSecs, Name: fmt.Sprintf("%s %s", r.Firstname, r.Surname), Date: date,
		})
	}
	return results, nil
}
