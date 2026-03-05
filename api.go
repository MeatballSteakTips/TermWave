package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

type Station struct {
	Name    string `json:"name"`
	URL     string `json:"url"`
	Tags    string `json:"tags"`
	Country string `json:"country"`
	Image   string `json:"facvicon"`
}

func StationSearch(searchTerm string) ([]Station, error) {
	safeQuery := url.QueryEscape(searchTerm)
	endPoint := fmt.Sprintf("http://de1.api.radio-browser.info/json/stations/search?name=%s&limit=10", safeQuery)

	resp, err := http.Get(endPoint)
	if err != nil {
		return nil, fmt.Errorf("Request failed: %w", err)
	}

	defer resp.Body.Close() //Defer: Run before function is done

	var stations []Station
	err = json.NewDecoder(resp.Body).Decode(&stations)
	if err != nil {
		return nil, fmt.Errorf("Failed to parse JSON: %w", err)
	}
	return stations, nil
}

