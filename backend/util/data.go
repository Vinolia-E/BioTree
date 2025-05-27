package util

import (
	"regexp"
	"strconv"
)

type DataPoint struct {
	Value float64 `json:"value"`
	Unit  string  `json:"unit"`
}

/*
GetData extracts numerical values and their associated measurement units from the input text.

It searches the text using a regular expression for patterns that match a number
(including optional decimal points and negative signs) followed by an optional unit.

Supported units include:

µg/m³, ppm, °C, °F, mm, in, ha, vehicles/hr, count/month, permits, vehicles

If a unit is not found next to a number, the unit will be recorded as "(none)".

Parameters:
- text: a string containing text with embedded data values.

Returns:
- []DataPoint: a slice of DataPoint structs, each containing a float64 value and a unit string.

Example:
  input := "Air quality is 23.5 µg/m³, temperature is 30.2 °C, and rainfall is 5 mm"
  output := []DataPoint{
      {Value: 23.5, Unit: "µg/m³"},
      {Value: 30.2, Unit: "°C"},
      {Value: 5.0, Unit: "mm"},
  } */

func GetData(text string) []DataPoint {
	var results []DataPoint
	pattern := regexp.MustCompile(`(?i)(-?\d+(?:\.\d+)?)\s*(µg/m³|ppm|°C|°F|mm|in|ha|vehicles/hr|count/month|permits|vehicles)?`)
	matches := pattern.FindAllStringSubmatch(text, -1)

	for _, match := range matches {
		value, err := strconv.ParseFloat(match[1], 64)
		if err != nil {
			continue
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
