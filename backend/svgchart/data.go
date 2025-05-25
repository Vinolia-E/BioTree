package svgchart

import (
	"fmt"
	"sort"
	"strconv"
)

type DataPoint struct {
	Label string
	Value float64
}

// ChartData represents the data to be visualized
type ChartData []DataPoint

func ConvertData(data interface{}) (ChartData, error) {
	switch d := data.(type) {
	case map[string]float64:
		return mapToDataPoints(d), nil
	case []DataPoint:
		return d, nil
	case []Point:
		return pointsToDataPoints(d), nil
	case [][]interface{}:
		return sliceToDataPoints(d)
	default:
		return nil, fmt.Errorf("unsupported data type")
	}
}

func mapToDataPoints(data map[string]float64) ChartData {
	var points ChartData
	for label, value := range data {
		points = append(points, DataPoint{Label: label, Value: value})
	}
	sort.Slice(points, func(i, j int) bool {
		return points[i].Label < points[j].Label
	})
	return points
}

func pointsToDataPoints(points []Point) ChartData {
	var data ChartData
	for _, p := range points {
		data = append(data, DataPoint{Label: p.X, Value: p.Y})
	}
	return data
}

func sliceToDataPoints(data [][]interface{}) (ChartData, error) {
	var points ChartData
	for _, item := range data {
		if len(item) < 2 {
			continue
		}

		label, ok := item[0].(string)
		if !ok {
			label = fmt.Sprintf("%v", item[0])
		}

		var value float64
		switch v := item[1].(type) {
		case float64:
			value = v
		case int:
			value = float64(v)
		default:
			strVal := fmt.Sprintf("%v", v)
			val, err := strconv.ParseFloat(strVal, 64)
			if err != nil {
				return nil, fmt.Errorf("invalid value type: %v", v)
			}
			value = val
		}

		points = append(points, DataPoint{Label: label, Value: value})
	}
	return points, nil
}

type Point struct {
	X string
	Y float64
}