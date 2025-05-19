package svgchart

import "sort"

// Chart dimensions and styling constants
const (
	Width           = 600
	Height          = 400
	PaddingTop      = 40
	PaddingLeft     = 60
	PaddingRight    = 40
	PaddingBottom   = 60
	GridColor       = "#dddddd"
	LineColor       = "#3366cc"
	AxisColor       = "#333333"
	TextColor       = "#333333"
	BackgroundColor = "#ffffff"
	FontSize        = 12
	PointRadius     = 3
)

// Point represents a data point in the chart
type Point struct {
	X float64
	Y float64
}

// LineChartFromPairs generates an SVG line chart from a slice of x,y pairs
func LineChartFromPairs(data []Point, xLabel, yLabel string) string {
	if len(data) == 0 {
		// return generateEmptyChart(xLabel, yLabel)
		return ""
	}

	// Sort data by X value
	sort.Slice(data, func(i, j int) bool {
		return data[i].X < data[j].X
	})
	return ""
}
