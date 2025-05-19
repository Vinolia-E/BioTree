package chart

import (
	"math"
	"sort"
)

const (
	width     = 400
	height    = 200
	padding   = 40
	axisColor = "#333"
	lineColor = "#40e0d0"
	bgColor   = "#ffffff"
	fontSize  = 10
	gridColor = "#eee"
)

// Point represents a 2D coordinate
type Point struct {
	X float64
	Y float64
}

// GenerateLineChartSVG creates an SVG string for a line chart
func GenerateLineChartSVG(data map[string]float64) string {
	if len(data) == 0 {
		return `<svg xmlns="http://www.w3.org/2000/svg" width="400" height="200"><text x="10" y="20" font-size="14" fill="red">No data available</text></svg>`
	}

	// Sort the keys (categorical x-axis)
	keys := make([]string, 0, len(data))
	for k := range data {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	// Compute min/max for scaling
	minY, maxY := math.Inf(1), math.Inf(-1)
	for _, v := range data {
		if v < minY {
			minY = v
		}
		if v > maxY {
			maxY = v
		}
	}
	if minY == maxY {
		minY -= 1
		maxY += 1
	}

	return ""
}
