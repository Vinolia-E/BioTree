package svgchart

import (
	"fmt"
	"sort"
	"strconv"

	"github.com/Vinolia-E/BioTree/backend/util"
)

// DataPoint represents a single data point with a label and value
type DataPoint struct {
	Label string  `json:"label"`
	Value float64 `json:"value"`
	Unit  string  `json:"unit,omitempty"`
}

// ChartData represents the data to be visualized
type ChartData []DataPoint

func ConvertData(data interface{}) (ChartData, error) {
	switch d := data.(type) {
	case map[string]float64:
		return mapToDataPoints(d), nil
	case []util.DataPoint:
		return convertUtilDataPoints(d), nil
	case []Point:
		return pointsToDataPoints(d), nil
	case [][]interface{}:
		return sliceToDataPoints(d)
	default:
		return nil, fmt.Errorf("unsupported data type")
	}
}

func convertUtilDataPoints(data []util.DataPoint) ChartData {
	result := make(ChartData, len(data))
	for i, d := range data {
		result[i] = DataPoint{
			Label:  fmt.Sprintf("%.2f %s", d.Value, d.Unit),
			Value:  d.Value,
			Unit:   d.Unit,
		}
	}
	return result
}

func mapToDataPoints(data map[string]float64) ChartData {
	var points ChartData
	for label, value := range data {
		points = append(points, DataPoint{Label: label, Value: value})
	}
	sort.Slice(points, func(i, j int) bool {
		return points[i].Value < points[j].Value
	})
	return points
}

func pointsToDataPoints(points []Point) ChartData {
	var result ChartData
	for _, p := range points {
		result = append(result, DataPoint{Label: p.X, Value: p.Y})
	}
	return result
}

func sliceToDataPoints(data [][]interface{}) (ChartData, error) {
	var points ChartData
	for _, row := range data {
		if len(row) != 2 {
			return nil, fmt.Errorf("each row must have exactly 2 values")
		}

		// Get label
		label, ok := row[0].(string)
		if !ok {
			label = fmt.Sprintf("%v", row[0])
		}

		// Get value
		var value float64
		switch v := row[1].(type) {
		case float64:
			value = v
		case string:
			var err error
			value, err = strconv.ParseFloat(v, 64)
			if err != nil {
				return nil, fmt.Errorf("invalid value: %v", v)
			}
		default:
			return nil, fmt.Errorf("unsupported value type: %T", v)
		}

		points = append(points, DataPoint{Label: label, Value: value})
	}

	sort.Slice(points, func(i, j int) bool {
		return points[i].Value < points[j].Value
	})

	return points, nil
}

type Point struct {
	X string
	Y float64
}
