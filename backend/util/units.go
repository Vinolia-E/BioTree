package util

import (
	"encoding/json"
	"fmt"
	"os"
)

// GetUnitsFromFile reads a JSON file and returns a slice of unique units found in the data.
func GetUnitsFromFile(filePath string) ([]string, error) {
	dataBytes, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("reading file: %w", err)
	}

	var data []DataPoint
	if err := json.Unmarshal(dataBytes, &data); err != nil {
		return nil, fmt.Errorf("unmarshaling JSON: %w", err)
	}

	unitSet := make(map[string]struct{})
	for _, dp := range data {
		unitSet[dp.Unit] = struct{}{}
	}

	// Convert map keys to slice
	var units []string
	for unit := range unitSet {
		units = append(units, unit)
	}

	return units, nil
}

// GetDataByUnitFromFile returns a JSON string of all entries from the file with the specified unit.
func GetDataByUnitFromFile(filePath string, unit string) (string, error) {
	dataBytes, err := os.ReadFile(filePath)
	if err != nil {
		return "", fmt.Errorf("reading file: %w", err)
	}

	var data []DataPoint
	if err := json.Unmarshal(dataBytes, &data); err != nil {
		return "", fmt.Errorf("unmarshaling JSON: %w", err)
	}

	var filtered []DataPoint
	for _, dp := range data {
		if dp.Unit == unit {
			filtered = append(filtered, dp)
		}
	}

	result, err := json.MarshalIndent(filtered, "", "  ")
	if err != nil {
		return "", fmt.Errorf("marshaling result: %w", err)
	}

	return string(result), nil
}
