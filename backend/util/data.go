package util

import (
	"regexp"
	"strconv"
)

type DataPoint struct {
	Value float64 `json:"value"`
	Unit  string  `json:"unit"`
}

var results []DataPoint

func GetData(text string) []DataPoint {
	pattern := regexp.MustCompile(`(?i)(-?\d+(?:\.\d+)?)\s*(µg/m³|ppm|°C|°F|mm|in|ha|vehicles/hr|count/month|permits|vehicles)?`)
	matches := pattern.FindAllStringSubmatch(text, -1)

	for _, match := range matches {
		value, err := strconv.ParseFloat(match[1], 64)
		if err != nil {
			value = 0
		}
		unit := match[2]
		if unit == "" {
			unit = "(none)"
		}

		results = append(results, DataPoint{
			Value: value,
			Unit:  unit,
		})
	}

	return results
}
