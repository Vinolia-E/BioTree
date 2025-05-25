package svgchart

import (
	"fmt"
	"sort"
	"strconv"

	"github.com/Vinolia-E/BioTree/backend/util"
)

// ChartData represents the data to be visualized
type ChartData []util.DataPoint

func ConvertData(data interface{}) (ChartData, error) {
	switch d := data.(type) {
	case map[string]float64:
		return mapToDataPoints(d), nil
	case []util.DataPoint:
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
		points = append(points, util.DataPoint{Unit: label, Value: value})
	}
	sort.Slice(points, func(i, j int) bool {
		return points[i].Unit < points[j].Unit
	})
	return points
}

func pointsToDataPoints(points []Point) ChartData {
	var data ChartData
	for _, p := range points {
		data = append(data, util.DataPoint{Unit: p.X, Value: p.Y})
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

		points = append(points, util.DataPoint{Unit: label, Value: value})
	}
	return points, nil
}

type Point struct {
	X string
	Y float64
}
