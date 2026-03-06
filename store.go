package main

import (
	"encoding/json"
	"os"
	"path/filepath"
)

func getConfigPath() (string, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return "Error: ", err
	}

	appDir := filepath.Join(configDir, "termwave")
	err = os.MkdirAll(appDir, 0755)
	if err != nil {
		return "Error: ", err
	}

	return filepath.Join(appDir, "saved_stations.json"), nil
}

func saveStations(stations []Station) error {
	path, err := getConfigPath()
	if err != nil {
		return err
	}

	data, err := json.MarshalIndent(stations, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0644)
}

func loadStations() ([]Station, error) {
	path, err := getConfigPath()
	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return []Station{}, nil
		}
		return nil, err
	}

	var stations []Station
	err = json.Unmarshal(data, &stations)
	if err != nil {
		return nil, err
	}

	return stations, nil
}
